extends Control

@export var UsernameLabel:Label
@export var DisplaynameLabel:Label

func _ready() -> void:
	DisplaynameLabel.text = KAccount._instance.MyAccount.display_name
	UsernameLabel.text = "@"+KAccount._instance.MyAccount.username
