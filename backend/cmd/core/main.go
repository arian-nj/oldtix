package main

import (
	"log"
	"os"

	"github.com/arian-nj/master-card/back/internal/server"
)

type Application struct {
	*server.CommonGlobals
	ReleaseMode string
}

func main() {
	globalStructs, poll, err := server.NewCommonGlobals()
	if err != nil {
		log.Fatal(err)
	}
	defer poll.Close()

	app := Application{
		CommonGlobals: globalStructs,
		ReleaseMode:   os.Getenv("RELEASE_MODE"),
	}
	if app.ReleaseMode == "" {
		globalStructs.Logger.Error("RELEASE_MODE is empty")
		return
	}
	app.Config.HTTPPort = 4444

	chiM := app.profileRoutes()
	err = app.ServeHTTP(chiM, app.Config.HTTPPort)
	if err != nil {
		app.Logger.Error(err.Error())
	}
}
