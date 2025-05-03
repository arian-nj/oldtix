class_name ChooseHokmPanel extends Control

signal HokmChoosed(Hokm:CardData.CardSuites)
# @export var sutiesSprite:AtlasTexture

@export var spadeSuite:SuiteButton
@export var heartSuite:SuiteButton
@export var clubSuite:SuiteButton
@export var diamondSuite:SuiteButton

@export var final_hokm_control:Control

@export var animation_time_span:float = 0.5

var last_hokm_suite:SuiteButton

var all_suits_btns:Array[SuiteButton]

func _ready() -> void:
	self.visible = true
	all_suits_btns.append_array([spadeSuite,heartSuite,clubSuite,diamondSuite])

	for btn:SuiteButton in all_suits_btns:
		btn.set_disolve(0)

	spadeSuite.SuitePressed.connect(_on_spades_button_pressed)
	heartSuite.SuitePressed.connect(_on_heart_button_pressed)
	clubSuite.SuitePressed.connect(_on_club_button_pressed)
	diamondSuite.SuitePressed.connect(_on_dimond_button_pressed)

func _on_spades_button_pressed(sb:SuiteButton) -> void:
	print("spade")
	HokmChoosed.emit(CardData.CardSuites.Spade)
	make_others_burn(sb)
	

func _on_heart_button_pressed(sb:SuiteButton) -> void:
	print("heart")
	HokmChoosed.emit(CardData.CardSuites.Heart)
	make_others_burn(sb)


func _on_club_button_pressed(sb:SuiteButton) -> void:
	print("club")
	HokmChoosed.emit(CardData.CardSuites.Club)
	make_others_burn(sb)


func _on_dimond_button_pressed(sb:SuiteButton) -> void:
	print("diamond")
	HokmChoosed.emit(CardData.CardSuites.Diamond)
	make_others_burn(sb)

func make_others_burn(sb:SuiteButton)->void:
	for btn:SuiteButton in all_suits_btns:
		if btn == sb:
			btn.move_hokm_to_position(final_hokm_control)
			last_hokm_suite = sb
			btn.selected = true
			continue
		btn.disolve_sprite()
		btn.locked = true
	
func reset_all_suites()->Signal:
	for btn:SuiteButton in all_suits_btns:
		btn.reset(animation_time_span)
		btn.selected = false
		btn.locked = false

	last_hokm_suite = null
	return get_tree().create_timer(animation_time_span).timeout

func find_btn_from_suite(suite: CardData.CardSuites)->SuiteButton:
	var btn :SuiteButton
	if suite == CardData.CardSuites.Club:
		btn = clubSuite
	elif suite == CardData.CardSuites.Spade:
		btn = spadeSuite
	elif suite == CardData.CardSuites.Diamond:
		btn = diamondSuite
	elif suite == CardData.CardSuites.Heart:
		btn = heartSuite
	return btn

func _on_show_choosed_hokm_button_pressed() -> void:
	await reset_all_suites()
	var btn :SuiteButton = find_btn_from_suite(CardData.CardSuites.Spade)
	btn.pressed.emit()


func _on_show_choosed_hokm_button_2_pressed() -> void:
	await reset_all_suites()
	var btn :SuiteButton = find_btn_from_suite(CardData.CardSuites.Club)
	btn.pressed.emit()

func _on_reset_button_pressed() -> void:
	reset_all_suites()
