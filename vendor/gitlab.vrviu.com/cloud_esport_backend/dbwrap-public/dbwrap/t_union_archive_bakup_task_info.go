package dbwrap

import (
	"reflect"
	"time"
)

type ArchiveBakupState int

const (
	BAKUPINIT   ArchiveBakupState = iota
	BAKUPDOING  ArchiveBakupState = 1
	BAKUPSUCC   ArchiveBakupState = 2
	BAKUPFAILED ArchiveBakupState = 3
	BAKUPIGNORE ArchiveBakupState = 4

	BAKUPRETRY ArchiveBakupState = 100
)

type DBUnionArchiveBakupTaskInfo struct {
	UABTID           uint64            `json:"uabtid"`
	UUID             uint64            `json:"uuid"`
	UGID             uint64            `json:"ugid"`
	AreaType         int               `json:"area_type"`
	State            ArchiveBakupState `json:"state"`
	DoneAreaVersion  int               `json:"done_area_version"`
	DoneUnionVersion int               `json:"done_union_version"`
	Size             uint64            `json:"size"`
	BakupAreaType    int               `json:"backup_area_type"`
	DoneTime         time.Time         `json:"done_time"`
	CreateTime       time.Time         `json:"create_time"`
	UpdateTime       time.Time         `json:"update_time"`
	ModifyTime       time.Time         `json:"modify_time"`
}

type UABTIWrap struct {
	DBWrap
}

func CreateUABTIWrap(host string, callerID int, flowID string) *UABTIWrap {
	return &UABTIWrap{
		DBWrap{
			_table:    _tUnionArchiveBakupTaskInfo,
			_host:     host,
			_callerID: callerID,
			_flowID:   flowID,
			_typ:      reflect.TypeOf(DBUnionArchiveBakupTaskInfo{}),
		},
	}
}

func (iDBWrap *UABTIWrap) Query(query string, sortby interface{}, ascending interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.Query(query, sortby, ascending)
}

func (iDBWrap *UABTIWrap) QueryAll(query string, sortby interface{}, ascending interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.QueryAll(query, sortby, ascending)
}

func (iDBWrap *UABTIWrap) Insert(info interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.Insert(info)
}

func (iDBWrap *UABTIWrap) Update(key interface{}, info interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.Update(key, info)
}
