class_name KatanaSocket extends Node

signal DisconnectedSig
signal NewEventSig(e:KEvent.Event)

var pause :bool = true
var state:int
var events_queue:Array[KEvent.Event] = []

# WebSocket instance
var socket:WebSocketPeer = WebSocketPeer.new()

func _ready() -> void:
	process_mode = Node.PROCESS_MODE_ALWAYS

func hold_events()->void:
	pause = true

func open_events()->void:
	pause = false
	
func send_event(event_type: String, event_data: String="") -> void:
	var event:KEvent.Event = KEvent.Event.new()
	event.type = event_type
	event.data = event_data
	socket.send_text(event.to_json())

func connect_to_game(coin_amount:int) -> void:
	print("Starting connection...")
	var ws_url:String = Katana._instance.Hokm4WsUrl + "/ws?auth_token=" + KClient._instance.Auth_Token +"&coin_amount=" + str(coin_amount)
	var err:int = socket.connect_to_url(ws_url)
	if err != OK:
		print("Unable to connect err code " + str(err) +" " + error_string(err))
		# get_tree().create_timer()

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
	
	process_events(_delta)

var time_passed:float

func process_events(_delta: float) -> void:
	time_passed += _delta
	if pause:
		return
	if time_passed >= .5:
		return
	time_passed = 0	

	var event : KEvent.Event = events_queue.pop_front()
	if event == null:
		return
	# print(KClient._instance.MyAccount.username + " process " +event.type)
	NewEventSig.emit(event)	

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
	events_queue.push_back(event)

# Handles the closed state
func _handle_closed_state() -> void:
	var code:int = socket.get_close_code()
	var reason:String = socket.get_close_reason()
	print("WebSocket closed with code: %d, reason: %s" % [code, reason])
	Katana._instance.logger.error("WebSocket closed with code: %d, reason: %s" % [code, reason])
	set_process(false)
	DisconnectedSig.emit()

func _notification(what: int) -> void:
	if what == NOTIFICATION_APPLICATION_FOCUS_OUT:
		# print("out")
		pass
	
