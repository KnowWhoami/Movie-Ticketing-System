package api

type Response struct {
	Success      bool        `json:"success"`
	Data         interface{} `json:"data"`
	ErrorMessage string      `json:"error_message,omitempty"`
}

func ErrorResponse(errorMessage string, _ int) *Response {
	return &Response{
		Success:      false,
		ErrorMessage: errorMessage,
	}
}

func SuccessResponse(payload interface{}) *Response {
	return &Response{
		Success: true,
		Data:    payload,
	}
}
