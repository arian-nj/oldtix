extends SceneLevel


func _on_play_tic_2_button_pressed() -> void:
	manager_change_scene.emit(SceneManger.Levels.GameHokm4)
