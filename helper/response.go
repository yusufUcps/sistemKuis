package helper

func FormatResponse(message string, data any, paging any) map[string]any {
	var response = map[string]any{}
	response["message"] = message
	if data != nil {
		response["data"] = data
	}
	if paging != nil {
		response["token"] = paging
	}
	return response
}

