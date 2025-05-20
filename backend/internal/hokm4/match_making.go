package hokm4

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/arian-nj/master-card/back/internal/socket"
	"github.com/arian-nj/master-card/back/sqldb"
)

var validBetAmounts = []int{BET_NO_MONEY, BET_AMOUNT_ONE, BET_AMOUNT_TWO}

type Lobby struct {
	MatchmakingQueueGlobal chan *MatchmakingTicket

	MatchmakingQueueForBetOne chan *MatchmakingTicket
	MatchmakingQueueForBetTwo chan *MatchmakingTicket

	// MakingMatches map[int]*GameState
	UserGames map[int]*GameState

	Mu sync.Mutex
}

func NewLobby() *Lobby {
	return &Lobby{
		MatchmakingQueueGlobal:    make(chan *MatchmakingTicket),
		MatchmakingQueueForBetOne: make(chan *MatchmakingTicket),
		MatchmakingQueueForBetTwo: make(chan *MatchmakingTicket),
		// Games: make(map[int]*GameState),
		UserGames: map[int]*GameState{},
		Mu:        sync.Mutex{},
	}
}

type MatchmakingTicket struct {
	Player    *HumanPlayer
	BetAmount int
	Timestamp time.Time
}

func NewMatchmakingTicket(hplayer *HumanPlayer) *MatchmakingTicket {
	return &MatchmakingTicket{
		Player:    hplayer,
		Timestamp: time.Now(),
	}
}
func (game *GameState) AddHumanPlayerToGame(player *HumanPlayer) error {
	_, err := game.Queries.InsertGamePlayer(context.Background(), sqldb.InsertGamePlayerParams{
		PlayerID: player.UserId,
		GameID:   game.ID,
		Team:     int(player.TeamId),
	})
	game.Players = append(game.Players, player)
	player.Game = game
	return err
}

func (game *GameState) AddBotPlayerToGame() {
	player := NewBotPlayer()
	game.Players = append(game.Players, player)
}

func (app *ApplicationHokm4) FilterMatchMakingByCoin() {
	for {
		newReq := <-app.Lobby.MatchmakingQueueGlobal
		if newReq.BetAmount == BET_AMOUNT_ONE || newReq.BetAmount == BET_NO_MONEY {
			app.Lobby.MatchmakingQueueForBetOne <- newReq
		} else if newReq.BetAmount == BET_AMOUNT_TWO {
			app.Lobby.MatchmakingQueueForBetTwo <- newReq
		} else {
			app.Logger.Error("can'r filter this amount of coin " + strconv.Itoa(newReq.BetAmount))
		}
	}
}
func (app *ApplicationHokm4) MatchUsers(matchesChan chan *MatchmakingTicket, betting_amount int) {
	MAX_PLAYERS := 4
	for {
		game, err := app.NewGameState(betting_amount)
		if err != nil {
			app.ReportError(err)
			return
		}
		game.BetAmount = betting_amount

		foundHumanPlayers := app.WaitForPlayers(matchesChan)
		for _, humanPlayer := range foundHumanPlayers {
			err := game.AddHumanPlayerToGame(humanPlayer)
			if err != nil {
				app.ReportError(err)
				return
			}
		}
		numberOfBots := MAX_PLAYERS - len(foundHumanPlayers)
		if numberOfBots > 0 {
			for range numberOfBots {
				game.AddBotPlayerToGame()
			}
		}

		game.Players[0].SetTeamID(TeamOne)
		game.Players[1].SetTeamID(TeamTwo)
		game.Players[2].SetTeamID(TeamOne)
		game.Players[3].SetTeamID(TeamTwo)

		app.Lobby.Mu.Lock()
		for _, p := range game.GetHumanPlayers() {
			app.Lobby.UserGames[p.UserId] = game
		}
		app.Lobby.Mu.Unlock()

		for _, hplayer := range game.GetHumanPlayers() {
			err = app.Queries.AddCoinToPerson(context.Background(), sqldb.AddCoinToPersonParams{
				Coin: -1 * hplayer.BetAmount,
				ID:   hplayer.UserId,
			})
			if err != nil {
				app.ReportError(err)
				return
			}
			err = game.SendGameData(TypeMatchFound, hplayer)
			if err != nil {
				app.ReportError(err)
				return
			}
		}

		app.BackgroundTask(func() {
			game.RunGame()

		})

	}
}

func (app *ApplicationHokm4) WaitForPlayers(TicketChan chan *MatchmakingTicket) []*HumanPlayer {
	var foundPlayers []*HumanPlayer
	firstTicket := <-TicketChan
	foundPlayers = append(foundPlayers, firstTicket.Player)
	timer := time.NewTimer(5 * time.Second)

LOOP:
	for {
		select {
		case foundMatchRequest := <-TicketChan:
			foundPl := foundMatchRequest.Player
			if foundPl.Client.State != socket.OPEN {
				continue
			}
			foundPlayers = append(foundPlayers, foundPl)

			for index, p := range foundPlayers { // check others connection
				if p.Client.State != socket.OPEN {
					foundPlayers = append(foundPlayers[:index], foundPlayers[index+1:]...)
				}
			}
			if len(foundPlayers) == 4 {
				break LOOP
			}

		case <-timer.C:
			for index, p := range foundPlayers { // check others connection
				if p.Client.State != socket.OPEN {
					foundPlayers = append(foundPlayers[:index], foundPlayers[index+1:]...)
				}
			}
			break LOOP
		}
	}
	return foundPlayers
}

func (app *ApplicationHokm4) MatchUsers3(matchesChan chan *MatchmakingTicket, betting_amount int) {
	for {
		game, err := app.NewGameState(betting_amount)
		if err != nil {
			app.ReportError(err)
			return
		}
		game.BetAmount = betting_amount

		var foundPlayers []*HumanPlayer

		for len(foundPlayers) < 2 {
			foundMatchRequest := <-matchesChan
			foundPl := foundMatchRequest.Player
			if foundPl.Client.State != socket.OPEN {
				continue
			}
			foundPlayers = append(foundPlayers, foundPl)

			for index, p := range foundPlayers {
				if p.Client.State != socket.OPEN {
					foundPlayers = append(foundPlayers[:index], foundPlayers[index+1:]...)
				}
			}
		}

		for _, p := range foundPlayers {
			p.Game = game
			err = game.AddHumanPlayerToGame(p)
			if err != nil {
				app.ReportError(err)
				return
			}
		}

		// foundPlayer.IsPlayng = true

		game.AddBotPlayerToGame()
		game.AddBotPlayerToGame()

		game.Players[0].SetTeamID(TeamOne)
		game.Players[1].SetTeamID(TeamTwo)
		game.Players[2].SetTeamID(TeamOne)
		game.Players[3].SetTeamID(TeamTwo)

		app.Lobby.Mu.Lock()
		for _, p := range game.GetHumanPlayers() {
			app.Lobby.UserGames[p.UserId] = game
		}
		app.Lobby.Mu.Unlock()

		for _, hplayer := range game.GetHumanPlayers() {
			err = app.Queries.AddCoinToPerson(context.Background(), sqldb.AddCoinToPersonParams{
				Coin: -1 * hplayer.BetAmount,
				ID:   hplayer.UserId,
			})
			if err != nil {
				app.ReportError(err)
				return
			}
			err = game.SendGameData(TypeMatchFound, hplayer)
			if err != nil {
				app.ReportError(err)
				return
			}
		}
		app.BackgroundTask(func() {
			game.RunGame()

		})
	}
}
