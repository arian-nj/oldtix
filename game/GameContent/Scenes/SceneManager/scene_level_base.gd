class_name SceneLevel extends Control

signal manager_change_scene(to_level:SceneManger.Levels)

var level_parameters:Dictionary = {
}

func load_level_parameters(new_level_parameters:Dictionary) -> void:
	level_parameters = new_level_parameters

func OnLoaded() -> void:
	manager_change_scene.is_null()
	pass
	
func CleanUp() -> void:
	queue_free()
