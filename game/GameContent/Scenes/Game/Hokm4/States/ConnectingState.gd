extends State

@export var ws:KatanaSocket
@export var status_label:Label


func Enter()->void:
	ws.hold_events()
	status_label.text = "Connecting..."
	var state:int = ws.socket.get_ready_state()
	while state != WebSocketPeer.STATE_OPEN:
		await get_tree().create_timer(.2).timeout
		state = ws.socket.get_ready_state()
	Transition.emit(self,"match_making")