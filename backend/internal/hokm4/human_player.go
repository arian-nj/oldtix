package hokm4

import (
	cards "github.com/arian-nj/master-card/back/internal/card"
	"github.com/arian-nj/master-card/back/internal/socket"
	"github.com/arian-nj/master-card/back/pkg/hokm4engine"
)

type HumanPlayer struct {
	*hokm4engine.PlayerCommon
	UserId    int             `json:"user_id"`
	Client    *socket.Client  `json:"-"`
	IsPlayng  bool            `json:"is_playing"`
	Game      *GameState      `json:"-"`
	BetAmount int             `json:"-"`
	AllEvents []*socket.Event `json:"-"`
}

func NewHumanPlayer(userId int, client *socket.Client, cards []cards.Card, is_playng bool) *HumanPlayer {
	return &HumanPlayer{
		PlayerCommon: hokm4engine.NewPlayerCommon(cards),
		UserId:       userId,
		Client:       client,
		IsPlayng:     is_playng,
	}
}

func (hplayer *HumanPlayer) AddToEgress(e *socket.Event, write_to_events bool) {
	if hplayer.Client != nil && hplayer.Client.State != socket.OPEN {
		return
	}

	go func() {
		hplayer.Client.Egres <- *e
	}()

	if write_to_events {
		if e.Type != TypeNewCard &&
			e.Type != TypeNewCardOne &&
			e.Type != TypeRejoinMatch &&
			e.Type != TypeMatchFound {
			hplayer.AllEvents = append(hplayer.AllEvents, e)
		}
	}
}

// events come here if not used go to GameEventCh
func (app *ApplicationHokm4) BackgroundSocketHandlers(hplayer *HumanPlayer) { // game is nil

	app.BackgroundTask(func() {
		for {
			select {
			case new_event := <-hplayer.Client.NewEvents:
				switch new_event.Type {
				case TypeDisconnect:
					hplayer.Client.Close()
				case TypeMakeMatch:
					activeGame, wasInAGame := app.Lobby.UserGames[hplayer.UserId]
					if wasInAGame {
						err := activeGame.SendGameData(TypeRejoinMatch, hplayer)
						if err != nil {
							app.Logger.Error(err.Error())
						}

						err = activeGame.sendCardsToEgres(hplayer.Cards, hplayer, TypeNewCardOne)
						if err != nil {
							app.Logger.Error(err.Error())
						}

						for _, e := range hplayer.AllEvents {
							hplayer.AddToEgress(e, false)
						}
					} else {
						app.Lobby.MatchmakingQueueGlobal <- NewMatchmakingTicket(hplayer)
					}
				default:
					if hplayer.Game == nil {
						break
					}
					hplayer.Game.GameEventsCh <- NewGameEvent(&new_event, hplayer)
				}
			case <-hplayer.Client.CancelCtx.Done():
				return
			}
		}
	})

}
