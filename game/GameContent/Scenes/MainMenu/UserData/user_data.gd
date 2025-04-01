extends Control

@export var usernameLabel:Label
@export var displaynameLabel:Label

@export var profileButton:Button
@export var profilePanelContainer:PanelContainer


func _ready() -> void:
	profilePanelContainer.visible = false
	profileButton.pressed.connect(_on_profile_pressed)
	resized.connect(_on_resize)
	KAccount._instance.MeChanged.connect(set_user_display_name_label)
	set_user_display_name_label()


func set_user_display_name_label()->void:
	displaynameLabel.text = KAccount._instance.MyAccount.display_name
	usernameLabel.text = "@"+KAccount._instance.MyAccount.username

func _on_resize()->void:
	if profilePanelContainer.visible:
		_on_profile_pressed()

func _on_profile_pressed()->void:
	panel_go_up()

func _on_button_pressed() -> void:
	panel_go_down()

func panel_go_up()->void:
	profileButton.disabled = true
	profilePanelContainer.visible = true

	profilePanelContainer.size = self.size
	profilePanelContainer.global_position = self.global_position
	profilePanelContainer.global_position.y = self.size.y
	
	var tween := create_tween().set_ease(Tween.EASE_IN).set_trans(Tween.TRANS_CUBIC)
	tween.parallel().tween_property(profilePanelContainer,"global_position",self.global_position,.8)
	await tween.finished

func panel_go_down()->void:
	var target_position:Vector2 = self.global_position
	target_position.y += self.size.y
	
	var tween := create_tween().set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_BOUNCE)
	tween.parallel().tween_property(profilePanelContainer,"global_position",target_position,1)
	await tween.finished
	profileButton.disabled = false
	profilePanelContainer.visible = false


# func panel_go_up()->void:
# 	# profileButton.disabled = true
# 	betPanelContainer.visible = true
# 	betPanelContainer.size = self.size
# 	betPanelContainer.global_position = self.global_position
# 	betPanelContainer.global_position.y = self.size.y
	
# 	var tween := create_tween().set_ease(Tween.EASE_IN).set_trans(Tween.TRANS_CUBIC)
# 	tween.parallel().tween_property(betPanelContainer,"global_position",self.global_position,.8)
# 	await tween.finished

# func panel_go_down()->void:
# 	var tween := create_tween().set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_BOUNCE)
# 	var target_position:Vector2 = self.global_position
# 	target_position.y += self.size.y

# 	tween.parallel().tween_property(betPanelContainer,"global_position",target_position,1)
# 	await tween.finished
# 	# profileButton.disabled = false
# 	betPanelContainer.visible = false


# func _on_play_tic_2_button_pressed() -> void:
# 	panel_go_up()


# func _on_button_pressed() -> void:
# 	panel_go_down()
