extends State

@export var ws:KatanaSocket
@export var status_label:Label
@export var table:Game4Table


func Enter()->void:
	status_label.text = "Wait for State"
	ws.new_event.connect(_on_new_event)
	ws.open_events()

func _on_new_event(e:KEvent.Event)->void:
	if e.type == KEvent.TYPE_NEW_TRICK:
		ws.hold_events()
		table.parse_game_data(e.data)
		status_label.text = "new trick"
		Transition.emit(self,"choose_hokm")
	elif e.type == KEvent.TYPE_THE_END:
		ws.hold_events()
		Transition.emit(self,"the_end")



func Exit()->void:
	ws.hold_events()
	ws.new_event.disconnect(_on_new_event)
