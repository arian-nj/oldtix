package main

import (
	cards "github.com/arian-nj/master-card/back/internal/card"
	"github.com/arian-nj/master-card/back/internal/socket"
	"github.com/arian-nj/master-card/back/internal/utils"
)

func (app *ApplicationH2) MatchUsers() error {
	for {
		p1 := <-app.lobby.Queue
		p2 := <-app.lobby.Queue

		p1.TeamId = TeamOne
		p2.TeamId = TeamTwo

		game := GameState{
			ID:      utils.GenerateRandomString(16),
			Players: []*Player{p1, p2},
			// Current:      0,
			GameEventsCh: make(chan *GameEvent),
		}
		app.lobby.Mu.Lock()
		app.lobby.Games[game.ID] = &game
		app.lobby.Mu.Unlock()

		app.BackgroundTask(func() error {
			err := app.RunGame(&game)
			app.Logger.Error("Game Loop Ended")
			if err != nil {
				app.ReportError(err)
			}
			return err
		})
	}
}

func (app *ApplicationH2) AddUserToMatchMaking(event *socket.Event, client *socket.Client) error {
	p := Player{
		UserId: client.User.ID,
		Client: client,
		Cards:  []cards.Card{},
	}
	app.lobby.Queue <- &p
	return nil
}
