package saas

import (
	"cdp-admin-service/internal/helper"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAreaConfigInfoListReq struct {
	Offset   int64    `json:"offset"`
	Limit    int64    `json:"limit"`
	CondList []string `json:"cond_list"`
	Sorts    string   `json:"sorts"`
	Orders   string   `json:"orders"`
}
type GetAreaConfigInfoListResp struct {
	Code int                       `json:"code"`
	Msg  string                    `json:"msg"`
	Data GetAreaConfigInfoListBody `json:"data"`
}
type GetAreaConfigInfoListBody struct {
	Total int64            `json:"total"`
	List  []AreaConfigInfo `json:"list"`
}
type AreaConfigInfo struct {
	TEsportsAreaConfigInfo
	OuterSpecInfoList []AreaOuterSpecInfo `json:"spec_info_list"`
	SetDurationList   []AreaSetDuration   `json:"set_duration_list"`
}

type TEsportsAreaConfigInfo struct {
	Acid uint64 `orm:"column(acid);auto" description:"区域外部信息自增id" json:"acid"`

	AreaType int64 `orm:"column(area_type)" description:"区域信息" json:"area_type"`

	Name string `orm:"column(name)" description:"名字" json:"name"`

	Desc string `orm:"column(desc)" description:"描述" json:"desc"`

	SpecIdList string `orm:"column(spec_id_list)" description:"规格ID列表(以|分割)" json:"spec_id_list"`

	TotalInstancesList string `orm:"column(total_instances_list)" description:"资源总数列表(以|分割)" json:"total_instances_list"`

	LimitDurationConfig string `orm:"column(limit_duration_config)" description:"业务调度设置" json:"limit_duration_config"`

	CreateBy string `orm:"column(create_by)" description:"" json:"create_by"`

	UpdateBy string `orm:"column(update_by)" description:"" json:"update_by"`

	Remark string `orm:"column(remark)" description:"" json:"remark"`

	State int32 `orm:"column(state)" description:"状态,0-未知;1-有效" json:"state"`

	CreateTime time.Time `orm:"column(create_time)" description:"" json:"create_time"`

	UpdateTime time.Time `orm:"column(update_time)" description:"" json:"update_time"`

	ModifyTime time.Time `orm:"column(modify_time)" description:"" json:"modify_time"`
}

type AreaOuterSpecInfo struct {
	Id             int64  `json:"outer_spec_id"`
	Name           string `json:"name"`
	InnerSpecID    int64  `json:"inner_spec_id"`
	TotalInstances int64  `json:"total_instances"`
	UsedInstances  int64  `json:"used_instances"`
}

type AreaSetDuration struct {
	SetID            int64          `json:"set_id"`
	WeekDurationList []WeekDuration `json:"week_duration_list"`
}
type WeekDuration struct {
	WeekDay   string `json:"week_day"`
	HourBegin string `json:"hour_begin"`
	HourEnd   string `json:"hour_end"`
}
type CreateAreaConfigInfoReq struct {
	Name               string            `json:"name"`
	Desc               string            `json:"desc"`
	AreaType           int32             `json:"area_type"`
	SpecIDList         string            `json:"spec_id_list"`
	TotalInstancesList string            `json:"total_instances_list"`
	SetDurationList    []AreaSetDuration `json:"set_duration_list"`
	CreateBy           string            `json:"create_by"`
}

type CreateAreaConfigInfoResp struct {
	Code int            `json:"code"`
	Msg  string         `json:"msg"`
	Data AreaConfigInfo `json:"data"`
}

type UpdateAreaConfigInfoReq struct {
	Name               string `json:"name"`
	Desc               string `json:"desc"`
	SpecIDList         string `json:"spec_id_list"`
	TotalInstancesList string `json:"total_instances_list"`
	UpdateBy           string `json:"update_by"`
}

type UpdateAreaConfigInfoResp struct {
	Code int            `json:"code"`
	Msg  string         `json:"msg"`
	Data AreaConfigInfo `json:"data"`
}

// 查询区域信息
func (s *SaasServerService) GetAreaConfigInfoList(ctx context.Context, sessionId string, areaType uint64) (total int64, list []AreaConfigInfo, err error) {

	url := fmt.Sprintf("%s/area_config/getAreaConfigInfoList", s.Host)

	req := GetOuterSpecificationInfoListReq{
		Limit:    1,
		CondList: []string{fmt.Sprintf("area_type:%d", areaType)},
	}
	resp := new(GetAreaConfigInfoListResp)
	_, err = helper.HttpPost(ctx, sessionId, url, s.Timeout, req, resp)
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] GetAreaConfigInfoList post err  host[%s] err:%+v", sessionId, url, err)
		return 0, nil, errors.New("获取区域信息失败")
	}
	if resp.Code != 0 {
		logx.WithContext(ctx).Errorf("[%s] url[%s] resp.code err. saasResp:%+v", sessionId, url, resp)
		return 0, nil, errors.New(resp.Msg)
	}

	logx.WithContext(ctx).Debugf("[%s] GetAreaConfigInfoList success url:%s, resp:%+v", sessionId, url, resp)

	list = resp.Data.List
	total = resp.Data.Total
	return
}

