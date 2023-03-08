package request

import "time"

type Response struct {
	ErrorCode    string      `json:"error_code"`
	ErrorMessage string      `json:"error_message"`
	ServerTime   int64       `json:"server_time"`
	Count        int         `json:"count,omitempty"`
	Data         interface{} `json:"data,omitempty"`
}

func NewResponse() *Response {
	resp := new(Response)
	resp.ServerTime = time.Now().Unix()
	return resp
}

func (resp *Response) Code(code string) *Response {
	resp.ErrorCode = code
	return resp
}

func (resp *Response) Msg(msg string) *Response {
	resp.ErrorMessage = msg
	return resp
}
