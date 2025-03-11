extends Control
class_name Game4Table

signal MyCardPlayed(card:Card)

@export var drawer:CardDrawer
@export var right_drawer:CardDrawer
@export var top_drawer:CardDrawer
@export var left_drawer:CardDrawer

var all_drawers :Array[CardDrawer]
var table_draw_queue:Array[Callable]

signal GameDataUpdated(game_data:GameData)

var isMyTurn:bool = false

func set_turn(b:bool)->void:
	isMyTurn = b

var game_data:GameData 

@onready var timer := Timer.new()

func clear_cards()->void:
	push_callback(
		_clear_cards.bind()
	)

func _clear_cards()->void:
	for drawer_queue in all_drawers:
		for c:Card in drawer_queue.cards:
			c.queue_free()
		drawer_queue.cards = []
		

func remove_one_card(card:Card) -> void:
	push_callback(
		card.queue_free.bind()
	)

func push_callback(c:Callable)->void:
	table_draw_queue.push_back(c)

func _on_card_played(card:Card)->void:
	MyCardPlayed.emit(card)

func _ready() -> void:
	game_data = GameData.new()
	drawer.MyCardPlayed.connect(_on_card_played)
	game_data.current_trick = TrickData.new()
	
	all_drawers.append(drawer)
	all_drawers.append(right_drawer)
	all_drawers.append(top_drawer)
	all_drawers.append(left_drawer)
	

	for dra in all_drawers:
		dra.AddToQueue.connect(push_callback)

	timer.wait_time = .2
	timer.timeout.connect(run_actions)
	add_child(timer)
	timer.start()

func set_player_to_hand()->void:
	var meId := KAccount._instance.MyAccount.id

	var mePlayer : PlayerData = null
	var beforeMe : Array[PlayerData]
	var afterMe : Array[PlayerData]

	for player in game_data.players:
		if player.user_id == meId:
			print_debug("user id is ",player.user_id)
			mePlayer = player
			continue
		if mePlayer == null:
			beforeMe.append(player)
		else:
			afterMe.append(player)

	var player_order_from_me :Array[PlayerData]

	player_order_from_me.append_array(afterMe)
	player_order_from_me.append_array(beforeMe)

	if len(player_order_from_me)!= 3:
		print_debug("game data player from me is ",len(game_data.players))
		print_debug("player order from me is ",len(player_order_from_me))

	drawer.unique_string = mePlayer.player_unique
	right_drawer.unique_string = player_order_from_me[0].player_unique
	top_drawer.unique_string = player_order_from_me[1].player_unique
	left_drawer.unique_string = player_order_from_me[2].player_unique

func run_actions()->void:
	timer.timeout.disconnect(run_actions)
	timer.stop()
	while true:
		var action_variant:Variant = table_draw_queue.pop_front()
		if action_variant == null:
			break
		var action:Callable = action_variant

		await action.call()
	timer.timeout.connect(run_actions)
	timer.start()

func new_cards_event(e:KEvent.Event)->void:
	push_callback(_new_cards_event.bind(e))

func _new_cards_event(e:KEvent.Event)->void:
	drawer.new_cards_event(e)
	drawer.break_action()

	right_drawer.new_cards_event(e)
	right_drawer.break_action()

	top_drawer.new_cards_event(e)
	top_drawer.break_action()

	left_drawer.new_cards_event(e)
	left_drawer.break_action()

func parse_game_data(json_string:String)->void:
	game_data = JsonClassConverter.json_string_to_class(GameData,json_string)
	GameDataUpdated.emit(game_data)
