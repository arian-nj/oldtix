package main

import "github.com/arian-nj/master-card/back/internal/socket"

func (app *ApplicationH2) RunGame(game *Game) {
	// notify match found
	for _, p := range game.Players {
		p.Client.Egres <- *socket.NewEvent(socket.TypeMatchFound, socket.EventMessage("yes"))
	}

}
