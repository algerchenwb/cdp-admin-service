package errwrap

import "fmt"

func New(code int, msg string) *ErrorWrap {
	return &ErrorWrap{
		Code:    code,
		Message: msg,
	}
}

type ErrorWrap struct {
	Code    int
	Message string
}

func (e *ErrorWrap) Error() string {
	if e.Code == 0 {
		return e.Message
	}
	return fmt.Sprintf("errcode(%d) errmsg:%s", e.Code, e.Message)
}

func (e *ErrorWrap) Set(msg string) *ErrorWrap {
	return &ErrorWrap{Code: e.Code, Message: msg}
}
