extends Control

@export var table:Game4Table

signal fired

func _ready() -> void:
	fired.connect(_new_card)
	fired.emit()
	fired.emit()
	fired.emit()

func _new_card() -> void:
	var e := KEvent.Event.new()
	e.data = '{"cards":[{"suit":0,"value":3},{"suit":1,"value":9},{"suit":0,"value":2},{"suit":3,"value":5}]}'
	table.new_cards_event(e)

