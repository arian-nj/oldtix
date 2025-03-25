package hokm4

import (
	"context"
	"fmt"
	"net/http"

	cards "github.com/arian-nj/master-card/back/internal/card"
	"github.com/arian-nj/master-card/back/internal/server"
	"github.com/arian-nj/master-card/back/internal/socket"
	"github.com/arian-nj/master-card/back/sqldb"
)

func (app *ApplicationHokm4) WsUpgradeHandler(w http.ResponseWriter, r *http.Request) {
	user := server.ContextGetAuthenticatedUser(r)

	conn, err := socket.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		app.Logger.Error(err.Error())
		return
	}

	client, err := socket.NewClient(conn)
	if err != nil {
		app.ReportServerError(r, err)
		return
	}
	app.Logger.Debug("new ws connection established")

	activeGame, wasInAGame := app.Lobby.UserGames[user.ID]

	var player *HumanPlayer

	if wasInAGame {
		for _, p := range activeGame.GetHumanPlayers() {
			if p.UserId == user.ID {
				player = p
				break
			}
		}
		player.Client = client
		player.IsPlayng = true
	} else { // new game
		player = NewHumanPlayer(user.ID, client, []cards.Card{}, true)
	}

	ctx, cancel := context.WithCancel(context.Background())
	player.Client.CancelCtx = ctx
	player.Client.Cancel = cancel

	app.ReadWriteClientInBackground(player, user)

	if !wasInAGame { // wait for make match request event
		app.BackgroundTask(func() error {
			for {
				select {
				case new_event := <-client.NewEvents:
					if new_event.Type == socket.TypeMakeMatch {
						app.Lobby.Queue <- player
						return nil
					}
				case <-player.Client.CancelCtx.Done():
					return nil
				}
			}
		})
	} else {
		err := activeGame.SendGameData(socket.TypeRejoinMatch, player)
		if err != nil {
			app.ServerError(w, r, err)
			return
		}
		player.BackgroundSocketHandlers(activeGame)

	}
}

func (app *ApplicationHokm4) ReadWriteClientInBackground(player *HumanPlayer, user *sqldb.Person) {
	player.Client.State = socket.OPEN
	app.BackgroundTask(func() error { // read messages
		defer player.Client.Cancel()
		defer app.Logger.Error(fmt.Sprintf("read disconnect %d", user.ID))
		defer func() {
			if player != nil {
				player.IsPlayng = false
			}
		}()

		return player.Client.ReadMessage(app.Logger, player.Client.CancelCtx)
	})

	app.BackgroundTask(func() error { // write messages
		defer func() {
			if player != nil {
				player.IsPlayng = false
			}
		}()

		defer func() {
			player.Client.Cancel()
			err := player.Client.Close()
			if err != nil {
				app.Logger.Error(err.Error())
			} else {
				app.Logger.Error(fmt.Sprintf("write disconnect %d", user.ID))
			}
		}()

		return player.Client.WriteMessage(app.Logger, player.Client.CancelCtx)
	})
}
