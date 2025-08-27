package dbwrap

import (
	"reflect"
	"time"
)

// DBUnionUserThirdInfo 统一用户第三方信息
type DBUnionUserThirdInfo struct {
	UUTID      uint64    `json:"uutid"`
	UUID       uint64    `json:"uuid"`
	ThirdID    string    `json:"third_id"`
	ThirdType  int       `json:"third_type"`
	BizType    int       `json:"biz_type"`
	State      int       `json:"state"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
	ModifyTime time.Time `json:"modify_time"`
}

type UUTIWrap struct {
	DBWrap
}

func CreateUUTIWrap(host string, callerID int, flowID string) *UUTIWrap {
	return &UUTIWrap{
		DBWrap{
			_table:    _tUnionUserThirdInfo,
			_host:     host,
			_callerID: callerID,
			_flowID:   flowID,
			_typ:      reflect.TypeOf(DBUnionUserThirdInfo{}),
		},
	}
}

func (iDBWrap *UUTIWrap) Query(query string, sortby interface{}, ascending interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.Query(query, sortby, ascending)
}

func (iDBWrap *UUTIWrap) Insert(info interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.Insert(info)
}
