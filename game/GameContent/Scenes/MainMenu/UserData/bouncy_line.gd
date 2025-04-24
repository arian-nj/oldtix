class_name  BouncyLine extends Control


@export var MyLabel:Label
@export var OtherLabel:Label

@export var MyPanel:PanelContainer
@export var OtherPanel:PanelContainer
@export var boner:bool = false

@onready var last_other_size_x:float = OtherPanel.size.x

@onready var last_my_pos_x:float = MyPanel.position.x
@onready var last_my_size_x:float = MyPanel.size.x
var tween:Tween

func animate(my_score:float,opp_score:float) -> void:
	if boner:
		return
	if my_score + opp_score == 0:
		return
	boner = true

	if tween and tween.is_running():
		tween.kill()

	var panel_size := MyPanel.size
	 
	if my_score + opp_score == 0:
		return
	
	MyLabel.text = str(my_score)
	OtherLabel.text = str(opp_score)

	var available_space :float = self.size.x - (panel_size.x*2)
	var chunk_len:float = available_space / (my_score + opp_score)
	
	last_other_size_x = OtherPanel.size.x

	last_my_pos_x = MyPanel.position.x
	last_my_size_x = MyPanel.size.x

	tween = create_tween().set_trans(Tween.TRANS_BOUNCE)
	tween.tween_property(OtherPanel,"size:x",OtherPanel.size.x + chunk_len*opp_score,1)

	tween.parallel().tween_property(MyPanel,"position:x",MyPanel.position.x - chunk_len*my_score,1)
	tween.parallel().tween_property(MyPanel,"size:x",MyPanel.size.x + chunk_len*my_score,1)
	await tween.finished
	
func redo_animation()->void:
	if boner == false:
		return
	boner = false

	if tween and tween.is_running():
		tween.kill()

	tween = create_tween()
	tween.tween_property(OtherPanel,"size:x",last_other_size_x,0.5)

	tween.parallel().tween_property(MyPanel,"position:x",last_my_pos_x,0.5)
	tween.parallel().tween_property(MyPanel,"size:x",last_my_size_x,0.5)
	await tween.finished
