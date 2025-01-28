package main

import (
	"net/http"

	"github.com/arian-nj/master-card/back/internal/socket"
)

func (app *Application) eventRoutes() {
	app.eventRouter.RegisterEventHandler(socket.TypeStatus, app.BroadCastChatToAll)
	app.eventRouter.RegisterEventHandler(socket.TypeMakeMatch, app.AddUserToMatchMaking)
}

func (app *Application) WsStartHandler(w http.ResponseWriter, r *http.Request) {
	user := contextGetAuthenticatedUser(r)
	if user == nil {
		app.authenticationRequired(w, r)
		return
	}

	conn, err := socket.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		app.logger.Error(err.Error())
		return
	}
	client := socket.NewClient(conn, app.eventRouter)
	client.DefaultHandler = app.WsDefaultHandler
	app.logger.Debug("new ws connection established")

	go client.ReadMessage(app.logger)
	go client.WriteMessage()

	//		p := Player{
	//			UserId: user.ID,
	//			Client: client,
	//			Mode:   Tic2GameMode,
	//		}
	//		p.Client.Egres <- *socket.NewEvent(socket.EventTypeStatus, socket.StatusConnected)
	//		app.ReadyPlayers <- &p
	//	}
}

func (app *Application) WsDefaultHandler(event *socket.Event, client *socket.Client) error {
	app.logger.Info(string(event.Type) + " " + string(*event.Data))
	return nil
}

func (app *Application) AddUserToMatchMaking(event *socket.Event, client *socket.Client) error {
	return nil
}

func (app *Application) BroadCastChatToAll(event *socket.Event, client *socket.Client) error {
	return nil
}
