package main

import "net/http"

func (app *Application) WsStartHandler(w http.ResponseWriter, r *http.Request) {
}

// 	user := contextGetAuthenticatedUser(r)
// 	if user == nil {
// 		app.authenticationRequired(w, r)
// 		return
// 	}

// 	conn, err := socket.Upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	client := socket.NewClient(conn, app.eventRouter)

// 	log.Println("new ws connection established")
// 	go client.ReadMessage()
// 	go client.WriteMessage()

// 	p := Player{
// 		UserId: user.ID,
// 		Client: client,
// 		Mode:   Tic2GameMode,
// 	}
// 	p.Client.Egres <- *socket.NewEvent(socket.EventTypeStatus, socket.StatusConnected)
// 	app.ReadyPlayers <- &p
// }

// func (app *Application) BroadCastChatToAll(event *socket.Event, client *socket.Client) error {
// 	// app.ActiveGameParies
// 	return nil
// }
