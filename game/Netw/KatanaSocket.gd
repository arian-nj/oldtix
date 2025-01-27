class_name KatanaSocket extends Node

signal new_event(e:Event)

# WebSocket instance
var socket:WebSocketPeer = WebSocketPeer.new()

# Event class for encapsulating event data
class Event:
	var type: String
	var data: String

	func to_json() -> String:
		return JSON.stringify({
			"type": type,
			"data": data
		})

	func from_json(json_string: String) -> bool:
		var data_dict:Variant = JSON.parse_string(json_string)
		if data_dict == null:
			return false

		if "type" not in data_dict and "data" not in data_dict:
			return false
		
		type = data_dict["type"]
		data = data_dict["data"]

		return true

# Sends an event through the WebSocket
func send_event(event_type: String, event_data: String) -> void:
	var event:Event = Event.new()
	event.type = event_type
	event.data = event_data
	socket.send_text(event.to_json())

# Initializes the WebSocket connection
func _ready() -> void:
	print("Starting connection...")
	var ws_url:String = Katana.WsBaseUrl + "/ws?auth_token=" + Katana.Auth_Token
	var err:int = socket.connect_to_url(ws_url)
	if err != OK:
		print("Unable to connect")
		set_process(false)



# Processes incoming WebSocket data
func _process(_delta: float) -> void:
	socket.poll()

	match socket.get_ready_state():
		WebSocketPeer.STATE_OPEN:
			pass
			# _handle_open_state()

		WebSocketPeer.STATE_CLOSING:
			print("WebSocket is closing...")

		WebSocketPeer.STATE_CLOSED:
			_handle_closed_state()

# Handles the open state and processes incoming messages
func _handle_open_state() -> void:
	while socket.get_available_packet_count() > 0:
		var message:String = socket.get_packet().get_string_from_utf8()
		print("Received data: ", message)
		var event:Event = Event.new()
		if not event.from_json(message):
			print("Failed to parse message: ", message)
		else:
			_handle_event(event)

# Handles incoming events (to be extended)
func _handle_event(event: Event) -> void:
	print("Event received: ", event.type, " - ", event.data)
	new_event.emit(event)

# Handles the closed state
func _handle_closed_state() -> void:
	var code:int = socket.get_close_code()
	var reason:String = socket.get_close_reason()
	print("WebSocket closed with code: %d, reason: %s" % [code, reason])
	set_process(false)