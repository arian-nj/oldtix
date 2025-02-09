extends State


@export var ws:KatanaSocket
@export var game_table:Game4Table
@export var status_label:Label

func Enter()->void:
	status_label.text = "Turn Started"
	ws.new_event.connect(_on_new_event)

func Exit()->void:
	ws.new_event.disconnect(_on_new_event)

func _on_new_event(e:KEvent.Event)->void:
	print_debug(e.type)
	if e.type == KEvent.TYPE_TURN_START:
		game_table.parse_game_data(e.data)
	# elif e.type == KEvent.TYPE_TURN_START:
	# 	game_table.parse_game_data(e.data)