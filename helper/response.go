package helper

func FormatResponse(message string, data any, token string) map[string]any {
	var response = map[string]any{}
	response["message"] = message
	if data != nil {
		response["data"] = data
	}
	if token != "" {
		response["token"] = token
	}
	return response
}

