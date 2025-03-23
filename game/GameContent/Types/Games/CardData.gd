class_name CardData extends Resource


@export var suit:int
@export var value:int

enum CardSuites { # sync it with internal/card/card.go
	Diamond = 0,
	Heart = 1,
	Spade = 2,
	Club = 3

}