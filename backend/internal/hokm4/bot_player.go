package hokm4

import (
	"math/rand/v2"

	cards "github.com/arian-nj/master-card/back/internal/card"
)

type BotPlayer struct {
	*Player
}

func (game *GameState) BotPlayTurn(hand_cards []cards.Card) int {
	current_turn := game.CurrentTrick.CurrentTurn
	if len(current_turn.CardsPlayed) > 0 {
		return game._not_first_play(&hand_cards)
	}
	return game._first_play(&hand_cards)

}

func (game *GameState) _first_play(hand_cards *[]cards.Card) int {
	rand_card_index := rand.IntN(len(*hand_cards))
	return rand_card_index
}

func (game *GameState) _not_first_play(hand_cards *[]cards.Card) int {
	current_turn := game.CurrentTrick.CurrentTurn
	first_card_played := current_turn.CardsPlayed[0]

	if cards.HasSuit(first_card_played.Card.Suit, hand_cards) {
		choosed_card := cards.SelectHighestCard(hand_cards, first_card_played.Card.Suit)
		cardIndex, _ := cards.IsCardInCards(&choosed_card, hand_cards)
		return cardIndex
	}

	if cards.HasSuit(game.CurrentTrick.Hokm, hand_cards) {
		choosed_card := cards.SelectLowestCard(hand_cards, game.CurrentTrick.Hokm)
		cardIndex, _ := cards.IsCardInCards(&choosed_card, hand_cards)
		return cardIndex
	}
	extracted_suites := cards.ExtractDeckSuites(hand_cards)
	rand_suite_index := rand.IntN(len(extracted_suites))

	choosed_suite := extracted_suites[rand_suite_index]
	choosed_card := cards.SelectLowestCard(hand_cards, choosed_suite)
	cardIndex, _ := cards.IsCardInCards(&choosed_card, hand_cards)
	return cardIndex
}
