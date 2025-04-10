package hokm4

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	cards "github.com/arian-nj/master-card/back/internal/card"
	"github.com/arian-nj/master-card/back/internal/server"
	"github.com/arian-nj/master-card/back/internal/socket"
	"github.com/arian-nj/master-card/back/pkg/validator"
)

func (app *ApplicationHokm4) WsUpgradeHandler(w http.ResponseWriter, r *http.Request) {
	user := server.ContextGetAuthenticatedUser(r)

	// Coin Amount
	coin_amount_string := r.URL.Query().Get("coin_amount")
	coin_amount_int, err := strconv.Atoi(coin_amount_string)
	if err != nil {
		app.BadRequest(w, r, err)
		return
	}

	tmpValidator := validator.Validator{}
	tmpValidator.CheckField(coin_amount_string != "", "coin_amount", "not provided")

	tmpValidator.CheckField(err == nil, "coin_amount", "not a valid number")

	tmpValidator.CheckField(validator.In(coin_amount_int, validBetAmounts...), "coin_amount", "this amount is not allowed")

	tmpValidator.CheckField(coin_amount_int <= user.Coin, "coin_amount", "don't have enough coins")

	// if user.Coin >= BET_AMOUNT_ONE {
	// 	tmpValidator.CheckField(coin_amount_int != BET_NO_MONEY, "coin_amount", "you can't play a free game")
	// }

	if tmpValidator.HasErrors() {
		app.FailedValidation(w, r, tmpValidator)
		log.Println(tmpValidator.Errors)
		return
	}

	// Upgrade Connection To Websocket
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
	hplayer := app.CreatePlayer(client, user.ID, coin_amount_int)
	app.RunReadWriteInBackground(hplayer)
	app.Logger.Info("new ws connection established " + strconv.Itoa(coin_amount_int))
	app.BackgroundSocketHandlers(hplayer, coin_amount_int)
}

func (app *ApplicationHokm4) CreatePlayer(client *socket.Client, userID int, coin_amount int) *HumanPlayer {
	var player *HumanPlayer

	activeGame, wasInAGame := app.Lobby.UserGames[userID]

	if wasInAGame {
		for _, p := range activeGame.GetHumanPlayers() {
			if p.UserId == userID {
				player = p
				break
			}
		}
		player.Client = client
		player.IsPlayng = true
	} else { // new game
		player = NewHumanPlayer(userID, client, []cards.Card{}, true)
	}
	return player

}

func (app *ApplicationHokm4) RunReadWriteInBackground(player *HumanPlayer) {
	player.Client.State = socket.OPEN
	app.BackgroundTask(func() { // read messages
		defer player.Client.Close()
		defer app.Logger.Error(fmt.Sprintf("read disconnect %d", player.UserId))
		defer func() {
			if player != nil {
				player.IsPlayng = false
			}
		}()

		err := player.Client.ReadMessage(app.Logger, player.Client.CancelCtx)
		if err != nil {
			app.ReportError(err)
			return
		}
	})

	app.BackgroundTask(func() { // write messages
		defer func() {
			if player != nil {
				player.IsPlayng = false
			}
		}()

		defer func() {
			err := player.Client.Close()
			if err != nil {
				app.Logger.Error(err.Error())
			} else {
				app.Logger.Error(fmt.Sprintf("write disconnect %d", player.UserId))
			}
		}()

		err := player.Client.WriteMessage(app.Logger, player.Client.CancelCtx)
		if err != nil {
			app.ReportError(err)
			return
		}
	})
}
