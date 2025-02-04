class_name CardDrawer extends Control

@export var from: Control
@export var hand: Control

@export var rot_max: float = 10.0
@export var card_offset_x: float = 30.0
@export var card_scene: PackedScene
@export var card_movment_dur:float = .3
@export var two_card_movment_dur:float = .075

var cards:Array[Card] = []
var tween:Tween
var drawn:bool

var isDrawn:bool = false

func create_card(suit:Card.CardSuites,value:int)->void:
	var c:Card = card_scene.instantiate()
	c.suit = suit
	c.value = value
	add_child(c)
	cards.append(c)
	c.card_played.connect(_card_played)
	c.button_up.connect(func()->void:
		draw_cards()
		)

func _card_played(card:Card)->void:
	if isDrawn:
		print(card.suite_name()+"--"+card.value_name())
		card.queue_free()
		cards.erase(card)
		draw_cards.call_deferred()


# func _ready() -> void:
# 	for i:int in range(TotalCards):
# 		create_card(randi_range(0,3),randi_range(1,13))

	# draw_cards(from.global_position)
func draw_cards(from_pos: Vector2 = Vector2.ZERO) -> void:
	if len(cards) <= 0:
		return
	isDrawn = false
	sort_cards()
	
	if tween and tween.is_running():
		tween.kill()
	
	tween = create_tween().set_ease(Tween.EASE_IN_OUT).set_trans(Tween.TRANS_CUBIC)
	var deck_x_length:float = card_offset_x * (len(cards)-1) + cards[0].size.x
	var x_offset:float = deck_x_length/2 
	
	for i:int in len(cards):
		var instance: Card = cards[i]
		remove_child(instance)
		add_child(instance)
		
		instance.z_index = i

		if from_pos != Vector2.ZERO and !instance.in_hand:
			instance.global_position = from_pos
			instance.global_position -= instance.size
			instance.in_hand = true
			# instance.swingComponent.last_position = instance.global_position

		
		var final_pos: Vector2 = Vector2(card_offset_x * i , 0.0)
		final_pos.y -= instance.size.y/2
		final_pos.x -= x_offset
		final_pos += hand.global_position


		tween.parallel().tween_property(instance, "global_position", final_pos, card_movment_dur + (i * two_card_movment_dur))
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
