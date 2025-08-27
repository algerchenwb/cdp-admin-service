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

type StrategyDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStrategyDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StrategyDeleteLogic {
	return &StrategyDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StrategyDeleteLogic) StrategyDelete(req *types.StrategyDeleteReq) (resp *types.StrategyDeleteResp, err error) {
	sessionId := helper.GetSessionId(l.ctx)
	updateBy := helper.GetUserName(l.ctx)
	resp = new(types.StrategyDeleteResp)

	if req.Id == 0 || req.AreaId == 0 {
		return nil, errorx.NewDefaultError(errorx.ParamErrorCode)
	}

	if !helper.CheckAreaId(l.ctx, fmt.Sprintf("%d", req.AreaId)) {
		return nil, errorx.NewDefaultCodeError("用户无该区域的权限")
	}

	strategyInfo, _, err := table.T_TCdpResourceStrategyService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d$status:%d", req.Id, table.BizStrategyStatusValid), nil, nil)
	if err != nil && err != gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] T_TCdpResourceStrategyService Query err. strategyId[%d] err:%+v", sessionId, req.Id, err)
		return nil, errorx.NewDefaultError(errorx.ParamErrorCode)
	}
	if err == gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] T_TCdpResourceStrategyService Query err. ErrNotExist strategyId[%d] err:%+v", sessionId, req.Id, err)
		return nil, errorx.NewDefaultError(errorx.ListStrategyErrorCode)
	}

	bizStrategyInfo, _, err := table.T_TCdpBizStrategyService.Query(l.ctx, sessionId, fmt.Sprintf("inst_strategy_id:%d$status:%d", req.Id, table.BizStrategyStatusValid), nil, nil)
	if err != nil && err != gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] T_TCdpResourceStrategyService Query err. strategyId[%d] err:%+v", sessionId, req.Id, err)
		return nil, errorx.NewDefaultCodeError("查询合约策略失败")
	}
	if err == nil && bizStrategyInfo != nil {
		return nil, errorx.NewDefaultCodeError("该策略下有合约正在使用，无法删除")
	}

	newStrategy, _, err := table.T_TCdpResourceStrategyService.Update(l.ctx, sessionId, strategyInfo.Id, map[string]interface{}{
		"name":      fmt.Sprintf("%s-del-%s", strategyInfo.Name, time.Now().Format("20060102150405")),
		"status":    table.ResourceStrategyStatusInvalid,
		"update_by": updateBy})
	if err != nil {
		l.Logger.Errorf("[%s] Table T_TCdpResourceStrategyService err. strategyId[%d] err:%+v", sessionId, req.Id, err)
		return nil, errorx.NewDefaultError(errorx.ListStrategyErrorCode)
	}

	l.Logger.Infof("[%s] T_TCdpResourceStrategyService delete success. strategyId[%d] newStrategy:%s", sessionId, req.Id, helper.ToJSON(newStrategy))

	searchReq := &diskless.SearchPoolRequest{
		FlowId:     sessionId,
		AreaType:   int32(req.AreaId),
		ResourceId: &strategyInfo.Id,
		Conditions: []string{},
		Offset:     0,
		Length:     99,
		Order:      "",
		Sortby:     "",
	}

	searchPool, err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).SearchPool(sessionId, int64(req.AreaId), searchReq)
	if err != nil {
		l.Logger.Errorf("[%s] SearchPool searchReq:%s err: %v", sessionId, helper.ToJSON(searchReq), err)
		return nil, errorx.NewDefaultCodeError("查询无盘资源池失败")
	}

	if searchPool.Total != 0 {
		l.Logger.Errorf("[%s] SearchPool searchPool.Total !=0  searchReq:%s searchPool total: %d", sessionId, helper.ToJSON(searchReq), searchPool.Total)
		return nil, errorx.NewDefaultCodeError("该策略下有实例还没释放，请先释放后，再删除")
	}
	// 更新无盘资源，策略, todo 优先调用无盘接口
	state := int32(table.ResourceStrategyStatusInvalid)
	updateResReq := &diskless.UpdateResourceRequest{
		FlowId:     sessionId,
		AreaType:   int32(req.AreaId),
		ResourceId: &strategyInfo.Id,
		State:      &state,
	}
	if err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).UpdateResource(sessionId, int64(req.AreaId), updateResReq); err != nil {
		l.Logger.Errorf("[%s] diskless.UpdateResource AreaId[%d] listResReq:%s err:%+v", sessionId, req.AreaId, helper.ToJSON(updateResReq), err)
		return nil, errorx.NewDefaultCodeError("更新无盘策略失败")
	}

	return
}
