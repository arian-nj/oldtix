package main

import (
	"sync"

	cards "github.com/arian-nj/master-card/back/internal/card"
	"github.com/arian-nj/master-card/back/internal/socket"
)

type Team int

const (
	TeamOne Team = iota
	TeamTwo
)

type Player struct {
	UserId int32          `json:"user_id"`
	TeamId Team           `json:"team"`
	Client *socket.Client `json:"-"`
	Cards  []cards.Card   `json:"-"`
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
	ID             string          `json:"id"`
	Players        []*Player       `json:"players"`
	GameEventsCh   chan *GameEvent `json:"-"`
	CurrentTrick   *Trick          `json:"current_trick"`
	Tricks         []*Trick        `json:"-"`

	TeamOneTricksScore int `json:"team_one_trick_score"`
	TeamTwoTricksScore int `json:"team_two_trick_score"`
}

type Trick struct {
	Hokm               cards.Suite `json:"hokm"`
	HakemIndex         int         `json:"hakem_index"`
	StarterPlayerIndex int         `json:"-"`

	CurrentTurn *Turn   `json:"current_turn"`
	Turns       []*Turn `json:"-"`

	TeamOneTurnScore int `json:"team_one_turn_score"`
	TeamTwoTurnScore int `json:"team_two_turn_score"`
}

type Turn struct {
	CardsPlayed []*PlayerCardPlayed `json:"played_cards"`
}

func NewTrick() *Trick {
	return &Trick{}
}

func NewTurn() *Turn {
	return &Turn{
		CardsPlayed: []*PlayerCardPlayed{},
	}
}

type Lobby struct {
	Queue chan *Player
	Games map[string]*GameState
	Mu    sync.Mutex
}

type PlayerCardPlayed struct {
	Player *Player    `json:"player_id"`
	Card   cards.Card `json:"card"`
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
