extends Control

@export var coinLabel:Label

func _ready() -> void:
	set_coin_label()
	KAccount._instance.MeChanged.connect(set_coin_label)

func set_coin_label()-> void:
	coinLabel.text = str(KAccount._instance.MyAccount.coin)
