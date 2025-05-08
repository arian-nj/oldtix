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
		uid_string = await get_new_guest_uid()
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
	var new_token := await get_guest_token(uid_string)
	if new_token == "":
		print_debug("new token is empty")
		return
	ErrorBoard._instance.new_error("Logged in",ErrorBoard.SuccessLevel)
	KAccount._instance.set_token(new_token)



func load_auth_config() -> String:
	var load_err := authConfig.load(auth_config_address)
	if load_err == OK:
		return ""
	load_err = authConfig.save(auth_config_address)
	if load_err != OK:
		return "can't save file code "+str(load_err)
	return ""

func get_new_guest_uid() -> String:
	var http_req := Katana.NewHttpRequest()
	add_child(http_req)
	var err := http_req.request(Katana.CoreHttpUrl+Katana.CreateGuestRand)
	if err != OK:
		print_debug("request failed with err ", err)
		ErrorBoard._instance.new_error("uid request failed close app and try again",ErrorBoard.ErrorLevel)
		return ""
	var res:Variant = await http_req.request_completed
	http_req.queue_free()

	var result:int = res[0]
	if result != OK:
		print_debug("result is not ok result: ", str(err))
		ErrorBoard._instance.new_error("uid request failed close app and try again r:"+str(result),ErrorBoard.ErrorLevel)
	
	var response_code:int = res[1]
	if response_code != HTTPClient.RESPONSE_CREATED:
		print_debug("response code: ", str(response_code))
		ErrorBoard._instance.new_error("uid failed with rc: "+str(result),ErrorBoard.ErrorLevel)
	# var _headers = response[2] # <-- not used
	
	var body_byte:PackedByteArray = res[3]
	var body_string := body_byte.get_string_from_utf8()

	var guestUIDResponse :CreateGuestUID = JsonClassConverter.json_string_to_class(CreateGuestUID,body_string)
	return guestUIDResponse.uid_string

func get_guest_token(guest_uid:String) -> String:
	var http_req := Katana.NewHttpRequest()
	add_child(http_req)
	
	var request_data_json:Dictionary = {
		"uid_string":guest_uid
	}
	var request_data := JSON.stringify(request_data_json)
	var err := http_req.request(Katana.CoreHttpUrl+Katana.GetGuestToken,[],HTTPClient.METHOD_POST,request_data)
	if err != OK:
		print_debug("request failed with err ", err)
		ErrorBoard._instance.new_error("token request failed close app and try again",ErrorBoard.ErrorLevel)
		return ""
	
	var res:Variant = await http_req.request_completed
	http_req.queue_free()

	var result:int = res[0]
	if result != OK:
		print_debug("result is not ok result: ", str(err))
		ErrorBoard._instance.new_error("token request failed close app and try again r:"+str(result),ErrorBoard.ErrorLevel)
		return ""
	
	var response_code:int = res[1]
	if response_code != HTTPClient.RESPONSE_OK:
		print_debug("response code: ", str(response_code))
		ErrorBoard._instance.new_error("token failed with rc: "+str(result),ErrorBoard.ErrorLevel)
		return ""
	
	var body_byte:PackedByteArray = res[3]
	var body_string := body_byte.get_string_from_utf8()

	var newJwtToken :JwtToken = JsonClassConverter.json_string_to_class(JwtToken,body_string)
	return newJwtToken.token
