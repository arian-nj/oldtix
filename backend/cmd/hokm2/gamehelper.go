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
	// game.Logger.Debug("new card " + card.String())
	cardIndex, haveCard := isCardInCards(card, &player.Cards)
	if !haveCard {
		game.Logger.Debug("not in hand")
		return 0, false
	}
	currentTurn := game.CurrentTrick.CurrentTurn
	if len(currentTurn.CardsPlayed) > 0 { // if not first move
		// game.Logger.Debug("first card " + currentTurn.CardsPlayed[0].Card.String())
		if hasSuit(currentTurn.CardsPlayed[0].Card.Suit, &player.Cards) {
			if card.Suit != currentTurn.CardsPlayed[0].Card.Suit {
				game.Logger.Debug("no allowed")
				return 0, false
			}
		}
	}

	return cardIndex, true
}

// func (game *GameState) WhoWins() *PlayerCardPlayed {

// 	HokmSuite := game.CurrentTrick.Hokm

// 	Winner := game.CurrentTrick.CurrentTurn.CardsPlayed[0]

// 	for _, pc := range game.CurrentTrick.CurrentTurn.CardsPlayed {
// 		if Winner.Card.Suit == HokmSuite {
// 			if pc.Card.Suit == HokmSuite && pc.Card.Value >= Winner.Card.Value {
// 				Winner = pc
// 			}
// 		} else {
// 			if pc.Card.Suit == HokmSuite {
// 				Winner = pc
// 			} else if pc.Card.Suit == Winner.Card.Suit && pc.Card.Value > Winner.Card.Value {
// 				Winner = pc
// 			}
// 		}
// 	}
// 	return Winner

// }

// hokm
// zamineh
// number
