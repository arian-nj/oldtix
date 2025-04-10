class_name DownloadProgressBar extends ProgressBar

func start_downloading(domain:String,pack_url:String)->String:
	var err := 0
	var client := HTTPClient.new()
	err = client.connect_to_host(domain,443,TLSOptions.client())
	if err != OK:
		return "can't connect to domain err code => " +str(err)

	while client.get_status() == HTTPClient.STATUS_CONNECTING or client.get_status() == HTTPClient.STATUS_RESOLVING:
		client.poll()
		# print("connecting")
		await get_tree().process_frame
	
	if client.get_status() != HTTPClient.STATUS_CONNECTED:
		return "can't establish a connection"

	var headers := [
		"User-Agent: Pirulo/1.0 (Godot)",
		"Accept: */*"
	]

	err = client.request(HTTPClient.METHOD_GET,pack_url,headers)
	if err != OK:
		return "can't send request"
	
	while client.get_status() == HTTPClient.STATUS_REQUESTING:
		client.poll()
		# print("Requesing...")
		await get_tree().process_frame
	
	if (client.get_status() != HTTPClient.STATUS_BODY and client.get_status() != HTTPClient.STATUS_CONNECTED): # Make sure request finished well.
		return "error in requesting status: " + str(client.get_status())
	
	# print("response? ",client.has_response())

	if client.has_response():
		# var res_headers := client.get_response_headers_as_dictionary()
		# print("code: ", client.get_response_code()) # Show response code.
		# print("**headers:\\n", res_headers) # Show headers.
		
		if client.is_response_chunked():
			print("Response is chuncked")
		else:
			var bl := client.get_response_body_length()
			print("Response Length: ", bl)
		
		var rb := PackedByteArray()
		while client.get_status() == HTTPClient.STATUS_BODY:
			client.poll()
			var chunk := client.read_response_body_chunk()
			if chunk.size() == 0:
				await get_tree().process_frame
			else:
				rb = rb + chunk
				_send_loading_signal.bind(rb.size(),client.get_response_body_length()).call_deferred()
		_send_loading_signal.bind(rb.size(),client.get_response_body_length()).call_deferred()
		
		# print("bytes got:",rb.size())
		# var text := rb.get_string_from_ascii()
		# print("Text: ",text)
		var f := FileAccess.open("user://system.pck",FileAccess.WRITE)
		f.store_buffer(rb)
		f.close()
	return ""
	

func _send_loading_signal(l:int,t:int) -> void:
	var new_val := (float(l)/float(t)) * max_value
	# print(new_val)
	self.value = new_val
