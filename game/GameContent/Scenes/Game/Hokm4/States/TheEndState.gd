extends State

@export var ws:KatanaSocket
@export var status_label:Label
@export var table:Game4Table
@export var win_panel:WinPanel
@export var trick_score_panel:TrickScore


func Enter()->void:
	ws.set_process(false)
	status_label.text = "The End"
	win_panel.visible = true
	trick_score_panel.end_game()

func _process(_delta: float) -> void:
	pass

func Exit()->void:
	pass