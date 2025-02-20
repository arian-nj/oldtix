class_name ErrorClass extends CanvasLayer

@export var VContainer:VBoxContainer
@export var ErrorRowScene:PackedScene

@export var ErrorLevelColor:Color
@export var InfoLevelColor:Color
@export var SuccessLevelColor:Color

enum {
	ErrorLevel,
	InfoLevel,
	SuccessLevel,
}

func new_error(error_message:String,error_level:int)->void:

	var new_err:ErrorRow = ErrorRowScene.instantiate()
	new_err.label.text = error_message

	if error_level == ErrorLevel:
		new_err.color_rect.color = ErrorLevelColor
	if error_level == InfoLevel:
		new_err.color_rect.color = InfoLevelColor
	if error_level == SuccessLevel:
		new_err.color_rect.color = SuccessLevelColor
	
	VContainer.add_child(new_err)
	await get_tree().create_timer(3).timeout
	if new_err != null:
		new_err.queue_free()
