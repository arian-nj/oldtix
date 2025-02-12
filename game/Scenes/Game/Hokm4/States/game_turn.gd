extends State


@export var ws:KatanaSocket
@export var table:Game4Table
@export var status_label:Label

var last_card_played:Card = null

func Enter()->void:
	ws.new_event.connect(_on_new_event)
	table.CardPlayed.connect(_card_played)



func Exit()->void:
	ws.new_event.disconnect(_on_new_event)



func _on_new_event(e:KEvent.Event)->void:
	if e.type == KEvent.TYPE_TURN_START:
		status_label.text = "Turn Started"
		table.parse_game_data(e.data)

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
		last_card_played = null

		table.drawer.draw_cards()
	elif e.type == KEvent.TYPE_PLAY_TIMEOUT:
		status_label.text = "Timeout"
		table.isMyTurn = false
	
	elif e.type == KEvent.TYPE_GAME_DATA:
		status_label.text = "Turn Ended"
		if last_card_played != null:
			print("not null")
			last_card_played.queue_free()
		else:
			print("null")


func _card_played(card:Card)->void:
	print("card played ",card.card_data.suit , " ", card.card_data.value)
	if table.isMyTurn:
		send_card_playend_event(card)
		
	last_card_played = card
	

func send_card_playend_event(card:Card)->void:
	var card_played_string := JsonClassConverter.class_to_json_string(card.card_data)
	ws.send_event(KEvent.TYPE_PLAY_TURN,card_played_string)