// 创建区域信息
func (s *SaasServerService) CreateAreaConfigInfo(ctx context.Context, sessionId string, name, desc, specIdList, totalInstancesList, createBy string, areaType int32) (info AreaConfigInfo, err error) {

	url := fmt.Sprintf("%s/area_config/createAreaConfigInfo", s.Host)

	req := &CreateAreaConfigInfoReq{
		Name:               name,
		Desc:               desc,
		SpecIDList:         specIdList,
		TotalInstancesList: totalInstancesList,
		CreateBy:           createBy,
		AreaType:           areaType,
		SetDurationList:    nil,
	}
	resp := new(CreateAreaConfigInfoResp)

	_, err = helper.HttpPost(ctx, sessionId, url, s.Timeout, req, resp)
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] CreateAreaConfigInfo post err  host[%s] err:%+v", sessionId, url, err)
		return info, errors.New("创建区域信息失败")
	}
	if resp.Code != 0 {
		logx.WithContext(ctx).Errorf("[%s] url[%s] resp.code err. saasResp:%+v", sessionId, url, resp)
		return info, errors.New(resp.Msg)
	}

	logx.WithContext(ctx).Debugf("[%s] CreateAreaConfigInfo success url:%s, resp:%+v", sessionId, url, resp)
	info = resp.Data
	return
}

// 创建区域信息
func (s *SaasServerService) UpdateAreaConfigInfo(ctx context.Context, sessionId string, name, desc, specIdList, totalInstancesList, updateBy string, id uint64) (info AreaConfigInfo, err error) {

	url := fmt.Sprintf("%s/area_config/updateAreaConfigInfo/%d", s.Host, id)

	req := &UpdateAreaConfigInfoReq{
		Name:               name,
		Desc:               desc,
		SpecIDList:         specIdList,
		TotalInstancesList: totalInstancesList,
		UpdateBy:           updateBy,
	}
	resp := new(UpdateAreaConfigInfoResp)

	_, err = helper.HttpPut(ctx, sessionId, url, s.Timeout, req, resp)
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] UpdateAreaConfigInfo post err  host[%s] err:%+v", sessionId, url, err)
		return info, errors.New("更新区域信息失败")
	}
	if resp.Code != 0 {
		logx.WithContext(ctx).Errorf("[%s] url[%s] resp.code err. saasResp:%+v", sessionId, url, resp)
		return info, errors.New(resp.Msg)
	}

	logx.WithContext(ctx).Debugf("[%s] UpdateAreaConfigInfo success url:%s, resp:%+v", sessionId, url, resp)
	info = resp.Data
	return
}
