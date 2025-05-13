package main

import (
	"context"
	"log"
	"time"

	"github.com/arian-nj/master-card/back/internal/hokm4"
	"github.com/arian-nj/master-card/back/internal/server"
)

func main() {
	globalStructs, poll, err := server.NewCommonGlobals("HOKM4_HTTP_PORT")
	if err != nil {
		log.Panic(err)
	}
	defer poll.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	err = poll.Ping(ctx)
	if err != nil {
		panic(err)
	}
	cancel()

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
