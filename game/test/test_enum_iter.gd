extends Node

enum LOG_LEVEL {NONE, ERROR, WARNING, INFO, VERBOSE, DEBUG}

func _ready() -> void:
    for k:String in LOG_LEVEL:
        print(LOG_LEVEL[k])
        var what := "=== %s : %s === " % ["Katana", k]
        print(what)