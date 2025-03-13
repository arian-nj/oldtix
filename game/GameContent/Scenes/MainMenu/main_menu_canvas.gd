extends SceneLevel


func _on_play_tic_2_button_pressed() -> void:
	manager_change_scene.emit(SceneManager.Levels.GameHokm4)


func _on_continue_button_pressed() -> void:
	manager_change_scene.emit(SceneManager.Levels.GameHokm4)