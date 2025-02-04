extends State

@export var ws:KatanaSocket
@export var status_label:Label


func Enter()->void:
	status_label.text = "looking for opponent"
	ws.new_event.connect(_on_new_event)
	ws.send_event(KEvent.TYPE_MAKE_MATCH,"me")

func _on_new_event(e:KEvent.Event)->void:
	if e.type == KEvent.TYPE_MATCH_FOUND:
		Transition.emit(self,"choose_hokm")
