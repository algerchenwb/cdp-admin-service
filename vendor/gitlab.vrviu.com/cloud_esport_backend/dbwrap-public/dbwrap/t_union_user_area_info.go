package dbwrap

import (
	"reflect"
	"time"
)

// DBUnionUserAreaInfo 统一用用户区域信息
type DBUnionUserAreaInfo struct {
	UUAID      uint64    `json:"uuaid"`
	UUID       uint64    `json:"uuid"`
	AreaType   int       `json:"area_type"`
	UID        uint64    `json:"uid"`
	State      int       `json:"state"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
	ModifyTime time.Time `json:"modify_time"`
}

type UUAIWrap struct {
	DBWrap
}

func CreateUUAIWrap(host string, callerID int, flowID string) *UUAIWrap {
	return &UUAIWrap{
		DBWrap{
			_table:    _tUnionUserAreaInfo,
			_host:     host,
			_callerID: callerID,
			_flowID:   flowID,
			_typ:      reflect.TypeOf(DBUnionUserAreaInfo{}),
		},
	}
}

func (iDBWrap *UUAIWrap) Query(query string, sortby interface{}, ascending interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.Query(query, sortby, ascending)
}

func (iDBWrap *UUAIWrap) Insert(info interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.Insert(info)
}
