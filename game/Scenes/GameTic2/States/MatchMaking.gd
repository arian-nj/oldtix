extends State

@export var ws:KatanaSocket
@export var status_label:Label


func Enter()->void:
	status_label.text = "looking for opponent"
	ws.send_event(KEvent.TYPE_MAKE_MATCH,"me")
	ws.new_event.connect(_on_new_event)

func _on_new_event(e:KEvent.Event)->void:
	if e.type == KEvent.TYPE_MATCH_FOUND:
		status_label.text = "match Found"