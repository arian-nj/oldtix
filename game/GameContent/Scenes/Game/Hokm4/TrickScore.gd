class_name TrickScore extends Node

@export var MyTeamLabel:Label
@export var OtherTeamLabel:Label

@export var MyTeamPanel:PanelContainer
@export var OtherTeamPanel:PanelContainer

@export var table:Game4Table

@onready var view_size:Vector2

func update_view_size()->void:
	view_size = get_viewport().get_visible_rect().size

func _ready() -> void:
	table.GameDataUpdated.connect(_game_data_updated)
	update_view_size()
	get_viewport().size_changed.connect(update_view_size)

func get_scores()-> Array[int]:
	var my_score:int = 0
	var opp_score:int = 0

	if table.game_data.your_team == GameData.TeamOne:
		my_score = table.game_data.team_one_trick_score
		opp_score = table.game_data.team_two_trick_score
	else:
		my_score = table.game_data.team_two_trick_score
		opp_score = table.game_data.team_one_trick_score
	return [my_score,opp_score]

var sum_of_trick_scores := 0
func _game_data_updated(game_data:GameData)->void:
	if table.game_data.current_trick == null:
		return

	var r := get_scores()
	var my_score := r[0]
	var opp_score := r[1]
	
	MyTeamLabel.text = str(my_score)
	OtherTeamLabel.text = str(opp_score)
	var new_sum_of_trick_scores := my_score + opp_score
	if new_sum_of_trick_scores != sum_of_trick_scores:
		sum_of_trick_scores = new_sum_of_trick_scores
		data_updated_animation()

func end_game()->void:
	var panel_size := MyTeamPanel.size
	 
	var r := get_scores()
	var my_score := r[0]
	var opp_score := r[1]
	if my_score + opp_score == 0:
		return

	var available_space :float = view_size.x - (panel_size.x*2)
	var chunk_len:float = available_space / (my_score + opp_score)
	
	var tween := create_tween().set_trans(Tween.TRANS_BOUNCE)
	tween.tween_property(OtherTeamPanel,"size:x",OtherTeamPanel.size.y + chunk_len*opp_score,1)

	tween.parallel().tween_property(MyTeamPanel,"position:x",MyTeamPanel.position.x - chunk_len*my_score,1)
	tween.parallel().tween_property(MyTeamPanel,"size:x",MyTeamPanel.size.x + chunk_len*my_score,1)

func data_updated_animation()->void:
	var panel_size := MyTeamPanel.size
	 
	var r := get_scores()
	var my_score := r[0]
	var opp_score := r[1]
	if my_score + opp_score == 0:
		return

	var available_space :float = view_size.x - (panel_size.x*2)
	available_space = available_space/3
	var chunk_len:float = available_space / (my_score + opp_score)
	
	var last_other_size_x := OtherTeamPanel.size.x

	var last_my_pos_x := MyTeamPanel.position.x
	var last_my_size_x := MyTeamPanel.size.x

	var tween := create_tween().set_trans(Tween.TRANS_BOUNCE)
	tween.tween_property(OtherTeamPanel,"size:x",OtherTeamPanel.size.y + chunk_len*opp_score,1)

	tween.parallel().tween_property(MyTeamPanel,"position:x",MyTeamPanel.position.x - chunk_len*my_score,1)
	tween.parallel().tween_property(MyTeamPanel,"size:x",MyTeamPanel.size.x + chunk_len*my_score,1)
	await tween.finished
	
	
	tween = create_tween().set_trans(Tween.TRANS_BOUNCE)
	tween.tween_property(OtherTeamPanel,"size:x",last_other_size_x,1)

	tween.parallel().tween_property(MyTeamPanel,"position:x",last_my_pos_x,1)
	tween.parallel().tween_property(MyTeamPanel,"size:x",last_my_size_x,1)
	await tween.finished