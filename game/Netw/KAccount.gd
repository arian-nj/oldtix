class_name KAccountClass extends Node

# Global value
@onready var MyAccount:Account
@onready var Auth_Token:String = ""

## Signals
signal LoggedIn
signal GotMe

## Consts
const RegisterUrl:String = Katana.HttpBaseUrl + "/user/register"
const TokenUrl:String = Katana.HttpBaseUrl + "/user/token"

const MeUrl:String = Katana.HttpBaseUrl + "/user/me"


class Account:
	var id:int
	var username:String
	var display_name:String
	var bio:String
	var coin:int

func _ready() -> void:
	await LoggedIn
	await RefetchME()
	GotMe.emit()

# token
func set_token(token:String)->void:
	Auth_Token = token
	LoggedIn.emit()

# header
func AddAuthHeader(headers:PackedStringArray = [])->PackedStringArray:
	headers.append("Authorization: Bearer "+Auth_Token)
	return headers

func RefetchME()->Account:
	var me:Account = await _get_me()
	if me == null:
		print_debug("can't get user data")
		return
	MyAccount = me
	return MyAccount

func _get_me()->Account:
	var http_req_node:HTTPRequest = HTTPRequest.new()
	Katana.add_child(http_req_node)

	var err :int = http_req_node.request(MeUrl,AddAuthHeader(),HTTPClient.METHOD_GET)
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
	var body_content:Variant = JSON.parse_string(body_byte.get_string_from_utf8())

	var me :Account = Account.new()
	me.username = body_content["username"]
	me.display_name = body_content["display_name"]
	me.bio = body_content["bio"]
	me.coin = body_content["coin"]
	return me