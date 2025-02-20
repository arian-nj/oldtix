extends SceneLevel

func OnLoaded()->void:
	await KAccount.LoggedIn
	await KAccount.GotMe
	manager_change_scene.emit(SceneManger.Levels.MainMenu)
