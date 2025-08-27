package instance

import (
	"context"
	"fmt"
	"time"

	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/helper/diskless"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
)

type StrategyUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStrategyUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StrategyUpdateLogic {
	return &StrategyUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StrategyUpdateLogic) StrategyUpdate(req *types.StrategyUpdateReq) (resp *types.StrategyUpdateResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)
	updateBy := helper.GetUserName(l.ctx)

	if req.Id == 0 || req.Name == "" || req.InstPoolId == 0 {
		return nil, errorx.NewDefaultError(errorx.ParamErrorCode)
	}

	if !helper.CheckAreaId(l.ctx, fmt.Sprintf("%d", req.AreaId)) {
		return nil, errorx.NewDefaultCodeError("用户无该区域的权限")
	}

	strategyInfo, _, err := table.T_TCdpResourceStrategyService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d$status:%d", req.Id, table.BizStrategyStatusValid), nil, nil)
	if err != nil && err != gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] T_TCdpResourceStrategyService Query err. strategyId[%d] err:%+v", sessionId, req.Id, err)
		return nil, errorx.NewDefaultError(errorx.DalGetErrorCode)
	}
	if err == gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] T_TCdpResourceStrategyService Query err. ErrNotExist strategyId[%d] err:%+v", sessionId, req.Id, err)
		return nil, errorx.NewDefaultError(errorx.ListStrategyErrorCode)
	}

	if req.Name != strategyInfo.Name {
		strategyInfo, _, err := table.T_TCdpResourceStrategyService.Query(l.ctx, sessionId, fmt.Sprintf("area_id:%d$name:%s$status:%d", req.AreaId, req.Name, table.ResourceStrategyStatusValid), nil, nil)
		if err != nil && err != gopublic.ErrNotExist {
			l.Logger.Errorf("[%s] T_TCdpResourceStrategyService Query err. AreaId[%d] name[%s] err:%+v", sessionId, req.AreaId, req.Name, err)
			return nil, errorx.NewDefaultError(errorx.DalGetErrorCode)
		}

		if err == nil && strategyInfo != nil && strategyInfo.Id != 0 {
			l.Logger.Errorf("[%s] T_TCdpResourceStrategyService Query Name Exist name[%d] strategyInfo:%s", sessionId, req.Name, helper.ToJSON(strategyInfo))
			return nil, errorx.NewDefaultError(errorx.StrategyNameErrorCode)
		}
	}
	// todo 优先调用无盘接口
	newStrategy, _, err := table.T_TCdpResourceStrategyService.Update(l.ctx, sessionId, strategyInfo.Id, map[string]interface{}{
		"name":            req.Name,
		"inst_pool_id":    req.InstPoolId,
		"total_instances": req.TotalInstances,
		"remark":          req.Remark,
		"update_time":     time.Now(),
		"update_by":       updateBy})
	if err != nil {
		l.Logger.Errorf("[%s] Table T_TCdpResourceStrategyService err. strategyId[%d] err:%+v", sessionId, req.Id, err)
		return nil, errorx.NewDefaultError(errorx.UpdateResourceStrategyErrorCode)
	}

	l.Logger.Infof("[%s] T_TCdpResourceStrategyService Update success. strategyId[%d] newStrategy:%s", sessionId, req.Id, helper.ToJSON(newStrategy))

	// 更新无盘资源，策略
	vlan := int32(req.VlanId)
	init := int32(req.PreBootCount)
	totalInst := int32(req.TotalInstances)
	instPoolId := int64(req.InstPoolId)
	updateResReq := &diskless.UpdateResourceRequest{
		FlowId:        sessionId,
		AreaType:      int32(req.AreaId),
		ResourceId:    &strategyInfo.Id,
		Capacity:      &totalInst,
		Name:          &req.Name,
		Vlan:          &vlan,
		AssignConfig:  &req.PreBootSchemaIdInfo,
		Init:          &init,
		Specification: &instPoolId,
	}
	if err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).UpdateResource(sessionId, int64(req.AreaId), updateResReq); err != nil {
		l.Logger.Errorf("[%s] diskless.UpdateResource AreaId[%d] listResReq:%s err:%+v", sessionId, req.AreaId, helper.ToJSON(updateResReq), err)
		return nil, errorx.NewDefaultCodeError("更新无盘策略失败")
	}

	return
}
