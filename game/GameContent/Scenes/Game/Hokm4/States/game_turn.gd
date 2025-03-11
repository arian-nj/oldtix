extends State


@export var ws:KatanaSocket
@export var table:Game4Table
@export var status_label:Label

var last_card_played:Card = null
var others_card_played:Array[Card]

func Enter()->void:
	ws.new_event.connect(_on_new_event)
	ws.open_events()
	table.MyCardPlayed.connect(_card_played)
	for drawer in table.all_drawers:
		drawer.OtherCardPlayed.connect(_other_played_card)

func _other_played_card(card:Card)->void:
	print("other card played")
	others_card_played.append(card)

func Exit()->void:
	ws.hold_events()
	ws.new_event.disconnect(_on_new_event)
	table.MyCardPlayed.disconnect(_card_played)
	for drawer in table.all_drawers:
		drawer.OtherCardPlayed.disconnect(_other_played_card)

func _on_new_event(e:KEvent.Event)->void:
	if e.type == KEvent.TYPE_TURN_START:
		status_label.text = "Turn Started"
		table.parse_game_data(e.data)
		others_card_played = []

	elif e.type == KEvent.TYPE_YOUR_TURN:
		status_label.text = "Your Turn"
		table.isMyTurn = true
		if last_card_played != null:
			send_card_playend_event(last_card_played)

	elif e.type == KEvent.TYPE_VALID_PLAY:
		status_label.text = "Valid"
		table.isMyTurn = false
		table.drawer.cards.erase(last_card_played)

	elif e.type == KEvent.TYPE_INVALID_PLAY:
		status_label.text = "Invalid"
		if last_card_played != null:
			table.drawer.cards.append(last_card_played)
			#table.drawer.draw_cards()
		last_card_played = null
		table.drawer.draw_cards()

	elif e.type == KEvent.TYPE_PLAY_TIMEOUT:
		status_label.text = "Timeout"
		table.isMyTurn = false
	
	elif e.type == KEvent.TYPE_TURN_PLAYED:
		status_label.text = "other played a turn"
		other_turn_played_logic(e.data)
	
	elif e.type == KEvent.TYPE_END_TURN:
		status_label.text = "Turn Ended"
		table.parse_game_data(e.data)
		await get_tree().create_timer(1).timeout
		if last_card_played != null:
			table.remove_one_card(last_card_played)
		for card in others_card_played:
			table.remove_one_card(card)

	elif e.type == KEvent.TYPE_END_TRICK:
		status_label.text = "Trick Ended"
		table.parse_game_data(e.data)
		table.clear_cards()
		Transition.emit(self,"new_trick")


func _card_played(card:Card)->void:
	print("card played ",card.card_data.suit , " ", card.card_data.value)
	if table.isMyTurn:
		send_card_playend_event(card)
		
	last_card_played = card
	

func send_card_playend_event(card:Card)->void:
	var card_played_string := JsonClassConverter.class_to_json_string(card.card_data)
	ws.send_event(KEvent.TYPE_TURN_PLAYED,card_played_string)

func other_turn_played_logic(data:String)->void:
	# print(KAccount._instance.MyAccount.username," ==> ",data)
	var card_played :CardPlayedData= JsonClassConverter.json_string_to_class(CardPlayedData,data)
	# card_played.card
	for drawer in table.all_drawers:
		if drawer.unique_string == card_played.player.player_unique:
			drawer.play_others_card(card_played.card)
			break
