extends SceneLevel

func _on_main_button_pressed() -> void:
	manager_change_scene.emit(SceneManger.Levels.MainMenu)