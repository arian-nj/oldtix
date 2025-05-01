package hokm4

import (
	"fmt"
	"net/http"
	"strconv"

	cards "github.com/arian-nj/master-card/back/internal/card"
	"github.com/arian-nj/master-card/back/internal/server"
	"github.com/arian-nj/master-card/back/internal/socket"
	"github.com/arian-nj/master-card/back/pkg/validator"
	"github.com/arian-nj/master-card/back/sqldb"
)

func (app *ApplicationHokm4) WsUpgradeHandler(w http.ResponseWriter, r *http.Request) {
	user := server.ContextGetAuthenticatedUser(r)
	activeGame, wasInAGame := app.Lobby.UserGames[user.ID]

	if wasInAGame {
		app.RejonGame(w, r, user, activeGame)
	} else {
		app.JoinGame(w, r, user)
	}
}
func (app *ApplicationHokm4) RejonGame(w http.ResponseWriter, r *http.Request, person *sqldb.Person, activeGame *GameState) {
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
	// hplayer := app.CreatePlayer(client, user.ID, coin_amount_int)
	var Hplayer *HumanPlayer

	for _, p := range activeGame.GetHumanPlayers() {
		if p.UserId == person.ID {
			Hplayer = p
			break
		}
	}

	Hplayer.Client = client
	Hplayer.IsPlayng = true

	app.RunReadWriteInBackground(Hplayer)
	app.Logger.Info("new ws connection re-established " + strconv.Itoa(activeGame.BetAmount))

	app.BackgroundSocketHandlers(Hplayer)

}

func (app *ApplicationHokm4) JoinGame(w http.ResponseWriter, r *http.Request, person *sqldb.Person) {
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

	tmpValidator.CheckField(coin_amount_int <= person.Coin, "coin_amount", "don't have enough coins")

	if person.Coin >= BET_AMOUNT_ONE {
		tmpValidator.CheckField(coin_amount_int != BET_NO_MONEY, "coin_amount", "you can't play a free game")
	}

	if tmpValidator.HasErrors() {
		app.FailedValidation(w, r, tmpValidator)
		app.Logger.Error(fmt.Sprintln(tmpValidator.Errors))
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

	Hplayer := NewHumanPlayer(person.ID, client, []cards.Card{}, true)
	Hplayer.BetAmount = coin_amount_int

	app.RunReadWriteInBackground(Hplayer)
	app.Logger.Info("new ws connection established " + strconv.Itoa(coin_amount_int))
	app.BackgroundSocketHandlers(Hplayer)

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
