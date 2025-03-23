package main

import (
	"log"
	"sync"

	"github.com/arian-nj/master-card/back/internal/server"
)

type ApplicationH2 struct {
	*server.CommonGlobals
	lobby *Lobby
}

func main() {
	globalStructs, poll, err := server.NewCommonGlobals()
	if err != nil {
		log.Fatal(err)
	}
	defer poll.Close()

	app := ApplicationH2{
		CommonGlobals: globalStructs,
		lobby: &Lobby{
			Queue: make(chan *Player),
			// Games: make(map[int64]*GameState),
			UserGames: map[int64]*GameState{},
			Mu:        sync.Mutex{},
		},
	}
	app.Config.HTTPPort = 4445

	chiM := app.wsHokm2Router()
	// hashedPassword, err := password.Hash("arian123")
	// if err != nil {
	// 	app.Logger.Error(err.Error())
	// 	return
	// }
	// for i := range 50 {
	// 	app.Queries.InsertPerson(context.Background(), sqldb.InsertPersonParams{
	// 		Username:       "arian" + strconv.Itoa(i),
	// 		HashedPassword: hashedPassword,
	// 	})
	// }

	app.BackgroundTask(app.MatchUsers)
	err = app.ServeHTTP(chiM, app.Config.HTTPPort)
	if err != nil {
		app.Logger.Error(err.Error())
	}

}
