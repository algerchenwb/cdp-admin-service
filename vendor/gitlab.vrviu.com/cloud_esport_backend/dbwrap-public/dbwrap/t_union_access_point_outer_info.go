package dbwrap

import (
	"reflect"
	"time"
)

// AreaPointType 区域接入类型定于
type AreaPointType int

// 区域接入类型枚举
const (
	AreaPointTypeAccess          AreaPointType = iota
	AreaPointTypeInnerArchiver   AreaPointType = 1  // be内部使用
	AreaPointTypeInnerServo      AreaPointType = 2  // be内部使用
	AreaPointTypeStun            AreaPointType = 3  // 已废弃
	AreaPointTypeWebAccess       AreaPointType = 4  // 区域接入（web）
	AreaPointTypeInnerLiveAccess AreaPointType = 5  // be内部使用
	AreaPointTypeControl         AreaPointType = 6  // 手柄控制文件
	AreaPointTypeWebControl      AreaPointType = 7  // 手柄控制文件（web）
	AreaPointTypeInnerAccess     AreaPointType = 8  // be内部使用，区域接入
	AreaPointTypeConfigment      AreaPointType = 9  // be内部使用，区域配置代理地址
	AreaPointTypeInnerStorage    AreaPointType = 10 // be内部使用，区域存储调度
)

// DBUnionAccessPointOuterInfo 统一接入点外部信息
type DBUnionAccessPointOuterInfo struct {
	UAPOID           uint64        `json:"uapoid"`
	UAPID            uint64        `json:"uapid"`
	ZoneID           int           `json:"zone_id"`
	AreaType         int           `json:"area_type"`
	MgrState         int           `json:"mgr_state"`
	ISP              string        `json:"isp"`
	Country          string        `json:"country"`
	Province         string        `json:"province"`
	City             string        `json:"city"`
	OuterDomain      string        `json:"outer_domain"`
	OuterIpv4Address string        `json:"outer_ipv4_address"`
	OuterIpv6Address string        `json:"outer_ipv6_address"`
	OuterPort        uint16        `json:"outer_port"`
	Weight           int           `json:"weight"`
	Type             AreaPointType `json:"type"`
	CreateTime       time.Time     `json:"create_time"`
	UpdateTime       time.Time     `json:"update_time"`
	ModifyTime       time.Time     `json:"modify_time"`
}

type UAPOIWrap struct {
	DBWrap
}

func CreateUAPOIWrap(host string, callerID int, flowID string) *UAPOIWrap {
	return &UAPOIWrap{
		DBWrap{
			_table:    _tUnionAccessPointOuterInfo,
			_host:     host,
			_callerID: callerID,
			_flowID:   flowID,
			_typ:      reflect.TypeOf(DBUnionAccessPointOuterInfo{}),
		},
	}
}

func (iDBWrap *UAPOIWrap) Query(query string, sortby interface{}, ascending interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.Query(query, sortby, ascending)
}

func (iDBWrap *UAPOIWrap) QueryAll(query string, sortby interface{}, ascending interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.QueryAll(query, sortby, ascending)
}

func (iDBWrap *UAPOIWrap) Insert(info interface{}) (interface{}, int, error) {
	return iDBWrap.DBWrap.Insert(info)
}
