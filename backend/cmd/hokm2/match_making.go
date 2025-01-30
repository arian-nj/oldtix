package main

import (
	"github.com/arian-nj/master-card/back/internal/utils"
)

func (app *ApplicationH2) MatchUsers() error {
	for {
		p1 := <-app.lobby.Queue
		p2 := <-app.lobby.Queue
		game := Game{
			ID:      utils.GenerateRandomString(16),
			Players: [2]*Player{p1, p2},
			Current: 0,
		}
		app.lobby.Mu.Lock()
		app.lobby.Games[game.ID] = &game
		app.lobby.Mu.Unlock()

		app.RunGame(&game)
	}
}
