extends State

@export var ws:KatanaSocket
@export var status_label:Label


func Enter()->void:
	status_label.text = "Connecting..."
	var state:int = ws.socket.get_ready_state()
	if state != WebSocketPeer.STATE_OPEN:
		await get_tree().create_timer(.2).timeout
	Transition.emit(self,"match_making")
