class_name StatisticsSection extends Control

@export var WinLoseLine:BouncyLine
@export var TrickWinLoseLine:BouncyLine
@export var TurnWinLoseLine:BouncyLine

func get_animate_statistics() -> void:
	var result := await KClient._instance.GetUserStatistics(str(KClient._instance.MyAccount.id))
	var userStatistics:UserStatisticsData = result[0]
	var err:String = result[1]
	if err != "":
		Katana._instance.logger.error(err)
	if userStatistics == null:
		return
	
	await WinLoseLine.animate(userStatistics.win,userStatistics.lose)
	await TrickWinLoseLine.animate(userStatistics.tricks_won,userStatistics.tricks_lost)
	await TurnWinLoseLine.animate(userStatistics.turns_won,userStatistics.turns_lost)


func redo_animate_statistics() -> void:
	WinLoseLine.redo_animation()
	TrickWinLoseLine.redo_animation()
	TurnWinLoseLine.redo_animation()
