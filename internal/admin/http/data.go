package http

const (
	OK = 20000
)

type Response struct {
	Code int `json:"code"`
	Data any `json:"data"`
}

func Success(data any) *Response {
	return &Response{
		Code: OK,
		Data: data,
	}
}

func Error(code int, msg string) *Response {
	return &Response{
		Code: code,
		Data: map[string]string{"error": msg},
	}
}
