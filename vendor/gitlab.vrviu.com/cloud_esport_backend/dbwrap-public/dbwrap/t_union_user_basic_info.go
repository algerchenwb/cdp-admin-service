package dbwrap

import (
	"reflect"
	"time"
)

// DBUnionUserBasicInfo 统一用户基础信息
type DBUnionUserBasicInfo struct {
	UUID          uint64    `json:"uuid,omitempty"`
	Name          string    `json:"name,omitempty"`
	ThirdID       string    `json:"third_id,omitempty"`
	UseType       int       `json:"user_type,omitempty"`
	State         int       `json:"state,omitempty"`
	ActiveVersion int       `json:"active_version,omitempty"`
	CreateTime    time.Time `json:"create_time,omitempty"`
	UpdateTime    time.Time `json:"update_time,omitempty"`
	ModifyTime    time.Time `json:"modify_time,omitempty"`
}

type UUBIWrap struct {
	DBWrap
}

func CreateUUBIWrap(host string, callerID int, flowID string) *UUBIWrap {
	return &UUBIWrap{
		DBWrap{
			_table:    _tUnionUserBasicInfo,
			_host:     host,
			_callerID: callerID,
			_flowID:   flowID,
			_typ:      reflect.TypeOf(DBUnionUserBasicInfo{}),
		},
	}
}

func (iDBWrap *UUBIWrap) Query(query string, sortby interface{}, ascending interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.Query(query, sortby, ascending)
}

func (iDBWrap *UUBIWrap) Insert(info interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.Insert(info)
}
