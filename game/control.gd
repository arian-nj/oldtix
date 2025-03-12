extends Control

var table_draw_queue:Array[Callable]
@onready var timer := Timer.new()

func Greet(hname:String)->void:
	print("hello "+ hname)
	await get_tree().create_timer(1).timeout
	print("bye "+ hname)


func Greet2(hname:String)->void:
	print("hello "+ hname)
	var x := 0
	for i in range(10_000_000):
		x += 1
	print(x)
	print("bye "+ hname)

func push_callback(c:Callable)->void:
	table_draw_queue.push_back(c)

func _ready() -> void:
	# push_callback(Greet.bind("hiii1"))
	push_callback(Greet2.bind("hiii2"))
	# push_callback(Greet.bind("hiii3"))

	timer.wait_time = .2
	timer.timeout.connect(run_actions)
	add_child(timer)
	timer.start()

func run_actions()->void:
	while true:
		var action_variant:Variant = table_draw_queue.pop_front()
		if action_variant != null:			
			var action:Callable = action_variant
			print("start")
			await action.call()
			print("end")
		else:
			await get_tree().create_timer(.2).timeout
