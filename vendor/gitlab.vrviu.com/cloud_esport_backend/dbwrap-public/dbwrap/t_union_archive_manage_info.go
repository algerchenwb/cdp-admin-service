package dbwrap

import (
	"reflect"
	"time"
)

// DBUnionArchiveManageInfo 各区域存档库保存存档信息
type DBUnionArchiveManageInfo struct {
	UAMID            uint64    `json:"uamid,omitempty"`
	UUID             uint64    `json:"uuid,omitempty"`
	UGID             uint64    `json:"ugid,omitempty"`
	AreaType         int       `json:"area_type,omitempty"`
	State            int       `json:"state,omitempty"`
	DoneUnionVersion int       `json:"done_union_version,omitempty"`
	DoneAreaVersion  int       `json:"done_area_version,omitempty"`
	Name             string    `json:"name,omitempty"`
	Path             string    `json:"path,omitempty"`
	UID              uint64    `json:"uid,omitempty"`
	GID              uint64    `json:"gid,omitempty"`
	Desc             string    `json:"desc,omitempty"`
	Opeator          string    `json:"opeator,omitempty"`
	CreateTime       time.Time `json:"create_time,omitempty"`
	UpdateTime       time.Time `json:"update_time,omitempty"`
	ModifyTime       time.Time `json:"modify_time,omitempty"`
}

type UAMIWrap struct {
	DBWrap
}

func CreateUAMIWrap(host string, callerID int, flowID string) *UAMIWrap {
	return &UAMIWrap{
		DBWrap{
			_table:    _tUnionArchiveManageInfo,
			_host:     host,
			_callerID: callerID,
			_flowID:   flowID,
			_typ:      reflect.TypeOf(DBUnionArchiveManageInfo{}),
		},
	}
}

func (iDBWrap *UAMIWrap) Query(query string, sortby interface{}, ascending interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.Query(query, sortby, ascending)
}

func (iDBWrap *UAMIWrap) QueryAll(query string, sortby interface{}, ascending interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.QueryAll(query, sortby, ascending)
}

func (iDBWrap *UAMIWrap) Insert(info interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.Insert(info)
}

func (iDBWrap *UAMIWrap) Update(key interface{}, info interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.Update(key, info)
}
