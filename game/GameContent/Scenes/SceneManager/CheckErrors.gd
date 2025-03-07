# ------------------------------------------------------------------------------
# Running this code at the command line does not show warnings.  Which I found 
# odd.  But if you make a scene and add this as the script and run it from the 
# editor it will:
#   Set all debug-gdscript warnings to "warn"
#   Load all the .gd files it can find (recursively)
# You get all the warnings you want in the Debugger tab
# ------------------------------------------------------------------------------
class_name ErrorChecker
extends Node
var include_subdirectories := true


# ------------------------------------------------------------------------------
# Gets all the files in a directory and all subdirectories if include_subdirectories
# is true.  The files returned are all sorted by name.
# ------------------------------------------------------------------------------
func _get_files(path: String, prefix: String, suffix: String) -> Array[String]:
	var files: Array[String] = []
	var directories: Array[String] = []

	var d := DirAccess.open(path)
	# true parameter tells list_dir_begin not to include "." and ".." directories.
	d.list_dir_begin() # TODO 4.0 fill missing arguments https://github.com/godotengine/godot/pull/40547

	# Traversing a directory is kinda odd.  You have to start the process of listing
	# the contents of a directory with list_dir_begin then use get_next until it
	# returns an empty string.  Then I guess you should end it.
	var fs_item := d.get_next()
	var full_path := ''
	while fs_item != '':
		full_path = path.path_join(fs_item)

		#file_exists returns fasle for directories
		if d.file_exists(full_path):
			if fs_item.begins_with(prefix) and fs_item.ends_with(suffix):
				files.append(full_path)
		elif include_subdirectories and d.dir_exists(full_path):
			directories.append(full_path)

		fs_item = d.get_next()
	d.list_dir_end()

	for dir in range(directories.size()):
		var dir_files := _get_files(directories[dir], prefix, suffix)
		for i in range(dir_files.size()):
			files.append(dir_files[i])

	files.sort()
	return files


func set_all_gdscript_warnings_to_warning() -> void:
	var props := ProjectSettings.get_property_list()
	for prop in props:
		var prop_name: String = prop.name
		var prop_hint_string: String = prop.hint_string
		if prop_name.begins_with('debug/gdscript/warnings'):
			if prop_hint_string == 'Ignore,Warn,Error':
				# print(prop_name, ' -> ', 1)
				ProjectSettings.set_setting(prop_name, 1)


func load_all_scripts() -> void:
	var files := _get_files('res://', '', '.gd')
	for f in files:
		# print(f)
		var _thing := load(f)


func _ready() -> void:
	# set_all_gdscript_warnings_to_warning()
	load_all_scripts()