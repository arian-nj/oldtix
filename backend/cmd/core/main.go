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
	gb, poll, err := server.NewCommonGlobals()
	if err != nil {
		log.Fatal(err)
	}
	defer poll.Close()

	app := Application{
		CommonGlobals: gb,
	}
	release_mode := os.Getenv("RELEASE_MODE")
	if release_mode == "" {
		log.Fatal("RELEASE_MODE is empty")
	}
	app.ReleaseMode = release_mode

	chiM := app.profileRoutes()
	err = gb.ServeHTTP(chiM, 4444)
	if err != nil {
		app.Logger.Error(err.Error())
	}

}
