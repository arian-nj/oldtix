package cards

import (
	"fmt"
	"math/rand"
)

const (
	HEART  = "H"
	SPADE  = "S"
	DIMOND = "D"
	CLUBS  = "C"
)
const (
	ACE = iota + 1
	N2
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
)

type Suite string
type CardValue int

var Suits = []Suite{HEART, SPADE, DIMOND, CLUBS}
var Values = []CardValue{ACE, FK, FQ, FJ, N10, N9, N8, N7, N6, N5, N4, N3, N2}

type Card struct {
	Suit  Suite
	Value CardValue
}

var AllCards = []Card{}

func init() {
	for _, s := range Suits {
		for _, v := range Values {
			AllCards = append(AllCards, Card{Suit: s, Value: v})
		}
	}

}

func GiveRandomCards(numberOfCards int, availableCards []Card) (randomCards []Card, remaningCards []Card, err error) {
	for range numberOfCards {
		if len(availableCards) <= 0 {
			return randomCards, availableCards, fmt.Errorf("available cards is %d cant give more", len(availableCards))
		}
		random_index := rand.Intn(len(availableCards))
		random_card := availableCards[random_index]
		availableCards = append(availableCards[:random_index], availableCards[random_index+1:]...)
		randomCards = append(randomCards, random_card)
	}
	return randomCards, availableCards, nil
}
