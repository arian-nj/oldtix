package hokm4

import (
	cards "github.com/arian-nj/master-card/back/internal/card"
	"github.com/arian-nj/master-card/back/internal/randutils"
	"github.com/arian-nj/master-card/back/internal/socket"
)

type HumanPlayer struct {
	UserId       int64          `json:"user_id"`
	PlayerUnique string         `json:"player_unique"`
	TeamId       Team           `json:"team"`
	Client       *socket.Client `json:"-"`
	Cards        []cards.Card   `json:"-"`
	IsPlayng     bool           `json:"is_playing"`
}

func NewPlayer(UserId int64, Client *socket.Client, Cards []cards.Card, IsPlayng bool) *HumanPlayer {
	randString := randutils.GenerateRandomString(16)

	return &HumanPlayer{
		UserId:       UserId,
		PlayerUnique: randString,
		Client:       Client,
		Cards:        Cards,
		IsPlayng:     IsPlayng,
	}
}

func (p *HumanPlayer) AddToEgress(e *socket.Event) {
	if p.Client != nil && p.Client.State != socket.OPEN {
		return
	}
	go func() {
		p.Client.Egres <- *e

	}()
}
