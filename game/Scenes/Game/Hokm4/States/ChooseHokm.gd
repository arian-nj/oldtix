extends State

@export var ws:KatanaSocket
@export var status_label:Label
@export var card_drawer:CardDrawer
@export var game_table:Game4Table


func Enter()->void:
	status_label.text = "Choose_Hokm"
	ws.send_event(KEvent.TYPE_GET_DATA,"game")
	ws.new_event.connect(_on_new_event)

func _on_new_event(e:KEvent.Event)->void:
	if e.type == KEvent.TYPE_GAME_DATA:
		game_table.parse_game_data(e.data)
	# if e.type == KEvent.TYPE_NEW_CARD:
	# 	var json_data:Variant = JSON.parse_string(e.data)
	# 	var cards_json:Variant = json_data["cards"]
	# 	for card_json:Variant in cards_json:
	# 		card_drawer.create_card(card_json["suit"],card_json["value"])
	# 	card_drawer.draw_cards(card_drawer.from.global_position)
		
