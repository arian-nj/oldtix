extends State

@export var ws:KatanaSocket
@export var status_label:Label
@export var table:Game4Table
@export var win_panel:WinPanel
@export var trick_score_panel:TrickScore


func Enter()->void:
	ws.set_process(false)
	status_label.text = "The End"
	ws.NewEventSig.connect(_on_new_event)
	ws.open_events()
	win_panel.visible = true
	trick_score_panel.end_game()

func _on_new_event(_e:KEvent.Event)->void:
	pass

func Exit()->void:
	ws.NewEventSig.disconnect(_on_new_event)
