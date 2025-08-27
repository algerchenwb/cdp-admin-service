package gopublic

import (
	"errors"
)

// 错误码error定义
var (
	ErrOK           = errors.New("ok")
	ErrNotExist     = errors.New("item not exist")
	ErrMarshal      = errors.New("marshal json fail")
	ErrUnmarshal    = errors.New("unmarshal json fail")
	ErrAlreadyExist = errors.New("item already exist")
	ErrTimeout      = errors.New("timeout")
)
