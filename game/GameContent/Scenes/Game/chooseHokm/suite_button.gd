class_name SuiteButton extends Button

signal SuitePressed(sb:SuiteButton)

var selected:bool = false # to use when reseting
var locked:bool = false # locking from clicking

var initial_position:Vector2 # start of scnene
var initial_scale:Vector2

func _ready() -> void:
	initial_position = self.position
	initial_scale = scale
	self.pressed.connect(_on_pressed)

func _on_pressed()->void:
	if locked:
		return
	self.SuitePressed.emit(self)

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

func redo_scale(time_span:float)->void:
	var tween := create_tween()
	tween.parallel().tween_property(self,"scale",initial_scale,time_span)
	return 

func redo_global_position(time_span:float)->void:
	var tween := create_tween()
	tween.parallel().tween_property(self,"position",initial_position,time_span)
	return

func _scale_up(time_span:float)->Signal:
	var tween:Tween = create_tween()
	self.scale = Vector2(0.0,0.0)
	tween.parallel().tween_property(self,"scale",initial_scale,time_span)
	return tween.finished

func reset(time_span:float) -> void:
	self.set_disolve(1.0)
	self.redo_global_position(time_span)
	if selected:
		self.redo_scale(time_span)
	else:
		self._scale_up(time_span)


# func go_up(final_hokm_control:Control,time_span:float)->void:
# 	position = Vector2(0,0)
# 	self.set_disolve(1)
# 	await _scale_up(time_span)
# 	await get_tree().create_timer(time_span).timeout
# 	move_hokm_to_position(final_hokm_control)
