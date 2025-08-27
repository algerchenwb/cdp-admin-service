package biz

import (
	"context"
	"fmt"
	"strings"
	"time"

	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigBizStrategyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfigBizStrategyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigBizStrategyLogic {
	return &ConfigBizStrategyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfigBizStrategyLogic) ConfigBizStrategy(req *types.ConfigBizStrategyReq) (resp *types.ConfigBizStrategyResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)
	updateBy := helper.GetUserName(l.ctx)

	if len(req.BizStrategyList) == 0 || req.BizId == 0 {
		return nil, errorx.NewDefaultError(errorx.ParamErrorCode)
	}

	bizInfo, _, err := table.T_TCdpBizInfoService.Query(l.ctx, sessionId, fmt.Sprintf("biz_id:%d$status__ex:0", req.BizId), nil, nil)
	if err != nil { // todo : 判断 err != gopublic.ErrNotExist 不数据不存在单独返回
		l.Logger.Errorf("[%s] T_TCdpBizInfoService Query err. bizId[%d] err:%+v", sessionId, req.BizId, err)
		return nil, errorx.NewDefaultError(errorx.ParamErrorCode)
	}

	existBizStrategys, _, err := table.T_TCdpBizStrategyService.QueryAll(l.ctx, sessionId,
		fmt.Sprintf("biz_id:%d$status:%d", req.BizId, table.BizStrategyStatusValid), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpBizStrategyService QueryAll err. bizId[%d] err:%+v", sessionId, req.BizId, err)
		return nil, errorx.NewDefaultError(errorx.QueryBizStrategyFailedErrorCode)
	}

	var instStrategyIds []string
	for _, tgyInfo := range req.BizStrategyList {
		instStrategyIds = append(instStrategyIds, fmt.Sprintf("%d", tgyInfo.StrategyId))
	}

	for _, existBizStrategy := range existBizStrategys {
		instStrategyIds = append(instStrategyIds, fmt.Sprintf("%d", existBizStrategy.InstStrategyId))
	}

	instStrategys, _, err := table.T_TCdpResourceStrategyService.QueryAll(l.ctx, sessionId,
		fmt.Sprintf("id__in:%s$status:%d", strings.Join(instStrategyIds, ","), table.ResourceStrategyStatusValid), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] Table Query err. bizId[%d] err:%+v", sessionId, req.BizId, err)
		return nil, errorx.NewDefaultError(errorx.QueryBizStrategyFailedErrorCode)
	}

	var instStrategyMap = make(map[int64]table.TCdpResourceStrategy)
	for _, instStrategy := range instStrategys {
		instStrategyMap[instStrategy.Id] = instStrategy
	}

	var existPoolMap = make(map[int64]struct{})
	for _, existBizStrategy := range existBizStrategys {
		existPoolMap[instStrategyMap[existBizStrategy.InstStrategyId].InstPoolId] = struct{}{}
	}

	var ErrMsgs []string

	for _, tgyInfo := range req.BizStrategyList {
		instStrategy, ok := instStrategyMap[tgyInfo.StrategyId]
		if !ok {
			ErrMsgs = append(ErrMsgs, fmt.Sprintf("算力策略[%d]不存在", tgyInfo.StrategyId))
			continue
		}

		if _, ok := existPoolMap[instStrategy.InstPoolId]; ok {
			ErrMsgs = append(ErrMsgs, fmt.Sprintf("算力策略[%s]绑定的算力池已分配给当前合约", instStrategy.Name))
			continue
		}

		newBizStrategyInfo, _, err := table.T_TCdpBizStrategyService.Insert(l.ctx, sessionId, &table.TCdpBizStrategy{
			BizId:          req.BizId,
			InstStrategyId: tgyInfo.StrategyId,
			AreaId:         int64(bizInfo.AreaId),
			Status:         1,
			CreateBy:       updateBy,
			UpdateBy:       updateBy,
			CreateTime:     time.Now(),
			UpdateTime:     time.Now(),
			ModifyTime:     time.Now()})
		if err != nil {
			l.Logger.Errorf("[%s] T_TCdpBizStrategyService Insert err. bizId[%d] err:%+v", sessionId, req.BizId, err)
			ErrMsgs = append(ErrMsgs, fmt.Sprintf("%d:合约[%d]策略绑定失败", tgyInfo.StrategyId, req.BizId))
		}
		existPoolMap[instStrategy.InstPoolId] = struct{}{}
		l.Logger.Infof("[%s] T_TCdpBizStrategyService Insert  success. bizId[%d] newBizStrategyInfo:%s", sessionId, req.BizId, helper.ToJSON(newBizStrategyInfo))
	}

	resp = &types.ConfigBizStrategyResp{}
	resp.ErrMsgs = append(resp.ErrMsgs, ErrMsgs...)
	return
}
