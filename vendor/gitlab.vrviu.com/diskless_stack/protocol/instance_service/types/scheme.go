package types

import (
	"time"
)

type ListSchemesRequest struct {
	Offset    int    `json:"offset" form:"offset"`   // 偏移
	Length    int    `json:"length" form:"length"`   // 长度 0-所有
	SortBy    string `json:"sort_by" form:"sort_by"` // 排序字段 id/create_time/modify_time
	Order     string `json:"order" form:"order"`     // asc or desc
	SchemeId  int64  `json:"id" form:"id"`
	Name      string `json:"name" form:"name"`
	OsImageId string `json:"os_image_id" form:"os_image_id"`
}

type DataImage struct {
	Id      int64  `json:"id"`
	Type    int    `json:"type"`
	Path    string `json:"path"`
	Version int    `json:"version"`
	UUID    string `json:"uuid"`
	Remark  string `json:"remark"`
	State   int    `json:"state"`
}

type InstanceScheme struct {
	Id          int64        `json:"id"`
	Name        string       `json:"name"`
	OsImageId   string       `json:"os_image_id"`
	DataImages  []*DataImage `json:"data_image,omitempty"`
	StorageType int          `json:"storage_type,omitempty"` // 0-默认iscsi; 1-游戏盘使用自研存储
	State       int          `json:"state"`                  // 0-初始 100-有效 900-失效
	CreateTime  time.Time    `json:"create_time"`
	UpdateTime  time.Time    `json:"update_time"`
	ModifyTime  time.Time    `json:"modify_time"`
}

type UpdateableInstanceScheme struct {
	Name         string `json:"name,omitempty"`
	OsImageId    string `json:"os_image_id,omitempty"`
	DataImageIds string `json:"data_image_ids,omitempty"`
	StorageType  *int   `json:"storage_type,omitempty"`
}

type ListSchemesResponse struct {
	Schemes []InstanceScheme `json:"schemes"` // 实例编排列表
	Total   int              `json:"total"`
}

type GetSchemeRequest struct {
	SchemeId   int64 `json:"id" form:"id"`
	InstanceId int64 `json:"instance_id" form:"instance_id"`
}

type GetSchemeResponse struct {
	Scheme InstanceScheme `json:"scheme"`
}

type CreateSchemeRequest struct {
	FlowId string `json:"flow_id"` // 流水ID
	UpdateableInstanceScheme
}

type CreateSchemeResponse struct {
	SchemeId int64 `json:"scheme_id"`
}

type UpdateSchemeRequest struct {
	FlowId   string `json:"flow_id"` // 流水ID
	SchemeId int64  `json:"scheme_id"`
	UpdateableInstanceScheme
}

type UpdateSchemeResponse struct {
}

type DeleteSchemeRequest struct {
	FlowId   string `json:"flow_id"` // 流水ID
	SchemeId int64  `json:"scheme_id"`
}

type DeleteSchemeResponse struct {
}
