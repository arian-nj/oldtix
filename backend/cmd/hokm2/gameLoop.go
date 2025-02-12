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

func (app *ApplicationH2) RunGame(game *GameState) error {
	defer func() {
		for _, p := range game.Players {
			p.Client.Conn.Close()
		}
	}()

	// choose hakem
	err := app.GameInitialize(game)
	if err != nil {
		return err
	}

	for range 1 {
		new_trick := NewTrick()
		game.CurrentTrick = new_trick
		game.Tricks = append(game.Tricks, game.CurrentTrick)

		new_trick.HakemIndex = rand.Intn(len(game.Players))
		for _, p := range game.Players {
			err := app.SendGameData(game, p)
			if err != nil {
				return err
			}
		}

		all_cards := cards.NewAllCards()
		all_cards = app.sendCards(5, all_cards, game.Players)

		app.WaitToChooseHokm(game)       // put hokm in game.CurrentTurn.Hokm
		for _, p := range game.Players { // update hokm data
			app.SendGameData(game, p)
		}

		// send rest of cards
		all_cards = app.sendCards(4, all_cards, game.Players)
		app.sendCards(4, all_cards, game.Players)

		for range 10 {
			err = app.RunTurn(game)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func (app *ApplicationH2) GameInitialize(game *GameState) error {
	for _, p := range game.Players {
		app.BackgroundTask(func() error {
			app.socketHandlers(game, p)
			return nil
		})
	}

	return nil
}

func (app *ApplicationH2) WaitToChooseHokm(game *GameState) {

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
				app.Logger.Info("no data")
				continue
			}
			app.Logger.Info("here setting")
			hokm_int, err := strconv.Atoi(string(*hokm_data))
			if err != nil {
				app.Logger.Error(fmt.Sprintf("trying to set %s as hokm", string(*hokm_data)))
				continue
			}
			new_hokm := cards.Suite(hokm_int)
			game.CurrentTrick.Hokm = new_hokm
			app.Logger.Info(fmt.Sprintf("new hokm is choosed by hakem %d ", hokm_int))
			return
		case <-choose_hokm_ticker.C:
			rand_index := rand.Intn(4)
			new_hokm := cards.AllSuits[rand_index]
			game.CurrentTrick.Hokm = new_hokm
			app.Logger.Info(fmt.Sprintf("new hokm is choosed by server %d ", int(new_hokm)))
			return
		}
	}
}

func (app *ApplicationH2) RunTurn(game *GameState) error {
	// new Turn
	game.CurrentTrick.CurrentTurn = NewTurn()
	time.Sleep(time.Second * 1)

	// game starts
	for _, p := range game.Players {
		game_data, err := json.Marshal(game)
		if err != nil {
			return err
		}
		p.Client.Egres <- *socket.NewEvent(socket.TypeTurnStart, socket.EventMessage(game_data))
	}

	to_play_order := []*Player{}
	starter_player_index := game.CurrentTrick.HakemIndex
	after_ward := []*Player{}
	for ind, p := range game.Players {
		if ind != starter_player_index {
			after_ward = append(after_ward, p)
		} else {
			to_play_order = append(to_play_order, p)
		}
	}
	to_play_order = append(to_play_order, after_ward...)

	for _, p := range to_play_order {
		err := app.PlayerDoMove(game, p)
		if err != nil {
			return err
		}
	}

	for _, p := range to_play_order {
		app.SendGameData(game, p)
	}

	return nil
}
func (app *ApplicationH2) PlayerDoMove(game *GameState, player *Player) error {
	player.Client.Egres <- *socket.NewEvent(socket.TypeYourTurn, socket.EventMessage(""))
	NewTicker := time.NewTicker(time.Second * 60)
	var card_played cards.Card
	var cardIndex int

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
				app.Logger.Debug(err.Error())
				continue
			}

			var isValid bool
			cardIndex, isValid = game.ValidateAndDoMove(new_game_event.Player, &card_played)

			app.Logger.Debug(card_played.String())
			if !isValid {
				new_game_event.Player.Client.Egres <- *socket.NewEvent(socket.TypeInvalidPlay, socket.EventMessage(""))
				// app.Logger.Debug("move not valid")
				continue
			}

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
		Card:   &player.Cards[cardIndex],
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
func (app *ApplicationH2) socketHandlers(game *GameState, p *Player) {
	for {
		new_event := <-p.Client.NewEvents
		if new_event.Type == socket.TypeGetData {
			app.SendGameData(game, p)
		} else {
			game.GameEventsCh <- NewGameEvent(&new_event, p)
		}
	}
}

func (app *ApplicationH2) SendGameData(game *GameState, p *Player) error {
	// send game data
	game_data, err := json.Marshal(game)
	if err != nil {
		return err
	}
	p.Client.Egres <- *socket.NewEvent(socket.TypeGameData, socket.EventMessage(game_data))
	return nil
}

func (app *ApplicationH2) sendCards(number int, all_cards []cards.Card, players []*Player) []cards.Card {
	var remaining_cards []cards.Card = all_cards
	fmt.Println(remaining_cards)
	fmt.Println(" ")
	for _, p := range players {
		var randomCards []cards.Card
		var err error
		randomCards, remaining_cards, err = cards.GiveRandomCards(number, remaining_cards)
		fmt.Println(remaining_cards)
		// fmt.Println(randomCards)
		if err != nil {
			app.Logger.Error(err.Error())
		}
		fmt.Println(" ")

		var output struct {
			NewCards []cards.Card `json:"cards"`
		}
		output.NewCards = randomCards
		data_byte, err := json.Marshal(output)
		if err != nil {
			app.Logger.Error(err.Error())
		}
		p.Client.Egres <- *socket.NewEvent(socket.TypeNewCard, socket.EventMessage(data_byte))
		p.Cards = append(p.Cards, randomCards...)
	}
	return remaining_cards

}
