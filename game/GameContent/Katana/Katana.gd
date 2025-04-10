class_name Katana extends Node


const BaseUrl:String = "localhost"
# const HtppUserUrl:String = "https://user."+BaseUrl

# const WsBaseUrl:String = "wss://hokm4."+BaseUrl

const HttpPreUrl:String = "http://"
const WsPreUrl:String = "ws://"

const HtppUserUrl:String = HttpPreUrl + "192.168.188.205:4444"

const WsHokmUrl:String = WsPreUrl+"192.168.188.205:4445"
const HttpHokmUrl:String = HttpPreUrl+"192.168.188.205:4445"


## Consts
const RegisterUrl:String = HtppUserUrl + "/register"
const TokenUrl:String = HtppUserUrl + "/token"

const PersonUrl:String = HtppUserUrl + "/person/"
const MeUrl:String = HtppUserUrl + "/me"
const StatusUrl:String = HtppUserUrl + "/status"
const PersonStatisticsAfter = "/stat"

const ActiveGameUrl:String = HttpHokmUrl + "/active_game"

