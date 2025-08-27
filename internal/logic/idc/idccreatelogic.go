package idc

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"cdp-admin-service/internal/helper"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/helper/saas"

	"github.com/zeromicro/go-zero/core/logx"
	"gitlab.vrviu.com/cloud_esport_backend/saas_user_server/pb/saas_user"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
)

type IdcCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIdcCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IdcCreateLogic {
	return &IdcCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 一代规格 ，新增区域， 区域绑定规格
func (l *IdcCreateLogic) IdcCreate(req *types.CreateIdcReq) (resp *types.CreateIdcResp, err error) {
	sessionId := helper.GetSessionId(l.ctx)

	primary, err := saas.NewSaasServerService(l.svcCtx.Config.Saas.ServerHost, l.svcCtx.Config.Saas.Timeout).
		GetPrimaryInfo(l.ctx, sessionId, int64(req.PrimaryId))
	if err != nil {
		l.Logger.Errorf("sessionId: %s, create agent info error: %v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.GetPrimaryInfoErrorCode)
	}

	rpcPrimary, err := l.svcCtx.PrimaryRpc.GetPrimaryInfo(l.ctx, &saas_user.GetPrimaryInfoReq{
		SessionId: sessionId,
		PrimaryId: int64(req.PrimaryId),
	})
	if err != nil {
		l.Logger.Errorf("sessionId: %s, create agent info error: %v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.GetPrimaryInfoErrorCode)
	}
	editSpecInstances := make([]*saas_user.EditSpecInstance, 0)
	for _, specInstance := range primary.SpecInstances {
		editSpecInstances = append(editSpecInstances, &saas_user.EditSpecInstance{
			OuterSpecId:    int64(specInstance.OuterSpecID),
			TotalInstances: int64(specInstance.TotalInstances),
		})
	}
	var areaTypes string
	if rpcPrimary.GetPrimaryInfo().GetAreaTypes() == "" {
		areaTypes = fmt.Sprint(req.AreaId)
	} else {
		areaTypes = strings.Join(helper.ArrayUniqueValue(append(strings.Split(rpcPrimary.GetPrimaryInfo().GetAreaTypes(), ","), fmt.Sprint(req.AreaId))), ",")
	}
	_, err = l.svcCtx.PrimaryRpc.UpdatePrimaryInfo(l.ctx, &saas_user.UpdatePrimaryInfoReq{
		SessionId:                sessionId,
		PrimaryId:                int64(primary.Id),
		Name:                     primary.Name,
		TemplateId:               int64(primary.TemplateID),
		TotalInstances:           int64(primary.TotalInstances),
		ExpectedEffectiveDate:    primary.ExpectedEffectiveDate.Format("2006-01-02T15:04:05-07:00"),
		Remark:                   primary.Remark,
		AuthoritySet:             primary.AuthoritySet,
		Email:                    primary.Email,
		BelongId:                 int64(primary.BelongID),
		TemplatePriority:         int64(primary.TemplatePriority),
		SpecInstances:            editSpecInstances,
		UpdateBy:                 primary.UpdateBy,
		State:                    int32(primary.State),
		PlatformShareNum:         primary.PlatformShareNum,
		PlatformShardModel:       primary.PlatformShardModel,
		BulletinBoardAuthorities: primary.BulletinBoardAuthorities,
		Tag:                      int32(primary.Tag),
		AreaTypes:                &areaTypes,
	})
	if err != nil {
		l.Logger.Errorf("sessionId: %s, update primary info error: %v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.UpdatePrimaryInfoErrorCode)
	}

	//  查询区域与区域配置信息
	_, areaConfigs, err := saas.NewSaasServerService(l.svcCtx.Config.Saas.ServerHost, l.svcCtx.Config.Saas.Timeout).
		GetAreaConfigInfoList(l.ctx, sessionId, uint64(req.AreaId))
	if err != nil {
		l.Logger.Errorf("sessionId: %s, get area config info list error: %v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.GetAreaConfigInfoListErrorCode)
	}
	areaConfig, err := l.configAreaConfig(sessionId, areaConfigs, primary, req.AreaId)
	if err != nil {
		l.Logger.Errorf("sessionId: %s, config area config info error: %v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.ConfigAreaConfigInfoErrorCode)
	}

	// 创建二代
	agentResp, err := l.createAgent(sessionId, areaConfig, primary, req.Name)
	if err != nil {
		l.Logger.Errorf("sessionId: %s, create agent info error: %v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.CreateAgentInfoErrorCode)
	}
	// 创建默认算力池
	err = l.insertDefaultResourcePool(sessionId, int32(req.AreaId))
	if err != nil {
		l.Logger.Errorf("sessionId: %s, insert default resource pool error: %v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.InsertDefaultResourcePoolErrorCode)
	}

	areaInfo, err := l.insertAreaInfo(sessionId, &table.TCdpAreaInfo{
		AreaId:            int32(req.AreaId),
		PrimaryId:         int64(req.PrimaryId),
		AgentId:           int64(agentResp.AgentInfo.AgentId),
		Name:              req.Name,
		RegionId:          int32(req.RegionId),
		DeploymentType:    int32(req.DeploymentType),
		SchemaConfig:      req.SchemaConfig,
		ResetSchemaConfig: req.ResetSchemaConfig,
		Status:            1,
		Remark:            req.Remark,
		CreateTime:        time.Now(),
		UpdateTime:        time.Now(),
		ModifyTime:        time.Now(),
	})
	if err != nil {
		l.Logger.Errorf("sessionId: %s, insert area info error: %v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.InsertAreaInfoErrorCode)
	}
	// 更新超管账号节点
	err = l.updateSuperAdminNode(sessionId, int32(req.AreaId))
	if err != nil {
		l.Logger.Errorf("sessionId: %s, update super admin node error: %v", sessionId, err)
		return nil, errorx.NewDefaultCodeError("更新超管账号失败")
	}

	return &types.CreateIdcResp{
		Id:      int(areaInfo.Id),
		AgentId: int(agentResp.AgentInfo.AgentId),
	}, nil
}

// areaConfig 区域配置信息
func (l *IdcCreateLogic) configAreaConfig(sessionId string, areaConfigs []saas.AreaConfigInfo, primary saas.PrimaryInfoBody, areaId int32) (areaConfigInfo saas.AreaConfigInfo, err error) {
	// 如果区域配置为空，报错
	if len(areaConfigs) == 0 {
		return areaConfigInfo, errors.New("area config is empty")
	}
	var specMap = make(map[int64]int64)
	for _, specInstance := range areaConfigs[0].OuterSpecInfoList {
		specMap[specInstance.Id] += specInstance.TotalInstances
	}
	for _, specInstance := range primary.SpecInstances {
		specMap[specInstance.OuterSpecID] += specInstance.TotalInstances
	}
	var specIdList string
	var totalInstancesList string
	for specId, totalInstances := range specMap {
		specIdList += fmt.Sprintf("|%d", specId)
		totalInstancesList += fmt.Sprintf("|%d", totalInstances)
	}
	if len(specIdList) > 0 {
		specIdList = specIdList + "|"
		totalInstancesList = totalInstancesList + "|"
	}
	areaConfigInfo, err = saas.NewSaasServerService(l.svcCtx.Config.Saas.ServerHost, l.svcCtx.Config.Saas.Timeout).
		UpdateAreaConfigInfo(l.ctx, sessionId, areaConfigs[0].Name, areaConfigs[0].Remark, specIdList, totalInstancesList, "skyhash", uint64(areaConfigs[0].TEsportsAreaConfigInfo.Acid))
	if err != nil {
		l.Logger.Errorf("sessionId: %s, update area config info error: %v", sessionId, err)
		return areaConfigInfo, err
	}
	return areaConfigInfo, nil

}

// 二代新增
func (l *IdcCreateLogic) createAgent(sessionId string, areaConfig saas.AreaConfigInfo, primary saas.PrimaryInfoBody, name string) (agentResp *saas_user.CreateAgentInfoResp, err error) {
	agentResp, err = l.svcCtx.AgentRpc.CreateAgentInfo(l.ctx, &saas_user.CreateAgentInfoReq{
		SessionId:                sessionId,
		AuthoritySet:             "",
		AreaType:                 int64(areaConfig.AreaType),
		PrimaryId:                int64(primary.Id),
		Name:                     name,
		TemplateId:               1,
		TemplatePriority:         2,
		TotalInstances:           0,
		ExpectedEffectiveDate:    time.Now().Format("2006-01-02T15:04:05-07:00"),
		CreateBy:                 "",
		Remark:                   "",
		Email:                    "",
		PrimarySharePercent:      0,
		BulletinBoardAuthorities: 254,
		Tag:                      0,
		RegionCode:               "110101",
	})
	if err != nil {
		l.Logger.Errorf("sessionId: %s, create agent info error: %v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.CreateAgentInfoErrorCode)
	}
	_, err = l.svcCtx.AgentRpc.UpdateAgentInfo(l.ctx, &saas_user.UpdateAgentInfoReq{
		SessionId:                sessionId,
		AgentId:                  agentResp.AgentInfo.AgentId,
		AuthoritySet:             "",
		AreaType:                 int64(areaConfig.AreaType),
		Name:                     name,
		TemplateId:               1,
		TemplatePriority:         2,
		TotalInstances:           0,
		ExpectedEffectiveDate:    time.Now().Format("2006-01-02T15:04:05-07:00"),
		UpdateBy:                 "",
		Remark:                   "",
		State:                    1, // 生效
		Email:                    "",
		PrimarySharePercent:      0,
		BulletinBoardAuthorities: 254,
		Tag:                      0,
	})
	if err != nil {
		l.Logger.Errorf("sessionId: %s, update agent info error: %v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.UpdateAgentInfoErrorCode)
	}
	return
}

// 更新二代资源池
func (l *IdcCreateLogic) updateResourcePool(sessionId string, areaConfig saas.AreaConfigInfo, primary saas.PrimaryInfoBody, agentId int64) (err error) {
	var editSpecInstances []saas.EditSpecInstanceInfoRequest
	var editResourcePools []saas.EditResourcePoolInfoRequest

	for _, specInstance := range primary.SpecInstances {
		editSpecInstances = append(editSpecInstances, saas.EditSpecInstanceInfoRequest{
			OuterSpecID:    specInstance.OuterSpecID,
			TotalInstances: l.svcCtx.Config.Saas.Specification.DefaultTotalInstances,
		})
		editResourcePools = append(editResourcePools, saas.EditResourcePoolInfoRequest{
			OuterSpecID:    specInstance.OuterSpecID,
			Number:         l.svcCtx.Config.Saas.Specification.DefaultTimeResourcePoolNum,
			ValidityPeriod: time.Now().AddDate(l.svcCtx.Config.Saas.Specification.DefaultValidityPeriodYears, 0, 0),
			ConfigID:       "",
			Type:           1,
		})
		editResourcePools = append(editResourcePools, saas.EditResourcePoolInfoRequest{
			OuterSpecID:    specInstance.OuterSpecID,
			Number:         l.svcCtx.Config.Saas.Specification.DefaultFrequencyResourcePoolNum,
			ValidityPeriod: time.Now().AddDate(l.svcCtx.Config.Saas.Specification.DefaultValidityPeriodYears, 0, 0),
			ConfigID:       "",
			Type:           2,
		})
	}

	err = saas.NewSaasServerService(l.svcCtx.Config.Saas.ServerHost, l.svcCtx.Config.Saas.Timeout).
		UpdateResourcePoolByAgent(l.ctx, sessionId, agentId, uint64(areaConfig.AreaType), areaConfig.Name, editResourcePools, editSpecInstances)
	if err != nil {
		l.Logger.Errorf("sessionId: %s, update resource pool error: %v", sessionId, err)
		return err
	}
	return nil
}

// 更新cdp area 信息
func (l *IdcCreateLogic) insertAreaInfo(sessionId string, areaInfo *table.TCdpAreaInfo) (area *table.TCdpAreaInfo, err error) {
	area, _, err = table.T_TCdpAreaInfoService.Insert(l.ctx, sessionId, areaInfo)
	if err != nil {
		l.Logger.Errorf("sessionId: %s, insert area info error: %v", sessionId, err)
		return nil, err
	}
	return area, nil
}

// 插入默认资源池信息
func (l *IdcCreateLogic) insertDefaultResourcePool(sessionId string, areaId int32) (err error) {
	pool, _, err := table.T_TCdpInstancePoolService.Query(l.ctx, sessionId, fmt.Sprintf("area_id:%d$pool_id:%d", areaId, l.svcCtx.Config.Instance.DefaultPoolId), nil, nil)
	if err != nil && !errors.Is(err, gopublic.ErrNotExist) {
		l.Logger.Errorf("sessionId: %s, insert default resource pool error: %v", sessionId, err)
		return err
	}
	if pool != nil {
		l.Logger.Debugf("sessionId: %s, default resource pool already exists", sessionId)
		return nil
	}
	pool, _, err = table.T_TCdpInstancePoolService.Insert(l.ctx, sessionId, &table.TCdpInstancePool{
		AreaId:       int64(areaId),
		PoolId:       l.svcCtx.Config.Instance.DefaultPoolId,
		InstPoolName: fmt.Sprintf("默认资源池-%d", areaId),
		Status:       1,
		Remark:       "新建节点创建默认资源池",
		CreateBy:     "",
		UpdateBy:     "",
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
		ModifyTime:   time.Now(),
	})
	if err != nil {
		l.Logger.Errorf("sessionId: %s, insert default resource pool error: %v", sessionId, err)
		return err
	}
	l.Logger.Debugf("sessionId: %s, insert default resource pool success[%v]", sessionId, pool)

	return nil
}

// 更新超管账号节点
func (l *IdcCreateLogic) updateSuperAdminNode(sessionId string, areaId int32) (err error) {
	users, _, err := table.T_TCdpSysUserService.QueryAll(l.ctx, sessionId, fmt.Sprintf("is_admin:%d$status:%d", table.IsAdminYes, table.SysUserStatusEnable), nil, nil)
	if err != nil {
		l.Logger.Errorf("sessionId: %s, update super admin node error: %v", sessionId, err)
		return err
	}
	for _, user := range users {
		areaIds := strings.Split(user.AreaIds, ",")
		areaIds = append(areaIds, fmt.Sprint(areaId))
		_, _, err = table.T_TCdpSysUserService.Update(l.ctx, sessionId, user.Id, map[string]any{
			"area_ids": strings.Join(areaIds, ","),
		})
		if err != nil {
			l.Logger.Errorf("sessionId: %s, update super admin node error: %v", sessionId, err)
			return err
		}
	}

	return nil
}
