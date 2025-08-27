package dbwrap

import (
	"reflect"
	"time"
)

// DBUnionAccessPointGeographyScheduleInfo -
type DBUnionAccessPointGeographyScheduleInfo struct {
	UAPGSID    uint64    `json:"uapgsid"`
	BizType    int       `json:"biz_type"`
	ISP        string    `json:"isp"`
	Country    string    `json:"country"`
	Province   string    `json:"province"`
	City       string    `json:"city"`
	UAPID      int       `json:"uapid"`
	State      int       `json:"state"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
	ModifyTime time.Time `json:"modify_time"`
}

type UAPGSIWrap struct {
	DBWrap
}

func CreateUAPGSIWrap(host string, callerID int, flowID string) *UAPGSIWrap {
	return &UAPGSIWrap{
		DBWrap{
			_table:    _tUnionAccessPointGeographyScheduleInfo,
			_host:     host,
			_callerID: callerID,
			_flowID:   flowID,
			_typ:      reflect.TypeOf(DBUnionAccessPointGeographyScheduleInfo{}),
		},
	}
}

func (iDBWrap *UAPGSIWrap) Query(query string, sortby interface{}, ascending interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.Query(query, sortby, ascending)
}
