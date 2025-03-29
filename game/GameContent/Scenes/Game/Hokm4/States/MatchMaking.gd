extends State

@export var ws:KatanaSocket
@export var status_label:Label
@export var table:Game4Table


func Enter()->void:
	status_label.text = "looking for opponent"
	ws.new_event.connect(_on_new_event)
	ws.open_events()
	ws.send_event(KEvent.TYPE_MAKE_MATCH,"me")

func _on_new_event(e:KEvent.Event)->void:
	if e.type == KEvent.TYPE_MATCH_FOUND:
		table.parse_game_data(e.data)
		table.set_player_to_hand()

		status_label.text = "waiting new trick"
		ws.hold_events()
		Transition.emit(self,"wait_state")

	elif e.type == KEvent.TYPE_REJOIN_MATCH:
		print("rehoining request")
		status_label.text = "rejoining"
		table.parse_game_data(e.data)
		table.set_player_to_hand()
		# table.
		ws.send_event(KEvent.TYPE_GET_MY_CARDS)
	
	elif e.type == KEvent.TYPE_GET_MY_CARDS: # move it to game turn
		status_label.text = "Fetching cards"
		table.new_cards_event(e)
		status_label.text = "Wait for new turn"
		Transition.emit(self,"game_turn")


func Exit()->void:
	ws.new_event.disconnect(_on_new_event)
