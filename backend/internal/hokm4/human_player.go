package hokm4

import (
	cards "github.com/arian-nj/master-card/back/internal/card"
	"github.com/arian-nj/master-card/back/internal/randutils"
	"github.com/arian-nj/master-card/back/internal/socket"
)

type HumanPlayer struct {
	*Player
	UserId   int64          `json:"user_id"`
	Client   *socket.Client `json:"-"`
	IsPlayng bool           `json:"is_playing"`
}

func NewHumanPlayer(UserId int64, Client *socket.Client, Cards []cards.Card, IsPlayng bool) *HumanPlayer {
	randString := randutils.GenerateRandomString(16)

	return &HumanPlayer{
		Player: &Player{
			PlayerUnique: randString,
			Cards:        Cards,
		},
		UserId:   UserId,
		Client:   Client,
		IsPlayng: IsPlayng,
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
