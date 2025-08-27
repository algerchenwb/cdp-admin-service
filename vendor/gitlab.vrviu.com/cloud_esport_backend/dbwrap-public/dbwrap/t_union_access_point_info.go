package dbwrap

import (
	"reflect"
	"time"
)

// DBUnionAccessPointInfo 统一接入点外部信息
type DBUnionAccessPointInfo struct {
	UAPID       int       `json:"uapid"`
	AreaType    int       `json:"area_type"`
	ZoneID      int       `json:"zone_id"`
	OnlineState int       `json:"online_state"`
	MgrState    int       `json:"mgr_state"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
	ModifyTime  time.Time `json:"modify_time"`
}

type UAPIWrap struct {
	DBWrap
}

func CreateUAPIWrap(host string, callerID int, flowID string) IDBWrap {
	return &UAPIWrap{
		DBWrap{
			_table:    _tUnionAccessPointInfo,
			_host:     host,
			_callerID: callerID,
			_flowID:   flowID,
			_typ:      reflect.TypeOf(DBUnionAccessPointInfo{}),
		},
	}
}

func (iDBWrap *UAPIWrap) Query(query string, sortby interface{}, ascending interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.Query(query, sortby, ascending)
}

func (iDBWrap *UAPIWrap) QueryAll(query string, sortby interface{}, ascending interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.QueryAll(query, sortby, ascending)
}

func (iDBWrap *UAPIWrap) Insert(info interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.Insert(info)
}
