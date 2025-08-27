package types

import (
	"fmt"
	"mime/multipart"
	"strings"
	"time"
)

const (
	DEFAULT_PAGING_LENGTH = 100 // 默认分页查询记录数

	INSTANCE_BOOTTYPE_NORECOVER        = 0 // 不恢复
	INSTANCE_BOOTTYPE_RECOVER          = 1 // 恢复，老的调用image_service, Iaas合研使用
	INSTANCE_BOOTTYPE_LOCAL            = 2 // 本地盘启动，用于并行启动
	INSTANCE_BOOTTYPE_BINDNORECOVER    = 3 // 不恢复，被串流调度绑定并且启动时绑定的客户机
	INSTANCE_BOOTTYPE_DISKLESS_UPGRADE = 4 // 恢复，新的调用diskless_image_service
	INSTANCE_BOOTTYPE_DRIVER           = 5 // 驱动模式(只恢复系统镜像), 老的调用image_service, Iaas合研使用

	SCHEME_STORAGETYPE_ISCSI      = 0 // os/data iscsi
	SCHEME_STORAGETYPE_RBDMAPDATA = 1 // os iscsi + data rbd-map
	SCHEME_STORAGETYPE_RBDMAPALL  = 2 // os/data rbd-map
	// SCHEME_STORAGETYPE_LIBRBD      = 3 // os/data ceph librbd
	// SCHEME_STORAGETYPE_QCOW2BDEV   = 4 // imgqcow2 + rawbdev
	SCHEME_STORAGETYPE_COMBINATION = 5   // 组合
	SCHEME_STORAGETYPE_NETZONE     = 900 // netzone兼容模式

	////// 业务类型
	// 系统盘和数据盘区分不同type, 目的是页面展示/操作逻辑不同，比如上传下载
	// 私有盘和数据盘区分不同type, 目的是页面展示/更新内容writable的逻辑不同
	IMG_IMAGETYPE_OSRBD          = 0   // ceph系统盘
	IMG_IMAGETYPE_PRIRBD         = 1   // 私有ceph盘
	IMG_IMAGETYPE_OSIMGQCOW2     = 2   // 本地系统盘
	IMG_IMAGETYPE_DATARAWBDEV    = 3   // 裸盘数据盘
	IMG_IMAGETYPE_PRIIMGQCOW2    = 4   // 私有本地盘
	IMG_IMAGETYPE_DATARBDCOMSNAP = 5   // ceph数据盘母盘（统一快照）
	IMG_IMAGETYPE_DATARBDSNAP    = 6   // ceph数据盘指定快照
	IMG_IMAGETYPE_COWCOMSNAP     = 7   // 自研cow快照盘母盘（统一快照）
	IMG_IMAGETYPE_COWSNAP        = 8   // 自研cow快照盘指定快照
	IMG_IMAGETYPE_OSRBDSNAP      = 9   // ceph系统盘指定快照
	IMG_IMAGETYPE_PERSONAL       = 10  // 个人盘,对应bstype=6
	IMG_IMAGETYPE_NETZONE        = 900 // netzone磁盘

	////// 真实存储类型
	IMG_BSTYPE_INIT      = 0   // "rbd"
	IMG_BSTYPE_RBD       = 1   // "rbd"
	IMG_BSTYPE_IMGQCOW2  = 2   // "image"
	IMG_BSTYPE_RBDSNAP   = 3   // "rbd快照"
	IMG_BSTYPE_RAWBDEV   = 4   // "raw"
	IMG_BSTYPE_COW       = 5   // "cow"
	IMG_BSTYPE_DIRECTRBD = 6   // "direct rbd，直接独占读写"
	IMG_BSTYPE_NZ_BASE   = 900 // nezone基础盘
	IMG_BSTYPE_NZ_TPL    = 901 // nezone磁盘模板
	IMG_BSTYPE_NZ_DIFF   = 902 // nezone差异盘

	///// 自研挂载接口参数 bsPath + srcType
	IMG_MNTTYPE_INIT           = 0 // "iscsi原始方式"
	IMG_MNTTYPE_ISCSI          = 1 // "iscsi原始方式"
	IMG_MNTTYPE_NBD            = 2 // "自研nbd：跟原始iscsi相同路径"
	IMG_MNTTYPE_RBD            = 3 // "自研rbd"
	IMG_MNTTYPE_LIBRBD         = 4 // "自研librbd"
	IMG_MNTTYPE_IMGQCOW2       = 5 // "自研image"
	IMG_MNTTYPE_RAWBDEV        = 6 // "自研raw"
	IMG_MNTTYPE_COW            = 7 // "自研cow"
	IMG_MNTTYPE_ISCSI_PERSONAL = 8 // "iscsi个人独占盘"

	///// 传给底层IO引擎的字符串描述
	IMG_MNTTYPEDES_ISCSI    = "iscsi"
	IMG_MNTTYPEDES_NBD      = "nbd"
	IMG_MNTTYPEDES_RBD      = "rbd"
	IMG_MNTTYPEDES_LIBRBD   = "librbd"
	IMG_MNTTYPEDES_IMGQCOW2 = "image"
	IMG_MNTTYPEDES_RAWBDEV  = "raw"
	IMG_MNTTYPEDES_COW      = "cow"

	INSTANCE_KEEPTYPE_NONE  = 0 // old mode
	INSTANCE_KEEPTYPE_READY = 1 // never reboot
	INSTANCE_KEEPTYPE_DONE  = 2 // reboot happened

	SCHEME_WRBACKTYPE_NONE    = 0
	SCHEME_WRBACKTYPE_LDISK   = 1
	SCHEME_WRBACKTYPE_BACKEND = 2

	BOOT_PROTOCOL_ISCSI    = 0 // iscsi
	BOOT_PROTOCOL_DISKLESS = 1 // diskless

	VOLUME_TYPE_RBD      = 1 // rbd
	VOLUME_TYPE_DISKLESS = 2 // diskless
	VOLUME_TYPE_NBD      = 3 // qcow2+nbd

	DEVICE_CLOUD_HOST    = 0
	DEVICE_LOCAL_HOST    = 1
	DEVICE_CLOUD_BOX     = 2
	DEVICE_LOCAL_BOX     = 3 // 废弃，合并到2中
	DEVICE_IAASPLUS_HOST = 1001
	DEVICE_ESHORE_HOST   = 2001

	INSTANCE_VOLUME_ENDPOINT_LOCALDISK = "localdisk"  // 本地盘启动
	INSTANCE_VOLUME_ENDPOINT_ISCSI     = "iscsi"      // 单个iscsi链接
	INSTANCE_VOLUME_ENDPOINT_V1        = "disklessv1" // dios
	INSTANCE_VOLUME_ENDPOINT_V2        = "disklessv2" // 无盘镜像管理调用iscsi存储引擎创建的
	INSTANCE_VOLUME_ENDPOINT_V3        = "disklessv3" // 无盘镜像管理调用自研存储引擎创建的

	WriteModeSave = 0 // 保存模式
	WriteModeWB   = 1 // 回写文件模式

)

