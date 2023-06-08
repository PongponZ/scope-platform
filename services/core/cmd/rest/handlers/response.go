package handlers

type Response struct {
	Success   bool        `json:"success"`
	ErrorCode string      `json:"error_code",omitempty`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
}
