package main

import (
	"context"
	"log"

	"github.com/arian-nj/master-card/back/internal/hokm4"
	"github.com/arian-nj/master-card/back/internal/server"
)

func main() {
	globalStructs, poll, err := server.NewCommonGlobals(4445)
	if err != nil {
		log.Panic(err)
	}
	defer poll.Close()

	err = poll.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	app := hokm4.NewHokm4Application(globalStructs)
	app.Config.HTTPPort = 4445

	chiRouter := app.Hokm4Router()

	app.BackgroundTask(func() {
		app.FilterMatchMkingByCoin()
	})

	app.BackgroundTask(func() {
		app.MatchUsers(app.Lobby.MatchmakingQueueForBetOne, hokm4.BET_AMOUNT_ONE)
	})

	app.BackgroundTask(func() {
		app.MatchUsers(app.Lobby.MatchmakingQueueForBetTwo, hokm4.BET_AMOUNT_TWO)
	})

	err = app.ServeHTTP(chiRouter, app.Config.HTTPPort)
	if err != nil {
		app.Logger.Error(err.Error())
	}

}
