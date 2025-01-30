extends Node
class_name State

signal Transition(state:State,new_state_name:String)

func Enter()->void:
    Transition.is_null()
    pass

func Exit()->void:
    pass


func Update(_delta:float)->void:
    pass

func Physics_Update(_delta:float)->void:
    pass