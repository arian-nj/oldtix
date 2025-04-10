### DONT TOUCH THIS IF YOU DON'T KNOW HOW IT'S WORKING i don't have no clue
extends State

@export var ws:KatanaSocket
@export var status_label:Label
@export var table:Game4Table
@export var chooseHokmPanel:ChooseHokmPanel
var is_hakem := false

var got_cards_time:int

func Enter()->void:
	got_cards_time = 0
	status_label.text = "Choose Hokm"
	ws.new_event.connect(_on_new_event)
	ws.open_events()
	
	var hakem_player := table.game_data.players[table.game_data.current_trick.hakem_index]
	# print("my account ",KAccount.MyAccount.id," hakem account id ",hakem_player.UserId)
	if KAccount._instance.MyAccount.id == hakem_player.user_id:
		# HokmChoosed.connect(_on_hokm_choosed)
		chooseHokmPanel.HokmChoosed.connect(_on_hokm_choosed)
		chooseHokmPanel.reset_all_suites()
		is_hakem = true

func Exit()->void:
	ws.new_event.disconnect(_on_new_event)
	if chooseHokmPanel.HokmChoosed.is_connected(_on_hokm_choosed):
		chooseHokmPanel.HokmChoosed.disconnect(_on_hokm_choosed)

func _on_new_event(e:KEvent.Event)->void:
	if e.type == KEvent.TYPE_NEW_CARD:
		table.new_cards_event(e)		
		got_cards_time += 1
		if got_cards_time == 3:
			ws.hold_events()
			Transition.emit(self,"game_turn")

	elif e.type == KEvent.TYPE_NEW_HOKM:
		table.parse_game_data(e.data)
		status_label.text = "Hokm Choosed"
		if is_hakem:
			var btn :SuiteButton = chooseHokmPanel.find_btn_from_suite(table.game_data.current_trick.hokm)
			btn.pressed.emit()
		else:
			chooseHokmPanel.come_up(table.game_data.current_trick.hokm)

	# else:
		# print_debug(e.type)



# only runs if hakem panel is shown and hakem choosed hokm
func _on_hokm_choosed(Hokm:CardData.CardSuites)->void:
	ws.send_event(KEvent.TYPE_PLAYER_SELECTED_HOKM,str(Hokm))
