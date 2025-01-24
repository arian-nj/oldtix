extends Node
class_name State

signal Transition(state:State,new_state_name:String)

func Enter()->void:
    pass

func Exit()->void:
    pass


func Update(delta:float)->void:
    pass

func Physics_Update(delta:float)->void:
    pass