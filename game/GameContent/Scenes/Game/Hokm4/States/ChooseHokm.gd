### DONT TOUCH THIS IF YOU DON'T KNOW HOW IT'S WORKING i don't have no clue
extends State

@export var ws:KatanaSocket
@export var status_label:Label
@export var table:Game4Table
@export var choose_hokm_scene:PackedScene


var choose_hokm_instance:ChooseHokmPanel = null
var got_cards_time:int

func Enter()->void:
	got_cards_time = 0
	status_label.text = "Choose Hokm"
	ws.new_event.connect(_on_new_event)
	ws.open_events()
	
	var hakem_player := table.game_data.players[table.game_data.current_trick.hakem_index]
	# print("my account ",KAccount.MyAccount.id," hakem account id ",hakem_player.UserId)
	if KAccount._instance.MyAccount.id == hakem_player.user_id:
		choose_hokm_instance = choose_hokm_scene.instantiate()
		table.add_child(choose_hokm_instance)
		choose_hokm_instance.HokmChoosed.connect(_on_hokm_choosed)

func Exit()->void:
	ws.new_event.disconnect(_on_new_event)

func _on_new_event(e:KEvent.Event)->void:
	if e.type == KEvent.TYPE_NEW_CARD:
		table.new_cards_event(e)
		got_cards_time += 1
		if got_cards_time == 3:
			Transition.emit(self,"game_turn")


	elif e.type == KEvent.TYPE_GAME_DATA:
		table.parse_game_data(e.data)
		status_label.text = "Hokm Choosed"
		if choose_hokm_instance != null:
			choose_hokm_instance.queue_free()
	else:
		print_debug(e.type)



# only runs if hakem panel is shown and hakem choosed hokm
func _on_hokm_choosed(Hokm:CardData.CardSuites)->void:
	ws.send_event(KEvent.TYPE_HOKM_CHOOSED,str(Hokm))
