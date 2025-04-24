extends SceneLevel

@export var ws:KatanaSocket

func _ready() -> void:
	# ws.Disconnected.connect(_on_disconnect)
	pass

func _on_home_button_pressed() -> void:
	manager_change_scene.emit(SceneManager.Levels.MainMenu)
	ws.send_event(KEvent.TYPE_DISCONNECT)
	
func OnLoaded()->void:
	print("start onload")
	connect_ws()
	print("end onload")

func _on_disconnect()-> void:
	await get_tree().create_timer(.5).timeout
	connect_ws()

func connect_ws()-> void:
	var bm :int = level_parameters.get(SharedBetAmount)
	if bm == null:
		print_debug("fialed to get" + SharedBetAmount)
		return
	ws.connect_to_game(bm)
