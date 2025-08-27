package dbwrap

import (
	"reflect"
	"time"
)

// DBUnionGameAreaInfo 统一用游戏区域信息
type DBUnionGameAreaInfo struct {
	UGAID      uint64    `json:"ugaid"`
	UGID       uint64    `json:"ugid"`
	AreaType   int       `json:"area_type"`
	GID        uint64    `json:"gid"`
	State      int       `json:"state"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
	ModifyTime time.Time `json:"modify_time"`
}

type UGAIWrap struct {
	DBWrap
}

func CreateUGAIWrap(host string, callerID int, flowID string) *UGAIWrap {
	return &UGAIWrap{
		DBWrap{
			_table:    _tUnionGameAreaInfo,
			_host:     host,
			_callerID: callerID,
			_flowID:   flowID,
			_typ:      reflect.TypeOf(DBUnionGameAreaInfo{}),
		},
	}
}

func (iDBWrap *UGAIWrap) Query(query string, sortby interface{}, ascending interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.Query(query, sortby, ascending)
}

func (iDBWrap *UGAIWrap) QueryAll(query string, sortby interface{}, ascending interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.QueryAll(query, sortby, ascending)
}

func (iDBWrap *UGAIWrap) Insert(info interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.Insert(info)
}
