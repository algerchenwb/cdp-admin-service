package types

const (
	POWER_ON     = "on"
	POWER_OFF    = "off"
	POWER_FOFF   = "foff"
	POWER_REBOOT = "reboot"
	POWER_STATUS = "status"
)

const (
	POWER_STATUS_OFF = 0
	POWER_STATUS_ON  = 1
)

type PowerIpmiInfo struct {
	Type    int    `json:"type,omitempty"`
	Slot    int    `json:"slot,omitempty"`
	Address string `json:"address,omitempty"`
	Param   string `json:"param,omitempty"`
	Mac     string `json:"mac,omitempty"`
}

type PowerCtrlRequest struct {
	FlowId     string         `json:"flow_id"` // 流水ID
	InstanceId int64          `json:"instance_id"`
	Mac        string         `json:"mac" form:"mac"`
	HostId     int64          `json:"host_id,omitempty"`
	Operation  string         `json:"operation"` // on/off/reboot
	IsAsync    bool           `json:"is_async,omitempty"`
	IpmiInfo   *PowerIpmiInfo `json:"ipmi_info,omitempty"` // 选择ipmi类型
}

type PowerCtrlResponse struct {
	FlowId string `json:"flow_id"` // 流水ID
	Status int    `json:"status"`  // 电源状态 0-off 1-on
}

type GetBootActionRequest struct {
	FlowId      string `json:"flow_id" form:"flow_id"` // 流水ID
	Mac         string `json:"mac" form:"mac"`
	BootSession string `json:"bootsession" form:"bootsession"`
	Ip          string `json:"ip" form:"ip"`
	Phase       int    `json:"phase" form:"phase"` // 阶段 0-等待 1-恢复启动
}

const (
	BOOT_ACTION_WAIT   = "ctn" // 等待
	BOOT_ACTION_RESUME = "brk" // 恢复启动
)
