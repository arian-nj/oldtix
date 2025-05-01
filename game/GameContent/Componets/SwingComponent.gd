class_name SwingComponent extends Node

@export var CardButton:Button

@export var spring:float = 200
@export var damp:float = 12
var displacement:float = 0

var x_velocity:float = 0
var last_position:Vector2

func swing(delta:float)->void:
	
	x_velocity = (CardButton.global_position.x - last_position.x)

	var dir :float= 1
	if CardButton.pivot_offset.y > CardButton.size.y/2:
		dir = -1

	last_position = CardButton.global_position

	var force:float = -spring * displacement - damp *x_velocity
	x_velocity += force * delta
	displacement += x_velocity * delta

	CardButton.rotation_degrees = displacement * 8 * dir
