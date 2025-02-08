extends State

@export var ws:KatanaSocket
@export var status_label:Label
@export var game_table:Game4Table
@export var choose_hokm_scene:PackedScene
@export var card_drawer:CardDrawer


var choose_hokm_instance:ChooseHokmPanel = null
var got_cards_time:int

func Enter()->void:
	got_cards_time = 0
	status_label.text = "Choose Hokm"
	ws.new_event.connect(_on_new_event)
	
	var hakem_player := game_table.game_data.players[game_table.game_data.hakem]
	# print("my account ",KAccount.MyAccount.id," hakem account id ",hakem_player.UserId)
	if KAccount.MyAccount.id == hakem_player.UserId:
		choose_hokm_instance = choose_hokm_scene.instantiate()
		game_table.add_child(choose_hokm_instance)
		choose_hokm_instance.HokmChoosed.connect(_on_hokm_choosed)

func Exit()->void:
	ws.new_event.disconnect(_on_new_event)

func _on_new_event(e:KEvent.Event)->void:
	if e.type == KEvent.TYPE_NEW_CARD:
		card_drawer.new_cards_event(e)
		got_cards_time +=1
		if got_cards_time == 2:
			Transition.emit(self,"game_turn")

	elif e.type == KEvent.TYPE_GAME_DATA:
		game_table.parse_game_data(e.data)
		status_label.text = "Hokm Choosed"
		if choose_hokm_instance != null:
			choose_hokm_instance.queue_free()


# only runs if 
func _on_hokm_choosed(Hokm:Card.CardSuites)->void:
	ws.send_event(KEvent.TYPE_HOKM_CHOOSED,str(Hokm))