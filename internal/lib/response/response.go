package response

const (
	StatusOk  = "Ok"
	StatusErr = "Error"
)

type Response struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func Ok(msg string) *Response {
	return &Response{
		Status:  StatusOk,
		Message: msg,
	}
}

func Err(msg string) *Response {
	return &Response{
		Status:  StatusErr,
		Message: msg,
	}
}
