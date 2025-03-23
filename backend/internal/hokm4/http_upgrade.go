package hokm4

import (
	"context"
	"fmt"
	"net/http"

	cards "github.com/arian-nj/master-card/back/internal/card"
	"github.com/arian-nj/master-card/back/internal/server"
	"github.com/arian-nj/master-card/back/internal/socket"
)

func (app *ApplicationHokm4) WsUpgradeHandler(w http.ResponseWriter, r *http.Request) {
	user := server.ContextGetAuthenticatedUser(r)

	conn, err := socket.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		app.Logger.Error(err.Error())
		return
	}

	client := socket.NewClient(conn)
	app.Logger.Debug("new ws connection established")

	activeGame, wasInAGame := app.Lobby.UserGames[user.ID]

	var player *Player

	if wasInAGame {
		for _, p := range activeGame.Players {
			if p.UserId == user.ID {
				player = p
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
		player = NewPlayer(user.ID, client, []cards.Card{}, false)
	}

	ctx, cancel := context.WithCancel(context.Background())
	player.CancelCtx = ctx
	player.Cancel = cancel

	if !wasInAGame { // wait for make match request event
		app.BackgroundTask(func() error {
			for {
				select {
				case new_event := <-client.NewEvents:
					if new_event.Type == socket.TypeMakeMatch {
						app.Lobby.Queue <- player
						return nil
					}
				case <-player.CancelCtx.Done():
					if player != nil {
						player.IsPlayng = false
					}
					return nil
				}
			}
		})
	}
	app.BackgroundTask( // read messages
		func() error {
			defer player.Cancel()
			defer app.Logger.Error(fmt.Sprintf("read disconnect %d", user.ID))
			defer func() {
				if player != nil {
					player.IsPlayng = false
				}
			}()

			return client.ReadMessage(app.Logger, player.CancelCtx)
		})

	app.BackgroundTask(func() error { // write messages
		defer player.Cancel()
		defer func() {
			if player != nil {
				player.IsPlayng = false
			}
		}()

		defer func() {
			err := client.Close()
			if err != nil {
				app.Logger.Error(err.Error())
			} else {
				app.Logger.Error(fmt.Sprintf("write disconnect %d", user.ID))
			}
		}()

		return client.WriteMessage(app.Logger, player.CancelCtx)
	})

}
