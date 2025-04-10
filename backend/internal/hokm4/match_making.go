package hokm4

import (
	"context"
	"strconv"
	"sync"
	"time"

	cards "github.com/arian-nj/master-card/back/internal/card"
	"github.com/arian-nj/master-card/back/internal/socket"
	"github.com/arian-nj/master-card/back/sqldb"
)

var validBetAmounts = []int{BET_AMOUNT_ONE, BET_AMOUNT_TWO}

type Lobby struct {
	MatchmakingQueueGlobal    chan *MatchmakingRequest
	MatchmakingQueueForBetOne chan *MatchmakingRequest
	MatchmakingQueueForBetTwo chan *MatchmakingRequest

	MakingMatches map[int]*GameState
	UserGames     map[int]*GameState
	Mu            sync.Mutex
}

func NewLobby() *Lobby {
	return &Lobby{
		MatchmakingQueueGlobal:    make(chan *MatchmakingRequest),
		MatchmakingQueueForBetOne: make(chan *MatchmakingRequest),
		MatchmakingQueueForBetTwo: make(chan *MatchmakingRequest),
		// Games: make(map[int]*GameState),
		UserGames: map[int]*GameState{},
		Mu:        sync.Mutex{},
	}
}

type MatchmakingRequest struct {
	Player    *HumanPlayer
	BetAmount int
	Timestamp time.Time
}

func NewMatchmakingRequest(hplayer *HumanPlayer, bet_amount int) *MatchmakingRequest {
	return &MatchmakingRequest{
		Player:    hplayer,
		BetAmount: bet_amount,
		Timestamp: time.Now(),
	}
}
func (game *GameState) AddHumanPlayerToGame(player *HumanPlayer, gameId int) error {
	_, err := game.Queries.InsertGamePlayer(context.Background(), sqldb.InsertGamePlayerParams{
		PlayerID: player.UserId,
		GameID:   gameId,
		Team:     int(player.TeamId),
	})
	game.Players = append(game.Players, player)
	return err
}

func (game *GameState) AddBotPlayerToGame(player *BotPlayer) {
	game.Players = append(game.Players, player)
}
func (app *ApplicationHokm4) FilterMatchMkingByCoin() {
	for {
		newReq := <-app.Lobby.MatchmakingQueueGlobal
		if newReq.BetAmount == BET_AMOUNT_ONE {
			app.Lobby.MatchmakingQueueForBetOne <- newReq
		} else if newReq.BetAmount == BET_AMOUNT_TWO {
			app.Lobby.MatchmakingQueueForBetTwo <- newReq
		} else {
			app.Logger.Error("can'r filter this amount of coin " + strconv.Itoa(newReq.BetAmount))
		}
	}
}

func (app *ApplicationHokm4) MatchUsers(matchesChan chan *MatchmakingRequest, betting_amount int) {
	for {
		game, err := app.NewGameState()
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
			err = game.AddHumanPlayerToGame(p, game.ID)
			if err != nil {
				app.ReportError(err)
				return
			}
		}

		// foundPlayer.IsPlayng = true

		p3 := NewBotPlayer([]cards.Card{})
		game.AddBotPlayerToGame(p3)

		p4 := NewBotPlayer([]cards.Card{})
		game.AddBotPlayerToGame(p4)

		game.Players[0].SetTeamID(TeamOne)
		game.Players[1].SetTeamID(TeamTwo)
		game.Players[2].SetTeamID(TeamOne)
		game.Players[3].SetTeamID(TeamTwo)

		app.Lobby.Mu.Lock()
		for _, p := range game.GetHumanPlayers() {
			app.Lobby.UserGames[p.UserId] = game
		}
		app.Lobby.Mu.Unlock()

		for _, p := range game.GetHumanPlayers() {
			err = app.Queries.AddCoinToPerson(context.Background(), sqldb.AddCoinToPersonParams{
				Coin: -1 * betting_amount,
				ID:   p.UserId,
			})
			if err != nil {
				app.ReportError(err)
				return
			}
			err = game.SendGameData(TypeMatchFound, p)
			if err != nil {
				app.ReportError(err)
				return
			}
		}
		app.RunGameInBackground(game)
	}
}
