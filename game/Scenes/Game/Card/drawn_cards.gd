class_name CardDrawer extends Control

@export var from_middle: Control
@export var hand: Control

@export var rot_max: float = 10.0
@export var card_offset_x: float = 30.0
@export var card_scene: PackedScene
@export var card_movment_dur:float = .3
@export var two_card_movment_dur:float = .075
@export var flip_card_dur:float = .075

var cards:Array[Card] = []
var tween:Tween
var drawn:bool

var isDrawn:bool = true

func _ready() -> void:
	get_window().size_changed.connect(draw_cards)

func new_cards_event(e:KEvent.Event)->void:
	var json_data:Variant = JSON.parse_string(e.data)
	var cards_json:Variant = json_data["cards"]
	for card_json:Variant in cards_json:
		create_card(card_json["suit"],card_json["value"])
	if tween and tween.is_running():
		await tween.finished
	while isDrawn == false:
		await get_tree().create_timer(.1).timeout
	draw_cards(from_middle.global_position)
	

func create_card(suit:Card.CardSuites,value:int)->void:
	var c:Card = card_scene.instantiate()
	c.suit = suit
	c.value = value
	add_child(c)
	cards.append(c)
	c.card_played.connect(_card_played)
	# c.card_unplayed.connect(_card_unplayed)
	c.button_up.connect(func()->void:
		draw_cards()
	)

func _card_played(card:Card)->void:
	if isDrawn:
		print(card.suite_name()+"--"+card.value_name())
		cards.erase(card)
		draw_cards.call_deferred()

# func _card_unplayed(card:Card)->void:
# 	if card.in_hand:
# 		print("undrawn ",isDrawn)
# 		cards.append(card)
# 		draw_cards.call_deferred()


func draw_cards(from_pos: Vector2 = Vector2.ZERO) -> void:
	if len(cards) <= 0:
		return
	
	isDrawn = false

	sort_cards()
	if tween and tween.is_running():
		tween.kill()
	
	tween = create_tween().set_ease(Tween.EASE_IN_OUT).set_trans(Tween.TRANS_CUBIC)
	var deck_x_length:float = card_offset_x * (len(cards)) + cards[0].size.x
	var x_offset:float = deck_x_length/2 
	
	var not_sorted_counter:int = 0
	for i:int in len(cards):
		var instance: Card = cards[i]
		# if instance.button_pressed:
		# 	continue

		# in tree order to handle input hiarchy correctly + Render Order
		remove_child(instance)
		add_child(instance)

		if from_pos == from_middle.global_position and !instance.in_hand:
			instance.global_position = from_pos
			instance.global_position -= instance.size
		

		var final_pos: Vector2 = Vector2(card_offset_x * i , 0.0)
		final_pos.y -= instance.size.y/2 # center y
		final_pos.x -= x_offset 
		final_pos += hand.global_position
			
		var dur := card_movment_dur
		var delay :float = 0
		if !instance.in_hand:
			delay= (not_sorted_counter * two_card_movment_dur)
			not_sorted_counter += 1
			if from_pos == from_middle.global_position:
				tween.finished.connect(func ()->void:
					instance.prespective3DShader.flip_y(flip_card_dur,dur+delay,instance.load_assets)
				)
			instance.in_hand = true
			tween.parallel().tween_property(instance, "global_position", final_pos, dur).set_delay(delay)

		elif final_pos != instance.global_position:
			if final_pos.distance_to(instance.global_position) < 2:
				instance.global_position = final_pos
			else:
				dur += (i * two_card_movment_dur/4)
				tween.parallel().tween_property(instance, "global_position", final_pos, dur).set_delay(delay)


	
	await tween.finished
	isDrawn = true


var in_deck_suites:Array[Card.CardSuites] = []

func sort_cards()->void:
	in_deck_suites = []

	for card:Card in cards:
		if in_deck_suites.has(card.suit) == false:
			in_deck_suites.append(card.suit)
	
	if len(cards) <= 2:
		return
	
	in_deck_suites = sort_deck_suits(in_deck_suites)

	# for card in cards
	cards.sort_custom(value_sort)
	cards.sort_custom(suite_sort)


# sort filters
func suite_sort(a:Card,b:Card)->bool:
	var a_suite_index:int = in_deck_suites.find(a.suit)
	var b_suite_index :int = in_deck_suites.find(b.suit)

	if a_suite_index < b_suite_index:
		return true
	return false

func value_sort(a:Card,b:Card)->bool:
	if a.value < b.value:
		return true
	return false

func sort_deck_suits(suits: Array[Card.CardSuites]) -> Array[Card.CardSuites]:
	var red_suits: Array[Card.CardSuites] = []
	var black_suits: Array[Card.CardSuites] = []
	
	# Separate red and black suits
	for suit:Card.CardSuites in suits:
		if suit == Card.CardSuites.Diamond or suit == Card.CardSuites.Heart:
			red_suits.append(suit)
		else:
			black_suits.append(suit)

	# Create a new array to hold the sorted suits
	var sorted_array: Array[Card.CardSuites] = []
	
	# Determine the maximum length for interleaving
	var max_length:int = max(red_suits.size(), black_suits.size())
	
	# Interleave red and black suits
	for i:int in range(max_length):
		if black_suits.size() > red_suits.size():
			if i < black_suits.size():
				sorted_array.append(black_suits[i])
			if i < red_suits.size():
				sorted_array.append(red_suits[i])
		else:
			if i < red_suits.size():
				sorted_array.append(red_suits[i])
			if i < black_suits.size():
				sorted_array.append(black_suits[i])
		
	return sorted_array




func _on_re_draw_button_pressed() -> void:
	draw_cards()
