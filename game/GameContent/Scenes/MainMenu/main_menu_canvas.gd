extends SceneLevel

@export var max_aspect:float
@export var aspect_ration_container:AspectRatioContainer

func _ready() -> void:
	_on_aspect()
	resized.connect(_on_aspect)

func _on_aspect()->void:
	var view_rect := get_viewport_rect()
	var aspect := view_rect.size.x/view_rect.size.y
	print(aspect)
	aspect = min(aspect,max_aspect)
	aspect_ration_container.ratio = aspect
	


func _on_play_tic_2_button_pressed() -> void:
	manager_change_scene.emit(SceneManager.Levels.GameHokm4)


func _on_continue_button_pressed() -> void:
	manager_change_scene.emit(SceneManager.Levels.GameHokm4)
