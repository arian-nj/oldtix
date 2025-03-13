extends Control


@export var StatusLabel:Label

@export var PlayButton:Button
@export var ContinueGameButton:Button

func set_label(text:String)->void:
	StatusLabel.text = text

func _ready() -> void:
	PlayButton.disabled = true
	ContinueGameButton.visible = false

	var http_req:HTTPRequest = HTTPRequest.new()
	add_child(http_req)
	http_req.request(Katana.ActiveGameUrl,KAccount._instance.AddAuthHeader())
	http_req.request_completed.connect(_on_active_game_request_completed)

	
	set_label("Waiting...")
	http_req = HTTPRequest.new()
	add_child(http_req)
	http_req.request(Katana.StatusUrl)
	http_req.request_completed.connect(_on_status_request_completed)

func _on_active_game_request_completed(_result: int, response_code: int, _headers: PackedStringArray, body: PackedByteArray)->void:
	if response_code != HTTPClient.RESPONSE_OK:
		ErrorBoard._instance.new_error("somthing went wrong getting active game.",ErrorBoard.ErrorLevel)
		print_debug(str(response_code) +" response code ")
		return
	var active_game:ActiveGameData = JsonClassConverter.json_string_to_class(ActiveGameData,body.get_string_from_utf8())
	if active_game.is_active:
		ContinueGameButton.visible = true
	else:
		PlayButton.disabled = false

func _on_status_request_completed(_result: int, response_code: int, _headers: PackedStringArray, body: PackedByteArray)->void:
	if response_code != HTTPClient.RESPONSE_OK:
		ErrorBoard._instance.new_error("somthing went wrong getting status.",ErrorBoard.ErrorLevel)
		print_debug(str(response_code) +" response code ")
		return
	var json_data:Variant = JSON.parse_string(body.get_string_from_utf8())
	set_label(json_data["status"])
