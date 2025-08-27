package dbwrap

import (
	"reflect"
	"time"
)

// DBUnionGameMetaInfo 统一游戏信息
type DBUnionGameMetaInfo struct {
	UGID          uint64    `json:"ugid"`
	Name          string    `json:"name"`
	EnName        string    `json:"en_name"`
	Developer     string    `json:"developer"`
	Publisher     string    `json:"publisher"`
	ReleaseDate   time.Time `json:"release_date"`
	IssueDate     time.Time `json:"issue_date"`
	Description   string    `json:"description"`
	Cover         string    `json:"cover"`
	Categories    string    `json:"categories"`
	Tags          string    `json:"tags"`
	Controllers   string    `json:"controllers"`
	State         int       `json:"state"`
	ActiveVersion int       `json:"active_version"`
	ExeName       string    `json:"exe_name"`
	PlayConfig    string    `json:"play_config"`
	VMTypeList    string    `json:"vm_type_list"`
	GameConfig    string    `json:"game_config"`
	CreateTime    time.Time `json:"create_time"`
	UpdateTime    time.Time `json:"update_time"`
	ModifyTime    time.Time `json:"modify_time"`
}

type UGMIWrap struct {
	DBWrap
}

func CreateUGMIWrap(host string, callerID int, flowID string) *UGMIWrap {
	return &UGMIWrap{
		DBWrap{
			_table:    _tUnionGameMetaInfo,
			_host:     host,
			_callerID: callerID,
			_flowID:   flowID,
			_typ:      reflect.TypeOf(DBUnionGameMetaInfo{}),
		},
	}
}

func (iDBWrap *UGMIWrap) Query(query string, sortby interface{}, ascending interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.Query(query, sortby, ascending)
}

func (iDBWrap *UGMIWrap) Insert(info interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.Insert(info)
}
