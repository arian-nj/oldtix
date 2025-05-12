extends State

@export var ksocket:KatanaSocket
@export var status_label:Label


func Enter()->void:
	status_label.text = "Connecting..."
	var state:int = ksocket._ws.get_ready_state()
	while state != WebSocketPeer.STATE_OPEN:
		await get_tree().create_timer(.1).timeout
		state = ksocket._ws.get_ready_state()
	StateTransition.emit(self,"match_making")