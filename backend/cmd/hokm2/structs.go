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

type GameState struct {
	ID      string    `json:"id"`
	Players []*Player `json:"players"`
	Current int       `json:"Current"`
	Hakem   int32     `json:"Hakem"`
}

type Lobby struct {
	Queue chan *Player
	Games map[string]*GameState
	Mu    sync.Mutex
}

type CurrentPlayer int

func (gs *GameState) NextCurrent() {
	lastPlayerIndex := len(gs.Players) - 1
	if gs.Current+1 <= lastPlayerIndex {
		gs.Current += 1
	} else if gs.Current+1 > lastPlayerIndex {
		gs.Current = 0
	}

}
