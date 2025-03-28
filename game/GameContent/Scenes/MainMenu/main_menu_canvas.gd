extends SceneLevel

@export var aspect_ration_container:CustomRatioAspectContainer

func _ready() -> void:
	resized.connect(aspect_ration_container._on_aspect)

func _on_play_tic_2_button_pressed() -> void:
	manager_change_scene.emit(SceneManager.Levels.GameHokm4)


func _on_continue_button_pressed() -> void:
	manager_change_scene.emit(SceneManager.Levels.GameHokm4)
