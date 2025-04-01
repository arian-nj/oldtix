package hokm4

import (
	"encoding/json"
	"log"

	cards "github.com/arian-nj/master-card/back/internal/card"
	"github.com/arian-nj/master-card/back/internal/socket"
)

type HumanPlayer struct {
	*PlayerCommon
	UserId   int            `json:"user_id"`
	Client   *socket.Client `json:"-"`
	IsPlayng bool           `json:"is_playing"`
	Game     *GameState     `json:"-"`
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
	if e.Type == TypeRejoinMatch {
		log.Println("sending rejoin")
	}
	go func() {
		hplayer.Client.Egres <- *e

	}()
}

// events come here if not used go to GameEventCh
func (app *ApplicationHokm4) BackgroundSocketHandlers(hplayer *HumanPlayer, bet_amount int) { // game is nil

	app.BackgroundTask(func() error {
		for {
			select {
			case new_event := <-hplayer.Client.NewEvents:
				switch new_event.Type {
				case TypeGetData:
					if hplayer.Game == nil {
						break
					}
					err := hplayer.Game.SendGameData(TypeGameData, hplayer)
					if err != nil {
						return err
					}
				case TypeGetMyCards:
					var output struct {
						NewCards []cards.Card `json:"cards"`
					}
					output.NewCards = hplayer.Cards
					data_byte, err := json.Marshal(output)
					if err != nil {
						app.Logger.Error(err.Error())
						continue
					}
					hplayer.AddToEgress(socket.NewEvent(TypeGetMyCards, socket.EventMessage(data_byte)))
				case TypeDisconnect:
					hplayer.Client.Close()
				case TypeMakeMatch:
					activeGame, wasInAGame := app.Lobby.UserGames[hplayer.UserId]
					if wasInAGame {
						err := activeGame.SendGameData(TypeRejoinMatch, hplayer)
						if err != nil {
							app.Logger.Error(err.Error())
						}
					} else {
						new_match_request := NewMatchmakingRequest(hplayer)
						if bet_amount == 0 {
							app.Lobby.MatchmakingQueueFor0 <- new_match_request
						} else if bet_amount == 50 {
							app.Lobby.MatchmakingQueueFor50 <- new_match_request
						}
					}
				default:
					if hplayer.Game == nil {
						break
					}
					hplayer.Game.GameEventsCh <- NewGameEvent(&new_event, hplayer)
				}
			case <-hplayer.Client.CancelCtx.Done():
				return nil
			}
		}
	})

}
