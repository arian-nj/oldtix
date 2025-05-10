class_name AuthPanel extends PanelContainer

@onready var authConfig := ConfigFile.new()
var auth_config_address := "user://AuthConfig.cfg"

func _extract_auth_config(args:PackedStringArray)->String:
	var arg_name:="--auth_config"
	for arg:String in args:
		if arg.begins_with(arg_name):
			return arg.substr(len(arg_name)+1)
	return ""


func _ready() -> void:
	Katana.change_debug_mode(true)

	var args := OS.get_cmdline_args()
	var new_auth_path := _extract_auth_config(args)
	if new_auth_path != "":
		auth_config_address = new_auth_path
	
	print(new_auth_path)
	var uid_string :String = ""
	var load_err := authConfig.load(auth_config_address)
	if load_err != OK:
		uid_string = await KAccount._instance.get_new_guest_uid()
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
		ErrorBoard._instance.new_error("no guest uid",ErrorBoard.ErrorLevel)
		return
	var is_ok := await KAccount._instance.setup_token(uid_string)
	if is_ok:
		ErrorBoard._instance.new_error("Logged in",ErrorBoard.SuccessLevel)



func load_auth_config() -> String:
	var load_err := authConfig.load(auth_config_address)
	if load_err == OK:
		return ""
	load_err = authConfig.save(auth_config_address)
	if load_err != OK:
		return "can't save file code "+str(load_err)
	return ""