var ImgMntDes = map[int]string{
	IMG_MNTTYPE_INIT:     IMG_MNTTYPEDES_ISCSI,
	IMG_MNTTYPE_ISCSI:    IMG_MNTTYPEDES_ISCSI,
	IMG_MNTTYPE_NBD:      IMG_MNTTYPEDES_NBD,
	IMG_MNTTYPE_RBD:      IMG_MNTTYPEDES_RBD,
	IMG_MNTTYPE_LIBRBD:   IMG_MNTTYPEDES_LIBRBD,
	IMG_MNTTYPE_IMGQCOW2: IMG_MNTTYPEDES_IMGQCOW2,
	IMG_MNTTYPE_RAWBDEV:  IMG_MNTTYPEDES_RAWBDEV,
	IMG_MNTTYPE_COW:      IMG_MNTTYPEDES_COW,
}

// bstype
func CheckOsBstype(bstype int32) bool {
	if bstype == IMG_BSTYPE_INIT ||
		bstype == IMG_BSTYPE_RBD ||
		bstype == IMG_BSTYPE_IMGQCOW2 {
		return true
	}
	return false
}

func CheckDataBstype(bstype int32) bool {
	if bstype == IMG_BSTYPE_INIT || // private
		bstype == IMG_BSTYPE_RBD || // private
		bstype == IMG_BSTYPE_IMGQCOW2 || // private
		bstype == IMG_BSTYPE_RBDSNAP ||
		bstype == IMG_BSTYPE_RAWBDEV ||
		bstype == IMG_BSTYPE_COW {
		return true
	}
	return false
}

