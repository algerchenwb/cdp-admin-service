package diskless

const (
	MANAGE_STATUS_AVAILABLE  = 1   // 允许使用
	MANAGE_STATUS_DISABLE    = 2   // 暂停使用
	MANAGE_STATUS_UPGRADABLE = 3   // 允许升级
	MANAGE_STATUS_DEBUG      = 4   // 允许调试
	MANAGE_STATUS_DESTROYED  = 900 // 已销毁
)

// from ： https://confluence.vrviu.com/pages/viewpage.action?pageId=1674489786
// assign_config格式
type SchemeGrayItem struct {
	Percentage int   `json:"percentage,omitempty"` // 所有百分比相加不得超过100
	SchemeId   int64 `json:"scheme_id,omitempty"`
}
type AssignInfo struct {
	SchemeId         int64            `json:"scheme_id,omitempty"`          // 默认编排方案
	SchemeGrayConfig []SchemeGrayItem `json:"scheme_gray_config,omitempty"` // 灰度配置 存在的时候忽略默认编排方案
}

// 不更新就不用填，要更新为0就填-1，要更新为空字符串就填"-"
const (
	DisklessZero        = -1
	DisklessEmptyString = "-"
)
