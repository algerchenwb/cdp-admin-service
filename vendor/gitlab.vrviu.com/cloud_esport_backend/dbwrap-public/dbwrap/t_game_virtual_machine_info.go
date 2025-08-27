package dbwrap

import (
	"reflect"
	"time"
)

// DBGameVirtualMachineInfo 虚拟接端口映射表
type DBGameVirtualMachineInfo struct {
	VMID              uint64      `json:"vmid,omitempty"`
	AreaType          int         `json:"area_type,omitempty"`
	SubAreaType       int         `json:"sub_area_type,omitempty"`
	InnerMacAddress   string      `json:"inner_mac_address,omitempty"`
	OuterMacAddress   string      `json:"outer_mac_address,omitempty"`
	HMID              uint64      `json:"hmid,omitempty"`
	VMType            int         `json:"vm_type,omitempty"`
	Name              string      `json:"name,omitempty"`
	HostDeviceSet     interface{} `json:"host_device_set,omitempty"`
	OnlineState       int         `json:"online_state,omitempty"`
	MgrState          int         `json:"mgr_state,omitempty"`
	OsState           int         `json:"os_state,omitempty"`
	AssignState       int         `json:"assign_state,omitempty"`
	AssignUID         uint64      `json:"assign_uid,omitempty"`
	AssignGID         uint64      `json:"assign_gid,omitempty"`
	Vmoiid            int         `json:"vmoiid,omitempty"`
	ImageVersion      int         `json:"image_version,omitempty"`
	TotalImageSize    int         `json:"total_image_size,omitempty"`
	MgrIpv4Address    string      `json:"mgr_ipv4_address,omitempty"`
	StreamIpv4Address string      `json:"stream_ipv4_address,omitempty"`
	StreamIpv4Port    int         `json:"stream_ipv4_port,omitempty"`
	StreamMethod      int         `json:"stream_method,omitempty"`
	GfeVersion        int         `json:"gfe_version,omitempty"`
	LsVersion         int         `json:"ls_version,omitempty"`
	CreateTime        time.Time   `json:"create_time,omitempty"`
	UpdateTime        time.Time   `json:"update_time,omitempty"`
	ModifyTime        time.Time   `json:"modify_time,omitempty"`
}

type GVMIWrap struct {
	AreaDBWrap
}

func CreateGVMIWrap(host string, callerID int, flowID string) IAreaDBWrap {
	return &GVMIWrap{
		AreaDBWrap{
			_table:    _tGameVirtualMachineInfo,
			_host:     host,
			_callerID: callerID,
			_flowID:   flowID,
			_typ:      reflect.TypeOf(DBGameVirtualMachineInfo{}),
		},
	}
}