// type
func CheckOsType(bstype, typ int32) bool {
	switch bstype {
	case IMG_BSTYPE_INIT, IMG_BSTYPE_RBD, IMG_BSTYPE_DIRECTRBD:
		if typ == IMG_IMAGETYPE_OSRBD {
			return true
		}
	case IMG_BSTYPE_IMGQCOW2:
		if typ == IMG_IMAGETYPE_OSIMGQCOW2 {
			return true
		}
	default:
		return false
	}
	return false
}

func CheckDataType(bstype, typ int32) bool {
	switch bstype {
	case IMG_BSTYPE_INIT, IMG_BSTYPE_RBD: //private
		if typ == IMG_IMAGETYPE_PRIRBD {
			return true
		}
	case IMG_BSTYPE_IMGQCOW2: //private
		if typ == IMG_IMAGETYPE_PRIIMGQCOW2 {
			return true
		}
	case IMG_BSTYPE_RBDSNAP:
		if typ == IMG_IMAGETYPE_DATARBDSNAP || typ == IMG_IMAGETYPE_DATARBDCOMSNAP {
			return true
		}
	case IMG_BSTYPE_RAWBDEV:
		if typ == IMG_IMAGETYPE_DATARAWBDEV {
			return true
		}
	case IMG_BSTYPE_COW:
		if typ == IMG_IMAGETYPE_COWCOMSNAP || typ == IMG_IMAGETYPE_COWSNAP {
			return true
		}
	default:
		return false
	}
	return false
}

// mnt type
func CheckOsMntType(bstype int32, mnttype int) bool {
	switch bstype {
	case IMG_BSTYPE_INIT, IMG_BSTYPE_RBD, IMG_BSTYPE_DIRECTRBD:
		if mnttype == IMG_MNTTYPE_INIT || mnttype == IMG_MNTTYPE_ISCSI || mnttype == IMG_MNTTYPE_RBD || mnttype == IMG_MNTTYPE_LIBRBD {
			return true
		}
	case IMG_BSTYPE_IMGQCOW2:
		if mnttype == IMG_MNTTYPE_IMGQCOW2 {
			return true
		}
	default:
		return false
	}
	return false
}

func CheckDataMntType(bstype int32, mnttype int) bool {
	switch bstype {
	case IMG_BSTYPE_INIT, IMG_BSTYPE_RBD: //private
		if mnttype == IMG_MNTTYPE_INIT || mnttype == IMG_MNTTYPE_ISCSI || mnttype == IMG_MNTTYPE_RBD || mnttype == IMG_MNTTYPE_LIBRBD {
			return true
		}
	case IMG_BSTYPE_IMGQCOW2: //private
		if mnttype == IMG_MNTTYPE_IMGQCOW2 {
			return true
		}
	case IMG_BSTYPE_RBDSNAP:
		if mnttype == IMG_MNTTYPE_INIT || mnttype == IMG_MNTTYPE_ISCSI || mnttype == IMG_MNTTYPE_RBD || mnttype == IMG_MNTTYPE_LIBRBD {
			return true
		}
	case IMG_BSTYPE_RAWBDEV:
		if mnttype == IMG_MNTTYPE_RAWBDEV {
			return true
		}
	case IMG_BSTYPE_COW:
		if mnttype == IMG_MNTTYPE_COW {
			return true
		}
	default:
		return false
	}
	return false
}

