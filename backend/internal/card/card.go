package cards

import (
	"fmt"
	"math/rand"
)

type Suite int
type CardValue int

const ( // sync it with card/card.gd
	HEART  Suite = iota
	CLUB         = 1
	DIMOND       = 2
	SPADE        = 3
)

// const
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
		if len(availableCards) <= 0 {
			return allRandomCards, availableCards, fmt.Errorf("available cards is %d cant give more", len(availableCards))
		}
		random_index := rand.Intn(len(availableCards))
		random_card := availableCards[random_index]
		availableCards = append(availableCards[:random_index], availableCards[random_index+1:]...)
		allRandomCards = append(allRandomCards, random_card)
	}
	return allRandomCards, availableCards, nil
}
