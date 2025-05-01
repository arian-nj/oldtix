extends Control

@export var table:Game4Table

func _ready() -> void:
	table.MyCardPlayed.connect(_on_card_played)

func _on_card_played(card:Card) -> void:
	print(card.card_data.suit, " -- " ,card.card_data.value)

func _on_add_card_button_pressed() -> void:
	var e := KEvent.Event.new()
	e.data = '{"cards":[{"suit":0,"value":3},{"suit":1,"value":9},{"suit":0,"value":2},{"suit":3,"value":5}]}'
	table.new_cards_event(e)

func _on_toggle_turn_button_pressed() -> void:
	table.isMyTurn = !table.isMyTurn
