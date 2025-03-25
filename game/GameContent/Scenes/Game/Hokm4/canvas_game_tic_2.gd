extends SceneLevel

@export var ws:KatanaSocket

func _on_home_button_pressed() -> void:
	manager_change_scene.emit(SceneManager.Levels.MainMenu)
	ws.send_event(KEvent.TYPE_DISCONNECT)
	
