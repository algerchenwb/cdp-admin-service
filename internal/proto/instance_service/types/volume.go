package types

const (
	VOLUME_OPTION_LIVE           = "live"           // 热挂载
	VOLUME_OPTION_NON_PERSISTENT = "non-persistent" // 非持久挂载
)

type AttachVolumeRequest struct {
	FlowId     string   `json:"flow_id"` // 流水ID
	InstanceId int64    `json:"instance_id"`
	Endpoint   string   `json:"endpoint"`
	MountPoint string   `json:"mount_point,omitempty"`
	Options    []string `json:"options,omitempty"`
	Mac        string   `json:"mac"` // 有mac就不使用instanceid
}
type AttachVolumeResponse struct {
	MountPoint string `json:"mount_point"`
}

type DetachVolumeRequest struct {
	FlowId     string   `json:"flow_id"` // 流水ID
	InstanceId int64    `json:"instance_id"`
	Endpoint   string   `json:"endpoint"`
	Options    []string `json:"options,omitempty"`
	Mac        string   `json:"mac"` // 有mac就不使用instanceid
}
type DetachVolumeResponse struct{}

// DynamicUpgradeRequest 指定机器动态挂载卸载数据盘
type DynamicUpgradeRequest struct {
	FlowID      string `json:"flow_id"` // 流水ID
	Mac         string `json:"mac"`     // 有mac就不使用instanceid
	InstanceID  int64  `json:"instance_id"`
	GameVersion int64  `json:"game_version,omitempty"` // 数据盘版本
	MountPoint  string `json:"mount_point,omitempty"`
}

// DynamicUpgradeResponse 指定机器动态挂载卸载数据盘
type DynamicUpgradeResponse struct {
	MountPoint string `json:"mount_point,omitempty"`
}
