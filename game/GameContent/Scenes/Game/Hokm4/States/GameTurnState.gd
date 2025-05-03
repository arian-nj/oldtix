extends State

@export var ws:KatanaSocket
@export var table:Game4Table
@export var status_label:Label

func Enter()->void:
	ws.NewEventSig.connect(_on_new_event)
	ws.open_events()
	table.me_drawer.MyCardPlayed.connect(on_me_card_played)

func Exit()->void:
	ws.hold_events()
	ws.NewEventSig.disconnect(_on_new_event)
	table.me_drawer.MyCardPlayed.disconnect(on_me_card_played)

var lock:bool = false

func _on_new_event(e:KEvent.Event)->void:
	if e.type == KEvent.TYPE_TURN_START:
		lock = false
		status_label.text = "Turn Started"
		table.parse_game_data(e.data)
		var last_card := table.last_card_played
		if last_card != null and is_instance_valid(last_card):
			table.me_drawer.append_card(last_card)
			table.me_drawer.draw_cards()
			table.last_card_played = null

	elif e.type == KEvent.TYPE_YOUR_TURN:
		status_label.text = "Your Turn"
		table.me_player_panel.start_shader(InternalSetting.PLAYER_PLAY_WAIT)
		table.isMyTurn = true
		if table.last_card_played != null:
			send_card_playend_event(table.last_card_played)

	elif e.type == KEvent.TYPE_VALID_PLAY:
		lock = true
		status_label.text = "Valid"
		table.me_player_panel.stop_shader()
		table.isMyTurn = false

	elif e.type == KEvent.TYPE_INVALID_PLAY:
		status_label.text = "Invalid"
		if table.last_card_played != null and is_instance_valid(table.last_card_played):
			table.me_drawer.append_card(table.last_card_played)
		table.last_card_played = null
		table.me_drawer.draw_cards_and_sort()
		

	elif e.type == KEvent.TYPE_PLAY_TIMEOUT:
		lock = true
		status_label.text = "Timeout"
		table.me_player_panel.stop_shader()
		table.isMyTurn = false
		# print("time out " + e.data)
		var new_card_data :CardData = JsonClassConverter.json_string_to_class(CardData,e.data)
		var played_card := table.me_drawer.play_me_card(new_card_data)
		if table.last_card_played != null and is_instance_valid(table.last_card_played):
			table.me_drawer.append_card(table.last_card_played)
			table.me_drawer.draw_cards()
		table.last_card_played = played_card

	
	elif e.type == KEvent.TYPE_TURN_PLAYED:
		status_label.text = "other played a turn"
		other_turn_played_logic(e.data)
	
	elif e.type == KEvent.TYPE_END_TURN:
		status_label.text = "Turn Ended"
		table.parse_game_data(e.data)
		table.remove_played_cards()
	
	elif e.type == KEvent.TYPE_END_TRICK:
		status_label.text = "Trick Ended"
		table.parse_game_data(e.data)
		table.clear_all_cards()
		ws.hold_events()
		Transition.emit(self,"wait_state")


func on_me_card_played(card:Card)->void:
	print("card played ",card.card_data.suit , " ", card.card_data.value)
	if lock:
		table.me_drawer.draw_cards()
		return
	
	table.me_drawer.cards.erase(card)
	if table.isMyTurn:
		send_card_playend_event(card)
	else:
		if table.last_card_played == null:
			table.last_card_played = card
		else:
			table.me_drawer.append_card(table.last_card_played)
			table.last_card_played = card
			table.me_drawer.draw_cards_and_sort()

	
func send_card_playend_event(card:Card)->void:
	table.last_card_played = card
	table.me_drawer.cards.erase(card)
	var card_played_string := JsonClassConverter.class_to_json_string(card.card_data)
	ws.send_event(KEvent.TYPE_TURN_PLAYED,card_played_string)

func other_turn_played_logic(data:String)->void:
	# print(KAccount._instance.MyAccount.username," ==> ",data)
	var card_played :CardPlayedData= JsonClassConverter.json_string_to_class(CardPlayedData,data)
	
	var played_card :Card = null
	for drawer in table.all_drawers:
		if drawer.unique_string == card_played.player.player_unique:
			played_card = drawer.play_random_card(card_played.card)
			break
	
	if played_card != null:
		table.others_card_played.append(played_card)
