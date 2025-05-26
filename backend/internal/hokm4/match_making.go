package hokm4

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/arian-nj/master-card/back/internal/socket"
	"github.com/arian-nj/master-card/back/pkg/hokm4engine"
	"github.com/arian-nj/master-card/back/sqldb"
)

var validBetAmounts = []int{BET_NO_MONEY, BET_AMOUNT_ONE}

type Lobby struct {
	MatchmakingQueueGlobal chan *MatchmakingTicket

	MatchmakingQueueForBetOne chan *MatchmakingTicket
	MatchmakingQueueForBetTwo chan *MatchmakingTicket

	// MakingMatches map[int]*GameState
	PrivateRoom map[string]*Room
	UserGames   map[int]*Room

	Mu sync.Mutex
}

func NewLobby() *Lobby {
	return &Lobby{
		MatchmakingQueueGlobal:    make(chan *MatchmakingTicket),
		MatchmakingQueueForBetOne: make(chan *MatchmakingTicket),
		MatchmakingQueueForBetTwo: make(chan *MatchmakingTicket),
		// Games: make(map[int]*Room),
		UserGames: map[int]*Room{},
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
func (game *Room) AddHumanPlayerToGame(player *HumanPlayer) error {
	_, err := game.Queries.InsertGamePlayer(context.Background(), sqldb.InsertGamePlayerParams{
		PlayerID: player.UserId,
		GameID:   game.ID,
		Team:     int(player.TeamId),
	})
	game.Players = append(game.Players, player)
	player.Game = game
	return err
}

func (game *Room) AddBotPlayerToGame() {
	player := NewBotPlayer()
	game.Players = append(game.Players, player)
}

func (app *ApplicationHokm4) FilterMatchMakingByCoin() {
	for {
		newReq := <-app.Lobby.MatchmakingQueueGlobal
		if newReq.BetAmount == BET_AMOUNT_ONE || newReq.BetAmount == BET_NO_MONEY {
			app.Lobby.MatchmakingQueueForBetOne <- newReq
		} else {
			app.Logger.Error("can'r filter this amount of coin " + strconv.Itoa(newReq.BetAmount))
		}
	}
}
func (app *ApplicationHokm4) MatchUsers(matchesChan chan *MatchmakingTicket, betting_amount int) {
	MAX_PLAYERS := 4
	for {
		newRoom, err := app.NewRoom(betting_amount)
		if err != nil {
			app.ReportError(err)
			return
		}
		newRoom.BetAmount = betting_amount

		foundHumanPlayers := app.WaitForPlayers(matchesChan)
		for _, humanPlayer := range foundHumanPlayers {
			err := newRoom.AddHumanPlayerToGame(humanPlayer)
			if err != nil {
				app.ReportError(err)
				return
			}
		}
		numberOfBots := MAX_PLAYERS - len(foundHumanPlayers)
		if numberOfBots > 0 {
			for range numberOfBots {
				newRoom.AddBotPlayerToGame()
			}
		}

		newRoom.Players[0].SetTeamID(hokm4engine.TeamOne)
		newRoom.Players[1].SetTeamID(hokm4engine.TeamTwo)
		newRoom.Players[2].SetTeamID(hokm4engine.TeamOne)
		newRoom.Players[3].SetTeamID(hokm4engine.TeamTwo)

		app.Lobby.Mu.Lock()
		for _, p := range newRoom.GetHumanPlayers() {
			app.Lobby.UserGames[p.UserId] = newRoom
		}
		app.Lobby.Mu.Unlock()

		for _, hplayer := range newRoom.GetHumanPlayers() {
			err = app.Queries.AddCoinToPerson(context.Background(), sqldb.AddCoinToPersonParams{
				Coin: -1 * hplayer.BetAmount,
				ID:   hplayer.UserId,
			})
			if err != nil {
				app.ReportError(err)
				return
			}
			err = newRoom.SendGameData(TypeMatchFound, hplayer)
			if err != nil {
				app.ReportError(err)
				return
			}
		}

		app.BackgroundTask(func() {
			newRoom.RunGame()

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
