class_name CardDrawer extends Control

signal CardPlayed(card:Card)

func _card_played(card:Card)->void:
	if isDrawn:
		CardPlayed.emit(card)

@export var from_middle: Control
@export var hand: Control

@export var rot_max: float = 10.0
@export var card_offset_x: float = 30.0
@export var card_scene: PackedScene

@export var card_movment_dur:float = .3
@export var two_card_movment_dur:float = .075

@export var flip_card_dur:float = .075
@export var two_card_dur:float = .075


var cards:Array[Card] = []
var tween:Tween
var drawn:bool

var isDrawn:bool = true



@export var IsDrawnLabel: Label
func _process(_delta: float) -> void:
	IsDrawnLabel.text = str(isDrawn)

var draw_queue :Array[Callable] = []

@onready var timer := Timer.new()

func _ready() -> void:
	get_window().size_changed.connect(draw_cards)
	timer.wait_time = .1
	timer.timeout.connect(run_actions)
	add_child(timer)
	timer.start()


func run_actions()->void:
	# print("running action")
	timer.timeout.disconnect(run_actions)
	while draw_queue.size() > 0:
		var action:Callable = draw_queue.pop_front()
		# print("action => ",action)
		action.call()
	timer.timeout.connect(run_actions)

func clear_cards()->void:
	push_callback(
		_clear_cards.bind()
	)

func _clear_cards()->void:
	for c:Card in cards:
		c.queue_free()
	cards = []


func remove_one_card(card:Card) -> void:
	push_callback(
		card.queue_free.bind()
	)
	push_callback(
		draw_cards.bind()
	)

func push_callback(c:Callable)->void:
	draw_queue.push_back(c)
	

func new_cards_event(e:KEvent.Event)->void:
	push_callback(_new_cards_event.bind(e))


func _new_cards_event(e:KEvent.Event)->void:
	var json_data:Variant = JSON.parse_string(e.data)
	var cards_json:Variant = json_data["cards"]
	for card_json:Variant in cards_json:
		create_card(card_json["suit"],card_json["value"])
	draw_cards()
	print("cards size is : ",cards.size())


func create_card(suite:CardData.CardSuites,value:int)->void:
	var c:Card = card_scene.instantiate()
	c.card_data = CardData.new()
	c.card_data.suit = suite
	c.card_data.value = value

	add_child(c)
	cards.append(c)

	c.card_played.connect(_card_played)
	c.not_inplace.connect(draw_cards)


# func draw_cards(from_pos: Vector2 = Vector2.ZERO) -> void:
func draw_cards() -> void:
	push_callback(
		_draw_cards.bind()
	)

func _draw_cards(_from_pos:= Vector2.ZERO) -> void:
	if cards.is_empty():
		return

	isDrawn = false
	sort_cards()

	if tween and tween.is_running(): # Kill any running tween and create a new one with easing settings.
		tween.kill()
	tween = create_tween().set_ease(Tween.EASE_IN_OUT).set_trans(Tween.TRANS_CUBIC)

	# Calculate deck width and centering offset.
	var deck_x_length: float = card_offset_x * cards.size() + cards[0].size.x
	var x_offset: float = deck_x_length / 2.0

	var newCardsCounter: int = 0

	for i in range(cards.size()):
		var card: Card = cards[i]

		# Update tree order to correctly handle input and render order.
		remove_child(card)
		add_child(card)

		# Compute the final position for this card.
		var final_pos: Vector2 = _calculate_final_position(i, card, x_offset)

		# Set up movement parameters.
		var movementDuration: float = card_movment_dur
		var delay: float = 0.0

		# If the card is new to the hand, animate it in with a delay.
		if not card.in_hand:
			card.in_hand = true
			
			card.global_position = from_middle.global_position - card.size
			delay = newCardsCounter * two_card_movment_dur
			
			# Capture the current counter value for the tween callback.
			var currentCounter: int = newCardsCounter
			newCardsCounter += 1

			card.prespective3DShader.flip_y(flip_card_dur, currentCounter * two_card_dur, card.load_assets)
			tween.parallel().tween_property(card, "global_position", final_pos, movementDuration).set_delay(delay)
		
		# If the card is already in hand and its final position has changed…
		elif final_pos != card.global_position: 
			# Snap to final position if very close.
			if final_pos.distance_to(card.global_position) < 2.0:
				card.global_position = final_pos
			else:
				movementDuration += i * (two_card_movment_dur / 4.0)
				tween.parallel().tween_property(card, "global_position", final_pos, movementDuration).set_delay(delay)

	await tween.finished
	isDrawn = true


# Helper function to calculate a card's final position.
func _calculate_final_position(index: int, card: Card, x_offset: float) -> Vector2:
	var pos: Vector2 = Vector2(card_offset_x * index, 0.0)
	pos.y -= card.size.y / 2.0 # Center the card vertically.
	pos.x -= x_offset # Center horizontally.
	pos += hand.global_position # Offset by the hand's global position.
	return pos



var in_deck_suites:Array[CardData.CardSuites] = []

func sort_cards()->void:
	in_deck_suites = []

	for card:Card in cards:
		if not in_deck_suites.has(card.card_data.suit):
			in_deck_suites.append(card.card_data.suit)
	
	if len(cards) <= 2:
		return
	
	in_deck_suites = sort_deck_suits(in_deck_suites)

	# for card in cards
	cards.sort_custom(value_sort)
	cards.sort_custom(suite_sort)


# sort filters
func suite_sort(a:Card,b:Card)->bool:
	var a_suite_index:int = in_deck_suites.find(a.card_data.suit)
	var b_suite_index :int = in_deck_suites.find(b.card_data.suit)

	if a_suite_index < b_suite_index:
		return true
	return false

func value_sort(a:Card,b:Card)->bool:
	if a.card_data.value < b.card_data.value:
		return true
	return false

func sort_deck_suits(suits: Array[CardData.CardSuites]) -> Array[CardData.CardSuites]:
	var red_suits: Array[CardData.CardSuites] = []
	var black_suits: Array[CardData.CardSuites] = []
	
	# Separate red and black suits
	for suit:CardData.CardSuites in suits:
		if suit == CardData.CardSuites.Diamond or suit == CardData.CardSuites.Heart:
			red_suits.append(suit)
		else:
			black_suits.append(suit)

	# Create a new array to hold the sorted suits
	var sorted_array: Array[CardData.CardSuites] = []
	
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
