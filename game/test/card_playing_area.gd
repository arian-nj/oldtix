class_name PlayArea extends Control

@export var staticBody:StaticBody2D
@export var collisionShape:CollisionShape2D
@export var colorRect:ColorRect

var rectCollisionShape : RectangleShape2D

@export var turnOnColor:Color
@export var max_trans:float

@export_range(0,1) var x_percentage:float = .7
@export_range(0,1) var y_percentage:float = .3

@export var offset:Vector2

func _ready() -> void:
	get_viewport().size_changed.connect(set_default_size)
	set_default_size()
	colorRect.color = turnOnColor

func _set_both_size(view_size:Vector2)->void:
	rectCollisionShape = collisionShape.shape
	var target_size:Vector2

	target_size.x = view_size.x * x_percentage
	target_size.y = view_size.y * y_percentage


	rectCollisionShape.size = target_size
	colorRect.size = target_size

	staticBody.position = (view_size/2) + offset
	colorRect.position =  (-target_size/2) 


func set_default_size() -> void:
	_set_both_size(self.size)

func turn_on() -> void:
	set_default_size()
	colorRect.color.a = max_trans

func turn_off() -> void:
	colorRect.color.a = 0

func _on_button_button_down() -> void:
	turn_on()

func _on_button_button_up() -> void:
	turn_off()
