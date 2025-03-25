package hokm4

import (
	cards "github.com/arian-nj/master-card/back/internal/card"
	"github.com/arian-nj/master-card/back/internal/randutils"
	"github.com/arian-nj/master-card/back/internal/socket"
)

type PlayerInterface interface {
	// for sake of de-coupling i created this Chunky interface even used setters and getters
	// feature me remove them if you can
	AddToEgress(e *socket.Event)
	// PlayCard(game *GameState) (cardIndex int, err error)
	BackgroundSocketHandlers(game *GameState)

	GetTeamID() Team
	SetTeamID(team_id Team)
	SetCards(player_cards []cards.Card)
	GetCards() []cards.Card
}

type PlayerCommon struct {
	PlayerUnique string       `json:"player_unique"`
	TeamId       Team         `json:"team"`
	Cards        []cards.Card `json:"-"`
}

func (pc *PlayerCommon) SetCards(player_cards []cards.Card) {
	pc.Cards = player_cards
}
func (pc *PlayerCommon) GetCards() []cards.Card {
	return pc.Cards
}

func (pc *PlayerCommon) GetTeamID() Team {
	return pc.TeamId
}
func (pc *PlayerCommon) SetTeamID(team_id Team) {
	pc.TeamId = Team(team_id)
}

func NewPlayerCommon(player_cards []cards.Card) *PlayerCommon {
	randString := randutils.GenerateRandomString(24)

	return &PlayerCommon{
		PlayerUnique: randString,
		Cards:        player_cards,
	}
}
