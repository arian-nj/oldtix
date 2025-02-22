extends SceneLevel

func OnLoaded()->void:
	await KAccount._instance.LoggedIn

	if await KAccount._instance.RefetchME() == false:
		print("can't get me")
	manager_change_scene.emit(SceneManger.Levels.MainMenu)
