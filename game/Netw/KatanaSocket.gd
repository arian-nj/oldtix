class_name KatanaSocket extends Node

signal new_event(e:KEvent.Event)
# func _on_new_event(e:KEvent.Event)->void:
# 	pass

# WebSocket instance
var socket:WebSocketPeer = WebSocketPeer.new()


# Sends an event through the WebSocket
func send_event(event_type: String, event_data: String) -> void:
	var event:KEvent.Event = KEvent.Event.new()
	event.type = event_type
	event.data = event_data
	socket.send_text(event.to_json())

# Initializes the WebSocket connection
func _ready() -> void:
	print("Starting connection...")
	var ws_url:String = Katana.WsBaseUrl + "/ws?auth_token=" + KAccount.Auth_Token
	var err:int = socket.connect_to_url(ws_url)
	if err != OK:
		print("Unable to connect")
		set_process(false)


var state:int
# Processes incoming WebSocket data
func _process(_delta: float) -> void:
	socket.poll()

	state = socket.get_ready_state()
	match state:
		WebSocketPeer.STATE_OPEN:
			_handle_open_state()

		WebSocketPeer.STATE_CLOSING:
			print("WebSocket is closing...")

		WebSocketPeer.STATE_CLOSED:
			_handle_closed_state()

# Handles the open state and processes incoming messages
func _handle_open_state() -> void:
	while socket.get_available_packet_count() > 0:
		var message:String = socket.get_packet().get_string_from_utf8()
		var event:KEvent.Event = KEvent.Event.new()
		if not event.from_json(message):
			print("Failed to parse message: ", message)
		else:
			_handle_event(event)

# Handles incoming events (to be extended)
func _handle_event(event: KEvent.Event) -> void:
	print("Event received: ", event.type, " - ", event.data)
	new_event.emit(event)

# Handles the closed state
func _handle_closed_state() -> void:
	var code:int = socket.get_close_code()
	var reason:String = socket.get_close_reason()
	print("WebSocket closed with code: %d, reason: %s" % [code, reason])
	ErrorBoard.new_error("WebSocket closed with code: %d, reason: %s" % [code, reason],ErrorClass.ErrorLevel)
	set_process(false)
	var node_parent:Node = get_parent()
	if node_parent is SceneLevel:
		var scene_parent:SceneLevel = node_parent
		await  get_tree().create_timer(1).timeout
		scene_parent.manager_change_scene.emit(SceneManger.Levels.MainMenu)

func _notification(what: int) -> void:
	if what == NOTIFICATION_APPLICATION_FOCUS_OUT:
		print("out")
	