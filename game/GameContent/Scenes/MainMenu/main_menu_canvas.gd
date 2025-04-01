extends SceneLevel

@export var aspectRationContainer:CustomRatioAspectContainer
@export var betPanel:BetPanel

func _ready() -> void:
	resized.connect(aspectRationContainer._on_aspect)
	betPanel.BetAmountChoosed.connect(_on_bet_amount_choosed)

func _on_bet_amount_choosed(coin_amount:int)->void:
	level_parameters[SharedBetAmount] = coin_amount
	manager_change_scene.emit(SceneManager.Levels.GameHokm4)


func _on_continue_button_pressed() -> void:
	manager_change_scene.emit(SceneManager.Levels.GameHokm4)
