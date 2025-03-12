package main

import (
	"context"
	"sync"

	cards "github.com/arian-nj/master-card/back/internal/card"
	"github.com/arian-nj/master-card/back/internal/socket"
	"github.com/arian-nj/master-card/back/internal/utils"
	"github.com/arian-nj/master-card/back/sqldb"
)

type Team int16

const (
	TeamOne Team = iota
	TeamTwo
)

type Player struct {
	UserId       int64          `json:"user_id"`
	PlayerUnique string         `json:"player_unique"`
	TeamId       Team           `json:"team"`
	Client       *socket.Client `json:"-"`
	Cards        []cards.Card   `json:"-"`
	IsPlayng     bool           `json:"is_playing"`
}

func NewPlayer(UserId int64, Client *socket.Client, Cards []cards.Card, IsPlayng bool) *Player {
	randString := utils.GenerateRandomString(16)

	return &Player{
		UserId:       UserId,
		PlayerUnique: randString,
		Client:       Client,
		Cards:        Cards,
		IsPlayng:     IsPlayng,
	}
}

func (p *Player) AddToEgress(e *socket.Event) {
	if !p.IsPlayng {
		return
	}
	p.Client.Egres <- *e

}

type GameEvent struct {
	Player *Player
	event  *socket.Event
}

func NewGameEvent(event *socket.Event, player *Player) *GameEvent {
	return &GameEvent{
		event:  event,
		Player: player,
	}
}

type GameState struct {
	*ApplicationH2 `json:"-"`
	ID             int64           `json:"id"`
	Players        []*Player       `json:"players"`
	GameEventsCh   chan *GameEvent `json:"-"`
	CurrentTrick   *Trick          `json:"current_trick"`
	Tricks         []*Trick        `json:"-"`

	TeamOneTricksScore int `json:"team_one_trick_score"`
	TeamTwoTricksScore int `json:"team_two_trick_score"`
}

func (app *ApplicationH2) NewGameState(players []*Player) (*GameState, error) {

	gameRow, err := app.Queries.InsertHokm4Game(context.Background())
	if err != nil {
		return nil, err
	}

	for _, player := range players {
		app.Queries.InsertGamePlayer(context.Background(), sqldb.InsertGamePlayerParams{
			PlayerID: player.UserId,
			GameID:   gameRow.ID,
			Team:     int16(player.TeamId),
		})
	}

	return &GameState{
		ID:            gameRow.ID,
		ApplicationH2: app,
		GameEventsCh:  make(chan *GameEvent),
		Players:       players,
	}, nil
}

type Trick struct {
	id               int64
	Hokm             cards.Suite `json:"hokm"`
	HakemIndex       int         `json:"hakem_index"`
	TurnStarterIndex int         `json:"turn_starter_index"`

	CurrentTurn *Turn   `json:"current_turn"`
	Turns       []*Turn `json:"-"`

	TeamOneTurnScore int `json:"team_one_turn_score"`
	TeamTwoTurnScore int `json:"team_two_turn_score"`

	WinnerTeam Team `json:"-"`
}

func (game *GameState) NewTrick(HakemIndex int) (*Trick, error) {
	trickRow, err := game.Queries.InsertTrick(context.Background(), sqldb.InsertTrickParams{
		GameID:     int64(game.ID),
		HakemIndex: int32(HakemIndex),
	})
	if err != nil {
		return nil, err
	}

	return &Trick{
		id:         trickRow.TrickID,
		HakemIndex: int(trickRow.HakemIndex),
	}, nil
}

func (game *GameState) Save() {

}

type Turn struct {
	CardsPlayed []*PlayerCardPlayed `json:"played_cards"`
}

func NewTurn() *Turn {
	// game
	return &Turn{
		CardsPlayed: []*PlayerCardPlayed{},
	}
}

type PlayerCardPlayed struct {
	Player *Player    `json:"player"`
	Card   cards.Card `json:"card"`
}

func NewPlayerCardPlayed(player *Player, card cards.Card) *PlayerCardPlayed {
	return &PlayerCardPlayed{
		Player: player,
		Card:   card,
	}
}

type Lobby struct {
	Queue chan *Player
	Games map[int64]*GameState
	Mu    sync.Mutex
}

// type CurrentPlayer int

// func (gs *GameState) NextCurrent() {
// 	lastPlayerIndex := len(gs.Players) - 1
// 	if gs.Current+1 <= lastPlayerIndex {
// 		gs.Current += 1
// 	} else if gs.Current+1 > lastPlayerIndex {
// 		gs.Current = 0
// 	}

// }
