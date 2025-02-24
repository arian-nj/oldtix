package main

import (
	"log"

	"github.com/arian-nj/master-card/back/internal/server"
)

type ApplicationH2 struct {
	*server.CommonGlobals
	lobby *Lobby
}

func main() {
	gb, poll, err := server.NewCommonGlobals()
	if err != nil {
		log.Fatal(err)
	}
	defer poll.Close()

	app := ApplicationH2{
		CommonGlobals: gb,
		lobby: &Lobby{
			Queue: make(chan *Player),
			Games: make(map[string]*GameState),
		},
	}

	chiM := app.wsHokm2Router()

	app.BackgroundTask(app.MatchUsers)
	err = gb.ServeHTTP(chiM, 4445)
	if err != nil {
		app.Logger.Error(err.Error())
	}

}
