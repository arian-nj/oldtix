package main

import (
	"context"
	"net/http"

	"github.com/arian-nj/master-card/back/internal/server"
	"github.com/arian-nj/master-card/back/internal/socket"
)

func (app *ApplicationH2) WsUpgradeHandler(w http.ResponseWriter, r *http.Request) {
	user := server.ContextGetAuthenticatedUser(r)
	if user == nil {
		app.AuthenticationRequired(w, r)
		return
	}
	conn, err := socket.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		app.Logger.Error(err.Error())
		return
	}
	client := socket.NewClient(conn, user)
	app.Logger.Debug("new ws connection established")

	ctx, cancel := context.WithCancel(context.Background())

	app.BackgroundTask(
		func() error {
			defer cancel()
			return client.ReadMessage(app.Logger, ctx)
		})

	app.BackgroundTask(func() error {
		defer cancel()
		defer func() {
			err := client.Conn.Close()
			if err != nil {
				app.Logger.Error(err.Error())
			} else {
				app.Logger.Error("closed")
			}
		}()
		return client.WriteMessage(app.Logger, ctx)
	})

	app.BackgroundTask(func() error {
		defer app.Logger.Info("Make Match Waiter Ended")

		for {
			select {
			case new_event := <-client.NewEvents:
				if new_event.Type == socket.TypeMakeMatch {
					app.AddUserToMatchMaking(&new_event, client)
					return nil
				}
			case <-ctx.Done():
				return nil
			}
		}
	})

	client.Egres <- *socket.NewEvent(socket.TypeStatus, socket.StatusConnected)

}
