package hokm4

import (
	"context"

	"github.com/arian-nj/master-card/back/internal/randutils"
	"github.com/arian-nj/master-card/back/internal/socket"
	"github.com/arian-nj/master-card/back/pkg/hokm4engine"
	"github.com/arian-nj/master-card/back/sqldb"
)

type GameEvent struct {
	Player *HumanPlayer
	event  *socket.Event
}

func NewGameEvent(event *socket.Event, player *HumanPlayer) *GameEvent {
	return &GameEvent{
		event:  event,
		Player: player,
	}
}

type Room struct {
	*ApplicationHokm4 `json:"-"`
	*hokm4engine.Hokm4EngineState
	ID           int             `json:"id"`
	UID          string          `json:"uid"`
	BetAmount    int             `json:"-"`
	GameEventsCh chan *GameEvent `json:"-"`
}

func (app *ApplicationHokm4) NewRoom(bet_amount int) (*Room, error) {
	gameRow, err := app.Queries.InsertHokm4Game(context.Background(), bet_amount)
	if err != nil {
		return nil, err
	}

	return &Room{
		ID:               gameRow.ID,
		UID:              randutils.GenerateRandomString(8),
		ApplicationHokm4: app,
		Hokm4EngineState: hokm4engine.NewHokm4Engine(),
		GameEventsCh:     make(chan *GameEvent),
	}, nil
}

func (gs *Room) SaveGameStateData() error {
	err := gs.Queries.UpdateHokm4Game(context.Background(), sqldb.UpdateHokm4GameParams{
		TeamOneTricksScore: gs.TeamOneTrickScore,
		TeamTwoTricksScore: gs.TeamTwoTrickScore,
		ID:                 gs.ID,
	})

	if err != nil {
		gs.Logger.Error(err.Error())
	}
	return err
}

func (gs *Room) GetHumanPlayers() (allHumanPlayers []*HumanPlayer) {
	for _, p := range gs.Players {
		hp, ok := p.(*HumanPlayer)
		if !ok {
			continue
		}
		allHumanPlayers = append(allHumanPlayers, hp)
	}
	return allHumanPlayers
}
