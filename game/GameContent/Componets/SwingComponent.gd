class_name SwingComponent extends Node

@export var CardButton:Button

@export var spring:float = 200
@export var damp:float = 12
var displacement:float = 0

var velocity:float = 0
var last_position:Vector2
var x:bool = false

func swing(delta:float)->void:
	
	velocity = (CardButton.global_position.x - last_position.x)

	var dir :float= 1
	if CardButton.pivot_offset.y > CardButton.size.y/2:
		dir = -1

	last_position = CardButton.global_position

	var force:float = -spring * displacement - damp *velocity
	velocity += force * delta
	displacement += velocity * delta
	CardButton.rotation_degrees = displacement * 10 * dir