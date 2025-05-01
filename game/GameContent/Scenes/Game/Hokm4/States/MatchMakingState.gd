extends State

@export var ws:KatanaSocket
@export var status_label:Label
@export var table:Game4Table


func Enter()->void:
	status_label.text = "looking for opponent"
	ws.NewEventSig.connect(_on_new_event)
	ws.open_events()
	ws.send_event(KEvent.TYPE_MAKE_MATCH,"me")

func _on_new_event(e:KEvent.Event)->void:
	if e.type == KEvent.TYPE_MATCH_FOUND:
		ws.hold_events()
		table.parse_game_data(e.data)
		table.set_player_to_hand()

		status_label.text = "waiting new trick"
		Transition.emit(self,"wait_state")

	elif e.type == KEvent.TYPE_REJOIN_MATCH:
		print("rejoining request")
		status_label.text = "rejoining"
		table.parse_game_data(e.data)
		table.set_player_to_hand()
		
	elif e.type == KEvent.TYPE_GET_MY_CARDS: # move it to game turn
		status_label.text = "Fetching cards"
		table.new_cards_event(e,true)
		status_label.text = "Wait for new turn"
		Transition.emit(self,"game_turn")

func Exit()->void:
	ws.hold_events()
	ws.NewEventSig.disconnect(_on_new_event)
