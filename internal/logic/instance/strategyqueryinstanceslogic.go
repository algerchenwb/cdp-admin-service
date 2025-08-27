package instance

import (
	"context"

	"cdp-admin-service/internal/helper"
	"cdp-admin-service/internal/helper/diskless"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type StrategyQueryInstancesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStrategyQueryInstancesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StrategyQueryInstancesLogic {
	return &StrategyQueryInstancesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StrategyQueryInstancesLogic) StrategyQueryInstances(req *types.StrategyQueryInstancesReq) (resp *types.StrategyQueryInstancesResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)

	resp = &types.StrategyQueryInstancesResp{}

	searchReq := &diskless.SearchPoolRequest{
		FlowId:     sessionId,
		AreaType:   req.AreaType,
		ResourceId: req.ResourceId,
		Conditions: req.Conditions,
		Offset:     req.Offset,
		Length:     req.Limit,
		Order:      req.Order,
		Sortby:     req.Sortby,
	}

	searchPool, err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).SearchPool(sessionId, int64(req.AreaType), searchReq)
	if err != nil {
		l.Logger.Errorf("[%s] SearchPool searchReq:%s err: %v", sessionId, helper.ToJSON(searchReq), err)
		return nil, errorx.NewDefaultCodeError("查询无盘资源池失败")
	}

	resp.Total = int(searchPool.Total)

	for _, v := range searchPool.Lists {

		resp.List = append(resp.List, types.PoolItem{
			AreaType:     v.AreaType,
			ResourceId:   helper.StringToInt64(v.ResourceId),
			InstanceId:   helper.StringToInt64(v.InstanceId),
			Mac:          v.Mac,
			Address:      v.Address,
			Flags:        helper.StringToInt64(v.Flags),
			PoolSource:   v.PoolSource,
			PoolOrder:    v.PoolOrder,
			PoolStatus:   v.PoolStatus,
			AssignSource: v.AssignSource,
			AssignOrder:  v.AssignOrder,
			AssignParam:  v.AssignParam,
			AssignStatus: v.AssignStatus,
			CreateTime:   v.CreateTime,
			UpdateTime:   v.UpdateTime,
			ModifyTime:   v.ModifyTime,
		})
	}

	return
}
