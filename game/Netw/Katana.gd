extends Node

var IPAddr:String = "192.168.121.205"
var PortAddr:String = "4444"
var BaseUrl:String = "http://"+IPAddr+":"+PortAddr

var StatusUrl:String = BaseUrl + "/status"

var RegisterUrl:String = BaseUrl + "/user/register"
var TokenUrl:String = BaseUrl + "/user/token"

var Auth_Token:String = ""

func AddAuthHeader(headers:PackedStringArray)->PackedStringArray:
	headers.append("Authorization: Bearer "+Auth_Token)
	return headers

func set_token(token:String)->void:
	Auth_Token = token
# func _on_status_request_completed(_result: int, _response_code: int, _headers: PackedStringArray, _body: PackedByteArray)->void:
