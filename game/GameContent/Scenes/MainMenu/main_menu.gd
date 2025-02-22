extends Control


@export var StatusLabel:Label

func set_label(text:String)->void:
	StatusLabel.text = text

func _ready() -> void:
	var http_req:HTTPRequest = HTTPRequest.new()
	add_child(http_req)
	http_req.request(KAccount.StatusUrl)
	set_label("Waiting...")
	http_req.request_completed.connect(_on_status_request_completed)

func _on_status_request_completed(_result: int, response_code: int, _headers: PackedStringArray, body: PackedByteArray)->void:
	if response_code != HTTPClient.RESPONSE_OK:
		ErrorBoard._instance.new_error("somthing went wrong.",ErrorBoard.ErrorLevel)
		print_debug(str(response_code) +" response code ")
		return
	var json_data:Variant = JSON.parse_string(body.get_string_from_utf8())
	set_label(json_data["status"])
