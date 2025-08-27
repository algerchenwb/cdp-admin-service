package saas

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpc"
)

type GetPrimaryInfoResp struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data PrimaryInfoBody `json:"data"`
}
type PrimaryInfoBody struct {
	TEsportsPrimaryInfo
	SpecInstances []TEsportsSpecInstancesInfo `json:"spec_instances"`
}
type UpdatePrimaryInfoReq struct {
	Name                  string    `json:"name"`
	AuthoritySet          string    `json:"authority_set"`
	TemplateID            int64     `json:"template_id"`
	TemplatePriority      int64     `json:"template_priority"`
	TotalInstances        int64     `json:"-"`
	ExpectedEffectiveDate time.Time `json:"expected_effective_date"`
	UpdateBy              string    `json:"update_by"`
	Remark                string    `json:"remark"`
	State                 int       `json:"state"`
	Email                 string    `json:"email"`

	SpecInstances []EditSpecInstanceInfoRequest `json:"spec_instances"`
}

type UpdatePrimaryInfoResp struct {
	Code int                 `json:"code"`
	Msg  string              `json:"msg"`
	Data TEsportsPrimaryInfo `json:"data"`
}
type TEsportsPrimaryInfo struct {
	Id                    int64     `orm:"column(primary_id);auto" description:"自增ID" json:"primary_id"`
	AuthoritySet          string    `orm:"column(authority_set)" description:"权限集合" json:"authority_set"`
	Name                  string    `orm:"column(name)" description:"代理商名称" json:"name"`
	WallpaperID           int64     `orm:"column(wallpaper_id)" description:"壁纸ID" json:"wallpaper_id"`
	TemplateID            int64     `orm:"column(template_id)" description:"模板ID" json:"template_id"`
	TemplatePriority      int64     `orm:"column(template_priority)" description:"模板优先级" json:"template_priority"`
	TotalInstances        int64     `orm:"column(total_instances)" description:"实例数" json:"total_instances"`
	EffectiveDate         time.Time `orm:"column(effective_date);type(datetime)" description:"实际生效日期" json:"effective_date"`
	ExpectedEffectiveDate time.Time `orm:"column(expected_effective_date);type(datetime)" description:"期待生效日期" json:"expected_effective_date"`
	BelongID              int64     `orm:"column(belong_id)" description:"归属ID" json:"belong_id"`
	CreateBy              string    `orm:"column(create_by)" description:"创建者" json:"create_by"`
	UpdateBy              string    `orm:"column(update_by)" description:"修改者" json:"update_by"`
	Remark                string    `orm:"column(remark)" description:"备注" json:"remark"`
	State                 int       `orm:"column(state)" description:"状态: 0-无效, 1-有效, 2-删除, 3-暂停合作" json:"state"`
	Email                 string    `orm:"column(email)" description:"邮箱" json:"email"`
	Tag                   int32     `orm:"column(tag)" description:"标签 0-无 1-测试 2-POC 3-建设中 4-商用" json:"tag"`
	RegionCode            string    `orm:"column(region_code)" description:"行政编码" json:"region_code"`
	AreaTypes             string    `orm:"column(area_types)" description:"区域类型" json:"area_types"`
	CreateTime            time.Time `orm:"column(create_time);type(datetime);auto_now_add" json:"create_time"`
	UpdateTime            time.Time `orm:"column(update_time);type(datetime);auto_now_add" json:"update_time"`
	ModifyTime            time.Time `orm:"column(modify_time);type(datetime);auto_now" json:"modify_time"`

	SpecInstances []TEsportsSpecInstancesInfo `json:"spec_instances" orm:"-"`

	PlatformShareNum           uint64 `json:"platform_share_num" orm:"platform_share_num"`
	PlatformShardModel         int32  `json:"platform_shard_model" orm:"platform_shard_model"`
	BulletinBoardAuthorities   int64  `orm:"column(bulletin_board_authorities)" description:"公告板权限" json:"bulletin_board_authorities"`
	BulletinBoardAuthorityList []int  `json:"bulletin_board_authority_list" orm:"-"`
}

