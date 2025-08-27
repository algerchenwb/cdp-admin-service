package biz

import (
	"context"
	"fmt"

	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteBizStrategyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteBizStrategyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteBizStrategyLogic {
	return &DeleteBizStrategyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteBizStrategyLogic) DeleteBizStrategy(req *types.DeleteBizStrategyReq) (resp *types.DeleteBizStrategyResp, err error) {
	sessionId := helper.GetSessionId(l.ctx)
	updateBy := helper.GetUserName(l.ctx)
	resp = &types.DeleteBizStrategyResp{}

	if req.Id == 0 {
		return nil, errorx.NewDefaultError(errorx.ParamErrorCode)
	}
	bizStrategy, _, err := table.T_TCdpBizStrategyService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d$status:%d", req.Id, table.BizStrategyStatusValid), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpBizStrategyService Query err. id[%d] err:%+v", sessionId, req.Id, err)
		return nil, errorx.NewDefaultError(errorx.QueryBizStrategyFailedErrorCode)
	}

	resouceStrategy, _, err := table.T_TCdpResourceStrategyService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d$status:%d", bizStrategy.InstStrategyId, table.ResourceStrategyStatusValid), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpResourceStrategyService Query err. id[%d] err:%+v", sessionId, req.Id, err)
		return nil, errorx.NewDefaultError(errorx.QueryBizStrategyFailedErrorCode)
	}
	if err := l.checkBizStrategy(bizStrategy.BizId, resouceStrategy.Id); err != nil {
		return nil, err
	}

	_, _, err = table.T_TCdpBizStrategyService.Update(l.ctx, sessionId, req.Id, map[string]interface{}{
		"status":    table.BizStrategyStatusInvalid,
		"update_by": updateBy,
	})
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpBizStrategyService Update err. id[%d] err:%+v", sessionId, req.Id, err)
		return nil, errorx.NewDefaultError(errorx.DeleteBizStrategyFailedErrorCode)
	}

	l.Logger.Debugf("[%s] DeleteBizStrategy success. id[%d]", sessionId, req.Id)

	return
}

func (l *DeleteBizStrategyLogic) checkBizStrategy(bizId int64, id int64) (err error) {
	sessionId := helper.GetSessionId(l.ctx)
	// 主算力策略
	total, _, _, err := table.T_TCdpCloudclientInfoService.QueryPage(l.ctx, sessionId, fmt.Sprintf("biz_id:%d$first_strategy_id:%d$status:%d", bizId, id, table.CloudClientStatusValid), 0, 0, nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpCloudclientInfoService QueryPage err. id[%d] err:%+v", sessionId, id, err)
		return errorx.NewDefaultError(errorx.QueryBizStrategyFailedErrorCode)
	}
	if total > 0 {
		return errorx.NewDefaultError(errorx.DeleteBizStrategyRefusedErrorCode)
	}
	// 备算力策略
	total, _, _, err = table.T_TCdpCloudclientInfoService.QueryPage(l.ctx, sessionId, fmt.Sprintf("biz_id:%d$second_strategy_id:%d$status:%d", bizId, id, table.CloudClientStatusValid), 0, 0, nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpCloudclientInfoService QueryPage err. id[%d] err:%+v", sessionId, id, err)
		return errorx.NewDefaultError(errorx.QueryBizStrategyFailedErrorCode)
	}
	if total > 0 {
		return errorx.NewDefaultError(errorx.DeleteBizStrategyRefusedErrorCode)
	}
	// 主算力策略
	total, _, _, err = table.T_TCdpCloudboxInfoService.QueryPage(l.ctx, sessionId, fmt.Sprintf("biz_id:%d$first_strategy_id:%d$status:%d", bizId, id, table.CloudClientStatusValid), 0, 0, nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpCloudboxInfoService QueryPage err. id[%d] err:%+v", sessionId, id, err)
		return errorx.NewDefaultError(errorx.QueryBizStrategyFailedErrorCode)
	}
	if total > 0 {
		return errorx.NewDefaultError(errorx.DeleteBizStrategyRefusedErrorCode)
	}
	// 备算力策略
	total, _, _, err = table.T_TCdpCloudboxInfoService.QueryPage(l.ctx, sessionId, fmt.Sprintf("biz_id:%d$second_strategy_id:%d$status:%d", bizId, id, table.CloudClientStatusValid), 0, 0, nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpCloudboxInfoService QueryPage err. id[%d] err:%+v", sessionId, id, err)
		return errorx.NewDefaultError(errorx.QueryBizStrategyFailedErrorCode)
	}
	if total > 0 {
		return errorx.NewDefaultError(errorx.DeleteBizStrategyRefusedErrorCode)
	}
	return nil
}
