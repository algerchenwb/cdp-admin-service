package types

const (
	E_OK                   = 0
	E_NO_PARTITION         = -5
	E_MOUNTPOINT_NOT_MATCH = -6
	E_PARAM                = -14120001
	E_SYS                  = -14120002
	E_BUSY                 = -14120003
	E_CONFLICT             = -14120004
	E_NOT_EXISTS           = -14120005
)

var (
	EM_MSG = map[int]string{
		E_PARAM:                "参数错误",
		E_SYS:                  "系统错误, 请稍后重试",
		E_BUSY:                 "系统繁忙, 请稍后重试",
		E_CONFLICT:             "执行冲突, 请稍后重试",
		E_NOT_EXISTS:           "记录不存在",
		E_NO_PARTITION:         "磁盘分区不存在",
		E_MOUNTPOINT_NOT_MATCH: "盘符不匹配",
	}
)

func wrapErrorMsg(header *HTTPCommonHead) string {
	if prefix, ok := EM_MSG[header.Code]; ok && header.Code != E_OK {
		return prefix
	}
	return header.Msg
}
