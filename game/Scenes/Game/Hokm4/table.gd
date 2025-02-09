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