package saas

import (
	"cdp-admin-service/internal/helper"
	"cdp-admin-service/internal/svc"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResourcePoolInfoListReq struct {
	Offset   int64    `json:"offset"`
	Limit    int64    `json:"limit"`
	CondList []string `json:"cond_list"`
	Sorts    string   `json:"sorts"`
	Orders   string   `json:"orders"`
}

// {"total":2,"list":[{"pool_id":648,"outer_spec_id":26,"config_id":"372","number":315360000,"used_number":0,"type":2,"validity_period":"2026-04-27T20:30:14+08:00","primary_id":68,"agent_id":162,"biz_id":1018339,"create_by":"jeff","update_by":"jeff","remark":"","state":1,"create_time":"2025-04-27T20:30:18+08:00","update_time":"2025-04-27T20:30:18+08:00","modify_time":"2025-04-27T20:30:18+08:00","expired":false},{"pool_id":649,"outer_spec_id":26,"config_id":"372","number":10,"used_number":0,"type":1,"validity_period":"2026-04-27T20:30:14+08:00","primary_id":68,"agent_id":162,"biz_id":1018339,"create_by":"jeff","update_by":"jeff","remark":"","state":1,"create_time":"2025-04-27T20:30:18+08:00","update_time":"2025-04-27T20:30:18+08:00","modify_time":"2025-04-27T20:30:18+08:00","expired":false}
type TEsportsResourcePoolInfo struct {
	PoolID         int       `json:"pool_id"`
	OuterSpecID    int       `json:"outer_spec_id"`
	ConfigID       string    `json:"config_id"`
	Number         int       `json:"number"`
	UsedNumber     int       `json:"used_number"`
	Type           int       `json:"type"`
	ValidityPeriod time.Time `json:"validity_period"`
	PrimaryID      int       `json:"primary_id"`
	AgentID        int       `json:"agent_id"`
	BizID          int       `json:"biz_id"`
	CreateBy       string    `json:"create_by"`
	UpdateBy       string    `json:"update_by"`
	Remark         string    `json:"remark"`
	State          int       `json:"state"`
	CreateTime     time.Time `json:"create_time"`
	UpdateTime     time.Time `json:"update_time"`
	ModifyTime     time.Time `json:"modify_time"`
	Expired        bool      `json:"expired"`
}

const (
	ResourcePoolTypeHourly = 2
	ResourcePoolTypePerUse = 1
)

// {"code":0,"data":{"total":2,"list":[{"pool_id":652,"outer_spec_id":26,"config_id":"372","number":315360000,"used_number":0,"type":2,"validity_period":"2026-04-27T20:37:37+08:00","primary_id":68,"agent_id":162,"biz_id":1018341,"create_by":"jeff","update_by":"jeff","remark":"","state":1,"create_time":"2025-04-27T20:37:42+08:00","update_time":"2025-04-27T20:37:42+08:00","modify_time":"2025-04-27T20:37:42+08:00","expired":false},{"pool_id":653,"outer_spec_id":26,"config_id":"372","number":10,"used_number":0,"type":1,"validity_period":"2026-04-27T20:37:37+08:00","primary_id":68,"agent_id":162,"biz_id":1018341,"create_by":"jeff","update_by":"jeff","remark":"","state":1,"create_time":"2025-04-27T20:37:42+08:00","update_time":"2025-04-27T20:37:42+08:00","modify_time":"2025-04-27T20:37:42+08:00","expired":false}]},"msg":"成功"}
type ResourcePoolInfoList struct {
	Code int                      `json:"code"`
	Msg  string                   `json:"msg"`
	Data ResourcePoolInfoListBody `json:"data"`
}
type ResourcePoolInfoListBody struct {
	Total int64                      `json:"total"`
	List  []TEsportsResourcePoolInfo `json:"list"`
}

type EditResourcePoolInfoRequest struct {
	OuterSpecID    int64     `json:"outer_spec_id"`
	ConfigID       string    `json:"config_id"`
	Type           int       `json:"type"`
	Number         int64     `json:"number"`
	ValidityPeriod time.Time `json:"validity_period"`
}

type EditSpecInstanceInfoRequest struct {
	OuterSpecID    int64 `json:"outer_spec_id"`
	TotalInstances int64 `json:"total_instances"`
}

type UpdateResourcePoolByAgentReq struct {
	AgentID       int64                         `json:"agent_id"`
	AreaType      int64                         `json:"area_type"`
	UpdateBy      string                        `json:"update_by"`
	ResourcePools []EditResourcePoolInfoRequest `json:"resource_pools"`
	SpecInstances []EditSpecInstanceInfoRequest `json:"spec_instances"`
}

type UpdateResourcePoolByAgentResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type UpdateResourcePoolByUserReq struct {
	UserID        int64                         `json:"user_id"`
	UpdateBy      string                        `json:"update_by"`
	ResourcePools []EditResourcePoolInfoRequest `json:"resource_pools"`
	SpecInstances []EditSpecInstanceInfoRequest `json:"spec_instances"`
}
type UpdateResourcePoolByUserResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (s *SaasServerService) GetResourcePoolInfoList(ctx context.Context, sessionId string,
	condList []string) (total int64, list []TEsportsResourcePoolInfo, err error) {

	url := fmt.Sprintf("%s/resource_pool_info/getResourcePoolInfoList", s.Host)
	req := &ResourcePoolInfoListReq{
		Limit:    100000,
		CondList: condList,
	}
	resp := new(ResourcePoolInfoList)

	_, err = helper.HttpPost(ctx, sessionId, url, s.Timeout, req, resp)
	if err != nil {
		logx.WithContext(ctx).Errorf("helper.HttpPost err  host[%s] err:%+v", url, err)
		return 0, nil, errors.New("查询资源池失败")
	}

	logx.WithContext(ctx).Debugf("[%s] GetResourcePoolInfoList success url:%s, req:%s resp:%+v", sessionId, url, helper.ToJSON(req), helper.ToJSON(resp))
	total = resp.Data.Total
	list = resp.Data.List
	return
}

// 更新二代的资源池信息
func (s *SaasServerService) UpdateResourcePoolByAgent(ctx context.Context, sessionId string, agentId int64, areaType uint64, updateBy string, resourcePools []EditResourcePoolInfoRequest, specInstances []EditSpecInstanceInfoRequest) (err error) {

	url := fmt.Sprintf("%s/resource_pool_info/updateResourcePoolByAgent", s.Host)
	req := &UpdateResourcePoolByAgentReq{
		AgentID:       agentId,
		AreaType:      int64(areaType),
		UpdateBy:      updateBy,
		ResourcePools: resourcePools,
		SpecInstances: specInstances,
	}

	resp := new(UpdateResourcePoolByAgentResp)

	_, err = helper.HttpPost(ctx, sessionId, url, s.Timeout, req, resp)
	if err != nil {
		logx.WithContext(ctx).Errorf("UpdateResourcePoolByAgent helper.HttpPost err  host[%s] err:%+v", url, err)
		return errors.New("更新二代资源池失败")
	}
	if resp.Code != 0 {
		logx.WithContext(ctx).Errorf("[%s] UpdateResourcePoolByAgent url[%s] resp.code err. saasResp:%+v", sessionId, url, resp)
		return errors.New(resp.Msg)
	}

	logx.WithContext(ctx).Debugf("[%s] UpdateResourcePoolByAgent success url:%s, req:%s resp:%+v", sessionId, url, helper.ToJSON(req), helper.ToJSON(resp))

	return
}

// 更的资源池信息
func (s *SaasServerService) UpdateResourcePoolByUser(ctx context.Context, sessionId string, BizId int64, updateBy string, resourcePools []EditResourcePoolInfoRequest, specInstances []EditSpecInstanceInfoRequest) (err error) {

	url := fmt.Sprintf("%s/resource_pool_info/updateResourcePoolByUser", s.Host)
	req := &UpdateResourcePoolByUserReq{
		UserID:        BizId,
		UpdateBy:      updateBy,
		ResourcePools: resourcePools,
		SpecInstances: specInstances,
	}

	resp := new(UpdateResourcePoolByUserResp)

	_, err = helper.HttpPost(ctx, sessionId, url, s.Timeout, req, resp)
	if err != nil {
		logx.WithContext(ctx).Errorf("UpdateResourcePoolByUser helper.HttpPost err  host[%s] err:%+v", url, err)
		return errors.New("更新租户资源池失败")
	}
	if resp.Code != 0 {
		logx.WithContext(ctx).Errorf("[%s] UpdateResourcePoolByUser url[%s] resp.code err. saasResp:%+v", sessionId, url, resp)
		return errors.New(resp.Msg)
	}

	logx.WithContext(ctx).Debugf("[%s] UpdateResourcePoolByUser success url:%s, req:%s resp:%+v", sessionId, url, helper.ToJSON(req), helper.ToJSON(resp))

	return
}

func (s *SaasServerService) GetBizPoolInfo(ctx context.Context, svcCtx *svc.ServiceContext, sessionId string, bizId int64) (mapConfigId map[int]TEsportsResourcePoolInfo, err error) {

	var condList []string
	condList = append(condList, fmt.Sprintf("biz_id__eq:%d", bizId))
	condList = append(condList, "state__eq:1")
	_, resPoolInfoList, err := s.GetResourcePoolInfoList(ctx, sessionId, condList)
	if err != nil {
		logx.WithContext(ctx).Errorf("GetEsportUserInfoList bizId[%d] err:%+v", bizId, err)
		return nil, err
	}

	if len(resPoolInfoList) == 0 {
		return nil, errors.New("门店资源池为空")
	}

	for _, rescPool := range resPoolInfoList {
		mapConfigId[rescPool.OuterSpecID] = rescPool
	}

	return
}
