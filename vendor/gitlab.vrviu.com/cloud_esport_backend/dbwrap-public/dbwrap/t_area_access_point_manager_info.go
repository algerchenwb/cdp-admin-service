package dbwrap

import (
	"reflect"
	"time"
)

// DBAreaAccessPointManagerInfo 统一接入点外部信息
type DBAreaAccessPointManagerInfo struct {
	AAPMID      uint64    `json:"aapmid"`
	AreaType    int       `json:"area_type"`
	UserISP     string    `json:"user_isp"`
	AccessISP   string    `json:"access_isp"`
	MgrState    int       `json:"mgr_state"`
	Ipv4State   int       `json:"ipv4_state"`
	Ipv6State   int       `json:"ipv6_state"`
	DomainState int       `json:"domain_state"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
	ModifyTime  time.Time `json:"modify_time"`
}

type AAPMIWrap struct {
	DBWrap
}

func CreateAAPMIWrap(host string, callerID int, flowID string) *AAPMIWrap {
	return &AAPMIWrap{
		DBWrap{
			_table:    _tAreaAccessPointManagerInfo,
			_host:     host,
			_callerID: callerID,
			_flowID:   flowID,
			_typ:      reflect.TypeOf(DBAreaAccessPointManagerInfo{}),
		},
	}
}

func (iDBWrap *AAPMIWrap) Query(query string, sortby interface{}, ascending interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.Query(query, sortby, ascending)
}

func (iDBWrap *AAPMIWrap) QueryAll(query string, sortby interface{}, ascending interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.QueryAll(query, sortby, ascending)
}