// valid:配置合法
// upgrade:需要升级
// system-drvflag:给驱动的超管标志
func CheckBootKeepType(bootType, keepType int) (bool, bool, bool) {
	switch bootType {
	case INSTANCE_BOOTTYPE_NORECOVER, INSTANCE_BOOTTYPE_BINDNORECOVER:
		if keepType == INSTANCE_KEEPTYPE_READY {
			return true, true, true
		} else if keepType == INSTANCE_KEEPTYPE_DONE {
			return true, false, true
		} else {
			return false, false, false
		}
	case INSTANCE_BOOTTYPE_RECOVER, INSTANCE_BOOTTYPE_DISKLESS_UPGRADE:
		if keepType == INSTANCE_KEEPTYPE_NONE {
			return true, true, false
		} else {
			return false, false, false
		}
	case INSTANCE_BOOTTYPE_DRIVER:
		if keepType == INSTANCE_KEEPTYPE_READY {
			return true, true, false
		} else {
			return true, false, false
		}
	case INSTANCE_BOOTTYPE_LOCAL:
		return false, false, false
	default:
		return false, false, false
	}
}

func GetDbRbdBsPath(bspath string) (string, error) {
	if len(bspath) == 0 {
		return "", fmt.Errorf("bspath len is 0")
	}

	if strings.Contains(bspath, ":") {
		srcs := strings.Split(bspath, ":")
		if len(srcs) < 2 {
			return "", fmt.Errorf("bspath=%v split invalid len=%v < 2", bspath, len(srcs))
		}
		return srcs[1], nil
	}

	return bspath, nil
}

type NetInfo struct {
	Ip       string `json:"ip,omitempty"`
	Netmask  string `json:"netmask,omitempty"`
	Gateway  string `json:"gateway,omitempty"`
	Dns      string `json:"dns,omitempty"`
	Hostname string `json:"hostname,omitempty"`
	Vlan     int    `json:"vlan,omitempty"` // 虚拟子网
}

type InstanceOptions struct {
	ShareVolume  int `json:"share_volume,omitempty"`   // 1-使用只读共享卷
	OsVolumeOnly int `json:"os_volume_only,omitempty"` // 1-只挂载系统盘
	DebugAll     int `json:"debug_all,omitempty"`      // 1-所有盘调试
}

type DefaultConfig struct {
	NetInfo
}

type AssigParams struct {
	Uid      int64  `json:"uid,omitempty"`
	Gid      int64  `json:"gid,omitempty"`
	UUid     int64  `json:"uuid,omitempty"`
	UGid     int64  `json:"ugid,omitempty"`
	DeviceID string `json:"device_id,omitempty"`
}

type UpdateableInstanceInfo struct {
	SchemeId           *int64           `json:"scheme_id,omitempty"` // 业务编排Id
	OsVolume           string           `json:"os_volume,omitempty"`
	OsVolumeId         int64            `json:"os_volume_id,omitempty"` // os volume id
	DataVolumes        []string         `json:"data_volumes,omitempty"`
	NetInfo            *NetInfo         `json:"net_info,omitempty"`
	BootType           *int             `json:"boot_type,omitempty"`            // 0/3/5/6-keep 1-discard 2-local 4-normal
	KeepType           *int             `json:"keep_type,omitempty"`            // 0-normal 1-ready-save 2-saving
	BootAction         *int             `json:"boot_action,omitempty"`          // 0-normal 1-ipxe wait 2-ipxe resume
	Tags               []string         `json:"tags,omitempty"`                 // 标签
	Remark             *string          `json:"remark,omitempty"`               // 备注
	Specification      *int64           `json:"specification,omitempty"`        // 规格id
	Options            *InstanceOptions `json:"options,omitempty"`              // 选项 json格
	CleanupDataVolumes bool             `json:"cleanup_data_volumes,omitempty"` // 清空数据盘
	DefaultConfig      *DefaultConfig   `json:"default_config,omitempty"`       // 默认配置 json格式
	// 硬件信息
	IpmiType     int    `json:"ipmi_type"`
	IpmiAddress  string `json:"ipmi_address"`
	IpmiSlot     int    `json:"ipmi_slot"`
	Vlan         int    `json:"vlan"`
	SwitchId     int64  `json:"switch_id"`
	SwitchPort   *int   `json:"switch_port"`
	PowerAddress string `json:"power_address"`
}

