package main

import (
	"log"

	"github.com/arian-nj/master-card/back/internal/hokm4"
	"github.com/arian-nj/master-card/back/internal/server"
)

func main() {
	globalStructs, poll, err := server.NewCommonGlobals()
	if err != nil {
		log.Fatal(err)
	}
	defer poll.Close()

	app := hokm4.NewHokm4Application(globalStructs)
	app.Config.HTTPPort = 4445

	chiRouter := app.Hokm4Router()

	app.BackgroundTask(func() error {
		return app.MatchUsers(app.Lobby.MatchmakingQueueFor0, 0)
	})

	app.BackgroundTask(func() error {
		return app.MatchUsers(app.Lobby.MatchmakingQueueFor50, 50)
	})

	err = app.ServeHTTP(chiRouter, app.Config.HTTPPort)
	if err != nil {
		app.Logger.Error(err.Error())
	}

}
