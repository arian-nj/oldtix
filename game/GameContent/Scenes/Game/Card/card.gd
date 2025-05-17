@tool
class_name Card extends Button

signal card_played(card:Card)

signal NotInplaceSig(card:Card)

@export var cardDrawingTextureRect:TextureRect
@export var cardBackgroundTextureRect:TextureRect

@export var prespective3DShader:CardShader
@export var cardDrawingAtlasTexture:Texture2D

@export var cardBackgroundTexture:Texture2D

@export var cardArea:Area2D
@export var swingComponent:SwingComponent

var local_pos_on_press:Vector2 
var start_choosed_pos:Vector2 

@export var card_data:CardData

var default_scale:Vector2
var big_scale:Vector2

var in_hand:bool
var is_touching_area:bool
@export var can_be_selected:bool = true

func _ready()->void:
	prespective3DShader.all_textures.append_array([cardDrawingTextureRect,cardBackgroundTextureRect])
	set_process(false)
	default_scale = scale
	big_scale = scale * 1.1

	if can_be_selected:
		cardArea.body_entered.connect(_body_entered)
		cardArea.body_exited.connect(_body_exited)

		button_up.connect(_on_button_up)
		button_down.connect(_on_button_down)	
		
	set_process(true)

func _process(delta:float)->void:
	if can_be_selected and button_pressed:
		set_card_position()

		prespective3DShader.calcute_shader()
		swingComponent.swing(delta)

func set_card_position() -> void:
	var mouse_pos:Vector2 = get_global_mouse_position()
	global_position = mouse_pos-(local_pos_on_press * scale)

func _body_entered(_body:Node2D)->void:
	is_touching_area = true

func _body_exited(_body:Node2D)->void:
	is_touching_area = false

func _on_button_up() -> void:
	print("up")
	if is_touching_area:
		print("here1")
		card_played.emit(self)
	else :
		NotInplaceSig.emit(self)
	
	prespective3DShader.set_shader(0,0)

	change_scale(default_scale)
	
	rotation = 0


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
	change_scale(big_scale)

	swingComponent.last_position = global_position


### scale stuff
var scale_tween:Tween
func change_scale(new_scale:Vector2)->void:
	if scale_tween and scale_tween.is_running():
		scale_tween.kill()
	scale_tween= create_tween().set_ease(Tween.EASE_IN)
	scale_tween.tween_property(self,"scale",new_scale,.3)

## assets

func load_assets()->void:
	var region := Vector2((self.size.x * (card_data.value-2)),(self.size.y * (card_data.suit)))
	var croped_image := cardDrawingAtlasTexture.get_image().get_region(Rect2i(region,self.size))
	cardDrawingTextureRect.texture = ImageTexture.create_from_image(croped_image)
	cardBackgroundTextureRect.texture = cardBackgroundTexture
