package saas

import (
	"cdp-admin-service/internal/helper"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type TEsportsUserInfo struct {
	AuthoritySet          string    `orm:"column(authority_set)" description:"权限集合" json:"authority_set"`
	Id                    int64     `orm:"column(user_id);auto" description:"自增ID" json:"user_id"`
	Type                  int64     `orm:"column(type)" description:"类型: 1-biz, 2-Agent" json:"type"`
	AreaType              int64     `orm:"column(area_type)" description:"区域ID" json:"area_type"`
	SetID                 int64     `orm:"column(set_id)" description:"业务ID" json:"set_id"`
	PrimaryID             int64     `orm:"column(primary_id)" description:"一级代理商ID" json:"primary_id"`
	AgentID               int64     `orm:"column(agent_id)" description:"二级代理商ID" json:"agent_id"`
	BizID                 int64     `orm:"column(biz_id)" description:"租户Id" json:"biz_id"`
	UserName              string    `orm:"column(user_name)" description:"名称" json:"user_name"`
	WallpaperID           int64     `orm:"column(wallpaper_id)" description:"壁纸ID" json:"wallpaper_id"`
	TemplateID            int64     `orm:"column(template_id)" description:"模板ID" json:"template_id"`
	TemplatePriority      int64     `orm:"column(template_priority)" description:"模板优先级" json:"template_priority"`
	TotalInstances        int64     `orm:"column(total_instances)" description:"实例数" json:"total_instances"`
	EffectiveDate         time.Time `orm:"column(effective_date);type(datetime)" description:"实际生效日期" json:"effective_date"`
	ExpectedEffectiveDate time.Time `orm:"column(expected_effective_date);type(datetime)" description:"期待生效日期" json:"expected_effective_date"`
	Email                 string    `orm:"column(email)" description:"邮箱" json:"email"`
	AccessKey             string    `orm:"column(access_key)" description:"租户key" json:"access_key"`
	AccessSecret          string    `orm:"column(access_secret)" description:"租户secret" json:"access_secret"`
	EnableQualityReport   int       `orm:"column(enable_quality_report)" description:"质量上报开关: 0-关闭, 1-打开" json:"enable_quality_report"`
	UpgradeGreyType       int       `orm:"column(upgrade_grey_type)" description:"灰度方式: 0-余量不升级, 1-余量升级" json:"upgrade_grey_type"`
	Domain                string    `orm:"column(domain)" description:"域名" json:"domain"`
	VlanID                int64     `orm:"column(vlan_id)" description:"网络 vlanID" json:"vlan_id"`
	ShopMode              int       `orm:"column(shop_mode)" description:"门店模式，无盘平台用 1=1.0方案2=2.0方案" json:"shop_mode"`
	CreateBy              string    `orm:"column(create_by)" description:"创建者" json:"create_by"`
	UpdateBy              string    `orm:"column(update_by)" description:"修改者" json:"update_by"`
	Remark                string    `orm:"column(remark)" description:"备注" json:"remark"`
	State                 int       `orm:"column(state)" description:"状态: 0-无效, 1-有效, 2-删除, 3-暂停合作" json:"state"`
	//CoID                  int       `orm:"column(co_id)" description:"主体id: 0-云天, 1-瞳感" json:"co_id"`
	CreateTime time.Time `orm:"column(create_time);type(datetime);auto_now_add" json:"create_time"`
	UpdateTime time.Time `orm:"column(update_time);type(datetime);auto_now_add" json:"update_time"`
	ModifyTime time.Time `orm:"column(modify_time);type(datetime);auto_now" json:"modify_time"`

	SpecInstances []TEsportsSpecInstancesInfo `json:"spec_instances" orm:"-"`
	ResourcePools []TEsportsResourcePoolInfo  `json:"resource_pools" orm:"-"`
}

type TEsportsUserInfoList struct {
	Total int64              `json:"total"`
	List  []TEsportsUserInfo `json:"list"`
}

// 查询租户列表
type EsportUserInfoListResp struct {
	Code int                  `json:"code"`
	Msg  string               `json:"msg"`
	Data TEsportsUserInfoList `json:"data"`
}

// 单个查询租户
type EsportGettUserInfoResp struct {
	Code int                    `json:"code"`
	Msg  string                 `json:"msg"`
	Data EsportUserInfoListResp `json:"data"`
}

// 查询租户的列表
func (s *SaasServerService) GetEsportUserInfoList(ctx context.Context, sessionId string, req interface{}) (bizInfos TEsportsUserInfoList, err error) {

	url := fmt.Sprintf("%s/user/getUserInfoList", s.Host)

	resp := new(EsportUserInfoListResp)

	_, err = helper.HttpPost(ctx, sessionId, url, s.Timeout, req, resp)
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] helper.HttpPost err  host[%s] err:%+v", sessionId, url, err)
		return bizInfos, errors.New("查询 biz 信息 http 错误")
	}

	if resp.Code != 0 {
		logx.WithContext(ctx).Errorf("[%s] GetEsportUserInfoList httpc.Do err host[%s] err:%+v", sessionId, url, err)
		return bizInfos, errors.New(resp.Msg)
	}
	logx.WithContext(ctx).Infof("[%s] GetEsportUserInfoList  resp:%+v", sessionId, helper.ToJSON(resp))

	bizInfos = resp.Data
	return
}
