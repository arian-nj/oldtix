package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	cards "github.com/arian-nj/master-card/back/internal/card"
	"github.com/arian-nj/master-card/back/internal/socket"
	"github.com/arian-nj/master-card/back/sqldb"
	"github.com/jackc/pgx/v5/pgtype"
)

func (app *ApplicationH2) MatchUsers() error {
	for {
		p1 := <-app.lobby.Queue
		p2 := <-app.lobby.Queue
		p1.IsPlayng = true
		p2.IsPlayng = true

		p3 := NewPlayer(0, socket.NewClient(nil), []cards.Card{}, false)
		p4 := NewPlayer(0, socket.NewClient(nil), []cards.Card{}, false)

		p1.TeamId = TeamOne
		p2.TeamId = TeamTwo

		p3.TeamId = TeamOne
		p4.TeamId = TeamTwo

		players := []*Player{p1, p2, p3, p4}

		game, err := app.NewGameState(players)
		if err != nil {
			return err
		}

		app.lobby.Mu.Lock()
		app.lobby.Games[game.ID] = game
		app.lobby.Mu.Unlock()

		app.BackgroundTask(func() error {
			err := game.RunGame()
			app.Logger.Error("Game Loop Ended")
			if err != nil {
				app.ReportError(err)
			}
			return err
		})
	}
}

func (app *ApplicationH2) AddUserToMatchMaking(event *socket.Event, client *socket.Client, userId int64) error {
	p := NewPlayer(userId, client, []cards.Card{}, false)
	app.lobby.Queue <- p
	return nil
}

func (game *GameState) GameInitialize() error {
	for _, player := range game.Players {
		game.BackgroundTask(func() error {
			game.BackgroundSocketHandlers(player)
			return nil
		})
	}
	return nil
}

func (game *GameState) RunGame() error {
	defer func() {
		for _, p := range game.Players {
			p.Client.Close()
		}
	}()

	// choose hakem
	err := game.GameInitialize()
	if err != nil {
		return err
	}

	for _, p := range game.Players {
		game.Logger.Info(p.PlayerUnique)
	}

	for _, p := range game.Players {
		err := game.SendGameData(socket.TypeMatchFound, p)
		if err != nil {
			return err
		}
	}

	for i := range 5 { // run tricks
		err = game.RunTrick(i)
		if err != nil {
			return err
		}
		if game.TeamOneTricksScore >= SETTING_WINNING_TRICK_SCORE || game.TeamTwoTricksScore >= SETTING_WINNING_TRICK_SCORE {
			break
		}
	}
	return nil
}

func (game *GameState) DeclareHakemIndex(trick_number int) int {
	var HakemIndex int
	if trick_number == 0 { // if first trick
		HakemIndex = rand.Intn(len(game.Players))
	} else {
		if game.CurrentTrick.WinnerTeam != game.Players[game.CurrentTrick.HakemIndex].TeamId {
			HakemIndex = game.CurrentTrick.HakemIndex
			if HakemIndex < len(game.Players)-1 {
				HakemIndex += 1
			} else {
				HakemIndex = 0
			}
		}
	}
	return HakemIndex
}

func (game *GameState) RunTrick(trick_number int) error {
	var err error

	HakemIndex := game.DeclareHakemIndex(trick_number)

	game.CurrentTrick, err = game.NewTrick(HakemIndex)
	if err != nil {
		return err
	}

	game.Tricks = append(game.Tricks, game.CurrentTrick)

	game.CurrentTrick.TurnStarterIndex = game.CurrentTrick.HakemIndex

	for _, p := range game.Players {
		err := game.SendGameData(socket.TypeNewTrick, p)
		if err != nil {
			return err
		}
	}

	all_cards := cards.NewAllCards()
	all_cards = game.sendCards(5, all_cards, game.Players)

	game.WaitToChooseHokm()          // put hokm in game.CurrentTurn.Hokm
	for _, p := range game.Players { // update hokm data
		game.SendGameData(socket.TypeGameData, p)
	}
	err = game.Queries.UpdateHokmTrick(context.Background(), sqldb.UpdateHokmTrickParams{
		Hokm:    pgtype.Int4{Int32: int32(game.CurrentTrick.Hokm), Valid: true},
		TrickID: game.CurrentTrick.id,
	})
	if err != nil {
		return err
	}

	// send rest of cards
	all_cards = game.sendCards(4, all_cards, game.Players)
	game.sendCards(4, all_cards, game.Players)

	// time.Sleep(5 * time.Second)
	for range 13 {
		err := game.RunTurn()
		if err != nil {
			return err
		}
		if game.CurrentTrick.TeamOneTurnScore >= SETTING_WINNIG_TURN_SCORE || game.CurrentTrick.TeamTwoTurnScore >= SETTING_WINNIG_TURN_SCORE {
			if game.CurrentTrick.TeamOneTurnScore >= SETTING_WINNIG_TURN_SCORE {
				game.TeamOneTricksScore += 1
				game.CurrentTrick.WinnerTeam = TeamOne
			} else {
				game.TeamTwoTricksScore += 1
				game.CurrentTrick.WinnerTeam = TeamTwo
			}
			break
		}
	}

	// notify winners and end the game
	for _, p := range game.Players { // update hokm data
		err := game.SendGameData(socket.TypeEndTrick, p)
		if err != nil {
			return err
		}
	}

	err = game.Queries.UpdateTrickScores(context.Background(), sqldb.UpdateTrickScoresParams{
		TeamOneTricksScore: int32(game.TeamOneTricksScore),
		TeamTwoTricksScore: int32(game.TeamTwoTricksScore),
		ID:                 game.ID,
	})
	if err != nil {
		return err
	}
	return nil
}

