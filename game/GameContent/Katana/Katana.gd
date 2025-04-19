class_name Katana extends Node


const BaseUrl:String = "localhost"
# const HtppUserUrl:String = "https://user."+BaseUrl

# const WsBaseUrl:String = "wss://hokm4."+BaseUrl

const HttpPreUrl:String = "http://"
const WsPreUrl:String = "ws://"

# const HtppUserUrl:String = "https://"+ "hokm.filelord.ir"
const HtppUserUrl:String = "http://"+ "127.0.0.1:4444"

const WsHokmUrl:String = WsPreUrl+"127.0.0.1:4445"
const HttpHokmUrl:String = HttpPreUrl+"127.0.0.1:4445"


## Consts
const RegisterUrl:String = HtppUserUrl + "/register"
const TokenUrl:String = HtppUserUrl + "/token"

const PersonUrl:String = HtppUserUrl + "/person/"
const MeUrl:String = HtppUserUrl + "/me"
const StatusUrl:String = HtppUserUrl + "/status"
const PersonStatisticsAfter = "/stat"

const ActiveGameUrl:String = HttpHokmUrl + "/active_game"

