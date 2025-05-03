package cards

import (
	"fmt"

	"github.com/arian-nj/master-card/back/internal/randutils"
)

type Suite int
type CardValue int

const ( // sync it with card/card.gd
	HEART  Suite = iota
	CLUB         = 1
	DIMOND       = 2
	SPADE        = 3
)

const (
	N2 = iota + 2
	N3
	N4
	N5
	N6
	N7
	N8
	N9
	N10
	FJ
	FQ
	FK
	ACE
)

var AllSuits = []Suite{HEART, SPADE, DIMOND, CLUB}
var AllValues = []CardValue{ACE, FK, FQ, FJ, N10, N9, N8, N7, N6, N5, N4, N3, N2}

type Card struct {
	Suit  Suite     `json:"suit"`
	Value CardValue `json:"value"`
}

func (c *Card) String() string {
	return fmt.Sprintf("%d %d", c.Suit, c.Value)
}

// var allCards = []Card{}

func NewAllCards() []Card {
	var allCards = []Card{}
	for _, s := range AllSuits {
		for _, v := range AllValues {
			allCards = append(allCards, Card{Suit: s, Value: v})
		}
	}
	return allCards
}

func GiveRandomCards(numberOfCards int, availableCards []Card) (allRandomCards []Card, remaningCards []Card, err error) {
	for range numberOfCards {
		if len(availableCards) == 0 {
			return allRandomCards, availableCards, fmt.Errorf("available cards is %d cant give more", len(availableCards))
		}
		random_index := randutils.GenerateRandomNumber(len(availableCards))
		random_card := availableCards[random_index]
		availableCards = append(availableCards[:random_index], availableCards[random_index+1:]...)
		allRandomCards = append(allRandomCards, random_card)
	}
	return allRandomCards, availableCards, nil
}

func IsCardInCards(card *Card, cards *[]Card) (int, bool) {
	for index, inHandCard := range *cards {
		if inHandCard.Suit == card.Suit && inHandCard.Value == card.Value {
			return index, true
		}
	}
	return -1, false
}

func HasSuit(wanted_suite Suite, cards *[]Card) bool {
	for _, inHandCard := range *cards {
		if inHandCard.Suit == wanted_suite {
			return true
		}
	}
	return false
}
func ExtractSuiteCards(cards *[]Card, wanted_suite Suite) []Card {
	extractedCards := []Card{}
	for _, c := range *cards {
		if c.Suit == wanted_suite {
			extractedCards = append(extractedCards, c)
		}
	}
	return extractedCards
}

func SelectHighestCard(cards *[]Card, wanted_suite Suite) Card {
	playedSuiteCards := ExtractSuiteCards(cards, wanted_suite)
	highest := playedSuiteCards[0]
	for _, c := range playedSuiteCards {
		if c.Value > highest.Value {
			highest = c
		}
	}
	return highest
}

func SelectLowestCard(cards *[]Card, wanted_suite Suite) Card {
	playedSuiteCards := ExtractSuiteCards(cards, wanted_suite)
	lowest := playedSuiteCards[0]
	for _, c := range playedSuiteCards {
		if c.Value > lowest.Value {
			lowest = c
		}
	}
	return lowest
}

func ExtractDeckSuites(cards *[]Card) []Suite {
	suites_map := make(map[Suite]bool)
	for _, c := range *cards {
		suites_map[c.Suit] = true
	}
	keys := make([]Suite, 0, len(suites_map))
	for s := range suites_map {
		keys = append(keys, s)
	}
	return keys
}
