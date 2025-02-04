extends SceneLevel

func _on_home_button_pressed() -> void:
	manager_change_scene.emit(SceneManger.Levels.MainMenu)
