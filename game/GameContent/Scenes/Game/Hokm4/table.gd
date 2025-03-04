extends Control
class_name Game4Table

@export var drawer:CardDrawer
@export var right_drawer:CardDrawer
@export var top_drawer:CardDrawer
@export var left_drawer:CardDrawer

signal GameDataUpdated(game_data:GameData)

var isMyTurn:bool = false

func set_turn(b:bool)->void:
	isMyTurn = b
signal CardPlayed(card:Card)

func _card_played(card:Card)->void:
	CardPlayed.emit(card)

var game_data:GameData = GameData.new()

@onready var timer := Timer.new()
var all_adraw_queue:Array[Variant]

func _ready() -> void:
	drawer.CardPlayed.connect(_card_played)
	game_data.current_trick = TrickData.new()
	
	all_adraw_queue.append(drawer.draw_queue)
	all_adraw_queue.append(right_drawer.draw_queue)
	all_adraw_queue.append(top_drawer.draw_queue)
	all_adraw_queue.append(left_drawer.draw_queue)

	timer.wait_time = .1
	timer.timeout.connect(run_actions)
	add_child(timer)
	timer.start()

func run_actions()->void:
	# print("running action")
	timer.timeout.disconnect(run_actions)
	timer.stop()
	for draw_queue:Array[Callable] in all_adraw_queue:
		while draw_queue.size() > 0:
			var action:Callable = draw_queue.pop_front()
			var result :Variant = await action.call()
			if result is String:
				if result == "break":
					break
	timer.timeout.connect(run_actions)
	timer.start()

func new_cards_event(e:KEvent.Event)->void:
	# print("")
	drawer.new_cards_event(e)
	drawer.push_callback(break_action)
	# await get_tree().create_timer(1).timeout
	right_drawer.new_cards_event(e)
	right_drawer.push_callback(break_action)
	# await get_tree().create_timer(1).timeout
	left_drawer.new_cards_event(e)
	left_drawer.push_callback(break_action)
	# await get_tree().create_timer(1).timeout
	top_drawer.new_cards_event(e)
	top_drawer.push_callback(break_action)

func break_action()->String:
	return "break"

func parse_game_data(json_string:String)->void:
	var json := JSON.new()
	var err := json.parse(json_string)
	if err != OK:
		print("JSON Parse Error: ", json.get_error_message(), " in ", json_string, " at line ", json.get_error_line())
		return
	
	game_data = JsonClassConverter.json_to_class(GameData,json.data)
	GameDataUpdated.emit(game_data)
