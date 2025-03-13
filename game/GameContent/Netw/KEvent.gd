class_name KEvent extends Node 

# Event class for encapsulating event data
class Event:
	var type: String
	var data: String

	func to_json() -> String:
		return JSON.stringify({
			"type": type,
			"data": data
		})

	func from_json(json_string: String) -> bool:
		var data_dict:Variant = JSON.parse_string(json_string)
		if data_dict == null:
			return false

		if "type" not in data_dict and "data" not in data_dict:
			return false
		
		type = data_dict["type"]
		data = data_dict["data"]

		return true


const TYPE_CHAT:String = "chat"
const TYPE_GAME_DATA:String = "game_data"
const TYPE_GET_DATA:String = "get_data"
const TYPE_NEW_CARD:String = "new_card"
const TYPE_HOKM_CHOOSED:String = "hokm_choosed"

const TYPE_MAKE_MATCH:String = "make_match"
const TYPE_MATCH_FOUND:String = "found_match"
const TYPE_REJOIN_MATCH:String = "rejoin_match"

const TYPE_GET_MY_CARDS:String = "get_my_cards"

const TYPE_NEW_TRICK:String = "new_trick"
const TYPE_END_TRICK:String = "end_trick"

const TYPE_TURN_START:String = "turn_start"
const TYPE_YOUR_TURN:String = "your_turn"
const TYPE_VALID_PLAY:String = "valid_play"
const TYPE_INVALID_PLAY:String = "invalid_play"
const TYPE_PLAY_TIMEOUT:String = "play_timeout"
const TYPE_TURN_PLAYED:String = "turn_played"
const TYPE_END_TURN:String = "turn_end"