func (s *SaasServerService) GetPrimaryInfo(ctx context.Context, sessionId string, primaryId int64) (info PrimaryInfoBody, err error) {

	url := fmt.Sprintf("%s/primary/getPrimaryInfo/%d", s.Host, primaryId)

	logx.WithContext(ctx).Debugf("sessionId:%s, httpc.Do url:%s, req:%+v", sessionId, url, nil)
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logx.WithContext(ctx).Errorf("sessionId:%s, httpc.Do err  host[%s] err:%+v", sessionId, url, err)
		err = errors.New("获取一级代理商信息失败")
		return
	}
	r.Header.Set("X-System", "diskless-aggregator")
	httpResp, err := httpc.DoRequest(r)
	logx.WithContext(ctx).Debugf("sessionId:%s, httpc.Do url:%s, resp:%+v", sessionId, url, httpResp)
	if err != nil {
		logx.WithContext(ctx).Errorf("sessionId:%s, httpc.Do err  host[%s] err:%+v", sessionId, url, err)
		err = errors.New("获取一级代理商信息失败")
		return
	}
	if httpResp.StatusCode != 200 {
		logx.WithContext(ctx).Errorf("sessionId:%s, httpc.Do err url:%s, resp:%+v", sessionId, url, httpResp)
		err = errors.New("获取一级代理商信息失败")
		return
	}

	resp := new(GetPrimaryInfoResp)
	body, _ := io.ReadAll(httpResp.Body)

	if err = json.Unmarshal(body, resp); err != nil {
		logx.WithContext(ctx).Errorf("sessionId:%s, json.Unmarshal err host[%s], err:%+v", sessionId, url, err)
		err = errors.New("反序列化失败")
		return
	}
	if resp.Code != 0 {
		logx.WithContext(ctx).Errorf("sessionId:%s, saasResp code err  host[%s], saasResp:%+v", sessionId, url, resp)
		err = errors.New(resp.Msg)
		return
	}

	logx.WithContext(ctx).Debugf("sessionId:%s, httpc.Do url:%s, saasResp:%+v", sessionId, url, resp)
	info = resp.Data
	return
}

func (s *SaasServerService) UpdatePrimaryInfo(ctx context.Context, sessionId string, info TEsportsPrimaryInfo, specInstances []EditSpecInstanceInfoRequest) (err error) {

	url := fmt.Sprintf("%s/v3/primary/updatePrimaryInfo/%d", s.Host, info.Id)

	logx.WithContext(ctx).Debugf("sessionId:%s, httpc.Do url:%s, req:%+v", sessionId, url, nil)
	req := &UpdatePrimaryInfoReq{}
	copier.Copy(&req, info)
	req.SpecInstances = specInstances
	b, _ := json.Marshal(req)
	r, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(b))
	if err != nil {
		logx.WithContext(ctx).Errorf("sessionId:%s, httpc.Do err  host[%s] err:%+v", sessionId, url, err)
		err = errors.New("更新一级代理商信息失败")
		return
	}
	r.Header.Set("X-System", "diskless-aggregator")
	httpResp, err := httpc.DoRequest(r)
	logx.WithContext(ctx).Debugf("sessionId:%s, httpc.Do url:%s, resp:%+v", sessionId, url, httpResp)
	if err != nil {
		logx.WithContext(ctx).Errorf("sessionId:%s, httpc.Do err  host[%s] err:%+v", sessionId, url, err)
		err = errors.New("更新一级代理商信息失败")
		return
	}
	if httpResp.StatusCode != 200 {
		logx.WithContext(ctx).Errorf("sessionId:%s, httpc.Do err url:%s, resp:%+v", sessionId, url, httpResp)
		err = errors.New("更新一级代理商信息失败")
		return
	}

	resp := new(GetPrimaryInfoResp)
	body, _ := io.ReadAll(httpResp.Body)

	if err = json.Unmarshal(body, resp); err != nil {
		logx.WithContext(ctx).Errorf("sessionId:%s, json.Unmarshal err host[%s], err:%+v", sessionId, url, err)
		err = errors.New("反序列化失败")
		return
	}
	if resp.Code != 0 {
		logx.WithContext(ctx).Errorf("sessionId:%s, saasResp code err  host[%s], saasResp:%+v", sessionId, url, resp)
		err = errors.New(resp.Msg)
		return
	}

	logx.WithContext(ctx).Debugf("sessionId:%s, httpc.Do url:%s, saasResp:%+v", sessionId, url, resp)
	return
}
