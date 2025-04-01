package hokm4

import (
	"context"

	cards "github.com/arian-nj/master-card/back/internal/card"
	"github.com/arian-nj/master-card/back/internal/socket"
	"github.com/arian-nj/master-card/back/sqldb"
)

type Team int

const (
	TeamOne Team = iota
	TeamTwo
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

type GameState struct {
	*ApplicationHokm4 `json:"-"`
	ID                int               `json:"id"`
	Players           []PlayerInterface `json:"players"`
	GameEventsCh      chan *GameEvent   `json:"-"`
	CurrentTrick      *Trick            `json:"current_trick"`
	Tricks            []*Trick          `json:"-"`
	BetAmount         int               `json:"-"`
	TeamOneTrickScore int               `json:"team_one_trick_score"`
	TeamTwoTrickScore int               `json:"team_two_trick_score"`
}

func (gs *GameState) GetHumanPlayers() (allHumanPlayers []*HumanPlayer) {
	for _, p := range gs.Players {
		hp, ok := p.(*HumanPlayer)
		if !ok {
			continue
		}
		allHumanPlayers = append(allHumanPlayers, hp)
	}
	return allHumanPlayers
}

func (app *ApplicationHokm4) NewGameState() (*GameState, error) {
	gameRow, err := app.Queries.InsertHokm4Game(context.Background())
	if err != nil {
		return nil, err
	}

	return &GameState{
		ID:               gameRow.ID,
		ApplicationHokm4: app,
		GameEventsCh:     make(chan *GameEvent),
	}, nil
}

type Trick struct {
	id               int
	Hokm             cards.Suite `json:"hokm"`
	HakemIndex       int         `json:"hakem_index"`
	TurnStarterIndex int         `json:"turn_starter_index"`

	CurrentTurn *Turn   `json:"current_turn"`
	Turns       []*Turn `json:"-"`

	TeamOneTurnScore int `json:"team_one_turn_score"`
	TeamTwoTurnScore int `json:"team_two_turn_score"`

	WinnerTeam Team `json:"-"`
}

func (game *GameState) NewTrick(hakem_index int) (*Trick, error) {
	trickRow, err := game.Queries.InsertTrick(context.Background(), sqldb.InsertTrickParams{
		GameID:     game.ID,
		HakemIndex: hakem_index,
	})
	if err != nil {
		return nil, err
	}

	return &Trick{
		id:         trickRow.TrickID,
		HakemIndex: int(trickRow.HakemIndex),
	}, nil
}

func (game *GameState) Save() {

}

type Turn struct {
	CardsPlayed []*PlayerCardPlayed `json:"played_cards"`
}

func NewTurn() *Turn {
	// game
	return &Turn{
		CardsPlayed: []*PlayerCardPlayed{},
	}
}

type PlayerCardPlayed struct {
	Player PlayerInterface `json:"player"`
	Card   cards.Card      `json:"card"`
}

func NewPlayerCardPlayed(player PlayerInterface, card cards.Card) *PlayerCardPlayed {
	return &PlayerCardPlayed{
		Player: player,
		Card:   card,
	}
}
