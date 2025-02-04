class_name ChooseHokmPanel extends PanelContainer

signal HokmChoosed(Hokm:Card.CardSuites)


func _on_heart_button_pressed() -> void:
	HokmChoosed.emit(Card.CardSuites.Heart)


func _on_spades_button_pressed() -> void:
	HokmChoosed.emit(Card.CardSuites.Spade)


func _on_club_button_pressed() -> void:
	HokmChoosed.emit(Card.CardSuites.Club)


func _on_dimond_button_pressed() -> void:
	HokmChoosed.emit(Card.CardSuites.Diamond)
