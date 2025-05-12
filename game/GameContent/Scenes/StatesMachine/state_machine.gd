extends Node

@export var initial_state:State

var current_state:State
var states:Dictionary = {}

func _ready() -> void:
	
	for child:Node in get_children():
		if child is State:
			child.set_process(false)
			# var nchild:State = child
			states[child.name.to_lower()] = child
	
	if initial_state:
		initial_state.Enter()
		current_state = initial_state
		current_state.StateTransition.connect(on_child_transition)


func _process(delta: float) -> void:
	if current_state:
		current_state.Update(delta)

func _physics_process(delta: float) -> void:
	if current_state:
		current_state.Physics_Update(delta)

func on_child_transition(state:State,new_state_name:String)->void:
	if state != current_state:
		return
	
	var new_state:State = states.get(new_state_name.to_lower())
	if !new_state:
		print_debug("no state named: ",new_state_name)
		return
	new_state.StateTransition.connect(on_child_transition)
	
	if current_state:
		current_state.StateTransition.disconnect(on_child_transition)
		current_state.set_process(false)
		current_state.Exit()
	
	new_state.Enter()
	new_state.set_process(true)
	current_state = new_state
