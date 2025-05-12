extends SceneLevel

@export var ws:KatanaSocket

func _ready() -> void:
	# ws.Disconnected.connect(_on_disconnect)
	pass

func _on_home_button_pressed() -> void:
	manager_change_scene.emit(SceneManager.Levels.MainMenu)
	ws.send_event(KEvent.TYPE_DISCONNECT)
	
func OnLoaded()->void:
	connect_ws()

func _on_disconnect()-> void:
	await get_tree().create_timer(.5).timeout
	connect_ws()

func connect_ws()-> void:
	var bm :Variant = level_parameters.get(SharedBetAmount)
	if bm == null:
		print_debug("fialed to get" + SharedBetAmount)
		return
	var bm_int:int = int(bm)
	ws.connect_to_game(bm_int,30)