type InstanceInfo struct {
	Id            int64           `json:"id"`                  // 实例id
	HostId        int64           `json:"host_id"`             // 硬件id
	Area          int             `json:"area,omitempty"`      // 区域id
	SchemeId      int64           `json:"scheme_id,omitempty"` // 业务编排Id
	OsVolume      string          `json:"os_volume,omitempty"`
	OsVolumeId    int64           `json:"os_volume_id,omitempty"` // os volume id
	DataVolumes   []string        `json:"data_volumes,omitempty"`
	NetInfo       NetInfo         `json:"net_info,omitempty"`       // 网络配置
	BootDev       string          `json:"boot_dev,omitempty"`       // 启动设备 默认net0
	BootMac       string          `json:"boot_mac"`                 // mac地址
	BootType      int             `json:"boot_type,omitempty"`      // 0/3/5/6-keep 1-discard 2-local 4-normal
	KeepType      int             `json:"keep_type,omitempty"`      // 0-normal 1-ready-save 2-saving
	BootAction    int             `json:"boot_action,omitempty"`    // 0-normal 1-ipxe wait 2-ipxe resume
	BootProtocol  int             `json:"boot_protocol,omitempty"`  // 0-iscsi 1-diskless
	BootScript    string          `json:"boot_script,omitempty"`    // 自定义启动脚本
	DefaultConfig DefaultConfig   `json:"default_config,omitempty"` // 默认配置 json格式
	Options       InstanceOptions `json:"options,omitempty"`        // 选项 json格式
	Specification int64           `json:"specification,omitempty"`  // 规格id
	State         int             `json:"state"`                    // 0-初始 100-已创建 900-已销毁
	Tags          []string        `json:"tags,omitempty"`           // 标签
	Remark        string          `json:"remark"`                   // 备注
	Chksum        string          `json:"chksum"`                   // 检查和
	CreateTime    time.Time       `json:"create_time"`
	UpdateTime    time.Time       `json:"update_time"`
	ModifyTime    time.Time       `json:"modify_time"`
}

type CreateInstanceRequest struct {
	FlowId string `json:"flow_id"` // 流水ID
	InstanceInfo
	IpmiType     int    `json:"ipmi_type"`
	IpmiAddress  string `json:"ipmi_address"`
	IpmiSlot     int    `json:"ipmi_slot"`
	DeviceType   int    `json:"device_type"` // 0-云主机 1-本地主机 2-云盒 3-本地盒子
	Vlan         int    `json:"vlan"`
	SwitchId     int64  `json:"switch_id"`
	SwitchPort   *int   `json:"switch_port"`
	PowerAddress string `json:"power_address"`
}

type CreateInstanceResponse struct {
	FlowId     string `json:"flow_id,omitempty"` // 流水ID
	InstanceId int64  `json:"instance_id"`
	Mac        string `json:"mac"`
}

type DestroyInstanceRequest struct {
	FlowId     string `json:"flow_id"` // 流水ID
	InstanceId int64  `json:"instance_id,omitempty"`
	Mac        string `json:"mac,omitempty"`
	Freeze     bool   `json:"freeze,omitempty"` // 冻结实例 (无法开机)
}

