package hokm4

import cards "github.com/arian-nj/master-card/back/internal/card"

type Player struct {
	PlayerUnique string       `json:"player_unique"`
	TeamId       Team         `json:"team"`
	Cards        []cards.Card `json:"-"`
}
