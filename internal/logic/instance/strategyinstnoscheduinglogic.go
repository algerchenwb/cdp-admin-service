package instance

import (
	"context"
	"fmt"

	"cdp-admin-service/internal/helper"
	"cdp-admin-service/internal/helper/diskless"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type StrategyInstNoScheduingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStrategyInstNoScheduingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StrategyInstNoScheduingLogic {
	return &StrategyInstNoScheduingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StrategyInstNoScheduingLogic) StrategyInstNoScheduing(req *types.StrategyInstNoScheduingReq) (resp *types.StrategyInstNoScheduingResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)
	resp = new(types.StrategyInstNoScheduingResp)

	for _, item := range req.Items {
		if item.AreaId == 0 || item.InstanceId == 0 || item.Source == "" || item.FlowId == "" {
			return nil, errorx.NewDefaultCodeError("参数错误")
		}

		if !helper.CheckAreaId(l.ctx, fmt.Sprintf("%d", item.AreaId)) {
			l.Logger.Errorf("[%s] diskless.CheckAreaId 用户无该区域的权限 AreaId[%d] instalnceId[%d] err:%+v", sessionId, item.AreaId, item.InstanceId, err)
			resp.FailedInstanceIds = append(resp.FailedInstanceIds, item.InstanceId)
			continue
		}
		// 释放资源池中的实例
		status := int32(1000) // 禁止调度
		releaseReq := &diskless.UpdatePoolItemRequest{
			FlowId:     item.FlowId,
			Source:     item.Source,
			AreaType:   int32(item.AreaId),
			InstanceId: fmt.Sprintf("%d", item.InstanceId),
			Status:     &status,
		}
		if err = diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).UpdatePoolItem(sessionId, int64(item.AreaId), releaseReq); err != nil {
			l.Logger.Errorf("[%s] diskless.UpdatePoolItem AreaId[%d] instalnceId[%d] err:%+v", sessionId, item.AreaId, item.InstanceId, err)
			resp.FailedInstanceIds = append(resp.FailedInstanceIds, item.InstanceId)
			//return nil, errorx.NewHandlerError(errorx.ServerErrorCode, "策略禁止实例调度失败")
		}
		resp.SuccessInstanceIds = append(resp.SuccessInstanceIds, item.InstanceId)
	}

	return
}
