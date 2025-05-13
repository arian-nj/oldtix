### DONT TOUCH THIS IF YOU DON'T KNOW HOW IT'S WORKING i don't have no clue
## I have been touching myself. know i can see.
extends State

@export var ksocket:KatanaSocket
@export var status_label:Label
@export var table:Game4Table
@export var chooseHokmPanel:ChooseHokmPanel
var is_hakem := false

var got_cards_time:int

func Enter()->void:
	is_hakem = false
	got_cards_time = 0
	status_label.text = "Choose Hokm"

	var hakem_player := table.game_data.players[table.game_data.current_trick.hakem_index]
	if KClient._instance.MyAccount.id == hakem_player.user_id and table.rejoining == false:
		# HokmChoosed.connect(_on_hokm_choosed)
		chooseHokmPanel.HokmChoosed.connect(_on_hokm_choosed)
		chooseHokmPanel.reset_all_suites()
		is_hakem = true
		table.me_player_panel.start_shader(InternalSetting.SETTING_PLAYER_CHOOSE_HOKM_WAIT)

func Exit()->void:
	pass

func _process(_delta: float) -> void:
	var new_event := ksocket.get_latest_event()
	if new_event == null:
		return
	
	if new_event.type == KEvent.TYPE_NEW_CARD:
		table.new_cards_event(new_event)		
		got_cards_time += 1
		if got_cards_time == 3:
			StateTransition.emit(self,"game_turn")

	elif new_event.type == KEvent.TYPE_NEW_CARD_ONE:
		table.new_cards_event(new_event,true)		
		StateTransition.emit(self,"game_turn")
	
	elif new_event.type == KEvent.TYPE_NEW_HOKM:
		table.parse_game_data(new_event.data)
		table.me_player_panel.stop_timer_shader()
		
		status_label.text = "Hokm Choosed"

		if is_hakem:
			var btn :SuiteButton = chooseHokmPanel.find_btn_from_suite(table.game_data.current_trick.hokm)
			btn.pressed.emit()
		else:
			await chooseHokmPanel.reset_all_suites()
			var btn :SuiteButton = chooseHokmPanel.find_btn_from_suite(table.game_data.current_trick.hokm)
			btn.pressed.emit()

	else:
		ksocket.push_event(new_event)

# only runs if hakem panel is shown and hakem choosed hokm
func _on_hokm_choosed(Hokm:CardData.CardSuites)->void:
	ksocket.send_event(KEvent.TYPE_PLAYER_SELECTED_HOKM,str(Hokm))
