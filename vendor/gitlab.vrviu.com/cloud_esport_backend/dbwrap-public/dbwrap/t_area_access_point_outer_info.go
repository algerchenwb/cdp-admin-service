package dbwrap

import (
	"reflect"
	"time"
)

// DBAreaAccessPointOuterInfo 区域接入点对外信息
type DBAreaAccessPointOuterInfo struct {
	AAPOID           uint64    `json:"aapoid"`
	AAPID            uint64    `json:"aapid"`
	BizType          int       `json:"biz_type"`
	AreaType         int       `json:"area_type"`
	MgrState         int       `json:"mgr_state"`
	ISP              string    `json:"isp"`
	Country          string    `json:"country"`
	Province         string    `json:"province"`
	City             string    `json:"city"`
	OuterIpv4Address string    `json:"outer_ipv4_address"`
	OuterIpv6Address string    `json:"outer_ipv6_address"`
	OuterDomain      string    `json:"outer_domain"`
	CreateTime       time.Time `json:"create_time"`
	UpdateTime       time.Time `json:"update_time"`
	ModifyTime       time.Time `json:"modify_time"`
}

type AAPOIWrap struct {
	AreaDBWrap
}

func CreateAAPOIWrap(host string, callerID int, flowID string) IAreaDBWrap {
	return &AAPOIWrap{
		AreaDBWrap{
			_table:    _tAreaAccessPointOuterInfo,
			_host:     host,
			_callerID: callerID,
			_flowID:   flowID,
			_typ:      reflect.TypeOf(DBAreaAccessPointOuterInfo{}),
		},
	}
}
