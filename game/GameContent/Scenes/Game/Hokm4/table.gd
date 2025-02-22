class_name Game4Table extends Control

signal GameDataUpdated(game_data:GameData)

var isMyTurn:bool = false

func set_turn(b:bool)->void:
	isMyTurn = b

signal CardPlayed(card:Card)

func _card_played(card:Card)->void:
	CardPlayed.emit(card)

var game_data:GameData = GameData.new()
@export var drawer:CardDrawer

func _ready() -> void:
	drawer.CardPlayed.connect(_card_played)
	game_data.current_trick = TrickData.new()

func parse_game_data(json_string:String)->void:
	var json := JSON.new()
	var err := json.parse(json_string)
	if err != OK:
		print("JSON Parse Error: ", json.get_error_message(), " in ", json_string, " at line ", json.get_error_line())
		return
	
	game_data = JsonClassConverter.json_to_class(GameData,json.data)
	GameDataUpdated.emit(game_data)