type DestroyInstanceResponse struct {
	FlowId string `json:"flow_id"` // 流水ID
}

type GetInstanceRequest struct {
	InstanceId int64  `json:"id" form:"id"`
	Mac        string `json:"mac" form:"mac"`
	Vlan       int    `json:"vlan" form:"vlan"`
	Dhcp       bool   `json:"dhcp" form:"dhcp"`
}

type GetInstanceResponse InstanceInfo

type UpdateInstanceRequest struct {
	FlowId      string  `json:"flow_id"` // 流水ID
	InstanceId  int64   `json:"instance_id"`
	InstanceIds []int64 `json:"instance_ids"` // 当instance_id为0时才使用instance_ids
	UpdateableInstanceInfo
}

type UpdateInstanceResponse struct {
	FlowId      string  `json:"flow_id"` // 流水ID
	InstanceId  int64   `json:"instance_id"`
	InstanceIds []int64 `json:"instance_ids"`
}

// GetAdminRequest 获取
type GetAdminRequest struct {
	FlowID     string `json:"flow_id"`        // 流水ID
	Mac        string `json:"mac" form:"mac"` // 有mac就不使用instanceid
	InstanceID int64  `json:"id" form:"id"`
}

// UserMode 用户模式
type UserMode int

const (
	// AdminUser 超管模式
	AdminUser UserMode = iota
	// RegularUser 普通模式
	RegularUser
	// LocalBootUser 不走无盘的本地启动
	LocalBootUser
	// BindAdminUser 被调度绑定的超管模式
	BindAdminUser
	// RegularUser2 新版本普通模式
	RegularUser2
	// DriverUser 驱动模式 驱动正常工作+后台保存数据
	DriverUser
)

// GetAdminResponse 获取
type GetAdminResponse struct {
	UserMode UserMode `json:"user_mode"`
}

// SetAdminRequest 获取
type SetAdminRequest struct {
	FlowID     string   `json:"flow_id"` // 流水ID
	AppID      string   `json:"app_id"`  // 表示业务来源方
	Mac        string   `json:"mac"`     // 有mac就不使用instanceid
	InstanceID int64    `json:"id"`
	UserMode   UserMode `json:"user_mode"`
}

// SetAdminResponse 获取
type SetAdminResponse struct{}

type ListInstancesRequest struct {
	Offset      int    `json:"offset" form:"offset"` // 偏移
	Length      int    `json:"length" form:"length"` // 长度 0-所有
	Order       string `json:"order" form:"order"`   // 排序 asc/desc
	InstanceId  int64  `json:"id" form:"id"`
	Mac         string `json:"mac" form:"mac"`
	InstanceIds string `json:"ids" form:"ids"`   // 逗号分隔
	Macs        string `json:"macs" form:"macs"` //
	Tag         string `json:"tag" form:"tag"`   // 标签
	Vlan        int    `json:"vlan" form:"vlan"`
	DeviceType  int    `json:"device_type" form:"device_type"` // 设备类型
	Ips         string `json:"ips" form:"ips"`                 // 逗号分隔
}

