class_name KAccount extends Node

static var _instance:KAccount

static func instanciate()->KAccount:
	if _instance == null:
		_instance = KAccount.new()
	return _instance

# Global value
var MyAccount:AccountData
var Auth_Token:String = ""

## Signals
signal LoggedIn

## Consts
const RegisterUrl:String = Katana.UserBaseUrl + "/register"
const TokenUrl:String = Katana.UserBaseUrl + "/token"

const UserUrl:String = Katana.UserBaseUrl + "/get/"
const MeUrl:String = Katana.UserBaseUrl + "/me"
const StatusUrl:String = Katana.UserBaseUrl + "/status"

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
	var me:AccountData = await _get_user("me")
	if me == null:
		print_debug("can't get user data")
		return false
	MyAccount = me
	return true

func GetUser(user_id:String)->AccountData:
	var user:AccountData = await _get_user(user_id)
	return user

func _get_user(user_id:String)->AccountData:
	var http_req_node:HTTPRequest = HTTPRequest.new()
	add_child(http_req_node)
	var req_url:= UserUrl+user_id
	if user_id == "me":
		req_url = MeUrl
	
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
	me = JsonClassConverter.json_string_to_class(AccountData,body_byte.get_string_from_utf8()
)
	return me
