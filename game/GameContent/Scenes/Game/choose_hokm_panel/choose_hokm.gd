class_name ChooseHokmPanel extends Control

signal HokmChoosed(Hokm:CardData.CardSuites)
# @export var sutiesSprite:AtlasTexture

@export var spadeSuite:SuiteButton
@export var heartSuite:SuiteButton
@export var clubSuite:SuiteButton
@export var diamondSuite:SuiteButton

@export var final_hokm_control:Control

var choosed_hokm_suite:SuiteButton

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
	make_it_burn(sb)
	

func _on_heart_button_pressed(sb:SuiteButton) -> void:
	print("heart")
	HokmChoosed.emit(CardData.CardSuites.Heart)
	make_it_burn(sb)


func _on_club_button_pressed(sb:SuiteButton) -> void:
	print("club")
	HokmChoosed.emit(CardData.CardSuites.Club)
	make_it_burn(sb)


func _on_dimond_button_pressed(sb:SuiteButton) -> void:
	print("diamond")
	HokmChoosed.emit(CardData.CardSuites.Diamond)
	make_it_burn(sb)

func make_it_burn(sb:SuiteButton)->void:
	for btn:SuiteButton in all_suits_btns:
		if btn == sb:
			btn.move_hokm_to_position(final_hokm_control)
			choosed_hokm_suite = sb
			continue
		btn.disolve_sprite()
		btn.locked = true
	
func reset_all_suites()->void:
	for btn:SuiteButton in all_suits_btns:
		btn.go_to_place()

	choosed_hokm_suite = null


func come_up(suite: CardData.CardSuites)->void:
	if choosed_hokm_suite != null: 
		await choosed_hokm_suite.disolve_sprite()
		choosed_hokm_suite.redo_global_position()

	var btn :SuiteButton =find_btn_from_suite(suite)
	choosed_hokm_suite = btn
	btn.go_up(final_hokm_control)
	for other_btn:SuiteButton in all_suits_btns:
		if other_btn != btn:
			other_btn.set_disolve(0.0)

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