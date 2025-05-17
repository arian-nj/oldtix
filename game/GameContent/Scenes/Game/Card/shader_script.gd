class_name CardShader extends Control


var all_textures:Array[TextureRect]

# func _ready()->void:
# 	for texture in all_textures:
# 		var atlas_tex := texture as AtlasTexture
# 		if atlas_tex:
# 			var atlas_size := atlas_tex.atlas.get_size()
# 			var region := atlas_tex.region
# 			var offset := region.position / atlas_size
# 			var scaleb := region.size / atlas_size
# 			var mat:ShaderMaterial= texture.material

# 			mat.set_shader_parameter("atlas_offset", offset)
# 			mat.set_shader_parameter("atlas_scale", scaleb)

### 3d prespective Shader stuff
var rot_x_max:float = 10
var rot_y_max:float = 10

func calcute_shader()->void:
	var local_pos:Vector2 = get_local_mouse_position()
	
	var rot_x:float = -(rot_x_max*2/size.x)*(local_pos.x-size.x/2)
	var rot_y:float = (rot_y_max*2/size.y)*(local_pos.y-size.y/2)
	set_shader(rot_x,rot_y)

func set_shader(rot_x:float,rot_y:float)->void:
	var time_dur:float = .15
	var tween:Tween = create_tween()
	for texture in all_textures:
		var mat:ShaderMaterial= texture.material
		tween.parallel().tween_property(mat,"shader_parameter/y_rot",rot_x,time_dur)
		tween.parallel().tween_property(mat,"shader_parameter/x_rot",rot_y,time_dur)

func flip_y(duration:float,delay:float,callme:Callable)->void:
	var tween:Tween = create_tween()
	for texture in all_textures:
		var mat:ShaderMaterial= texture.material
		tween.parallel().tween_property(mat,"shader_parameter/y_rot",90,duration).set_delay(delay)
		
		if tween.finished.is_connected(callme) == false:
			tween.finished.connect(callme)
		tween.finished.connect(func()->void:
			var new_tween:Tween = create_tween()
			mat.set_shader_parameter("y_rot",-90)
			new_tween.tween_property(mat,"shader_parameter/y_rot",0,duration)
		)
