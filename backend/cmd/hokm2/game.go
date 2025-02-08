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

func (app *ApplicationH2) RunGame(game *GameState) error {
	// choose hakem
	game.HakemIndex = rand.Intn(len(game.Players))
	// game.Current = int(game.Hakem)

	for _, p := range game.Players {
		app.BackgroundTask(func() error {
			app.socketHandlers(game, p)
			return nil
		})
	}
	for _, p := range game.Players {
		err := app.SendGameData(game, p)
		if err != nil {
			return err
		}
	}

	all_cards := cards.NewAllCards()

	// give 5 car to all players
	all_cards = app.giveCards(5, all_cards, game.Players)
	_ = all_cards

	choose_hokm_ticker := time.NewTicker(10 * time.Second)

	// new Turn
	game.CurrentTurn = NewTurn(game.HakemIndex)

OuterLoop:
	for {

		select {
		case new_game_event := <-game.GameEventsCh:
			if new_game_event.event.Type != socket.TypeHokmChoosed {
				continue
			}
			if new_game_event.Player != game.Players[game.HakemIndex] {
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
			}
			new_hokm := cards.Suite(hokm_int)
			game.CurrentTurn.Hokm = new_hokm
			app.Logger.Info(fmt.Sprintf("new hokm is choosed by hakem %d ", hokm_int))
			break OuterLoop
		case <-choose_hokm_ticker.C:
			choose_hokm_ticker.Stop()
			rand_index := rand.Intn(4)
			new_hokm := cards.AllSuits[rand_index]
			game.CurrentTurn.Hokm = new_hokm
			app.Logger.Info(fmt.Sprintf("new hokm is choosed by server %d ", int(new_hokm)))
			break OuterLoop
		}
	}
	for _, p := range game.Players {
		app.SendGameData(game, p)
	}

	// give rest of cards
	all_cards = app.giveCards(4, all_cards, game.Players)
	app.giveCards(4, all_cards, game.Players)
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
	starter_player_index := game.HakemIndex
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
				_ = new_game_event
			case <-NewTicker.C:
				NewTicker.Stop()
			}
		}
	}

	return nil
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

func (app *ApplicationH2) giveCards(number int, all_cards []cards.Card, players []*Player) []cards.Card {
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
