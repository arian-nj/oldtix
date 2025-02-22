package main

import (
	"log"
	"os"

	"github.com/arian-nj/master-card/back/internal/server"
	"github.com/arian-nj/master-card/back/internal/version"
)

type Application struct {
	*server.CommonGlobals
	Version *version.Version
}

func main() {
	gb, err := server.NewCommonGlobals()
	if err != nil {
		log.Fatal(err)
	}

	app := Application{
		CommonGlobals: gb,
	}

	v, err := version.NewVersion(os.Getenv("VERSION"))
	if err != nil {
		app.Logger.Error(err.Error())
		return
	}
	app.Version = v

	chiM := app.profileRoutes()
	err = gb.ServeHTTP(chiM, 4444)
	if err != nil {
		app.Logger.Error(err.Error())
	}

}
