class_name Katana extends Node

const REQUEST_TIMEOUT := 5
static var _debug:bool = false

static var CoreHttpUrl:String = "https://" + "core.filelord.ir"

static var Hokm4WsUrl:String =  "wss://" + "hokm4.filelord.ir"
static var Hokm4HttpUrl:String = "https://"+"hokm4.filelord.ir"

static func change_debug_mode(new_debug:bool) -> void:
	_debug = new_debug
	if new_debug:
		var local_ip := "192.168.216.205"
		CoreHttpUrl = "http://" + local_ip +":4444"
		Hokm4WsUrl =  "ws://" + local_ip + ":4445"
		Hokm4HttpUrl =  "http://" + local_ip + ":4445"

## Consts
static var RegisterUrl:= "/register"
static var TokenUrl:=  "/token"

static var PersonUrl:= "/person/"
static var MeUrl:= "/me"
static var StatusUrl := "/status"
static var PersonStatisticsAfter := "/stat"

static var ActiveGameUrl:String = "/active_game"

static func NewHttpRequest()->HTTPRequest:
	var http_req := HTTPRequest.new()
	http_req.timeout = REQUEST_TIMEOUT
	return http_req
