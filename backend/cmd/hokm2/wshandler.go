package main

import (
	cards "github.com/arian-nj/master-card/back/internal/card"
	"github.com/arian-nj/master-card/back/internal/socket"
)

func (app *ApplicationH2) WsDefaultHandler(event *socket.Event, client *socket.Client) error {
	app.Logger.Info(string(event.Type) + " " + string(*event.Data))
	return nil
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
