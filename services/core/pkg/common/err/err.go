package err

type Err struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func New(code string, message string) Err {
	return Err{
		Code:    code,
		Message: message,
	}
}

func (e Err) Error() string {
	return e.Message
}
