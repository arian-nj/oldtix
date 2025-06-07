class_name Katana extends Node

const REQUEST_TIMEOUT := 5
var _debug:bool = false


var CoreHttpUrl:String = "https://" + "core.filelord.ir"

var Hokm4WsUrl:String =  "wss://" + "hokm4.filelord.ir"
var Hokm4HttpUrl:String = "https://"+"hokm4.filelord.ir"

var logger := KatanaLogger.new(KatanaLogger.LOG_LEVEL.DEBUG)

static var _instance:Katana
func _ready() -> void:
	_instance = self

func change_debug_mode(new_debug:bool) -> void:
	_debug = new_debug
	if new_debug:
		var local_ip := "127.0.0.1"
		CoreHttpUrl = "https://" + local_ip +":4446"
		Hokm4WsUrl =  "wss://" + local_ip + ":4445"
		Hokm4HttpUrl =  "https://" + local_ip + ":4445"

## Consts
# static var RegisterUrl:= "/register"
# static var TokenUrl:=  "/token"

# var CreateGuestRand := "/auth/guest/create"
# var GetGuestToken := "/auth/guest/token"

var GetTelegramToken := "/auth/telegram/token"

var PersonUrl:= "/person/"
var MeUrl:= "/me"
var StatusUrl := "/status"
var PersonStatisticsAfter := "/stat"

var ActiveGameUrl:String = "/active_game"

static func NewHttpRequest()->HTTPRequest:
	var http_req := HTTPRequest.new()
	http_req.timeout = REQUEST_TIMEOUT
	return http_req
