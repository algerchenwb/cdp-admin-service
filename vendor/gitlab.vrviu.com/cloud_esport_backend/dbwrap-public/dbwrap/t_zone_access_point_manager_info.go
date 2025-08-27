package dbwrap

import (
	"reflect"
	"time"
)

// DBZoneAccessPointManagerInfo 统一接入点外部信息
type DBZoneAccessPointManagerInfo struct {
	ZAPMID      uint64    `json:"zapmid"`
	ZoneID      int       `json:"zone_id"`
	UserISP     string    `json:"user_isp"`
	AccessISP   string    `json:"access_isp"`
	Country     string    `json:"country"`
	Province    string    `json:"province"`
	City        string    `json:"city"`
	MgrState    int       `json:"mgr_state"`
	Ipv4State   int       `json:"ipv4_state"`
	Ipv6State   int       `json:"ipv6_state"`
	DomainState int       `json:"domain_state"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
	ModifyTime  time.Time `json:"modify_time"`
}

type ZAPMIWrap struct {
	DBWrap
}

func CreateZAPMIWrap(host string, callerID int, flowID string) IDBWrap {
	return &ZAPMIWrap{
		DBWrap{
			_table:    _tZoneAccessPointManagerInfo,
			_host:     host,
			_callerID: callerID,
			_flowID:   flowID,
			_typ:      reflect.TypeOf(DBZoneAccessPointManagerInfo{}),
		},
	}
}

func (iDBWrap *ZAPMIWrap) Query(query string, sortby interface{}, ascending interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.Query(query, sortby, ascending)
}

func (iDBWrap *ZAPMIWrap) QueryAll(query string, sortby interface{}, ascending interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.QueryAll(query, sortby, ascending)
}
