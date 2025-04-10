class_name SuiteButton extends Button

signal SuitePressed(sb:SuiteButton)
var locked:bool = false

var initial_position:Vector2
var initial_scale:Vector2


func _ready() -> void:
	initial_position = self.position
	initial_scale = scale
	self.pressed.connect(_on_pressed)

func _on_pressed()->void:
	if locked:
		return
	SuitePressed.emit(self)
	locked = true

func move_hokm_to_position(final_hokm_control:Control) ->void:
	var tween := create_tween()
	tween.parallel().tween_property(self,"global_position",final_hokm_control.global_position,1.5)
	tween.parallel().tween_property(self,"scale",Vector2(2.0,2.0),2.0)

func disolve_sprite()->Signal:
	var tween:Tween = create_tween()
	var mat:ShaderMaterial= material
	tween.tween_property(mat,"shader_parameter/dissolve_value",0.0,2)
	return tween.finished
func set_disolve(disolve_value:float)->void:
	var mat:ShaderMaterial= material
	mat.set_shader_parameter("dissolve_value",disolve_value)

func redo_scale()->void:
	var tween := create_tween()
	tween.parallel().tween_property(self,"scale",initial_scale,0.5)

func redo_global_position()->void:
	var tween := create_tween()
	tween.parallel().tween_property(self,"position",initial_position,0.5)

func _scale_up()->Signal:
	var tween:Tween = create_tween()
	self.scale = Vector2(0.0,0.0)
	tween.parallel().tween_property(self,"scale",initial_scale,0.5)
	return tween.finished

func go_to_place() -> void:
	set_disolve(1.0)
	redo_global_position()
	if locked:
		redo_scale()
	else:
		_scale_up()
	locked = false

func go_up(final_hokm_control:Control)->void:
	position = Vector2(0,0)
	self.set_disolve(1)
	await _scale_up()
	await get_tree().create_timer(.5).timeout
	move_hokm_to_position(final_hokm_control)
