package helper

func FormatResponse(message string, data any) map[string]any {
	var response = map[string]any{}
	response["message"] = message
	if data != nil {
		response["data"] = data
	}
	return response
}

func FormatResponseJWT(message string, data any, token string) map[string]any {
	var response = map[string]any{}
	response["message"] = message
	response["token"] = token
	if data != nil {
		response["data"] = data
	}
	return response
}
