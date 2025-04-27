extends Control

@export var table:Game4Table

# signal fired

func _ready() -> void:
	var e := KEvent.Event.new()
	e.data = '{"cards":[{"suit":0,"value":3},{"suit":1,"value":9},{"suit":0,"value":2},{"suit":3,"value":5}]}'
	table.new_cards_event(e)

	var cd : CardData = CardData.new()
	cd.suit = CardData.CardSuites.Heart
	cd.value = 10
	table.right_drawer.play_random_card(cd)
	table.top_drawer.play_random_card(cd)
	table.left_drawer.play_random_card(cd)

# func _new_card() -> void:
# 	var e := KEvent.Event.new()
# 	e.data = '{"cards":[{"suit":0,"value":3},{"suit":1,"value":9},{"suit":0,"value":2},{"suit":3,"value":5}]}'
# 	table.new_cards_event(e)

