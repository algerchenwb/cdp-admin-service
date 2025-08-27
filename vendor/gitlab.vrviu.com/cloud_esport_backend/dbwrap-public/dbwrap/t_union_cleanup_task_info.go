package dbwrap

import (
	"reflect"
	"time"
)

// DBUnionCleanupTaskInfo 清理旧存档任务信息
type DBUnionCleanupTaskInfo struct {
	UCTID               uint64      `json:"uctid"`
	UUID                uint64      `json:"uuid"`
	UGID                uint64      `json:"ugid"`
	Identity            int         `json:"identity"`
	SrcAreaType         int         `json:"src_area_type"`
	DoneUnionVersion    int         `json:"done_union_version"`
	CurDoneUnionVersion int         `json:"cur_done_union_version"`
	CleanupState        RemoveState `json:"cleanup_state"`
	CleanupDoneTime     time.Time   `json:"cleanup_done_time"`
	CreateTime          time.Time   `json:"create_time"`
	UpdateTime          time.Time   `json:"update_time"`
	ModifyTime          time.Time   `json:"modify_time"`
}

type UCTIWrap struct {
	DBWrap
}

func CreateUCTIWrap(host string, callerID int, flowID string) *UCTIWrap {
	return &UCTIWrap{
		DBWrap{
			_table:    _tUnionCleanupTaskInfo,
			_host:     host,
			_callerID: callerID,
			_flowID:   flowID,
			_typ:      reflect.TypeOf(DBUnionCleanupTaskInfo{}),
		},
	}
}

func (iDBWrap *UCTIWrap) Insert(info interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.Insert(info)
}
