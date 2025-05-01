class_name Game4Table extends Control

signal MyCardPlayed(card:Card)

@export var me_drawer:CardDrawer
@export var right_drawer:CardDrawer
@export var top_drawer:CardDrawer
@export var left_drawer:CardDrawer

@export var right_player_panel:UserComponent
@export var left_player_panel:UserComponent
@export var top_player_panel:UserComponent
@export var me_player_panel:UserComponent

@export var allowedPlayArea:PlayArea

var all_drawers :Array[CardDrawer]
var table_draw_queue:Array[Callable]

var last_card_played:Card = null
var others_card_played:Array[Card]




signal GameDataUpdated(game_data:GameData)

var isMyTurn:bool = false

var game_data:GameData 

func _ready() -> void:
	game_data = GameData.new()
	me_drawer.MyCardPlayed.connect(_on_card_played)

	me_drawer.DrawerCardUp.connect(_on_card_up)
	me_drawer.DrawerCardDown.connect(_on_card_down)

	game_data.current_trick = TrickData.new()
	
	all_drawers.append(me_drawer)
	all_drawers.append(right_drawer)
	all_drawers.append(top_drawer)
	all_drawers.append(left_drawer)
	

	for dra in all_drawers:
		dra.AddToQueue.connect(push_callback)

	run_actions()

func _on_card_up()->void:
	allowedPlayArea.turn_off()

func _on_card_down()->void:
	allowedPlayArea.turn_on()

func clear_cards()->void:
	push_callback(
		_clear_cards.bind()
	)

func _clear_cards()->void:
	for drawer_queue in all_drawers:
		for c:Card in drawer_queue.cards:
			if is_instance_valid(c):
				remove_one_card(c)
		drawer_queue.cards = []
		

func remove_one_card(card:Card) -> void:
	push_callback(
		# card.disconnect
		card.queue_free.call_deferred.bind()
	)

func push_callback(c:Callable)->void:
	table_draw_queue.push_back(c)

func _on_card_played(card:Card)->void:
	MyCardPlayed.emit(card)


func set_player_to_hand()->void:
	var meId := KAccount._instance.MyAccount.id

	var mePlayer : PlayerData = null
	var beforeMe : Array[PlayerData]
	var afterMe : Array[PlayerData]

	for player in game_data.players:
		if player.user_id == meId:
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

	me_drawer.unique_string = mePlayer.player_unique
	right_drawer.unique_string = player_order_from_me[0].player_unique
	top_drawer.unique_string = player_order_from_me[1].player_unique
	left_drawer.unique_string = player_order_from_me[2].player_unique

	me_player_panel.user_id = mePlayer.user_id
	right_player_panel.user_id = player_order_from_me[0].user_id
	top_player_panel.user_id = player_order_from_me[1].user_id
	left_player_panel.user_id = player_order_from_me[2].user_id

func run_actions()->void:
	while true:
		var action_variant:Variant = table_draw_queue.pop_front()
		if action_variant != null:			
			var action:Callable = action_variant
			await action.call()
		else:
			await get_tree().create_timer(.2).timeout
		

func new_cards_event(e:KEvent.Event,no_animation:bool=false)->void:
	push_callback(_new_cards_event.bind(e,no_animation))

func _new_cards_event(e:KEvent.Event,no_animation:bool=false)->void:
	me_drawer.new_cards_event(e,no_animation)
	me_drawer.break_action()

	right_drawer.new_cards_event(e,no_animation)
	right_drawer.break_action()

	top_drawer.new_cards_event(e,no_animation)
	top_drawer.break_action()

	left_drawer.new_cards_event(e,no_animation)
	left_drawer.break_action()

func parse_game_data(json_string:String)->void:
	game_data = JsonClassConverter.json_string_to_class(GameData,json_string)
	GameDataUpdated.emit(game_data)

func remove_played_cards()->void:
	var cards_trash : Array[Card]
	cards_trash.append(last_card_played)
	cards_trash.append_array(others_card_played)
	last_card_played = null
	others_card_played = []

	push_callback(_remove_played_cards.bind(cards_trash))

func _remove_played_cards(cards_trash:Array[Card])->void:
	for c in cards_trash:
		if is_instance_valid(c) and c:
			c.queue_free()
	for dr in all_drawers:
		dr.draw_cards_and_sort()
