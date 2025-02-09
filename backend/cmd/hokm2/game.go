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

		err = app.RunTurn(game)
		if err != nil {
			return err
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
		p.Client.Egres <- *socket.NewEvent(socket.TypeYourTurn, socket.EventMessage(""))
		NewTicker := time.NewTicker(time.Second * 10)
		for {
			select {
			case new_game_event := <-game.GameEventsCh:
				if new_game_event.event.Type != socket.TypePlayTurn {
					continue
				}
				if new_game_event.Player.UserId != p.UserId {
					continue
				}
				if new_game_event.event.Data == nil {
					continue
				}

				var card_played cards.Card
				err := json.Unmarshal([]byte(*new_game_event.event.Data), &card_played)
				if err != nil {
					continue
				}

			case <-NewTicker.C:
				NewTicker.Stop()

			}
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
	for _, p := range players {
		var randomCards []cards.Card
		var err error
		randomCards, all_cards, err = cards.GiveRandomCards(number, all_cards)
		if err != nil {
			app.Logger.Error(err.Error())
		}
		p.Cards = append(p.Cards, randomCards...)

		var output struct {
			NewCards []cards.Card `json:"cards"`
		}
		output.NewCards = randomCards
		data_byte, err := json.Marshal(output)
		if err != nil {
			app.Logger.Error(err.Error())
		}
		p.Client.Egres <- *socket.NewEvent(socket.TypeNewCard, socket.EventMessage(data_byte))
	}
	return all_cards

}
