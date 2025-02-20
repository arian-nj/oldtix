class_name GameData extends Resource


@export var id:String
@export var your_team:int
@export var players:Array[PlayerData]

@export var current_trick:TrickData

@export var team_one_trick_score:int
@export var team_two_trick_score:int

enum {
	TeamOne = 0,
	TeamTwo = 1
}