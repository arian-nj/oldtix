package hokm4

import (
	"encoding/json"

	cards "github.com/arian-nj/master-card/back/internal/card"
	"github.com/arian-nj/master-card/back/internal/socket"
)

// have card
// same suite as first card if not first move
func (game *GameState) ValidateAndDoMove(player *Player, played_card *cards.Card) (int, bool) {
	// game.Logger.Debug("new card " + card.String())
	cardIndex, haveCard := cards.IsCardInCards(played_card, &player.Cards)
	if !haveCard {
		game.Logger.Debug("not in hand")
		return 0, false
	}
	currentTurn := game.CurrentTrick.CurrentTurn
	if len(currentTurn.CardsPlayed) > 0 { // if not first move
		// game.Logger.Debug("first card " + currentTurn.CardsPlayed[0].Card.String())
		first_card_played_suite := currentTurn.CardsPlayed[0].Card.Suit
		if cards.HasSuit(first_card_played_suite, &player.Cards) {
			if played_card.Suit != first_card_played_suite {
				game.Logger.Debug("no allowed")
				return 0, false
			}
		}
	}

	return cardIndex, true
}

func (game *GameState) WhoWins() *PlayerCardPlayed {

	HokmSuite := game.CurrentTrick.Hokm

	Winner := game.CurrentTrick.CurrentTurn.CardsPlayed[0]
	for _, pc := range game.CurrentTrick.CurrentTurn.CardsPlayed {
		if Winner.Card.Suit == HokmSuite {
			if pc.Card.Suit == HokmSuite && pc.Card.Value >= Winner.Card.Value {
				Winner = pc
			}
		} else {
			if pc.Card.Suit == HokmSuite {
				Winner = pc
			} else if pc.Card.Suit == Winner.Card.Suit && pc.Card.Value > Winner.Card.Value {
				Winner = pc
			}
		}
	}

	return Winner
}

// hokm
// zamineh
// number

type GameStateOut struct {
	*GameState
	YourTeam Team `json:"your_team"`
}

func (game *GameState) SendGameData(MessageTurn socket.EventType, p *Player) error {
	gsOut := GameStateOut{
		GameState: game,
		YourTeam:  p.TeamId,
	}
	// send game data
	game_data, err := json.Marshal(gsOut)
	if err != nil {
		return err
	}
	p.AddToEgress(socket.NewEvent(MessageTurn, socket.EventMessage(game_data)))
	return nil
}

func (game *GameState) sendCards(number int, all_cards []cards.Card, players []*Player) []cards.Card {
	var remaining_cards []cards.Card = all_cards

	for _, p := range players {
		var randomCards []cards.Card
		var err error
		randomCards, remaining_cards, err = cards.GiveRandomCards(number, remaining_cards)
		if err != nil {
			game.Logger.Error(err.Error())
		}

		var output struct {
			NewCards []cards.Card `json:"cards"`
		}
		output.NewCards = randomCards
		data_byte, err := json.Marshal(output)
		if err != nil {
			game.Logger.Error(err.Error())
		}
		p.AddToEgress(socket.NewEvent(socket.TypeNewCard, socket.EventMessage(data_byte)))
		p.Cards = append(p.Cards, randomCards...)
	}
	return remaining_cards

}
