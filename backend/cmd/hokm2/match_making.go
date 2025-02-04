package main

import (
	"github.com/arian-nj/master-card/back/internal/utils"
)

func (app *ApplicationH2) MatchUsers() error {
	for {
		p1 := <-app.lobby.Queue
		p2 := <-app.lobby.Queue
		game := GameState{
			ID:      utils.GenerateRandomString(16),
			Players: []*Player{p1, p2},
			Current: 0,
		}
		app.lobby.Mu.Lock()
		app.lobby.Games[game.ID] = &game
		app.lobby.Mu.Unlock()

		app.RunGame(&game)
	}
}
