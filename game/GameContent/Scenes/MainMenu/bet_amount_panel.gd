class_name BetPanel extends Control

@export var betVBoxContainer:PanelContainer

@export var firstBetButton:Button
@export var firstBetLabel:Label
@export var secondBetButton:Button

@export var playButton:Button
@export var continueGameButton:Button

signal BetAmountChoosed(coin_amount:int)


func _ready() -> void:
	betVBoxContainer.visible = false
	secondBetButton.pressed.connect(fire_play_button.bind(250))
	if KAccount._instance.MyAccount.coin < 10:
		firstBetLabel.text = "0 -> 10"
		firstBetButton.pressed.connect(fire_play_button.bind(0))
		secondBetButton.disabled = true

	else:
		firstBetButton.pressed.connect(fire_play_button.bind(10))
	
	playButton.disabled = true
	continueGameButton.visible = false

	send_active_game_request()

func send_active_game_request()->void:
	var http_req:HTTPRequest = Katana.NewHttpRequest()
	add_child(http_req)
	http_req.request(Katana.Hokm4HttpUrl + Katana.ActiveGameUrl,KAccount._instance.AddAuthHeader())
	http_req.request_completed.connect(_on_active_game_request_completed)
	await http_req.request_completed
	http_req.queue_free()

func _on_active_game_request_completed(_result: int, response_code: int, _headers: PackedStringArray, body: PackedByteArray)->void:
	if response_code != HTTPClient.RESPONSE_OK:
		ErrorBoard._instance.new_error("خطا در اتصال به بازی" + str(response_code),ErrorBoard.ErrorLevel)
		ErrorBoard._instance.new_error("در حال تلاش دوباره",ErrorBoard.InfoLevel)
		print_debug(str(response_code) +" response code ")
		get_tree().create_timer(5).timeout.connect(send_active_game_request)
		return
	
	var active_game:ActiveGameData = JsonClassConverter.json_string_to_class(ActiveGameData,body.get_string_from_utf8())
	if active_game.is_active:
		continueGameButton.visible = true
	else:
		playButton.disabled = false	
		get_tree().create_timer(5).timeout.connect(send_active_game_request)


func fire_play_button(coin_amount:int)->void:
	BetAmountChoosed.emit(coin_amount)

func panel_go_up()->void:
	playButton.disabled = true
	betVBoxContainer.visible = true
	
	betVBoxContainer.size = self.size
	betVBoxContainer.size.y = self.size.y/2

	betVBoxContainer.global_position = self.global_position
	betVBoxContainer.global_position.y += self.size.y

	var target_position := self.global_position
	target_position.y += self.size.y/2

	var tween := create_tween().set_ease(Tween.EASE_IN).set_trans(Tween.TRANS_LINEAR)
	tween.parallel().tween_property(betVBoxContainer,"global_position",target_position,.6)
	await tween.finished

func panel_go_down()->void:
	var target_position:Vector2 = self.global_position 
	target_position.y += self.size.y

	var tween := create_tween().set_ease(Tween.EASE_OUT).set_trans(Tween.TRANS_BACK)
	tween.parallel().tween_property(betVBoxContainer,"global_position",target_position,0.7)
	await tween.finished
	betVBoxContainer.visible = false
	playButton.disabled = false


func _on_play_tic_2_button_pressed() -> void:
	panel_go_up()


func _on_button_pressed() -> void:
	panel_go_down()