func (game *GameState) WaitToChooseHokm() {

	hakem := game.Players[game.CurrentTrick.HakemIndex]
	game.Logger.Info(fmt.Sprintln("hakem is ", hakem.PlayerUnique, game.CurrentTrick.HakemIndex))
	var choose_hokm_ticker *time.Ticker
	choose_hokm_ticker = time.NewTicker(SETTING_BOT_CHOOSE_HOKM_WAIT)
	if hakem.IsPlayng {
		choose_hokm_ticker = time.NewTicker(SETTING_PLAYER_CHOOSE_HOKM_WAIT)
	}
	defer choose_hokm_ticker.Stop()
	for {

		select {
		case new_game_event := <-game.GameEventsCh:
			if new_game_event.event.Type != socket.TypeHokmChoosed {
				continue
			}
			if new_game_event.Player != hakem {
				continue
			}
			hokm_data := new_game_event.event.Data
			if new_game_event.event.Data == nil {
				game.Logger.Info("no data")
				continue
			}
			game.Logger.Info("here setting")
			hokm_int, err := strconv.Atoi(string(*hokm_data))
			if err != nil {
				game.Logger.Error(fmt.Sprintf("trying to set %s as hokm", string(*hokm_data)))
				continue
			}
			new_hokm := cards.Suite(hokm_int)
			game.CurrentTrick.Hokm = new_hokm
			game.Logger.Info(fmt.Sprintf("new hokm is choosed by hakem %d ", hokm_int))
			return
		case <-choose_hokm_ticker.C:
			rand_index := rand.Intn(4)
			new_hokm := cards.AllSuits[rand_index]
			game.CurrentTrick.Hokm = new_hokm
			game.Logger.Info(fmt.Sprintf("new hokm is choosed by server %d ", int(new_hokm)))
			return
		}
	}
}

