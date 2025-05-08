extends Control

@export var cards:Array[Card]

func _ready() -> void:
	get_tree().create_timer(1).timeout.connect(flip_cards)

func flip_cards()->void:
	print("fliping")
	var new_val := randi_range(11,13)
	for card_index in len(cards):
		var card := cards[card_index]
		card.card_data.value = new_val
		card.prespective3DShader.flip_y(.1, card_index/2.0, card.load_assets)
