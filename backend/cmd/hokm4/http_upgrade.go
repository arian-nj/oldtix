package main

import (
	"context"
	"fmt"
	"net/http"

	cards "github.com/arian-nj/master-card/back/internal/card"
	"github.com/arian-nj/master-card/back/internal/server"
	"github.com/arian-nj/master-card/back/internal/socket"
)

func (app *ApplicationH2) WsUpgradeHandler(w http.ResponseWriter, r *http.Request) {
	user := server.ContextGetAuthenticatedUser(r)

	conn, err := socket.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		app.Logger.Error(err.Error())
		return
	}
	client := socket.NewClient(conn)
	app.Logger.Debug("new ws connection established")

	ctx, cancel := context.WithCancel(context.Background())

	app.BackgroundTask(
		func() error {
			defer cancel()
			defer app.Logger.Error(fmt.Sprintf("read disconnect %d", user.ID))
			err := client.ReadMessage(app.Logger, ctx)
			if err != nil {
				playinggame, exist := app.lobby.UserGames[user.ID]
				if exist {
					for _, p := range playinggame.Players {
						if p.UserId == user.ID {
							p.IsPlayng = false
						}
					}
				}

			}
			return nil
		})

	app.BackgroundTask(func() error {
		defer cancel()
		defer func() {
			err := client.Close()
			if err != nil {
				app.Logger.Error(err.Error())
			} else {
				app.Logger.Error(fmt.Sprintf("write disconnect %d", user.ID))
			}
		}()
		return client.WriteMessage(app.Logger, ctx)
	})

	activeGame, wasInAGame := app.lobby.UserGames[user.ID]

	if wasInAGame {
		for _, p := range activeGame.Players {
			if p.UserId == user.ID {
				player := p
				player.Client = client
				player.IsPlayng = true
				err := activeGame.SendGameData(socket.TypeRejoinMatch, p)
				if err != nil {
					app.ServerError(w, r, err)
					return
				}
				activeGame.BackgroundSocketHandlers(player)

				// p.AddToEgress(e)
			}
		}
	} else { // new game
		player := NewPlayer(user.ID, client, []cards.Card{}, false)

		app.BackgroundTask(func() error {
			// defer app.Logger.Info("Make Match Waiter Ended")
			for {
				select {
				case new_event := <-client.NewEvents:
					if new_event.Type == socket.TypeMakeMatch {
						app.lobby.Queue <- player
						return nil
					}
				case <-ctx.Done():
					return nil
				}
			}
		})
	}

}
