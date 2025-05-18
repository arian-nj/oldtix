class_name KatanaSocket
extends Node

var _ws := WebSocketPeer.new()
var _ws_last_state := WebSocketPeer.STATE_CLOSED
var _timeout : int = 30
var _start : float = 0
var logger := KatanaLogger.new()

## emitted when socket is connected.
signal connected()

## emitted when socket is disconnected.
signal closed()

## emitted when the socket has an error when connecting.
signal received_error_signal(p_exception:int)

## emitted when socket receives a message.
# signal received(p_bytes:PackedByteArray) 

## check is socket connected.
func is_connected_to_host() -> bool:
	return _ws.get_ready_state() == WebSocketPeer.STATE_OPEN

## Is Socket connecting to host.
func is_connecting_to_host() -> bool:
	return _ws.get_ready_state() == WebSocketPeer.STATE_CONNECTING

## close socket
func close() ->void:
	_ws.close()

var events_queue:Array[KEvent.Event] = []

func _ready() -> void:
	process_mode = Node.PROCESS_MODE_ALWAYS

func get_latest_event() -> KEvent.Event:
	return events_queue.pop_back()

func push_event(e:KEvent.Event) -> void:
	events_queue.push_front(e)

## coin_amount : can be null or int
func connect_to_game(coin_amount:Variant,p_timeout : int) -> void: 
	_timeout = p_timeout
	_start = Time.get_unix_time_from_system()
	var ws_url:String = Katana._instance.Hokm4WsUrl + "/ws?auth_token=" + KClient._instance.Auth_Token
	if coin_amount != null:
		ws_url += "&coin_amount=" + str(coin_amount)
	var err := _ws.connect_to_url(ws_url)
	if err != OK:
		Katana._instance.logger.error("Error connecting to host %s" % ws_url)
		call_deferred("emit_signal", "received_error_signal", err)
		received_error_signal.emit.bind(err).call_deferred()
		return
	_ws_last_state = WebSocketPeer.STATE_CLOSED

# func send(p_buffer : PackedByteArray) -> int:
# 	return _ws.send(p_buffer, WebSocketPeer.WRITE_MODE_TEXT)

func send_event(event_type: String, event_data: String="") -> void:
	var event:KEvent.Event = KEvent.Event.new()
	event.type = event_type
	event.data = event_data
	_ws.send_text(event.to_json())


func _process(_delta: float) -> void:
	if _ws.get_ready_state() != WebSocketPeer.STATE_CLOSED:
		_ws.poll()

	var current_state := _ws.get_ready_state()
	if _ws_last_state != current_state:
		_ws_last_state = current_state
		if current_state == WebSocketPeer.STATE_OPEN:
			connected.emit()
		elif current_state == WebSocketPeer.STATE_CLOSED:
			var code:int = _ws.get_close_code()
			var reason:String = _ws.get_close_reason()
			print("WebSocket closed with code: %d, reason: %s" % [code, reason])
			Katana._instance.logger.error("WebSocket closed with code: %d, reason: %s" % [code, reason])
			set_process(false)
			closed.emit()

	if current_state == WebSocketPeer.STATE_CONNECTING:
		if _start + _timeout < Time.get_unix_time_from_system():
			logger.debug("Timeout when connecting to socket")
			received_error_signal.emit(ERR_TIMEOUT)
			_ws.close()

	while _ws.get_ready_state() == WebSocketPeer.STATE_OPEN and _ws.get_available_packet_count():
		var message:String = _ws.get_packet().get_string_from_utf8()
		var event:KEvent.Event = KEvent.Event.new()
		if not event.from_json(message):
			print("Failed to parse message: ", message)
		else:
			self.push_event(event)

func _notification(what: int) -> void:
	if what == NOTIFICATION_APPLICATION_FOCUS_OUT:
		# print("out")
		pass
