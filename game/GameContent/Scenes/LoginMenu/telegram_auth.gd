extends Control

func _ready() -> void:
	# Katana._instance.change_debug_mode(true)
	KClient._instance.setup_telegram_token()
