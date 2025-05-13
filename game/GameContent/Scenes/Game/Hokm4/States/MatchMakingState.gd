extends State

@export var ksocket:KatanaSocket
@export var status_label:Label
@export var table:Game4Table


func Enter()->void:
	status_label.text = "looking for opponent"
	ksocket.send_event(KEvent.TYPE_MAKE_MATCH,"me")

func _process(_delta:float)->void:
	var new_event := ksocket.get_latest_event()
	if new_event == null:
		return
	
	if new_event.type == KEvent.TYPE_MATCH_FOUND:
		table.parse_game_data(new_event.data)
		table.set_player_to_hand()

		status_label.text = "waiting new trick"
		StateTransition.emit(self,"wait_state")

	elif new_event.type == KEvent.TYPE_REJOIN_MATCH:
		table.rejoined = true
		print("rejoining request")
		table.parse_game_data(new_event.data)
		table.set_player_to_hand()

		status_label.text = "rejoining"
		StateTransition.emit(self,"wait_state")
	else:
		ksocket.push_event(new_event)
