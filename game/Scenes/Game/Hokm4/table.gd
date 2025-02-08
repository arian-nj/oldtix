class_name Game4Table extends Control

class Player:
	var UserId:int
	var Team:int

class TableCard:
	var suite:Card.CardSuites
	var value:int

class Turn:
	var player_index:int
	var cards_played:Array[TableCard]
	var hokm:Card.CardSuites


class GameData:
	var id:String
	var players:Array[Player]

	var hakem:int
	var turn:Turn
	

var game_data:GameData = GameData.new()
func _ready() -> void:
	game_data.turn = Turn.new()

func parse_game_data(json_string:String)->bool:
	var json := JSON.new()
	var err := json.parse(json_string)
	if err != OK:
		print("JSON Parse Error: ", json.get_error_message(), " in ", json_string, " at line ", json.get_error_line())
		return false
	
	var json_data:Variant = json.data
	
	game_data.id = json_data["id"]
	game_data.hakem = json_data["hakem"]
	# print("Hokm is ",game_data.hokm)
	for player:Variant in json_data["players"]:
		var p := Player.new()
		p.UserId = player["user_id"]
		p.Team = player["team"]
		game_data.players.append(p)
	
	# turn stuff
	# var v:Variant = json_data.get("turn")
	if json_data.get("turn") != null:
		var turn_json:Variant = json_data["turn"]
		game_data.turn.player_index = turn_json["player_index"]
		game_data.turn.hokm = turn_json["hokm"]
		for card_json:Variant in turn_json["cards"]:
			var tb := TableCard.new()
			tb.suite = card_json["suit"]
			tb.value = card_json["value"]
			game_data.turn.cards_played.append(tb)
		print("got turn data")
	else:
		print("no turn data")
	
	return true
