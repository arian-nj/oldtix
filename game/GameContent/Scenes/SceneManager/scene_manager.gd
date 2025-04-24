extends Node
class_name SceneManager

@export var current_level:SceneLevel
@export var anim :AnimationPlayer

enum Levels{
	GameHokm4,
	MainMenu,
}

func _ready() -> void:
	KAccount.instanciate()
	self.add_child(KAccount._instance)
	ErrorBoard.instanciate()
	
	current_level.manager_change_scene.connect(handle_level_change)
	current_level.OnLoaded()


func handle_level_change(change_level_to:Levels)->void:
	# var start_time:float = Time.get_unix_time_from_system() 
	var next_level_name:String = ""
	match change_level_to:
		Levels.GameHokm4:
			next_level_name = "Game/Hokm4/OnlineHokm4"
		Levels.MainMenu:
			next_level_name = "MainMenu/MainMenu"
		_:
			print_debug("not match found")
			return

	var scene_address :String = "res://GameContent/Scenes/"+ next_level_name +".tscn"
	var next_level_packed:PackedScene = load(scene_address)
	if next_level_packed == null:
		print_debug("next_level is null can't find :\n"+scene_address)
		return
	var next_level:SceneLevel = next_level_packed.instantiate()
	current_level.visible = false
	anim.play("fade_in")
	await anim.animation_finished
	add_child(next_level)
	transfer_data_between_scenes(current_level,next_level)
	next_level.visible = true
	anim.play("fade_out")
	await anim.animation_finished
	
	print(scene_address)
	current_level.CleanUp()
	current_level = next_level
	current_level.manager_change_scene.connect(handle_level_change)
	current_level.OnLoaded()
	# var end_time:float = Time.get_unix_time_from_system() 
	# print(end_time-start_time - 1)

	

func transfer_data_between_scenes(old_scene:SceneLevel,new_scene:SceneLevel) -> void:
	new_scene.load_level_parameters(old_scene.level_parameters)

