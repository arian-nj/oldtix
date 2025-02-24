extends State

@export var ws:KatanaSocket
@export var status_label:Label
@export var table:Game4Table


func Enter()->void:
	status_label.text = "New Trick"
	ws.new_event.connect(_on_new_event)
	ws.open_events()
	ws.send_event(KEvent.TYPE_MAKE_MATCH,"me")

func _on_new_event(e:KEvent.Event)->void:
	if e.type == KEvent.TYPE_NEW_TRICK:
		table.parse_game_data(e.data)
		status_label.text = "new trick"
		ws.hold_events()
		Transition.emit(self,"choose_hokm")


func Exit()->void:
	ws.new_event.disconnect(_on_new_event)
