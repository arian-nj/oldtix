package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/arian-nj/master-card/back/internal/server"
	telegrambot "github.com/arian-nj/master-card/back/telegram"
)

func main() {
	globalStructs, poll, err := server.NewCommonGlobals("WEBAPP_HTTP_PORT")
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

	ctx, cancel = signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	app := telegrambot.NewTelegramApplication(globalStructs)

	go func() {
		TgRoutes := app.WebsiteRoutes()
		err := app.ServeHTTP(TgRoutes, app.Config.HTTPPort)
		if err != nil {
			panic(err)
		}
	}()

	err = app.StartTelegramBot(ctx)
	if err != nil {
		panic(err)
	}
}
