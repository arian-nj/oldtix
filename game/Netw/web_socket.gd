class_name KatanaSocket extends Node

var socket:WebSocketPeer = WebSocketPeer.new()

class Event:
	var type:String
	var data:String

	func get_string()->String:
		var data_dic:Dictionary = {
			"type":type,
			"data":data
		}
		return JSON.stringify(data_dic)

func send_event(type:String,data:String)->void:
	var new_event:Event = Event.new()
	new_event.data = data 
	new_event.type = type
	socket.send_text(new_event.get_string())

func _ready() -> void:
	print("starting connction")
	var err:int = socket.connect_to_url(Katana.WsBaseUrl+"/ws?auth_token="+Katana.Auth_Token)
	if err != OK:
		print("Unable to connect")
		set_process(false)
		return

	# Wait for the socket to connect.
	await get_tree().create_timer(1).timeout

func _process(_delta:float)->void:
	socket.poll()

	var state:int = socket.get_ready_state()

	if state == WebSocketPeer.STATE_OPEN:
		while socket.get_available_packet_count():
			print("Got data from server: ", socket.get_packet().get_string_from_utf8())

	elif state == WebSocketPeer.STATE_CLOSING:
		pass
	
	elif state == WebSocketPeer.STATE_CLOSED:
		var code:int = socket.get_close_code()
		print("WebSocket closed with code: %d. Clean: %s" % [code, code != -1])
		set_process(false) # Stop processing.
