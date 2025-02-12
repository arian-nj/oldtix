class_name CardData extends Resource


@export var suit:int
@export var value:int

enum CardSuites { # sync it with internal/card/card.go
	Heart = 0,
	Club = 1,
	Diamond = 2,
	Spade = 3
}