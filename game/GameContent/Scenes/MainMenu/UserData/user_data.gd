extends Control

@export var usernameLabel:Label
@export var displaynameLabel:Label
@export var profileButton:Button
@export var panel_container:PanelContainer


func _ready() -> void:
	profileButton.pressed.connect(_on_profile_pressed)
	resized.connect(_on_resize)
	KAccount._instance.MeChanged.connect(set_user_display_name_label)
	set_user_display_name_label()


func set_user_display_name_label()->void:
	displaynameLabel.text = KAccount._instance.MyAccount.display_name
	usernameLabel.text = "@"+KAccount._instance.MyAccount.username

func _on_resize()->void:
	if panel_container.visible:
		_on_profile_pressed()
func _on_profile_pressed()->void:
	profileButton.disabled = true
	panel_container.visible = true
	panel_container.size = self.size
	panel_container.global_position = self.global_position
	panel_container.global_position.y = self.size.y
	var tween := create_tween().set_ease(Tween.EASE_IN).set_trans(Tween.TRANS_CUBIC)
	tween.parallel().tween_property(panel_container,"global_position",self.global_position,.8)
	await tween.finished

func _on_button_pressed() -> void:
	var tween := create_tween().set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_BOUNCE)
	var target_position:Vector2 = self.global_position
	target_position.y += self.size.y

	tween.parallel().tween_property(panel_container,"global_position",target_position,1)
	await tween.finished
	profileButton.disabled = false
	panel_container.visible = false
