extends State

@export var ksocket:KatanaSocket
@export var status_label:Label
@export var table:Game4Table


func Enter()->void:
	status_label.text = "Wait for State"

func _process(_delta: float) -> void:
	var new_event := ksocket.get_latest_event()
	if new_event == null:
		return
	
	if new_event.type == KEvent.TYPE_NEW_TRICK:
		table.parse_game_data(new_event.data)
		status_label.text = "new trick"
		StateTransition.emit(self,"choose_hokm")
		
	elif new_event.type == KEvent.TYPE_THE_END:
		StateTransition.emit(self,"the_end")
	else:
		ksocket.push_event(new_event)


func Exit()->void:
	pass
