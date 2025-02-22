extends Node

@export var MyTeamLabel:Label
@export var OtherTeamLabel:Label

@export var table:Game4Table

func _ready() -> void:
	table.GameDataUpdated.connect(_game_data_updated)


func _game_data_updated(game_data:GameData)->void:
	if table.game_data.current_trick == null:
		return
	
	var my_score:int
	var opp_score:int
	if table.game_data.your_team == GameData.TeamOne:
		my_score = table.game_data.team_one_trick_score
		opp_score = table.game_data.team_two_trick_score
	else:
		my_score = table.game_data.team_two_trick_score
		opp_score = table.game_data.team_one_trick_score
	
	MyTeamLabel.text = str(my_score)
	OtherTeamLabel.text = str(opp_score)
