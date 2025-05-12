extends Node
class_name State

signal StateTransition(state:State,new_state_name:String)

func _init() -> void:
    set_process(false)

func Enter()->void:
    StateTransition.is_null()
    pass

func Exit()->void:
    pass


func Update(_delta:float)->void:
    pass

func Physics_Update(_delta:float)->void:
    pass