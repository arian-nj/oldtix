extends Control

@export var UsernameLabel:Label

func _ready() -> void:
	UsernameLabel.text = KAccount._instance.MyAccount.username
