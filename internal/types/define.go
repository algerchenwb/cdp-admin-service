package types

const (
	InstanceEventOpen  = "open"
	InstanceEventClose = "close"
)

const MaxFileSize = 1024 * 1024 * 2 // 2 MB
const FIXED_MAX_CSV_RECORD_NUM = 300

var ManageStatusDescMap = map[int]string{
	0: "超管中",
	1: "普通",
	3: "超管中",
	4: "普通",
}

type CloudBoxImport struct {
	Name       string `json:"name"         validate:"required,min=0,max=15"   label:"云盒名称"   csv:"云盒名称"`    // 云盒名称
	Mac        string `json:"mac"          validate:"required,mac"            label:"MAC地址"    csv:"MAC地址"` // 云盒 mac 地址
	IP         string `json:"ip"           validate:"required,ipv4"           label:"IP地址"     csv:"IP地址"`  // 云盒IP
	Mode       int    `json:"mode"         validate:"required,number,gte=0"   label:"启动方式"    csv:"启动方式"`   // 启动方式
	SpecInstId int64  `json:"specInstId"   validate:"required,number,gte=0"   label:"规格ID"   csv:"规格ID"`    //     规格ID                                           // 规格ID
	ConfigId   string `json:"configId"`                                                                     // 配置方案ID
}

type CloudClientImport struct {
	Name         string `json:"name"         validate:"required,min=0,max=64"   label:"云主机名"      csv:"云主机名"`    // 云主机名称
	Mac          string `json:"mac"          validate:"required,mac"            label:"MAC地址"       csv:"MAC地址"` // 云主机 mac 地址
	IP           string `json:"ip"           validate:"required,ipv4"           label:"IP地址"        csv:"IP地址"`  // 云主机 IP
	SpecInstId   int64  `json:"specInstId"   validate:"required,number,gte=0"   label:"规格ID"        csv:"规格ID"`  // 规格ID
	ConfigId     string `json:"configId"`                                                                        // 配置方案ID
	DeviceNumber string `json:"deviceNumber"`                                                                    // 设备编号
	PoolId       int64  `json:"poolId"`                                                                          // 资源池ID

}

type ShopInfoImportExt struct {
	PrimaryId   int64  `json:"primaryId,omitempty"`   // 一级代理商Id
	AgentId     int64  `json:"agentId,omitempty"`     // 二级代理商Id
	BizId       int64  `json:"bizId,omitempty"`       // 租户ID
	AccountName string `json:"accountName,omitempty"` // 登录帐户名
	AreaId      int64  `json:"areaId,omitempty"`      // 节点区域ID
	VlanId      int64  `json:"vlanId,omitempty"`      // 网络vlanId
	Gateway     string `json:"gateway,omitempty"`     // 网关
	SubNetMask  string `json:"subNetmask,omitempty"`  // 子网埯码
	PreDNS      string `json:"preDNs,omitempty"`      // 首选DNS
	BackupDNS   string `json:"backupDNS,omitempty"`   // 备用DNS
}
