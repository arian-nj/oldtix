package main

import (
	"fmt"

	cards "github.com/arian-nj/master-card/back/internal/card"
)

func isCardInCards(card *cards.Card, cards *[]cards.Card) (int, bool) {
	for index, inHandCard := range *cards {
		if inHandCard.Suit == card.Suit && inHandCard.Value == card.Value {
			return index, true
		}
	}
	return -1, false
}

func hasSuit(wanted_suite cards.Suite, cards *[]cards.Card) bool {
	for _, inHandCard := range *cards {
		if inHandCard.Suit == wanted_suite {
			return true
		}
	}
	return false
}

// have card
// same suite as first card if not first move
func (game *GameState) ValidateAndDoMove(player *Player, card *cards.Card) (int, bool) {
	cardIndex, haveCard := isCardInCards(card, &player.Cards)
	if !haveCard {
		fmt.Println("not in hand")
		return 0, false
	}
	currentTurn := game.CurrentTrick.CurrentTurn
	if len(currentTurn.CardsPlayed) > 0 { // if not first move
		if hasSuit(currentTurn.CardsPlayed[0].Card.Suit, &player.Cards) {
			if card.Suit != currentTurn.CardsPlayed[0].Card.Suit {
				fmt.Println("no allowed")
				return 0, false
			}
		}
	}

	return cardIndex, true
}

// func (game *GameState) WhoWins() {
// 	effective_playing_cards := []PlayerCardPlayed{}
// 	for _, pc := range game.CurrentTrick.CurrentTurn.CardsPlayed {

// 	}
// }
