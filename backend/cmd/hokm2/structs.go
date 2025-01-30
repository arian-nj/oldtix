package main

import (
	"sync"

	cards "github.com/arian-nj/master-card/back/internal/card"
	"github.com/arian-nj/master-card/back/internal/socket"
)

type Player struct {
	UserId int32
	Client *socket.Client
	Cards  []cards.Card
}

type Game struct {
	ID      string
	Players [2]*Player
	Current int
}

type Lobby struct {
	Queue chan *Player
	Games map[string]*Game
	Mu    sync.Mutex
}
