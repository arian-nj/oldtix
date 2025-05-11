extends RefCounted
class_name KatanaLogger

enum LOG_LEVEL {NONE, ERROR, WARNING, INFO, VERBOSE, DEBUG}

var _level := LOG_LEVEL.ERROR
var _module := "Katana"

func _init(p_level : LOG_LEVEL = LOG_LEVEL.ERROR) -> void:
	_level = p_level

func _log(level : LOG_LEVEL, msg:String) -> void:
	if level <= _level:
		if level == LOG_LEVEL.ERROR:
			printerr("=== %s : ERROR === %s" % [_module, str(msg)])
		else:
			var what := "=== UNKNOWN === "
			for k:String in LOG_LEVEL:
				if level == LOG_LEVEL[k]:
					what = "=== %s : %s === " % [_module, k]
					break
			print(what + str(msg))

func error(msg:String) -> void:
	_log(LOG_LEVEL.ERROR, msg)
	ErrorBoard._instance.new_error(msg,ErrorBoard.ErrorLevel)

func warning(msg:String) -> void:
	_log(LOG_LEVEL.WARNING, msg)
	ErrorBoard._instance.new_error(msg,ErrorBoard.ErrorLevel)

func info(msg:String) -> void:
	_log(LOG_LEVEL.INFO, msg)
	ErrorBoard._instance.new_error(msg,ErrorBoard.InfoLevel)

## no _log
func success(msg:String) -> void:
	ErrorBoard._instance.new_error(msg,ErrorBoard.SuccessLevel)

func verbose(msg:String) -> void:
	_log(LOG_LEVEL.VERBOSE, msg)

func debug(msg:String) -> void:
	_log(LOG_LEVEL.DEBUG, msg)
