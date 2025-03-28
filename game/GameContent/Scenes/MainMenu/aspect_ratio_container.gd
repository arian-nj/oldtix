class_name CustomRatioAspectContainer extends AspectRatioContainer


@export var max_aspect:float

func _ready() -> void:
	_on_aspect()
	resized.connect(_on_aspect)

func _on_aspect()->void:
	var view_rect := get_viewport_rect()
	var aspect := view_rect.size.x/view_rect.size.y
	aspect = min(aspect,max_aspect)
	self.ratio = aspect
	