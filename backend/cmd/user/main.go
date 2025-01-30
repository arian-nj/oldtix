package main

import (
	"log"

	"github.com/arian-nj/master-card/back/internal/server"
)

type Application struct {
	*server.CommonGlobals
}

func main() {
	gb, err := server.NewCommonGlobals()
	if err != nil {
		log.Fatal(err)
	}

	app := Application{
		CommonGlobals: gb,
	}

	chiM := app.profileRoutes()

	err = gb.ServeHTTP(chiM, 4444)
	if err != nil {
		app.Logger.Error(err.Error())
	}

}
