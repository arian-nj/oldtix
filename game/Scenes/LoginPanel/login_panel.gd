extends PanelContainer

@export var usernameField:LineEdit
@export var passwordField:LineEdit

func _ready()->void:
	_on_login_button_pressed()

func _on_register_button_pressed() -> void:
	var username:String = usernameField.text
	var password:String = passwordField.text

	var json_data:Dictionary = {
		"username":username,
		"password":password,
	}
	var json_data_strin:String = JSON.stringify(json_data)

	var http_req:HTTPRequest = HTTPRequest.new()
	add_child(http_req)
	http_req.request_completed.connect(_on_register_request_completed)

	var err:int = http_req.request(Katana.RegisterUrl,[],HTTPClient.METHOD_POST,json_data_strin)
	if err != OK:
		print_debug(err)
		ErrorBoard.new_error("error sending register request",ErrorClass.ErrorLevel)

	await http_req.request_completed
	http_req.queue_free()

func _on_register_request_completed(_result: int, response_code: int, _headers: PackedStringArray, body: PackedByteArray)->void:
	if response_code == HTTPClient.RESPONSE_CREATED:
		ErrorBoard.new_error("account registered",ErrorClass.SuccessLevel)
	elif response_code == HTTPClient.RESPONSE_UNPROCESSABLE_ENTITY:
		var body_content:Variant = JSON.parse_string(body.get_string_from_utf8())
		if body_content.has("FieldErrors"):
			for k:String in body_content["FieldErrors"]:
				ErrorBoard.new_error(body_content["FieldErrors"][k],ErrorClass.InfoLevel)

	else:
		ErrorBoard.new_error("failed",ErrorClass.ErrorLevel)
		print_debug(str(_result)+" " + str(response_code)," ",body.get_string_from_utf8())


func _on_login_button_pressed() -> void:
	var username:String = usernameField.text
	var password:String = passwordField.text

	var json_data:Dictionary = {
		"username":username,
		"password":password,
	}
	var json_data_strin:String = JSON.stringify(json_data)

	var http_req:HTTPRequest = HTTPRequest.new()
	add_child(http_req)
	http_req.request_completed.connect(_on_token_request_completed)

	var err:int = http_req.request(Katana.TokenUrl,[],HTTPClient.METHOD_POST,json_data_strin)
	if err != OK:
		print_debug(err)
		ErrorBoard.new_error("error sending login request",ErrorClass.ErrorLevel)
	
	await http_req.request_completed
	http_req.queue_free()


func _on_token_request_completed(_result: int, response_code: int, _headers: PackedStringArray, body: PackedByteArray)->void:

	if response_code == HTTPClient.RESPONSE_OK:
		ErrorBoard.new_error("You're In",ErrorClass.SuccessLevel)
		var tokenBodyJson :Variant = JSON.parse_string(body.get_string_from_utf8())
		var new_token:String = tokenBodyJson["AuthenticationToken"]	
		Katana.set_token(new_token)
		self.visible = false

	elif response_code == HTTPClient.RESPONSE_UNPROCESSABLE_ENTITY:
		var body_content:Variant = JSON.parse_string(body.get_string_from_utf8())
		if body_content.has("FieldErrors"):
			for k:String in body_content["FieldErrors"]:
				ErrorBoard.new_error(body_content["FieldErrors"][k],ErrorClass.InfoLevel)

	else:
		ErrorBoard.new_error("failed",ErrorClass.ErrorLevel)
		print_debug(str(_result)+" " + str(response_code)," ",body.get_string_from_utf8())
	

	
