class_name KAccount extends Node

static var _instance:KAccount

static func inst()->KAccount:
	if _instance == null:
		_instance = KAccount.new()
	return _instance

# Global value
var MyAccount:Account
var Auth_Token:String = ""

## Signals
signal LoggedIn

## Consts
const RegisterUrl:String = Katana.UserBaseUrl + "/register"
const TokenUrl:String = Katana.UserBaseUrl + "/token"

const UserUrl:String = Katana.UserBaseUrl + "/get/"
const MeUrl:String = Katana.UserBaseUrl + "/me"
const StatusUrl:String = Katana.UserBaseUrl + "/status"



class Account:
	var id:int
	var username:String
	var display_name:String
	var bio:String
	var coin:int

	func parse_from_json_string(json_string:String)->void:
		var body_content:Variant = JSON.parse_string(json_string)
		if body_content == null:
			print_debug("failed")
			return
		self.id = body_content["id"]
		self.username = body_content["username"]
		self.display_name = body_content["display_name"]
		self.bio = body_content["bio"]
		self.coin = body_content["coin"]

func _ready() -> void:
	await LoggedIn
	

# token
func set_token(token:String)->void:
	Auth_Token = token
	LoggedIn.emit()

# header
func AddAuthHeader(headers:PackedStringArray = [])->PackedStringArray:
	headers.append("Authorization: Bearer "+Auth_Token)
	return headers

func RefetchME()->bool:
	var me:Account = await _get_user("me")
	if me == null:
		print_debug("can't get user data")
		return false
	MyAccount = me
	return true

func GetUser(user_id:String)->Account:
	var user:Account = await _get_user(user_id)
	return user

func _get_user(user_id:String)->Account:
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

	var me :Account = Account.new()
	me.parse_from_json_string(body_byte.get_string_from_utf8())
	return me
