class_name LPError extends Control

@export var lpErrorPanelContainer:PanelContainer
@export var reTryButton:Button
@export var errorsLabel:Label

signal TryAgain

func _ready()->void:
	lpErrorPanelContainer.visible = false
	go_away()
	reTryButton.pressed.connect(go_away)
	reTryButton.pressed.connect(TryAgain.emit)
	self.resized.connect(_on_resize)

func _on_resize()->void:
	if lpErrorPanelContainer.visible == true:
		come_up()

func go_away()->void:	
	lpErrorPanelContainer.visible = false

func come_up() -> void:
	lpErrorPanelContainer.visible = true

func new_error(error_message:String)->void:
	errorsLabel.text = error_message
	if lpErrorPanelContainer.visible == false:
		come_up()
