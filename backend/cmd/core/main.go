package main

import (
	"context"
	"log"
	"time"

	core_api "github.com/arian-nj/master-card/back/internal/core"
	"github.com/arian-nj/master-card/back/internal/server"
)

func main() {

	globalStructs, poll, err := server.NewCommonGlobals("CORE_HTTP_PORT")
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

	app := core_api.NewApiApplication(globalStructs)
	if app.ReleaseMode == "" {
		globalStructs.Logger.Error("RELEASE_MODE is empty")
		return
	}
	app.Config.HTTPPort = 4444

	// for i := range 50 {
	// 	app.CreateBrandNewPerson("arian"+strconv.Itoa(i), "darian"+strconv.Itoa(i), "arian123")
	// }

	chiRouter := app.CoreRoutes()
	err = app.ServeHTTP(chiRouter, app.Config.HTTPPort)
	if err != nil {
		app.Logger.Error(err.Error())
	}
}
