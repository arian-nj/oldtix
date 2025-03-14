extends State

@export var ws:KatanaSocket
@export var status_label:Label
@export var table:Game4Table


func Enter()->void:
	status_label.text = "The End"
	ws.new_event.connect(_on_new_event)
	ws.open_events()

func _on_new_event(e:KEvent.Event)->void:
	pass

func Exit()->void:
	ws.new_event.disconnect(_on_new_event)
