package dbwrap

import (
	"reflect"
	"time"
)

// ScheduleState 分发任务状态类型
type ScheduleState int

// 分发任务状态类型枚举
const (
	_                   ScheduleState = iota
	GAMEPREPARE         ScheduleState = 1  // [1]  申请开始游戏
	GAMEDONE            ScheduleState = 2  // [2]  游戏完成,已产生游戏归档
	GAMESKIP            ScheduleState = 3  // [3]  跳过游戏,启动游戏失败
	DISPATCHPREPARE     ScheduleState = 4  // [4]  分发中...
	DISPATCHDONE        ScheduleState = 5  // [5]  分发已完成
	DISPATCHFAILED      ScheduleState = 6  // [6]  分发失败
	IMPORTDONE          ScheduleState = 7  // [7]  导入完成
	IMPORTFAILED        ScheduleState = 8  // [8]  导入失败
	SYNCDISPATCHPREPARE ScheduleState = 9  // [9]  实时迁移分发中...
	SYNCDISPATCHDONE    ScheduleState = 10 // [10] 实时迁移分发已完成
	SYNCDISPATCHFAILED  ScheduleState = 11 // [11] 实时迁移分发失败
)

// ScheduleDoneType 调度完成类型定义
type ScheduleDoneType int

// 调度完成类型枚举
const (
	_          ScheduleDoneType = iota
	DTGame     ScheduleDoneType = 1 // [1]  游戏类型
	DTDISPATCH ScheduleDoneType = 2 // [2]  分发类型
	DTIMPORT   ScheduleDoneType = 3 // [3]  导入类型
)

// DBUnionUserGameScheduleInfo 同一用户游戏区域调度信息
type DBUnionUserGameScheduleInfo struct {
	UUGAID           string           `json:"uugaid,omitempty"`
	UUID             uint64           `json:"uuid,omitempty"`
	UGID             uint64           `json:"ugid,omitempty"`
	AreaType         int              `json:"area_type,omitempty"`
	OpVersion        int              `json:"op_version"`
	OpState          ScheduleState    `json:"op_state"`
	DoneUnionVersion int              `json:"done_union_version"`
	DoneAreaVersion  int              `json:"done_area_version"`
	DoneType         ScheduleDoneType `json:"done_type"`
	DoneTime         time.Time        `json:"done_time"`
	CreateTime       time.Time        `json:"create_time"`
	UpdateTime       time.Time        `json:"update_time"`
	ModifyTime       time.Time        `json:"modify_time"`
}

type UUGASIWrap struct {
	DBWrap
}

func CreateUUGASIWrap(host string, callerID int, flowID string) *UUGASIWrap {
	return &UUGASIWrap{
		DBWrap{
			_table:    _tUnionUserGameAreaStorageInfo,
			_host:     host,
			_callerID: callerID,
			_flowID:   flowID,
			_typ:      reflect.TypeOf(DBUnionUserGameScheduleInfo{}),
		},
	}
}

func (iDBWrap *UUGASIWrap) Query(query string, sortby interface{}, ascending interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.Query(query, sortby, ascending)
}

func (iDBWrap *UUGASIWrap) QueryAll(query string, sortby interface{}, ascending interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.QueryAll(query, sortby, ascending)
}

func (iDBWrap *UUGASIWrap) Insert(info interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.Insert(info)
}

func (iDBWrap *UUGASIWrap) Update(key, info interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.Update(key, info)
}
