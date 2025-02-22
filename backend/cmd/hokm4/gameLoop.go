package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	cards "github.com/arian-nj/master-card/back/internal/card"
	"github.com/arian-nj/master-card/back/internal/socket"
)

func (game *GameState) GameInitialize() error {
	for _, p := range game.Players {
		game.BackgroundTask(func() error {
			game.socketHandlers(p)
			return nil
		})
	}

	return nil
}

func (game *GameState) RunGame() error {
	defer func() {
		for _, p := range game.Players {
			p.Client.Conn.Close()
		}
	}()

	// choose hakem
	err := game.GameInitialize()
	if err != nil {
		return err
	}

	for _, p := range game.Players {
		err := game.SendGameData(socket.TypeMatchFound, p)
		if err != nil {
			return err
		}
	}

	for i := range 5 { // run tricks
		if game.TeamOneTricksScore >= 3 || game.TeamTwoTricksScore >= 3 {
			break
		}
		game.CurrentTrick = NewTrick()
		game.Tricks = append(game.Tricks, game.CurrentTrick)

		if i == 0 { // if first trick
			game.CurrentTrick.HakemIndex = rand.Intn(len(game.Players))
		}

		game.CurrentTrick.StarterPlayerIndex = game.CurrentTrick.HakemIndex

		for _, p := range game.Players {
			err := game.SendGameData(socket.TypeNewTrick, p)
			if err != nil {
				return err
			}
		}

		err = game.RunTrick()
		if err != nil {
			return err
		}
	}
	return nil
}
func (game *GameState) RunTrick() error {

	all_cards := cards.NewAllCards()
	all_cards = game.sendCards(5, all_cards, game.Players)

	game.WaitToChooseHokm()          // put hokm in game.CurrentTurn.Hokm
	for _, p := range game.Players { // update hokm data
		game.SendGameData(socket.TypeGameData, p)
	}

	// send rest of cards
	all_cards = game.sendCards(4, all_cards, game.Players)
	game.sendCards(4, all_cards, game.Players)

	for range 13 {
		err := game.RunTurn()
		if err != nil {
			return err
		}
		if game.CurrentTrick.TeamOneTurnScore >= 2 || game.CurrentTrick.TeamTwoTurnScore >= 2 {
			if game.CurrentTrick.TeamOneTurnScore >= 2 {
				game.TeamOneTricksScore += 1
			} else {
				game.TeamTwoTricksScore += 1
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
	return nil
}

func (game *GameState) WaitToChooseHokm() {

	choose_hokm_ticker := time.NewTicker(10 * time.Second)
	defer choose_hokm_ticker.Stop()

	for {

		select {
		case new_game_event := <-game.GameEventsCh:
			if new_game_event.event.Type != socket.TypeHokmChoosed {
				continue
			}
			if new_game_event.Player != game.Players[game.CurrentTrick.HakemIndex] {
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
	time.Sleep(time.Second * 1)

	// game starts
	for _, p := range game.Players {
		err := game.SendGameData(socket.TypeTurnStart, p)
		if err != nil {
			game.Logger.Debug(err.Error())
		}
	}

	// actual game
	to_play_order := []*Player{}

	after_ward := []*Player{}
	for ind, p := range game.Players {
		if ind != game.CurrentTrick.StarterPlayerIndex {
			after_ward = append(after_ward, p)
		} else {
			to_play_order = append(to_play_order, p)
		}
	}
	to_play_order = append(to_play_order, after_ward...)

	for _, p := range to_play_order {
		err := game.PlayerPlayCard(p)
		if err != nil {
			return err
		}
	}

	// Decide who wins Turn
	Winner := game.WhoWins()
	if Winner.Player.TeamId == TeamOne {
		game.CurrentTrick.TeamOneTurnScore += 1
	} else {
		game.CurrentTrick.TeamTwoTurnScore += 1
	}

	for winnerIndex, winnerPlayer := range game.Players {
		if winnerPlayer.UserId == Winner.Player.UserId {
			game.CurrentTrick.StarterPlayerIndex = winnerIndex
		}
	}

	// End The Turn
	for _, p := range game.Players {
		err := game.SendGameData(socket.TypeTurnEnd, p)
		if err != nil {
			game.Logger.Debug(err.Error())
		}
	}

	return nil
}
func (game *GameState) PlayerPlayCard(player *Player) error {
	player.Client.Egres <- *socket.NewEvent(socket.TypeYourTurn, socket.EventMessage(""))
	NewTicker := time.NewTicker(time.Second * 60)
	var card_played cards.Card
	var cardIndex int = -1

OuterLoop:
	for {
		select {
		case new_game_event := <-game.GameEventsCh:
			if new_game_event.event.Type != socket.TypePlayTurn {
				// app.Logger.Debug("not same type")
				continue
			}
			if new_game_event.Player.UserId != player.UserId {
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
			choosen_card_by_bot := ""
			player.Client.Egres <- *socket.NewEvent(socket.TypePlayTimeout, socket.EventMessage(choosen_card_by_bot))
			return fmt.Errorf("time out")
		}
	}
	currentTurn := game.CurrentTrick.CurrentTurn
	new_card_player := &PlayerCardPlayed{
		Player: player,
		Card:   player.Cards[cardIndex],
	}

	currentTurn.CardsPlayed = append(currentTurn.CardsPlayed, new_card_player)
	player.Cards = append(player.Cards[:cardIndex], player.Cards[cardIndex+1:]...)

	played_data_byte, err := json.Marshal(new_card_player)
	if err != nil {
		return err
	}

	for _, otherp := range game.Players {
		if otherp.UserId != player.UserId {
			player.Client.Egres <- *socket.NewEvent(socket.TypeTurnPlayed, socket.EventMessage(played_data_byte))
		}
	}
	return nil
}

// events come here if not used go to GameEventCh
func (game *GameState) socketHandlers(p *Player) {
	for {
		new_event := <-p.Client.NewEvents
		if new_event.Type == socket.TypeGetData {
			game.SendGameData(socket.TypeGameData, p)
		} else {
			game.GameEventsCh <- NewGameEvent(&new_event, p)
		}
	}
}
