package diskless

type CreateImageFromAreaInstanceRequest struct {
	FlowId       string  `json:"flow_id,omitempty"`
	ImageId      string  `json:"image_id,omitempty"`      // 全局唯一ID(31020)
	Name         string  `json:"name,omitempty"`          // 镜像名：云电竞通用镜像win10
	OsVersion    string  `json:"os_version,omitempty"`    // 内置os版本描述：win10LTSC-64
	Remark       string  `json:"remark,omitempty"`        // 备注信息
	ManagerState int32   `json:"manager_state,omitempty"` // 1-启动， 2-停用
	AreaId       int64   `json:"area_id,omitempty"`       // 生成镜像的机房
	VmId         int64   `json:"vm_id,omitempty"`         // 生成镜像的实例
	OsVolumeId   float64 `json:"os_volume_id,omitempty"`  // 生成镜像的卷
	FlattenFlag  int32   `json:"flatten_flag,omitempty"`  // 是否拍平 0：要拍平，1：不拍平
}
