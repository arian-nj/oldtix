extends PanelContainer

# @export var usernameField:LineEdit
# @export var passwordField:LineEdit

# func _ready()->void:
# 	var r:int = randi()%50
# 	usernameField.text += str(r)
# 	# _on_login_button_pressed()

# func _on_register_button_pressed() -> void:
# 	var username:String = usernameField.text
# 	var password:String = passwordField.text

# 	var json_data:Dictionary = {
# 		"username":username,
# 		"password":password,
# 	}
# 	var json_data_string:String = JSON.stringify(json_data)

# 	var http_req:HTTPRequest = Katana.NewHttpRequest()
# 	add_child(http_req)
# 	http_req.request_completed.connect(_on_register_request_completed)

# 	var err:int = http_req.request(Katana.CoreHttpUrl + Katana.RegisterUrl,[],HTTPClient.METHOD_POST,json_data_string)
# 	if err != OK:
# 		print_debug(err)
# 		ErrorBoard._instance.new_error("error sending register request",ErrorBoard.ErrorLevel)

# 	await http_req.request_completed
# 	http_req.queue_free()

# func _on_register_request_completed(_result: int, response_code: int, _headers: PackedStringArray, body: PackedByteArray)->void:
# 	if response_code == HTTPClient.RESPONSE_CREATED:
# 		ErrorBoard._instance.new_error("account registered",ErrorBoard.SuccessLevel)
# 	elif response_code == HTTPClient.RESPONSE_UNPROCESSABLE_ENTITY:
# 		var body_content:Variant = JSON.parse_string(body.get_string_from_utf8())
# 		if body_content.has("FieldErrors"):
# 			for k:String in body_content["FieldErrors"]:
# 				ErrorBoard._instance.new_error(body_content["FieldErrors"][k],ErrorBoard.InfoLevel)

# 	else:
# 		ErrorBoard._instance.new_error("failed register",ErrorBoard.ErrorLevel)
# 		print_debug(str(_result)+" " + str(response_code)," ",body.get_string_from_utf8())


# func _on_login_button_pressed() -> void:
# 	var username:String = usernameField.text
# 	var password:String = passwordField.text

# 	var json_data:Dictionary = {
# 		"username":username,
# 		"password":password,
# 	}
# 	var json_data_string:String = JSON.stringify(json_data)

# 	var http_req:HTTPRequest = Katana.NewHttpRequest()
# 	http_req.timeout = 5
# 	add_child(http_req)
# 	http_req.request_completed.connect(_on_token_request_completed)

# 	var err:int = http_req.request(Katana.CoreHttpUrl + Katana.TokenUrl,[],HTTPClient.METHOD_POST,json_data_string)
# 	if err != OK:
# 		print_debug(err)
# 		ErrorBoard._instance.new_error("error sending login request",ErrorBoard.ErrorLevel)
# 	await http_req.request_completed
# 	http_req.queue_free()


# func _on_token_request_completed(_result: int, response_code: int, _headers: PackedStringArray, body: PackedByteArray)->void:
# 	if response_code == HTTPClient.RESPONSE_OK:
# 		ErrorBoard._instance.new_error("You're In",ErrorBoard.SuccessLevel)
# 		var tokenBodyJson :Variant = JSON.parse_string(body.get_string_from_utf8())
# 		var new_token:String = tokenBodyJson["token"]	
# 		KAccount._instance.set_token(new_token)
# 		self.visible = false


# 	elif response_code == HTTPClient.RESPONSE_UNPROCESSABLE_ENTITY:
# 		var body_content:Variant = JSON.parse_string(body.get_string_from_utf8())
# 		if body_content.has("FieldErrors"):
# 			for k:String in body_content["FieldErrors"]:
# 				ErrorBoard._instance.new_error(body_content["FieldErrors"][k],ErrorBoard.InfoLevel)
		
# 	else:
# 		ErrorBoard._instance.new_error("failed login r:" + str(_result)+" rc" + str(response_code),ErrorBoard.ErrorLevel)
# 		print_debug("token failed ",str(_result)+" " + str(response_code)," ",body.get_string_from_utf8())
	


# func _on_local_button_pressed() -> void:
# 	Katana.change_debug_mode(true)
