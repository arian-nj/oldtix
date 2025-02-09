class_name Game4Table extends Control


	

var game_data:GameData = GameData.new()

func _ready() -> void:
	game_data.current_trick = TrickData.new()

func parse_game_data(json_string:String)->void:
	var json := JSON.new()
	var err := json.parse(json_string)
	if err != OK:
		print("JSON Parse Error: ", json.get_error_message(), " in ", json_string, " at line ", json.get_error_line())
		return
	
	print(json.data)
	game_data = JsonClassConverter.json_to_class(GameData,json.data)
	
	# game_data.id = json_data["id"]
	# # print("Hokm is ",game_data.hokm)
	# for player:Variant in json_data["players"]:
	# 	var p := PlayerData.new()
	# 	p.user_id = player["user_id"]
	# 	p.team = player["team"]
	# 	game_data.players.append(p)
	
	# # turn stuff
	# var trick_json:Variant = json_data.get("current_trick")
	# if trick_json != null:
	# 	game_data.current_trick.hokm = trick_json["hokm"]
	# 	game_data.current_trick.hakem_index = trick_json["hakem_index"]

	# 	var current_turn:Variant = trick_json.get("current_turn")
	# 	if current_turn:
	# 		var new_turn := TurnData.new()
	# 		for card_json:Variant in current_turn["cards"]:
	# 			var tb := CardData.new()
	# 			tb.suite = card_json["suit"]
	# 			tb.value = card_json["value"]
	# 			game_data.current_trick.current_turn.cards_played.append(tb)
	# 		game_data.current_trick.current_turn = new_turn
	# 		print("got turn data")
	# 	else:
	# 		print("no turn data")
	
	# return true
