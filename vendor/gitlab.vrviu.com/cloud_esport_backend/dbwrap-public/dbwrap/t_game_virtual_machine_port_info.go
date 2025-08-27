package dbwrap

import (
	"reflect"
	"time"
)

// DBGameVirtualMachinePortInfo 虚拟接端口映射表
type DBGameVirtualMachinePortInfo struct {
	VMPID              uint64    `json:"vmpid"`
	AreaType           int       `json:"area_type"`
	VMID               uint64    `json:"vmid"`
	AAPID              uint64    `json:"aapid"`
	OuterIpv4Address   string    `json:"outer_ipv4_address"`
	OuterIpv4PortGroup int       `json:"outer_ipv4_port_group"`
	CreateTime         time.Time `json:"create_time"`
	UpdateTime         time.Time `json:"update_time"`
	ModifyTime         time.Time `json:"modify_time"`
}

type GVMPIWrap struct {
	AreaDBWrap
}

func CreateGVMPIWrap(host string, callerID int, flowID string) IAreaDBWrap {
	return &GVMPIWrap{
		AreaDBWrap{
			_table:    _tGameVirtualMachinePortInfo,
			_host:     host,
			_callerID: callerID,
			_flowID:   flowID,
			_typ:      reflect.TypeOf(DBGameVirtualMachinePortInfo{}),
		},
	}
}
