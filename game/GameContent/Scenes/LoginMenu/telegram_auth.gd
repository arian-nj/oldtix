extends Control

func _ready() -> void:
	var console := JavaScriptBridge.get_interface("console")
	console.log("hello from godot console bitch 2")

	var parent_window := JavaScriptBridge.get_interface("parent")
	var initdata : Variant= parent_window.getInitdata()
	console.log("from godot init data "+initdata)
	# console.log("from inside godot "+ webapp.initDataUnsafe.user)
