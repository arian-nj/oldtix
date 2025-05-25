package hokm4engine

import cards "github.com/arian-nj/master-card/back/internal/card"

type Hokm4EngineState struct {
	Players           []PlayerInterface `json:"players"`
	CurrentTrick      *Trick            `json:"current_trick"`
	TeamOneTrickScore int               `json:"team_one_trick_score"`
	TeamTwoTrickScore int               `json:"team_two_trick_score"`

	Tricks []*Trick `json:"-"`
}

func NewHokm4Engine() *Hokm4EngineState {
	return &Hokm4EngineState{}
}

type Trick struct {
	ID               int
	Hokm             cards.Suite `json:"hokm"`
	HakemIndex       int         `json:"hakem_index"`
	TurnStarterIndex int         `json:"turn_starter_index"`

	CurrentTurn *Turn   `json:"current_turn"`
	Turns       []*Turn `json:"-"`

	TeamOneTurnScore int `json:"team_one_turn_score"`
	TeamTwoTurnScore int `json:"team_two_turn_score"`

	WinnerTeam Team `json:"-"`
}

func (game *Hokm4EngineState) NewTrick(trick_id, hakem_index int) *Trick {
	return &Trick{
		ID:         trick_id,
		HakemIndex: hakem_index,
	}
}

func (game *Hokm4EngineState) Save() {

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
