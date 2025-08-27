package saas

import (
	"cdp-admin-service/internal/helper"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type TEsportsSpecInstancesInfo struct {
	Id             int64     `json:"spec_instance_id"`
	OuterSpecID    int64     `json:"outer_spec_id"`
	TotalInstances int64     `json:"total_instances"`
	PrimaryID      int64     `json:"primary_id"`
	AgentID        int64     `json:"agent_id"`
	BizID          int64     `json:"biz_id"`
	CreateBy       string    `json:"create_by"`
	UpdateBy       string    `json:"update_by"`
	Remark         string    `json:"remark"`
	CreateTime     time.Time `json:"create_time"`
	UpdateTime     time.Time `json:"update_time"`
	ModifyTime     time.Time `json:"modify_time"`

	OuterSpecName string `orm:"-" json:"outer_spec_name"`
}

type TEsportsAgentInfo struct {
	Id                    int64     `json:"agent_id"`                // 自增ID
	AuthoritySet          string    `json:"authority_set"`           // 权限集合
	AreaType              int64     `json:"area_type"`               // 区域ID
	PrimaryID             int64     `json:"primary_id"`              // 一级代理商ID
	Name                  string    `json:"name"`                    // 名称
	WallpaperID           int64     `json:"wallpaper_id"`            // 壁纸ID
	TemplateID            int64     `json:"template_id"`             // 模板ID
	TemplatePriority      int64     `json:"template_priority"`       // 模板优先级
	TotalInstances        int64     `json:"total_instances"`         // 实例数
	EffectiveDate         time.Time `json:"effective_date"`          // 实际生效日期
	ExpectedEffectiveDate time.Time `json:"expected_effective_date"` // 期待生效日期
	CreateBy              string    `json:"create_by"`               // 创建者
	UpdateBy              string    `json:"update_by"`               // 修改者
	Remark                string    `json:"remark"`                  // 备注
	State                 int       `json:"state"`                   // 状态: 0-无效, 1-有效, 2-删除, 3-暂停合作
	Email                 string    `json:"email"`                   // 邮箱
	CreateTime            time.Time `json:"create_time"`             // 创建时间
	UpdateTime            time.Time `json:"update_time"`             // 更新时间
	ModifyTime            time.Time `json:"modify_time"`             // 修改时间

	SpecInstances []TEsportsSpecInstancesInfo `json:"spec_instances" orm:"-"`
	ResourcePools []TEsportsResourcePoolInfo  `json:"resource_pools" orm:"-"`
}

type GetAgentInfoResp struct {
	TEsportsAgentInfo
}

// 单个二级代理商的信息
type EsportGetAgentInfoResp struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data GetAgentInfoResp `json:"data"`
}

func (s *SaasServerService) GetAgentInfo(ctx context.Context, sessionId string, agentId int64) (agentInfo *GetAgentInfoResp, err error) {
	url := fmt.Sprintf("%s/agent/getAgentInfo/%d", s.Host, agentId)

	resp := new(EsportGetAgentInfoResp)

	_, err = helper.HttpGet(ctx, sessionId, url, s.Timeout, resp)
	if err != nil {
		logx.WithContext(ctx).Errorf("helper.HttpPost err  host[%s] err:%+v", url, err)
		return nil, errors.New("查询agent信息http错误")
	}

	if resp.Code != 0 {
		logx.WithContext(ctx).Errorf("[%s] GetAgentInfo httpc.Do err agentId[%d] host[%s] err:%+v", sessionId, agentId, url, err)
		return nil, errors.New(resp.Msg)
	}

	agentInfo = &resp.Data

	return
}
