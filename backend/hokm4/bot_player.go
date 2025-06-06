package hokm4

import (
	cards "github.com/arian-nj/master-card/back/internal/card"
	"github.com/arian-nj/master-card/back/internal/socket"
	"github.com/arian-nj/master-card/back/pkg/hokm4engine"
)

type BotPlayer struct {
	*hokm4engine.PlayerCommon
}

func NewBotPlayer() *BotPlayer {
	return &BotPlayer{
		PlayerCommon: hokm4engine.NewPlayerCommon([]cards.Card{}),
	}
}

func (bplayer *BotPlayer) AddToEgress(e *socket.Event, write_to_events bool) {
}

func (bplayer *BotPlayer) BackgroundSocketHandlers(game *hokm4engine.Hokm4EngineState) {
}
