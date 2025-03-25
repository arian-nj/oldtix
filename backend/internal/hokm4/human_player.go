package hokm4

import (
	"encoding/json"

	cards "github.com/arian-nj/master-card/back/internal/card"
	"github.com/arian-nj/master-card/back/internal/socket"
)

type HumanPlayer struct {
	*PlayerCommon
	UserId   int            `json:"user_id"`
	Client   *socket.Client `json:"-"`
	IsPlayng bool           `json:"is_playing"`
}

func NewHumanPlayer(userId int, client *socket.Client, cards []cards.Card, is_playng bool) *HumanPlayer {
	return &HumanPlayer{
		PlayerCommon: NewPlayerCommon(cards),
		UserId:       userId,
		Client:       client,
		IsPlayng:     is_playng,
	}
}

func (hplayer *HumanPlayer) AddToEgress(e *socket.Event) {
	if hplayer.Client != nil && hplayer.Client.State != socket.OPEN {
		return
	}

	go func() {
		hplayer.Client.Egres <- *e

	}()
}

// events come here if not used go to GameEventCh
func (hplayer *HumanPlayer) BackgroundSocketHandlers(game *GameState) {

	game.BackgroundTask(func() error {
		for {
			select {
			case new_event := <-hplayer.Client.NewEvents:
				switch new_event.Type {
				case socket.TypeGetData:
					err := game.SendGameData(socket.TypeGameData, hplayer)
					if err != nil {
						return err
					}
				case socket.TypeGetMyCards:
					var output struct {
						NewCards []cards.Card `json:"cards"`
					}
					output.NewCards = hplayer.Cards
					data_byte, err := json.Marshal(output)
					if err != nil {
						game.Logger.Error(err.Error())
						continue
					}
					hplayer.AddToEgress(socket.NewEvent(socket.TypeGetMyCards, socket.EventMessage(data_byte)))
				default:
					game.GameEventsCh <- NewGameEvent(&new_event, hplayer)
				}
			case <-hplayer.Client.CancelCtx.Done():
				return nil
			}
		}
	})

}
