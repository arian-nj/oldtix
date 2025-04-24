extends Control
class_name SceneLevel

const SharedBetAmount = "bet_amount"

var level_parameters:Dictionary = {}

func load_level_parameters(new_level_parameters:Dictionary) -> void:
	level_parameters = new_level_parameters

func OnLoaded() -> void:
	pass
	
func CleanUp() -> void:
	queue_free()

signal manager_change_scene(to_level:SceneManager.Levels)

func do_nothing_not_call()->void: # only reseon to exist resolving editor warning
	manager_change_scene.is_null()
