class_name ChooseHokmPanel extends PanelContainer

signal HokmChoosed(Hokm:CardData.CardSuites)


func _on_heart_button_pressed() -> void:
	HokmChoosed.emit(CardData.CardSuites.Heart)


func _on_spades_button_pressed() -> void:
	HokmChoosed.emit(CardData.CardSuites.Spade)


func _on_club_button_pressed() -> void:
	HokmChoosed.emit(CardData.CardSuites.Club)


func _on_dimond_button_pressed() -> void:
	HokmChoosed.emit(CardData.CardSuites.Diamond)
