package main

import (
	"log"

	core_api "github.com/arian-nj/master-card/back/internal/core"
	"github.com/arian-nj/master-card/back/internal/server"
)

func main() {
	globalStructs, poll, err := server.NewCommonGlobals()
	if err != nil {
		log.Panic(err)
	}
	defer poll.Close()

	app := core_api.NewApiApplication(globalStructs)
	if app.ReleaseMode == "" {
		globalStructs.Logger.Error("RELEASE_MODE is empty")
		return
	}
	app.Config.HTTPPort = 4444

	chiRouter := app.CoreRoutes()
	err = app.ServeHTTP(chiRouter, app.Config.HTTPPort)
	if err != nil {
		app.Logger.Error(err.Error())
	}
}
