package diskless

// from proto/instance_scheduler/service.pb.go UpdateResourceRequest
// https://confluence.vrviu.com/pages/viewpage.action?pageId=1674489786
type UpdateResourceRequest struct {
	FlowId        string  `json:"flow_id,omitempty"` // 流水ID
	AreaType      int32   `json:"area_type,omitempty"`
	ResourceId    *int64  `json:"resource_id,omitempty"`
	Type          *int32  `json:"type,omitempty"`
	Name          *string `json:"name,omitempty"`
	Specification *int64  `json:"specification,omitempty"`
	Vlan          *int32  `json:"vlan,omitempty"`
	Mode          *int64  `json:"mode,omitempty"`
	Capacity      *int32  `json:"capacity,omitempty"`
	Buffer        *int32  `json:"buffer,omitempty"`
	Init          *int32  `json:"init,omitempty"`
	Concurrent    *int32  `json:"concurrent,omitempty"`
	Priority      *int32  `json:"priority,omitempty"`
	Preemptable   *int32  `json:"preemptable,omitempty"`
	AssignConfig  *string `json:"assign_config,omitempty"`
	Detail        *string `json:"detail,omitempty"`
	State         *int32  `json:"state,omitempty"`
}

type SearchPoolRequest struct {
	FlowId     string   `json:"flow_id,omitempty"` // 流水ID
	AreaType   int32    `json:"area_type,omitempty"`
	ResourceId *int64   `json:"resource_id,omitempty"`
	Conditions []string `json:"conditions,omitempty"`
	Offset     int32    `json:"offset,omitempty"`
	Length     int32    `json:"length,omitempty"`
	Order      string   `json:"order,omitempty"` // asc/desc
	Sortby     string   `json:"sortby,omitempty"`
}

type PoolItem struct {
	AreaType     int32  `json:"area_type,omitempty"`
	ResourceId   string `json:"resource_id,omitempty"`
	InstanceId   string `json:"instance_id,omitempty"`
	Mac          string `json:"mac,omitempty"`
	Address      string `json:"address,omitempty"`
	Flags        string `json:"flags,omitempty"`
	PoolSource   string `json:"pool_source,omitempty"`
	PoolOrder    string `json:"pool_order,omitempty"`
	PoolStatus   int32  `json:"pool_status,omitempty"`
	AssignSource string `json:"assign_source,omitempty"`
	AssignOrder  string `json:"assign_order,omitempty"`
	AssignParam  string `json:"assign_param,omitempty"`
	AssignResult string `json:"assign_result,omitempty"`
	AssignStatus int32  `json:"assign_status,omitempty"`
	CreateTime   string `json:"create_time,omitempty"`
	UpdateTime   string `json:"update_time,omitempty"`
	ModifyTime   string `json:"modify_time,omitempty"`
}

type SearchPoolBody struct {
	Total int32      `json:"total,omitempty"`
	Lists []PoolItem `json:"lists,omitempty"`
}

type SearchPoolResponse struct {
	Head HTTPCommonHead `json:"ret"`
	Body SearchPoolBody `json:"body,omitempty"`
}

// 释放资源
type ReleasePoolItemRequest struct {
	FlowId         string `json:"flow_id,omitempty"` // 流水ID ==> assign_order
	Source         string `json:"source,omitempty"`  // 来源 ==> assign_source
	InstanceId     string `json:"instance_id,omitempty"`
	AreaType       int32  `json:"area_type,omitempty"`
	ResourceConfig string `json:"resource_config,omitempty"`
}

// 重建资源池
type RebuildPoolRequest struct {
	AreaType   string `json:"area_type"`
	ResourceId string `json:"resource_id"`
}

type UpdatePoolItemRequest struct {
	FlowId     string `json:"flow_id,omitempty"` // 流水ID ==> assign_order
	Source     string `json:"source,omitempty"`  // 来源 ==> assign_source
	InstanceId string `json:"instance_id,omitempty"`
	AreaType   int32  `json:"area_type,omitempty"`
	Status     *int32 `json:"status,omitempty"` // 1000-禁止调度(使用完成后需要手动释放回资源池)
	Flags      *int64 `json:"flags,omitempty"`  // 0-取消调试 1-调试(串流结束后不释放回资源池 需要取消调试以后手动释放)
}
