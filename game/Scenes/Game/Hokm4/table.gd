class_name Game4Table extends Control

class Player:
	var UserId:int
	var Team:int

class GameData:
	var id:String
	var players:Array[Player]
	var current:int
	var hakem:int
	var hokm:Card.CardSuites

var game_data:GameData = GameData.new()

func parse_game_data(json_string:String)->bool:
	var json := JSON.new()
	var err := json.parse(json_string)
	if err != OK:
		print("JSON Parse Error: ", json.get_error_message(), " in ", json_string, " at line ", json.get_error_line())
		return false
	
	var data:Variant = json.data
	
	game_data.id = data["id"]
	game_data.current = data["current"]
	game_data.hakem = data["hakem"]
	game_data.hokm = data["hokm"]
	# print("Hokm is ",game_data.hokm)
	for player:Variant in data["players"]:
		var p := Player.new()
		p.UserId = player["user_id"]
		p.Team = player["team"]
		game_data.players.append(p)
	return true
