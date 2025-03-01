package main

import (
	"math/rand/v2"

	cards "github.com/arian-nj/master-card/back/internal/card"
)

func (game *GameState) BotPlayTurn(player *Player) int {
	current_turn := game.CurrentTrick.CurrentTurn
	if len(current_turn.CardsPlayed) > 0 {
		return game._not_first_play(player)
	}
	return game._first_play(player)

}

func (game *GameState) _first_play(player *Player) int {
	rand_card_index := rand.IntN(len(player.Cards))
	return rand_card_index
}

func (game *GameState) _not_first_play(player *Player) int {
	current_turn := game.CurrentTrick.CurrentTurn
	first_card_played := current_turn.CardsPlayed[0]

	if cards.HasSuit(first_card_played.Card.Suit, &player.Cards) {
		choosed_card := cards.SelectHighestCard(&player.Cards, first_card_played.Card.Suit)
		cardIndex, _ := cards.IsCardInCards(&choosed_card, &player.Cards)
		return cardIndex
	}

	if cards.HasSuit(game.CurrentTrick.Hokm, &player.Cards) {
		choosed_card := cards.SelectLowestCard(&player.Cards, game.CurrentTrick.Hokm)
		cardIndex, _ := cards.IsCardInCards(&choosed_card, &player.Cards)
		return cardIndex
	}
	extracted_suites := cards.ExtractDeckSuites(&player.Cards)
	rand_suite_index := rand.IntN(len(extracted_suites))

	choosed_suite := extracted_suites[rand_suite_index]
	choosed_card := cards.SelectLowestCard(&player.Cards, choosed_suite)
	cardIndex, _ := cards.IsCardInCards(&choosed_card, &player.Cards)
	return cardIndex
}