type ListInstancesRequestNew struct {
	Offset         int      `json:"offset" form:"offset"`                   // 偏移
	Length         int      `json:"length" form:"length"`                   // 长度 0-所有
	Order          string   `json:"order" form:"order"`                     // 排序 asc/desc
	InstanceIds    []int    `json:"ids" form:"ids"`                         //实例id
	Macs           []string `json:"macs" form:"macs"`                       //实例MAC
	HostIds        []int    `json:"hostids" form:"hostids"`                 //主机ids
	OsImageVersion []string `json:"os_image_ids" form:"os_image_ids"`       //镜像版本
	DataVersion    []string `json:"data_versions" form:"data_versions"`     //游戏镜像版本
	ManageStatus   []int    `json:"manager_status" form:"manager_status"`   //管理状态
	RunningStatus  []int    `json:"running_status" form:"running_status"`   //OS状态
	BusinessStatus []int    `json:"business_status" form:"business_status"` //应用状态

	DeviceTypes   []int    `json:"device_types" form:"device_types"`   // 设备类型
	AssignStatus  []int    `json:"assign_status" form:"assign_status"` // 分配状态
	PowerStatus   []int    `json:"power_status" form:"power_status"`   // 电源状态
	Specification []int    `json:"specification" form:"specification"` //规格
	Tag           []string `json:"tag" form:"tag"`                     // 标签
	Vlan          []int    `json:"vlan" form:"vlan"`
	Ips           []string `json:"ips" form:"ips"` // 逗号分隔
}

type HostInfo struct {
	Name          string `json:"name"`
	Arch          string `json:"arch,omitempty"`
	Cpu           string `json:"cpu,omitempty"`
	Gpu           string `json:"gpu,omitempty"`
	Net           string `json:"net,omitempty"`
	Mem           string `json:"mem,omitempty"`
	Disk0         string `json:"disk0,omitempty"` // 本地硬盘大小
	IpmiType      int    `json:"ipmi_type"`       // ipmi 类型 0-nodpwd.x 1-鑫誉v1 2-天翼云grpc
	IpmiAddress   string `json:"ipmi_address"`    // ipmi管理地址
	IpmiSlot      int    `json:"ipmi_slot"`       // ipmi槽位
	IpmiParam     string `json:"ipmi_param"`      // ipmi参数
	SwitchID      int64  `json:"switch_id"`       // 交换机id
	SwitchPort    int32  `json:"switch_port"`     // 交换机端口
	Vlan          int    `json:"vlan"`            // 默认Vlan
	DeviceType    int    `json:"device_type"`     // 0-云主机 1-本地主机 2-云盒 3-本地盒子
	Specification int    `json:"specification"`   // 规格
	Status        int    `json:"status"`          // 状态
	PowerAddress  string `json:"power_address"`   // 电源上报地址
}

type InstanceDetail struct {
	Id           int64           `json:"id"`                      // 实例id
	HostId       int64           `json:"host_id"`                 // 硬件id
	SchemeId     int64           `json:"scheme_id"`               // 编排方案id
	NetInfo      NetInfo         `json:"net_info,omitempty"`      // 网络配置
	BootMac      string          `json:"boot_mac"`                // mac地址
	BootType     int             `json:"boot_type,omitempty"`     // 0-normal 1-recover 2-本地盘启动
	KeepType     int             `json:"keep_type,omitempty"`     // 0-normal 1-ready 2-done
	BootAction   int             `json:"boot_action,omitempty"`   // 0-normal 1-ipxe wait 2-ipxe resume
	BootProtocol int             `json:"boot_protocol,omitempty"` // 0-iscsi 1-diskless
	BootScript   string          `json:"boot_script,omitempty"`   // 自定义启动脚本
	Options      InstanceOptions `json:"options,omitempty"`       // 选项 json格式
	State        int             `json:"state"`                   // 0-初始 100-已创建 900-已销毁
	Tags         []string        `json:"tags,omitempty"`          // 标签

	ActivityIp         string      `json:"activity_ip"`          // 实例当前IP
	BootTime           time.Time   `json:"boot_time"`            // 开机时间
	ManageStatus       int         `json:"manage_status"`        // 管理状态 0-normal 1-disable
	PowerStatus        int         `json:"power_status"`         // 电源状态 0-off 1-on
	RunningStatus      int         `json:"running_status"`       // 运行状态
	BootStatus         int         `json:"boot_status"`          // 启动状态
	BusinessStatus     int         `json:"business_status"`      // 业务状态
	AssignStatus       int         `json:"assign_status"`        // 分配状态
	ManageStatusDesc   string      `json:"manage_status_desc"`   // 管理状态描述
	PowerStatusDesc    string      `json:"power_status_desc"`    // 电源状态描述
	RunningStatusDesc  string      `json:"running_status_desc"`  // 运行状态描述
	BootStatusDesc     string      `json:"boot_status_desc"`     // 启动状态描述
	BusinessStatusDesc string      `json:"business_status_desc"` // 业务状态描述
	AssignStatusDesc   string      `json:"assign_status_desc"`   // 分配状态描述
	AssignSource       string      `json:"assign_source"`        // 分配来源
	AssignOrder        string      `json:"assign_order"`         // 分配订单
	Specification      int64       `json:"specification"`        // 规格id
	DeviceType         int         `json:"device_type"`          // 0-云主机 1-本地主机 2-云盒 3-本地盒子
	BootSession        string      `json:"boot_session"`         // 启动session
	AssignParam        AssigParams `json:"assign_param,omitempty"`
	StatusRemark       string      `json:"status_remark"`   // 状态备注
	InstanceRemark     string      `json:"instance_remark"` // 实例备注
	UserMode           UserMode    `json:"user_mode"`       // 0-超管模式 1-普通模式

	HostInfo      HostInfo      `json:"host_info"`                // 调度的网络信息
	DefaultConfig DefaultConfig `json:"default_config,omitempty"` // 分配出去的网络信息

	OsImage    string `json:"os_image"`
	OsVolumeId int64  `json:"os_volume_id,omitempty"` // os volume id
	DataImage  string `json:"data_image"`
}

