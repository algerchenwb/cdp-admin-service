package types

type HTTPCommonHead struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg,omitempty"`
	Detail string `json:"detail,omitempty"`
}

type HTTPResponse struct {
	Head HTTPCommonHead `json:"ret"`
	Body interface{}    `json:"body,omitempty"`
}

func NewResponse(header *HTTPCommonHead, body interface{}) *HTTPResponse {
	ret := HTTPCommonHead{
		Code:   header.Code,
		Detail: header.Msg,
		Msg:    wrapErrorMsg(header),
	}
	return &HTTPResponse{Head: ret, Body: body}
}