func (game *GameState) RunTurn() error {
	// new Turn
	game.CurrentTrick.CurrentTurn = NewTurn()
	time.Sleep(time.Second * 2)

	// game starts
	for _, p := range game.Players {
		err := game.SendGameData(socket.TypeTurnStart, p)
		if err != nil {
			game.Logger.Debug(err.Error())
		}
	}

	// actual game
	to_play_order := []*Player{}
	before_ward := []*Player{}
	after_ward := []*Player{}
	var starterFound = false

	game.Logger.Info(fmt.Sprintln(game.CurrentTrick.TurnStarterIndex))
	for player_index, p := range game.Players {
		game.Logger.Info(p.PlayerUnique)
		if player_index == game.CurrentTrick.TurnStarterIndex {
			starterFound = true
			to_play_order = append(to_play_order, p)
			continue
		}

		if starterFound {
			after_ward = append(after_ward, p)

		} else {
			before_ward = append(before_ward, p)

		}
	}
	to_play_order = append(to_play_order, after_ward...)
	to_play_order = append(to_play_order, before_ward...)

	game.Logger.Info(fmt.Sprintln("len of players", len(game.Players)))
	game.Logger.Info(fmt.Sprintln("len of to play", len(to_play_order)))
	for _, play_order := range to_play_order {
		game.Logger.Info(play_order.PlayerUnique)
	}

	for _, playing_player := range to_play_order {
		cardIndex, err := game.WaitForPlayerToPlayCard(playing_player)
		if err != nil {
			return err
		}
		currentTurn := game.CurrentTrick.CurrentTurn
		new_card_player := NewPlayerCardPlayed(playing_player, playing_player.Cards[cardIndex])

		// Brodcast played card
		b_data, err := json.Marshal(new_card_player)
		if err != nil {
			return err
		}

		turn_played_event := socket.NewEvent(socket.TypeTurnPlayed, socket.EventMessage(b_data))
		for _, player := range game.Players {
			if player != playing_player {
				player.AddToEgress(turn_played_event)
			}
		}

		currentTurn.CardsPlayed = append(currentTurn.CardsPlayed, new_card_player)
		playing_player.Cards = append(playing_player.Cards[:cardIndex], playing_player.Cards[cardIndex+1:]...)
	}

	// Decide who wins Turn
	Winner := game.WhoWins()
	if Winner.Player.TeamId == TeamOne {
		game.CurrentTrick.TeamOneTurnScore += 1
	} else {
		game.CurrentTrick.TeamTwoTurnScore += 1
	}

	for playerIndex, player := range game.Players {
		if player == Winner.Player {
			game.CurrentTrick.TurnStarterIndex = playerIndex
		}
	}

	// End The Turn
	for _, p := range game.Players {
		err := game.SendGameData(socket.TypeTurnEnd, p)
		if err != nil {
			game.Logger.Debug(err.Error())
		}
	}

	data_byte, err := json.Marshal(game.CurrentTrick.CurrentTurn.CardsPlayed)
	if err != nil {
		return err
	}
	_, err = game.Queries.InsertTurn(context.Background(), sqldb.InsertTurnParams{
		Moves:   string(data_byte),
		TrickID: game.CurrentTrick.id,
	})
	if err != nil {
		return err
	}
	err = game.Queries.UpdateTurnScores(context.Background(), sqldb.UpdateTurnScoresParams{
		TeamOneTurnScore: int32(game.CurrentTrick.TeamOneTurnScore),
		TeamTwoTurnScore: int32(game.CurrentTrick.TeamTwoTurnScore),
		TrickID:          game.CurrentTrick.id,
	})
	return err

}

func (game *GameState) WaitForPlayerToPlayCard(playing_player *Player) (cardIndex int, err error) {
	playing_player.AddToEgress(socket.NewEvent(socket.TypeYourTurn, socket.EventMessage("")))
	var NewTicker *time.Ticker

	NewTicker = time.NewTicker(SETTING_BOT_PLAY_WAIT)
	if playing_player.IsPlayng {
		NewTicker = time.NewTicker(SETTING_PLAYER_PLAY_WAIT)
	}
	var card_played cards.Card
	cardIndex = -1

OuterLoop:
	for {
		select {
		case new_game_event := <-game.GameEventsCh:
			if new_game_event.event.Type != socket.TypeTurnPlayed {
				// app.Logger.Debug("not same type")
				continue
			}
			if new_game_event.Player.UserId != playing_player.UserId {
				// app.Logger.Debug("not same user")
				continue
			}
			if new_game_event.event.Data == nil {
				// app.Logger.Debug("data is nil")
				continue
			}

			err := json.Unmarshal([]byte(*new_game_event.event.Data), &card_played)
			if err != nil {
				// app.Logger.Debug("can't marshal")
				game.Logger.Debug(err.Error())
				continue
			}

			newCardIndex, isValid := game.ValidateAndDoMove(new_game_event.Player, &card_played)

			// app.Logger.Debug(card_played.String())
			if !isValid {
				new_game_event.Player.Client.Egres <- *socket.NewEvent(socket.TypeInvalidPlay, socket.EventMessage(""))
				// app.Logger.Debug("move not valid")
				continue
			}
			cardIndex = newCardIndex
			new_game_event.Player.Client.Egres <- *socket.NewEvent(socket.TypeValidPlay, socket.EventMessage(""))
			break OuterLoop

		case <-NewTicker.C:
			NewTicker.Stop()
			cardIndex = game.BotPlayTurn(playing_player)
			choosen_card_by_bot := playing_player.Cards[cardIndex]
			data_byte, err := json.Marshal(choosen_card_by_bot)
			if err != nil {
				return -1, err
			}
			playing_player.AddToEgress(socket.NewEvent(socket.TypePlayTimeout, socket.EventMessage(data_byte)))
			break OuterLoop
		}
	}
	return cardIndex, nil
}

// events come here if not used go to GameEventCh
func (game *GameState) BackgroundSocketHandlers(p *Player) {
	for {
		new_event := <-p.Client.NewEvents
		if new_event.Type == socket.TypeGetData {
			game.SendGameData(socket.TypeGameData, p)
		} else {
			game.GameEventsCh <- NewGameEvent(&new_event, p)
		}
	}
}
