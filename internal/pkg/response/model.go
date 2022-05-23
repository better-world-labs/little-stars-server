package response

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func NewResponse(code int, msg string, data interface{}) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

func NewResponseOk(data interface{}) *Response {
	return NewResponse(StatusOK, "", data)
}

func NewResponseError(code int, msg string, data interface{}) *Response {
	return NewResponse(code, msg, data)
}
