extends State

@export var ws:KatanaSocket
@export var status_label:Label
@export var game_table:Game4Table
@export var card_drawer:CardDrawer


func Enter()->void:
	status_label.text = "looking for opponent"
	ws.new_event.connect(_on_new_event)
	ws.send_event(KEvent.TYPE_MAKE_MATCH,"me")

func _on_new_event(e:KEvent.Event)->void:
	if e.type == KEvent.TYPE_GAME_DATA:
		game_table.parse_game_data(e.data)
		# print("parsed game_data in match making")
	
	elif e.type == KEvent.TYPE_NEW_CARD:
		card_drawer.new_cards_event(e)
		Transition.emit(self,"choose_hokm")

func Exit()->void:
	ws.new_event.disconnect(_on_new_event)
