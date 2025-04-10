class_name UserComponent extends Button

var user_id:int:
	set(value):
		user_id = value
		request_user_data()

var userStat:UserStatisticsData = null

func _ready() -> void:
	pressed.connect(_on_pressed)

func request_user_data() -> void:
	if user_id == 0:
		self.text = "Bot"
		return
	var userData := await KAccount._instance.GetUser(str(user_id))
	if userData == null:
		self.text = "Bot"
	else:
		self.text = userData.display_name
	
	request_user_stats()

func request_user_stats()->void:
	self.userStat = await KAccount._instance.GetUserStatistics(str(user_id))

func _on_pressed()->void:
	if userStat == null:
		print_debug("user stat is null")
		return
	print(userStat.user_id ," ",userStat.win, " ",userStat.lose)

# send_request
# parse request
# set it in button text
