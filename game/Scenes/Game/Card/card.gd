class_name Card extends Button

signal card_played(card:Card)

@export var cardTexture:TextureRect
@export var cardArea:Area2D
@export var prespective3DShader:Prescpective3DShader
@export var swingComponent:SwingComponent

var local_pos_on_press:Vector2 
var start_choosed_pos:Vector2 

@export var suit :CardSuites
@export var value:int

var default_scale:Vector2
var in_hand:bool

enum CardSuites { # sync it with internal/card/card.go
	Heart = 0,
	Club = 1,
	Diamond = 2,
	Spade = 3
}


func suite_name() -> String:
	var sname:String = ""

	match suit:
		CardSuites.Club:
			sname ="7"
		CardSuites.Diamond:
			sname ="4"
		CardSuites.Heart:
			sname ="2"
		CardSuites.Spade:
			sname ="5"
	
	return sname

func value_name() -> String:
	match value:
		11:return "J"
		12:return "Q"
		13:return "K"
		1:return "A."
		_:return str(value)+"."

func get_assets_path() -> String:
	return "res://assets/cards/"+value_name()+suite_name()+".png"

func _ready()->void:
	# print(global_position)
	set_process(false)
	default_scale = scale

	# load assets
	var file_name:String = get_assets_path()
	var img:CompressedTexture2D = load(file_name)
	cardTexture.texture = img

	cardArea.body_entered.connect(_body_entered)
	button_up.connect(_on_button_up)
	button_down.connect(_on_button_down)	
	set_process(true)




func _body_entered(_body:Node2D)->void:
	card_played.emit(self)

func _process(delta:float)->void:
	if button_pressed:
		set_card_position()

		prespective3DShader.calcute_shader()
		swingComponent.swing(delta)


func set_card_position() -> void:
	var mouse_pos:Vector2 = get_global_mouse_position()
	global_position = mouse_pos-(local_pos_on_press * scale)

func _on_button_up() -> void:
	prespective3DShader.set_shader(0,0)

	change_scale(default_scale)
	
	rotation = 0
	
	# pivot_offset = Vector2(0.0,0.0)

func _on_button_down() -> void:
	start_choosed_pos = global_position
	local_pos_on_press = self.get_local_mouse_position()
	var old_z_index:int = self.z_index
	button_up.connect(func()->void:
		self.z_index = old_z_index
	)
	self.z_index = 20
	pivot_offset = local_pos_on_press
	rotation_degrees = 0

	prespective3DShader.calcute_shader()
	change_scale(default_scale*1.1)
	swingComponent.last_position = global_position

### scale stuff
var scale_tween:Tween
func change_scale(new_scale:Vector2)->void:
	if scale_tween and scale_tween.is_running():
		scale_tween.kill()
	scale_tween= create_tween().set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_BACK)
	scale_tween.tween_property(self,"scale",new_scale,.5)
