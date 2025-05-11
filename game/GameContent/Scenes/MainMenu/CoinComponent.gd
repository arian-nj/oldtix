extends Control

@export var coinLabel:Label

func _ready() -> void:
	KClient._instance.MeChanged.connect(set_coin_label)
	set_coin_label()

func set_coin_label()-> void:
	coinLabel.text = str(KClient._instance.MyAccount.coin)
