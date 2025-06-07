class_name KClient extends Node

static var _instance:KClient

func _ready() -> void:
	_instance = self

# Global value
var MyAccount:AccountData
var Auth_Token:String = ""

## Signals
signal LoggedIn
signal MeChanged

# token
func set_token(token:String)->void:
	if token == "":
		return
	Auth_Token = token
	LoggedIn.emit()

# header
func AddAuthHeader(headers:PackedStringArray = [])->PackedStringArray:
	headers.append("Authorization: Bearer "+Auth_Token)
	return headers

func RefetchME()->bool:
	var me:AccountData = await GetUser("me")
	if me == null:
		print_debug("can't get user data")
		return false
	MyAccount = me
	MeChanged.emit()
	return true

func GetUser(user_id:String)->AccountData:
	var http_req_node:HTTPRequest = Katana.NewHttpRequest()
	add_child(http_req_node)
	var req_url:= Katana._instance.CoreHttpUrl + Katana._instance.PersonUrl + user_id
	if user_id == "me":
		req_url = Katana._instance.CoreHttpUrl + Katana._instance.MeUrl
	
	var err :int = http_req_node.request(req_url,AddAuthHeader(),HTTPClient.METHOD_GET)
	
	if err != OK:
		print_debug("here1")
		return null
	var response:Variant = await http_req_node.request_completed
	http_req_node.queue_free()

	var result:int = response[0]
	if result != OK:
		print_debug(result)
	
	var response_code:int = response[1]
	if response_code != HTTPClient.RESPONSE_OK:
		print_debug(response_code)
		print_debug("here2")
		return null
	
	# var _headers = response[2] # <-- not used
	
	var body_byte:PackedByteArray = response[3]

	var me :AccountData = AccountData.new()
	me = JsonClassConverter.json_string_to_class(AccountData,body_byte.get_string_from_utf8())
	return me

func GetUserStatistics(user_id:String)->Array: # UserStatisticsData,String
	var http_req_node:HTTPRequest = Katana.NewHttpRequest()
	add_child(http_req_node)
	var req_url:= Katana._instance.CoreHttpUrl + Katana._instance.PersonUrl + user_id + Katana._instance.PersonStatisticsAfter
	# if user_id == "me":
	# 	req_url = Katana.MeUrl
	
	var err :int = http_req_node.request(req_url,AddAuthHeader(),HTTPClient.METHOD_GET)
	
	if err != OK:
		print_debug("here1")
		return [null,"request failed"]
	var response:Variant = await http_req_node.request_completed
	http_req_node.queue_free()

	var result:int = response[0]
	if result != OK:
		print_debug(result)
	
	var response_code:int = response[1]
	if response_code != HTTPClient.RESPONSE_OK:
		print_debug(response_code)
		return [null,"response code is " + str(response_code)]
	
	# var _headers = response[2] # <-- not used
	
	var body_byte:PackedByteArray = response[3]

	var me := UserStatisticsData.new()
	me = JsonClassConverter.json_string_to_class(UserStatisticsData,body_byte.get_string_from_utf8())
	return [me,""]

func get_new_guest_uid() -> String:
	var http_req := Katana.NewHttpRequest()
	add_child(http_req)
	var err := http_req.request(Katana._instance.CoreHttpUrl+Katana._instance.CreateGuestRand)
	if err != OK:
		print_debug("request failed with err ", err)
		Katana._instance.logger.error("uid request failed close app and try again")
		return ""
	var res:Variant = await http_req.request_completed
	http_req.queue_free()

	var result:int = res[0]
	if result != OK:
		print_debug("result is not ok result: ", str(err))
		Katana._instance.logger.error("uid request failed close app and try again r:"+str(result))
	
	# var _headers = response[2] # <-- not used
	
	var body_byte:PackedByteArray = res[3]
	var body_string := body_byte.get_string_from_utf8()
	
	var response_code:int = res[1]
	if response_code != HTTPClient.RESPONSE_CREATED:
		Katana._instance.logger.error("response code: "+ str(response_code) + " "+ body_string)
		Katana._instance.logger.error("uid failed with rc: "+str(result))
	

	var guestUIDResponse :CreateGuestUID = JsonClassConverter.json_string_to_class(CreateGuestUID,body_string)
	return guestUIDResponse.uid_string

func get_guest_token(guest_uid:String) -> String:
	var http_req := Katana.NewHttpRequest()
	add_child(http_req)
	
	var request_data_json:Dictionary = {
		"uid_string":guest_uid
	}
	var request_data := JSON.stringify(request_data_json)
	var err := http_req.request(Katana._instance.CoreHttpUrl+Katana._instance.GetGuestToken,[],HTTPClient.METHOD_POST,request_data)
	if err != OK:
		print_debug("request failed with err ", err)
		Katana._instance.logger.error("token request failed close app and try again")
		return ""
	
	var res:Variant = await http_req.request_completed
	http_req.queue_free()

	var result:int = res[0]
	if result != OK:
		print_debug("result is not ok result: ", str(err))
		Katana._instance.logger.error("token request failed close app and try again r:"+str(result))
		return ""
	
	var response_code:int = res[1]
	if response_code != HTTPClient.RESPONSE_OK:
		print_debug("response code: ", str(response_code))
		Katana._instance.logger.error("token failed with rc: "+str(result))
		return ""
	
	var body_byte:PackedByteArray = res[3]
	var body_string := body_byte.get_string_from_utf8()

	var newJwtToken :JwtToken = JsonClassConverter.json_string_to_class(JwtToken,body_string)
	return newJwtToken.token

# func setup_token(guest_uid:String) -> bool:
# 	var new_token := await KClient._instance.get_guest_token(guest_uid)
# 	if new_token == "":
# 		print_debug("new token is empty")
# 		return false
# 	KClient._instance.set_token(new_token)
# 	get_tree().create_timer(InternalSetting.JWT_EXPIARY_DURATION / 2.0).timeout.connect(self.setup_token.bind(guest_uid))
# 	return true
#
func get_telegram_token() -> String:
	var console := JavaScriptBridge.get_interface("console")
	
	var parent_window := JavaScriptBridge.get_interface("parent")
	var initdata : Variant= parent_window.getInitdata()
	console.log(initdata)

	var http_req_node := Katana.NewHttpRequest()
	add_child(http_req_node)
	var err := http_req_node.request(Katana._instance.CoreHttpUrl + Katana._instance.GetTelegramToken+"?"+initdata)
	if err != OK:
		Katana._instance.logger.error("telegram token request failed with error " + str(err))
		return ""

	var response:Variant = await http_req_node.request_completed
	http_req_node.queue_free()

	var result:int = response[0]
	if result != OK:
		Katana._instance.logger.error("telegram token response failed with error " + str(result))
		return ""
	var response_code:int = response[1]
	if response_code != HTTPClient.RESPONSE_OK:
		Katana._instance.logger.error("telegram token response is not ok " + str(response_code))
		return ""
	
	# var _headers = response[2] # <-- not used
	
	var body_byte:PackedByteArray = response[3]

	var jwtToken := JwtToken.new()
	jwtToken = JsonClassConverter.json_string_to_class(JwtToken,body_byte.get_string_from_utf8())
	return jwtToken.token

func setup_telegram_token() -> bool:
	var new_token := await KClient._instance.get_telegram_token()
	if new_token == "":
		Katana._instance.logger.error("empty token")
		return false
	KClient._instance.set_token(new_token)
	get_tree().create_timer(InternalSetting.JWT_EXPIARY_DURATION / 2.0).timeout.connect(self.setup_telegram_token.bind())
	return true
