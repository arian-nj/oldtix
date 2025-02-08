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
	ID           string          `json:"id"`
	Players      []*Player       `json:"players"`
	GameEventsCh chan *GameEvent `json:"-"`

	HakemIndex int `json:"hakem"`
	// Current     int         `json:"current"`
	CurrentTurn *Turn `json:"turn,omitempty"`
}

type Turn struct {
	PlayerIndex int          `json:"player_index"`
	CardsPlayed []cards.Card `json:"cards"`
	Hokm        cards.Suite  `json:"hokm"`
}

func NewTurn(player_index int) *Turn {
	return &Turn{
		PlayerIndex: player_index,
		CardsPlayed: []cards.Card{},
	}
}

type Lobby struct {
	Queue chan *Player
	Games map[string]*GameState
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
