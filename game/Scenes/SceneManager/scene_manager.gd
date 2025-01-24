class_name SceneManger extends Node

@export var current_level:SceneLevel
@export var anim :AnimationPlayer

enum Levels{
	GameTic2,
	MainMenu,
}

func _ready() -> void:
	current_level.manager_change_scene.connect(_handle_level_change)
	current_level.on_loaded()

func _handle_level_change(change_level_to:Levels) -> void:
	var next_level_name:String = ""
	
	match change_level_to:
		Levels.GameTic2:
			next_level_name = "GameTic2/GameTic2"
		Levels.MainMenu:
			next_level_name = "MainMenu/MainMenu"
		_:
			return
	
	var scene_address :String = "res://Scenes/"+ next_level_name +".tscn"
	var next_level_packed:PackedScene = load(scene_address)
	if next_level_packed == null:
		print_debug("next_level is null can't find :\n"+scene_address)
		return

	var next_level:SceneLevel = next_level_packed.instantiate()

	
	next_level.layer = -1
	anim.play("fade_in")
	await anim.animation_finished
	add_child(next_level)
	transfer_data_between_scenes(current_level,next_level)
	
	next_level.layer = 1
	anim.play("fade_out")
	await anim.animation_finished
	
	 
	current_level.cleanup()
	current_level = next_level
	current_level.manager_change_scene.connect(_handle_level_change)
	current_level.on_loaded()

	

func transfer_data_between_scenes(old_scene:SceneLevel,new_scene:SceneLevel) -> void:
	new_scene.load_level_parameters(old_scene.level_parameters)
