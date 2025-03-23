package hokm4

import "github.com/arian-nj/master-card/back/internal/socket"

type Player interface {
	AddToEgress(e *socket.Event)
	PlayCard(game *GameState) (cardIndex int, err error)
	ChooseHokm(game *GameState)
	BackgroundSocketHandlers(game *GameState)
}
