extends State

@export var ws:KatanaSocket
@export var status_label:Label


func Enter()->void:
	status_label.text = "looking for opponent"
	