type ListInstancesResponse struct {
	Instances []InstanceDetail `json:"instances"` // 实例信息列表
	Total     int              `json:"total"`
}

type ApplySchemeRequest struct {
	FlowId      string   `json:"flow_id"` // 流水ID
	InstanceId  int64    `json:"instance_id"`
	InstanceIds []int64  `json:"instance_ids"` // 批量升级使用
	SchemeId    int64    `json:"scheme_id"`
	NetInfo     *NetInfo `json:"net_info,omitempty"`
}

type ApplySchemeResponse struct {
	Status int `json:"status"`
}

type CancelSchemeRequest struct {
	FlowId     string `json:"flow_id"` // 流水ID
	InstanceId int64  `json:"instance_id"`
	Mac        string `json:"mac"`
}

type CancelSchemeResponse struct {
	Status int `json:"status"`
}

type RestoreInstanceRequest struct {
	FlowId           string `json:"flow_id"` // 流水ID
	InstanceId       int64  `json:"instance_id"`
	OsImageVersion   string `json:"os_image_version"`
	GameImageVersion string `json:"game_image_version,omitempty"`
	Type             int    `json:"type,omitempty"`
}

type RestoreInstanceResponse struct {
	TaskId string `json:"task_id,omitempty"`
}

type CreateInstanceItem struct {
	Hostname string `csv:"hostname"`
	Mac      string `csv:"mac"`
	Ip       string `csv:"ip"`
}

type BatchCreateInstancesRequest struct {
	FlowId     string                `form:"flow_id"` // 流水ID
	Vlan       int                   `form:"vlan"`
	DeviceType int                   `form:"device_type"` // 0-云主机 1-本地主机 2-云盒 3-本地盒子
	File       *multipart.FileHeader `form:"file"`        // csv format
}

type BatchCreateInstancesResponse struct {
	Total  int                      `json:"total"`
	Result []CreateInstanceResponse `json:"result"`
}

type InstanceExecRequest struct {
	FlowId     string `form:"flow_id"` // 流水ID
	InstanceId int64  `form:"instance_id"`
	Address    string `form:"address"`
	Order      string `form:"order"` // 订单号 用于保证命令只执行一次 为空则随机生成
	Shell      string `form:"shell"` // cmd/ps
	Async      bool   `form:"async"` // 异步执行
}

type InstanceExecResponse struct {
	Output string `json:"output"`
}
