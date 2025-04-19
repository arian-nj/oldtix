extends CanvasLayer

@export var oldVersionLabel:Label
@export var newVersionLabel:Label
@export var downloadProgressBar:DownloadProgressBar
@export var lpErrorBorad:LPErrorPanel

var server_url := "https://hokm.filelord.ir"
var patch_user_config_address := "user://version.cfg"
var config := ConfigFile.new()

func _ready() -> void:
	get_tree().change_scene_to_file("res://GameContent/Scenes/SceneManager/SceneManager.tscn")
	return
	lpErrorBorad.TryAgain.connect(do_request)
	do_request()

func do_request()->void:
	var err := await get_newest_version()
	if err != "":
		pass
		lpErrorBorad.new_error(err)
	
	
func get_newest_version()-> String:
	
	downloadProgressBar.indeterminate = true
	newVersionLabel.visible = false

	var ov_result := get_local_version()
	var old_version_string:String = ov_result[0]
	var err:String = ov_result[1]

	if err != "":
		return err

	oldVersionLabel.text = "V " + old_version_string

	var nv_result := await get_latest_version_number()
	var new_version_string:String = nv_result[0]
	err = nv_result[1]

	if err != "":
		return err
		
	if new_version_string != old_version_string:
		print("downloading new version")
		downloadProgressBar.indeterminate = false

		newVersionLabel.text = "V " + new_version_string
		newVersionLabel.visible = true
		
		var domain := "cgame.storage.c2.liara.space"
		var pack_url := "/patches/release/Tix_android_v"+new_version_string+".pck"
		pack_url = "/patches/dev/TixGame_android.pck"
		err = await downloadProgressBar.start_downloading(domain,pack_url)
		if err != "":
			return err
	config.set_value("player", "version",new_version_string)
	
	var success := ProjectSettings.load_resource_pack("user://system.pck")
	if !success:
		print_debug("loading failed")
		return "can't load resource pack"
	config.save(patch_user_config_address)
	var change_success := get_tree().change_scene_to_file("res://GameContent/Scenes/SceneManager/SceneManager.tscn")
	if change_success != OK:
		print_debug("change scene failed")
		return "changing scene failed"
	return ""

func get_latest_version_number() -> Array:
	var http_req_node:HTTPRequest = HTTPRequest.new()
	http_req_node.timeout = 5
	self.add_child(http_req_node)
	
	var err := http_req_node.request(server_url+"/version",PackedStringArray(),HTTPClient.METHOD_GET)
	
	if err != OK:
		return ["","can't connect to server"]
	
	var response:Variant = await http_req_node.request_completed
	http_req_node.queue_free()

	var result:int = response[0]
	if result != OK:
		return ["","result is not ok " + str(result)]
	
	var response_code:int = response[1]
	if response_code != HTTPClient.RESPONSE_OK:
		return ["","response is not ok " + str(response_code)
]	
	# var _headers = response[2] # <-- not used
	
	var body_byte:PackedByteArray = response[3]
	var body_string := body_byte.get_string_from_utf8()
	var nj := JSON.new()
	if nj.parse(body_string) != OK:
		print_debug("can't parse")
		return ["","can't parse version request body"]
	var version_string:String =  nj.data.get("version")
	return [version_string,""]

func get_local_version()->Array:
	var load_err := config.load(patch_user_config_address)
	if load_err != OK:
		load_err = config.save(patch_user_config_address)
		if load_err != OK:
			return ["","can't save file code "+str(load_err)]

	var old_version :String = config.get_value("player", "version","0.0.0")
	return [old_version,""]
