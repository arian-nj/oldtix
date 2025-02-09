package main

import (
	cards "github.com/arian-nj/master-card/back/internal/card"
)

func isCardInCards(card *cards.Card, cards *[]cards.Card) (int, bool) {
	for index, inHandCard := range *cards {
		if inHandCard.Suit == card.Suit && inHandCard.Value == card.Value {
			return index, true
		}
	}
	return 0, false
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
func (game *GameState) isMoveValid(player *Player, card *cards.Card) bool {
	cardIndex, haveCard := isCardInCards(card, &player.Cards)
	if !haveCard {
		return false
	}
	currentTurn := game.CurrentTrick.CurrentTurn
	if len(currentTurn.CardsPlayed) > 0 { // if not first move
		if hasSuit(currentTurn.CardsPlayed[0].Suit, &player.Cards) {
			if card.Suit != currentTurn.CardsPlayed[0].Suit {
				return false
			}
		}
	}
	currentTurn.CardsPlayed = append(currentTurn.CardsPlayed, player.Cards[cardIndex])
	player.Cards = append(player.Cards[:cardIndex], player.Cards[cardIndex+1:]...)
	return true
}
