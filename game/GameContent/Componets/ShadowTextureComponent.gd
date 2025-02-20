class_name ShadowComponent extends TextureRect

@export var MainObject:Control

var max_offset_shadow:float = -10

func handle_shaddow()->void:
	var center:Vector2 = get_viewport_rect().size
	var distancex:float = MainObject.global_position.x - (center.x/2) + (MainObject.size.x/2)
	var distancey:float = MainObject.global_position.y - center.y

	position.x = lerp(0.0,-sign(distancex)*max_offset_shadow,abs(distancex/(center.x)))
	position.y = lerp(0.0,-sign(distancey)*max_offset_shadow,abs(distancey/(center.y)))

func _physics_process(_delta:float)->void:
	handle_shaddow()