package dbwrap

import (
	"reflect"
	"time"
)

// DBUnionAreaZoneInfo 统一大区信息
type DBUnionAreaZoneInfo struct {
	UAZID      uint64    `json:"uazid"`
	ZoneID     int       `json:"zoneid"`
	AreaType   int       `json:"area_type"`
	Name       string    `json:"name"`
	State      int       `json:"state"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
	ModifyTime time.Time `json:"modify_time"`
}

type UAZIWrap struct {
	DBWrap
}

func CreateUAZIWrap(host string, callerID int, flowID string) *UAZIWrap {
	return &UAZIWrap{
		DBWrap{
			_table:    _tUnionAreaZoneInfo,
			_host:     host,
			_callerID: callerID,
			_flowID:   flowID,
			_typ:      reflect.TypeOf(DBUnionAreaZoneInfo{}),
		},
	}
}

func (iDBWrap *UAZIWrap) Query(query string, sortby interface{}, ascending interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.Query(query, sortby, ascending)
}
