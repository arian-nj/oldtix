class_name UserComponent extends Button

@export var user_id:int:
	set(value):
		user_id = value
		request_user_data()

var userStat:UserStatisticsData = null

@onready var mat: ShaderMaterial =self.material
@export var start_color:Color
@export var end_color:Color

var tween :Tween

func _ready() -> void:
	pressed.connect(_on_pressed)


func start_shader(timer_time:int) -> void:
	mat.set_shader_parameter("current_val",10)
	mat.set_shader_parameter("fill_color",start_color)
	
	tween = create_tween()	
	tween.finished.connect(stop_timer_shader)
	tween.parallel().tween_property(mat,"shader_parameter/current_val",90,timer_time)
	tween.parallel().tween_property(mat,"shader_parameter/fill_color",end_color,timer_time)	

func stop_timer_shader()->void:
	if tween:
		tween.kill()
	tween = create_tween()
	tween.tween_property(mat,"shader_parameter/current_val",0,.5)
	

func request_user_data() -> void:
	if user_id == 0:
		self.text = "Bot"
		return
	var userData := await KClient._instance.GetUser(str(user_id))
	if userData == null:
		self.text = "Bot"
	else:
		self.text = userData.display_name
	
	request_user_stats()

func request_user_stats()->void:
	var result := await KClient._instance.GetUserStatistics(str(user_id))
	self.userStat = result[0]
	var err:String = result[1]
	if err != "":
		return

func _on_pressed()->void:
	if userStat == null:
		print_debug("user stat is null")
		return
	print(userStat.user_id ," ",userStat.win, " ",userStat.lose)

# send_request
# parse request
# set it in button text
