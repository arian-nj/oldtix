extends Control

@export var UsernameLabel:Label

func _ready() -> void:
	UsernameLabel.text = KAccount.MyAccount.username
	