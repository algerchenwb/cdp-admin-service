package types

import (
	"time"
)

const (
	MANAGE_STATUS_DISABLE = 1 // 禁用
)

type InstanceStatusInfo struct {
	ID             int64     `json:"instance_id"`
	ManageStatus   int       `json:"manage_status"`   // 管理状态
	PowerStatus    int       `json:"power_status"`    // 电源状态 0-off 1-on
	RunningStatus  int       `json:"running_status"`  // 运行状态
	BusinessStatus int       `json:"business_status"` // 业务状态
	BootStatus     int       `json:"boot_status"`     // 启动状态
	BootTime       time.Time `json:"boot_time"`       // 开机时间
	Remark         string    `json:"remark"`
	AssignStatus   int       `json:"assign_status"` // 分配状态
	AssignSource   string    `json:"assign_source"` // 分配来源
	AssignOrder    string    `json:"assign_order"`  // 分配订单
}

// UpdateableInstanceStatusInfo 更新实例状态表
type UpdateableInstanceStatusInfo struct {
	ManageStatus   *int       `json:"manage_status,omitempty"` // 1允许使用，2 暂停使用 3 允许升级  4 允许调试
	PowerStatus    *int       `json:"power_status,omitempty"`
	RunningStatus  *int       `json:"running_status,omitempty"`
	BootStatus     *int       `json:"boot_status,omitempty"`
	BusinessStatus *int       `json:"business_status"`       // 业务状态
	BootTime       *time.Time `json:"boot_time,omitempty"`   // 开机时间
	ActivityIp     string     `json:"activity_ip,omitempty"` // IP
	Remark         string     `json:"remark,omitempty"`
}

// UpdateInstanceStatusRequest 实例状态表
type UpdateInstanceStatusRequest struct {
	FlowID      string  `json:"flow_id"` // 流水ID
	InstanceID  int64   `json:"id"`
	InstanceIDs []int64 `json:"ids"` // 当id为0时才使用ids
	UpdateableInstanceStatusInfo
}

// UpdateInstanceStatusResponse 实例状态表
type UpdateInstanceStatusResponse struct {
	FlowID      string  `json:"flow_id"` // 流水ID
	InstanceID  int64   `json:"instance_id"`
	InstanceIDs []int64 `json:"ids"` // 当id为0时才使用ids
}

// GetInstanceStatusRequest 实例状态
type GetInstanceStatusRequest struct {
	InstanceID int64  `json:"id" form:"id"`
	Mac        string `json:"mac" form:"mac"`
}

// GetInstanceStatusResponse 实例表结构
type GetInstanceStatusResponse InstanceStatusInfo
