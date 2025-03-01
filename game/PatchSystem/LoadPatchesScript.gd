extends CanvasLayer

@export var status_label:Label

var server_url := "http://192.168.171.205:4444"
var patch_user_config_address := "user://score.cfg"
var config := ConfigFile.new()

func _ready() -> void:
	var err := config.load(patch_user_config_address)
	if err != OK:
		err = config.save(patch_user_config_address)
		if err != OK:
			return

	var new_version_string:String = await check_version()
	if new_version_string == "":
		print_debug("version is empty")
		return
	print(new_version_string)
	var old_version :String = config.get_value("player", "version","0.0.0")
	config.set_value("player", "version", new_version_string)

	if new_version_string != old_version:
		print("downloading new version")
		var dl_success := await download("https://cgame.storage.c2.liara.space/patches/dev/GameContentV_"+new_version_string+".pck")
		if !dl_success:
			print_debug("download failed")
			return
		print("download success")	
	
	var success := ProjectSettings.load_resource_pack("user://system.pck")
	if !success:
		print_debug("loading failed")
		return
	var change_success := get_tree().change_scene_to_file("res://GameContent/Scenes/SceneManager/SceneManager.tscn")
	if change_success != OK:
		print_debug("change scene failed")
		return
	config.save(patch_user_config_address)


func download(link:String)->bool:
	var http_req := HTTPRequest.new()
	add_child(http_req)
	var request := http_req.request(link)
	if request != OK:
		print_debug("Http request error")
	status_label.text = "waiting"
	var response :Variant = await http_req.request_completed
	var result: int = response[0]
	var response_code: int = response[1]
	# var _headers: PackedStringArray = response[2]
	var body: PackedByteArray = response[3]
	
	if result != HTTPRequest.RESULT_SUCCESS:
		print_debug(result)
	
	if response_code != HTTPClient.RESPONSE_OK:
		print_debug(response_code)
		return false
	
	var f := FileAccess.open("user://system.pck",FileAccess.WRITE)
	f.store_buffer(body)
	f.close()
	return true
	
	

func check_version() -> String:
	var http_req_node:HTTPRequest = HTTPRequest.new()
	self.add_child(http_req_node)
	
	var err := http_req_node.request(server_url+"/version",PackedStringArray(),HTTPClient.METHOD_GET)
	
	if err != OK:
		print_debug("error code not ok ",err)
		return ""
	
	var response:Variant = await http_req_node.request_completed
	http_req_node.queue_free()

	var result:int = response[0]
	if result != OK:
		print_debug("result code not ok ",result)
	
	var response_code:int = response[1]
	if response_code != HTTPClient.RESPONSE_OK:
		print_debug("response code not ok ",response_code)
		return ""
	
	# var _headers = response[2] # <-- not used
	
	var body_byte:PackedByteArray = response[3]
	var body_string := body_byte.get_string_from_utf8()
	var nj := JSON.new()
	if nj.parse(body_string) != OK:
		print_debug("can't parse")
		return ""
	var version_string:String =  nj.data.get("version")
	return version_string
	
