class_name AuthPanel extends PanelContainer

@onready var authConfig := ConfigFile.new()
var auth_config_address := "user://AuthConfig.cfg"

var config_arg_name:="--auth_config"

func _extract_auth_config(args:PackedStringArray)->String:
	for arg:String in args:
		if arg.begins_with(config_arg_name):
			return arg.substr(len(config_arg_name)+1)
	return ""


func _ready() -> void:
	var os_args := OS.get_cmdline_args()
	for arg in os_args:
		if arg == "--dev":
			Katana._instance.change_debug_mode(true)

	var args := OS.get_cmdline_args()
	var new_auth_path := _extract_auth_config(args)
	if new_auth_path != "":
		auth_config_address = new_auth_path
	
	print(new_auth_path)
	var uid_string :String = ""
	var load_err := authConfig.load(auth_config_address)
	if load_err != OK:
		uid_string = await KClient._instance.get_new_guest_uid()
		if uid_string == "":
			print_debug("new guest uid is empty")
			return
		
		authConfig.set_value("p1","uid_string",uid_string)
		var err :=  authConfig.save(auth_config_address)
		if err != OK:
			print("can't save ", err)
			return
	else: # if auth config exist
		uid_string = authConfig.get_value("p1","uid_string")

	if uid_string == "":
		print_debug("no guest uid")
		Katana._instance.logger.error("no guest uid")
		return
	var is_ok := await KClient._instance.setup_token(uid_string)
	if is_ok:
		Katana._instance.logger.success("Logged in")



func load_auth_config() -> String:
	var load_err := authConfig.load(auth_config_address)
	if load_err == OK:
		return ""
	load_err = authConfig.save(auth_config_address)
	if load_err != OK:
		return "can't save file code "+str(load_err)
	return ""
