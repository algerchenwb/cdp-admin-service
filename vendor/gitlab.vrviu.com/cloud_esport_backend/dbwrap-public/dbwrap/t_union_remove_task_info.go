package dbwrap

import (
	"reflect"
	"time"
)

type RemoveState int

// DispatchState\CleanupState
const (
	StateBegin          RemoveState = iota
	StateDoing          RemoveState = 1
	StateSucc           RemoveState = 2
	StateFailed         RemoveState = 3
	StateIgnore         RemoveState = 4
	StateWaiting        RemoveState = 5
	StateUnable         RemoveState = 6
	StateNetErr         RemoveState = 7
	StateAreaVersionErr RemoveState = 8
	StateArchiverLarge  RemoveState = 9
)

// DBUnionRemoveTaskInfo 分发任务
type DBUnionRemoveTaskInfo struct {
	URTID            uint64      `json:"urtid"`
	UUID             uint64      `json:"uuid"`
	UGID             uint64      `json:"ugid"`
	Identity         int         `json:"identity"`
	SrcAreaType      int         `json:"src_area_type"`      // 源区域
	DstAreaType      int         `json:"dst_area_type"`      // 目的区域
	DoneUnionVersion int         `json:"done_union_version"` // 存档版本
	DispatchState    RemoveState `json:"dispatch_state"`
	CleanupState     RemoveState `json:"cleanup_state"`
	DispatchDoneTime time.Time   `json:"dispatch_done_time"`
	CleanupDoneTime  time.Time   `json:"cleanup_done_time"`
	CreateTime       time.Time   `json:"create_time"`
	ModifyTime       time.Time   `json:"modify_time"`
	UpdateTime       time.Time   `json:"update_time"`
}

type URTIWrap struct {
	DBWrap
}

func CreateURTIWrap(host string, callerID int, flowID string) *URTIWrap {
	return &URTIWrap{
		DBWrap{
			_table:    _tUnionRemoveTaskInfo,
			_host:     host,
			_callerID: callerID,
			_flowID:   flowID,
			_typ:      reflect.TypeOf(DBUnionRemoveTaskInfo{}),
		},
	}
}

func (iDBWrap *URTIWrap) Query(query string, sortby interface{}, ascending interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.Query(query, sortby, ascending)
}

func (iDBWrap *URTIWrap) QueryAll(query string, sortby interface{}, ascending interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.QueryAll(query, sortby, ascending)
}

func (iDBWrap *URTIWrap) QueryPage(query string, offset, limit int, sortby, ascending interface{}) (int, interface{}, int, error) {
	return iDBWrap.DBWrap.QueryPage(query, offset, limit, sortby, ascending)
}

func (iDBWrap *URTIWrap) Insert(info interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.Insert(info)
}

func (iDBWrap *URTIWrap) Update(key interface{}, info interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.Update(key, info)
}
