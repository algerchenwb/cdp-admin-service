package dbwrap

import (
	"reflect"
	"time"
)

// AreaModel 区域类型
type AreaModel int

// 分发任务状态类型枚举
const (
	CommonArea       AreaModel = 0 // [0]  普通区域
	StorageArea      AreaModel = 1 // [1]  存储区域
	CloudStorageArea AreaModel = 2 // [2]  云端存储区域
)

// DBUnionAreaInfo 统一区域信息
type DBUnionAreaInfo struct {
	UAID       uint64    `json:"uaid"`
	AreaType   int       `json:"area_type"`
	Name       string    `json:"name"`
	State      int       `json:"state"`
	Type       AreaModel `json:"type"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
	ModifyTime time.Time `json:"modify_time"`
}

type UAIWrap struct {
	DBWrap
}

func CreateUAIWrap(host string, callerID int, flowID string) *UAIWrap {
	return &UAIWrap{
		DBWrap{
			_table:    _tUnionAreaInfo,
			_host:     host,
			_callerID: callerID,
			_flowID:   flowID,
			_typ:      reflect.TypeOf(DBUnionAreaInfo{}),
		},
	}
}

func (iDBWrap *UAIWrap) QueryAll(query string, sortby interface{}, ascending interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.QueryAll(query, sortby, ascending)
}
