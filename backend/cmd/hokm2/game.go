package main

import (
	"encoding/json"
	"math/rand"

	cards "github.com/arian-nj/master-card/back/internal/card"
	"github.com/arian-nj/master-card/back/internal/socket"
)

func (app *ApplicationH2) RunGame(game *GameState) {
	// notify match found
	for _, p := range game.Players {
		p.Client.Egres <- *socket.NewEvent(socket.TypeMatchFound, socket.EventMessage("success"))
	}

	// choose hakem
	rand_hakem_index := rand.Int31n(int32(len(game.Players)))
	game.Hakem = game.Players[rand_hakem_index].UserId
	game.Current = int(game.Hakem)

	// send game data
	game_data, err := json.Marshal(game)
	if err != nil {
		app.Logger.Error(err.Error())
	}
	for _, p := range game.Players {
		p.Client.Egres <- *socket.NewEvent(socket.TypeGameData, socket.EventMessage(game_data))
	}

	// all_cards := cards.NewAllCards()

	// // give 5 car to all players
	// all_cards = app.giveCards(5, all_cards, game.Players)
	// // new cards
	// // new cards
	// time.Sleep(time.Second * 1)
	// all_cards = app.giveCards(4, all_cards, game.Players)
	// // new cards
	// time.Sleep(time.Second * 2)
	// app.giveCards(4, all_cards, game.Players)

}

// func (app *ApplicationH2) GameHandlers(game *GameState) {
// 	for {
// 	}
// }

func (app *ApplicationH2) giveCards(number int, all_cards []cards.Card, players []*Player) []cards.Card {
	for _, p := range players {
		var randomCards []cards.Card
		var err error
		randomCards, all_cards, err = cards.GiveRandomCards(number, all_cards)
		if err != nil {
			app.Logger.Error(err.Error())
		}

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
