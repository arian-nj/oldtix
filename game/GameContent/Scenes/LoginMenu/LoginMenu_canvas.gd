extends SceneLevel

func OnLoaded()->void:
	await KClient._instance.LoggedIn

	if await KClient._instance.RefetchME() == false:
		print("can't get me")
	manager_change_scene.emit(SceneManager.Levels.MainMenu)
