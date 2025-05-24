package hokm4

import (
	cards "github.com/arian-nj/master-card/back/internal/card"
	"github.com/arian-nj/master-card/back/internal/randutils"
	"github.com/arian-nj/master-card/back/internal/socket"
)

type BotPlayer struct {
	*PlayerCommon
}

func NewBotPlayer() *BotPlayer {
	return &BotPlayer{
		PlayerCommon: NewPlayerCommon([]cards.Card{}),
	}
}

func (bplayer *BotPlayer) AddToEgress(e *socket.Event, write_to_events bool) {
}

func (bplayer *BotPlayer) BackgroundSocketHandlers(game *GameState) {
}

func (game *GameState) BotPlayTurn(hand_cards []cards.Card) int {
	current_turn := game.CurrentTrick.CurrentTurn
	if len(current_turn.CardsPlayed) > 0 {
		return game._not_first_play(&hand_cards)
	}
	return game._first_play(&hand_cards)

}

func (game *GameState) _first_play(hand_cards *[]cards.Card) int {
	rand_card_index := randutils.GenerateRandomNumber(len(*hand_cards))
	return rand_card_index
}

func (game *GameState) _not_first_play(hand_cards *[]cards.Card) int {
	current_turn := game.CurrentTrick.CurrentTurn
	first_card_played := current_turn.CardsPlayed[0]
	turn_base_suite := first_card_played.Card.Suit
	led_card := game.WhoWinsLastTurn().Card
	led_suit := led_card.Suit
	isBureshed := (turn_base_suite != led_suit) && (led_suit == game.CurrentTrick.Hokm)

	if cards.HasSuit(turn_base_suite, hand_cards) {
		var choosed_card cards.Card
		if isBureshed {
			choosed_card = cards.SelectLowestCard(hand_cards, turn_base_suite)
		} else {
			choosed_card = cards.SelectHighestCard(hand_cards, turn_base_suite)
			if choosed_card.Value < led_card.Value {
				choosed_card = cards.SelectLowestCard(hand_cards, turn_base_suite)
			}
		}

		cardIndex, _ := cards.IsCardInCards(&choosed_card, hand_cards)
		return cardIndex
	}

	if cards.HasSuit(game.CurrentTrick.Hokm, hand_cards) {
		var choosed_card cards.Card

		if led_suit == game.CurrentTrick.Hokm {
			choosed_card = cards.SelectHighestCard(hand_cards, game.CurrentTrick.Hokm)
			if choosed_card.Value < led_card.Value {
				choosed_card = cards.SelectLowestCard(hand_cards, game.CurrentTrick.Hokm)
			}
		} else {
			choosed_card = cards.SelectLowestCard(hand_cards, game.CurrentTrick.Hokm)
		}

		cardIndex, _ := cards.IsCardInCards(&choosed_card, hand_cards)
		return cardIndex
	}
	extracted_suites := cards.ExtractDeckSuites(hand_cards)
	rand_suite_index := randutils.GenerateRandomNumber(len(extracted_suites))

	choosed_suite := extracted_suites[rand_suite_index]
	choosed_card := cards.SelectLowestCard(hand_cards, choosed_suite)
	cardIndex, _ := cards.IsCardInCards(&choosed_card, hand_cards)
	return cardIndex
}

func (game *GameState) WhoWinsLastTurn() *PlayerCardPlayed {
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
