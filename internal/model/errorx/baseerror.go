package errorx

import (
	"cdp-admin-service/internal/model/globalkey"
	"cdp-admin-service/internal/types"
	"encoding/json"
)

type CodeError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type CodeErrorResponse struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	RequestID string `json:"request_id,omitempty"`
}

func NewCodeError(code int, msg string) error {
	return &CodeError{Code: code, Msg: msg}
}

func NewDefaultError(code int) error {
	return NewCodeError(code, MapErrMsg(code))
}

func NewDefaultCodeError(msg string) error {
	return NewCodeError(ServerErrorCode, msg)
}

func NewHandlerError(code int, msg string) error {
	return NewCodeError(code, msg)
}

func NewSystemError(code int, msg string) error {
	if globalkey.SysShowSystemError {
		return NewCodeError(code, msg)
	} else {
		return NewCodeError(code, MapErrMsg(code))
	}
}

func (e *CodeError) Error() string {
	return e.Msg
}

func (e *CodeError) Data(requestId string) *types.CommonRet {
	return &types.CommonRet{
		Ret: types.Ret{
			Code:      e.Code,
			Msg:       e.Msg,
			RequestId: requestId,
		},
	}
}

func NewDefaultCommomError(code int) types.CommonRet {

	return types.CommonRet{
		Ret: types.Ret{
			Code: code,
			Msg:  MapErrMsg(code),
		},
	}
}

func NewSystemCommomError(msg string) types.CommonRet {

	return types.CommonRet{
		Ret: types.Ret{
			Code: ServerErrorCode,
			Msg:  msg,
		},
	}
}

func UnauthorizedError(sessionId string) string {
	ret := types.CommonRet{
		Ret: types.Ret{
			Code: UnauthorizedErrorCode,
			Msg:  MapErrMsg(UnauthorizedErrorCode),
		},
	}
	b, _ := json.Marshal(ret)
	return string(b)
}

func UnAccessError(sessionId string) string {
	ret := types.CommonRet{
		Ret: types.Ret{
			Code: UnAccessErrorCode,
			Msg:  MapErrMsg(UnAccessErrorCode),
		},
	}
	b, _ := json.Marshal(ret)
	return string(b)
}
