package dbwrap

import (
	"reflect"
	"time"
)

type DBUnionArchiveImportInfo struct {
	UAIID            uint64    `json:"uaiid"`
	UUID             uint64    `json:"uuid"`
	UGID             uint64    `json:"ugid"`
	SrcAreaType      int       `json:"src_area_type"`
	DstAreaType      int       `json:"dst_area_type"`
	State            int       `json:"state"`
	Name             string    `json:"name"`
	UID              uint64    `json:"uid"`
	GID              uint64    `json:"gid"`
	DoneAreaVersion  int       `json:"done_area_version"`
	DoneUnionVersion int       `json:"done_union_version"`
	DoneTime         time.Time `json:"done_time"`
	Opeator          string    `json:"opeator"`
	Desc             string    `json:"desc"`
	CreateTime       time.Time `json:"create_time"`
	UpdateTime       time.Time `json:"update_time"`
	ModifyTime       time.Time `json:"modify_time"`
}

type UAIIWrap struct {
	DBWrap
}

func CreateUAIIWrap(host string, callerID int, flowID string) *UAIIWrap {
	return &UAIIWrap{
		DBWrap{
			_table:    _tUnionArchiveImportInfo,
			_host:     host,
			_callerID: callerID,
			_flowID:   flowID,
			_typ:      reflect.TypeOf(DBUnionArchiveImportInfo{}),
		},
	}
}

func (iDBWrap *UAIIWrap) Query(query string, sortby interface{}, ascending interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.Query(query, sortby, ascending)
}

func (iDBWrap *UAIIWrap) Insert(info interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.Insert(info)
}

func (iDBWrap *UAIIWrap) Update(key interface{}, info interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.Update(key, info)
}
