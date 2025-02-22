extends CanvasLayer

@export var status_label:Label

var server_status_url := "http://192.168.174.205:4444"

var DownloadLnk := "https://cgame.storage.c2.liara.space/patches/GameContentV_0.2.0.pck"
func _ready() -> void:
	download(DownloadLnk)

func download(link:String)->void:
	var http := HTTPRequest.new()
	add_child(http)
	http.request_completed.connect(_download_completed)
	var request := http.request(link)
	if request != OK:
		print_debug("Http request error")
	status_label.text = "waiting"

func _download_completed(result: int, response_code: int, _headers: PackedStringArray, body: PackedByteArray )->void:
	if result != HTTPRequest.RESULT_SUCCESS:
		print_debug(result)
	
	if response_code != HTTPClient.RESPONSE_OK:
		print_debug(response_code)
		return
	
	var f := FileAccess.open("user://system.pck",FileAccess.WRITE)
	f.store_buffer(body)
	f.close()
	print("download completed")
	var success := ProjectSettings.load_resource_pack("user://system.pck")
	if success:
		get_tree().change_scene_to_file("res://GameContent/Scenes/SceneManager/SceneManager.tscn")
	else :
		print_debug("failed to load scene")
	

func check_version() -> void:
	var http_req_node:HTTPRequest = HTTPRequest.new()
	self.add_child(http_req_node)
	
	var err := http_req_node.request(server_status_url,PackedStringArray(),HTTPClient.METHOD_GET)
	
	if err != OK:
		print_debug("here1")
		return
	
	var response:Variant = await http_req_node.request_completed
	http_req_node.queue_free()

	var result:int = response[0]
	if result != OK:
		print_debug(result)
	
	var response_code:int = response[1]
	if response_code != HTTPClient.RESPONSE_OK:
		print_debug(response_code)
		print_debug("here2")
		return
	
	# var _headers = response[2] # <-- not used
	
	var body_byte:PackedByteArray = response[3]
	var body_string := body_byte.get_string_from_utf8()
	var nj := JSON.new()
	if nj.parse(body_string) != OK:
		print_debug("can't parse")
		return
	# var version_string:String =  nj.data.get("version")
	
