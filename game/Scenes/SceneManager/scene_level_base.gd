class_name SceneLevel extends CanvasLayer

signal manager_change_scene(to_level:SceneManger.Levels)

var level_parameters:Dictionary = {
}

func load_level_parameters(new_level_parameters:Dictionary) -> void:
	level_parameters = new_level_parameters

func on_loaded() -> void:
	pass
	
func cleanup() -> void:
	queue_free()
