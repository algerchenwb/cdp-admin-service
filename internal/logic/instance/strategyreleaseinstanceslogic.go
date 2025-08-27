package instance

import (
	"context"
	"fmt"

	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/helper/diskless"
	"cdp-admin-service/internal/model/errorx"
	instance_types "cdp-admin-service/internal/proto/instance_service/types"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type StrategyReleaseInstancesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStrategyReleaseInstancesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StrategyReleaseInstancesLogic {
	return &StrategyReleaseInstancesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StrategyReleaseInstancesLogic) StrategyReleaseInstances(req *types.StrategyReleaseInstancesReq) (resp *types.StrategyReleaseInstancesResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)

	for _, item := range req.Item {
		if item.AreaId == 0 || item.InstanceId == 0 || item.Source == "" || item.FlowId == "" {
			return nil, errorx.NewDefaultCodeError("参数错误")
		}

		if item.Mode == 2 {
			// 判断 实例所在的算力池，是不是最后一个，是则是让释放
			instListReq := &instance_types.ListInstancesRequestNew{
				Offset:      0,
				Length:      9999,
				InstanceIds: []int{int(item.InstanceId)},
			}
			InstanceDetails, err1 := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).SearcchInstanceList(int64(item.AreaId), sessionId, instListReq)
			if err1 != nil {
				l.Logger.Errorf("[%s] SearcchInstanceList areaId:%d, InstanceId:%d err:%+v", sessionId, item.AreaId, item.InstanceId, err1)
				return nil, errorx.NewDefaultError(errorx.QueryInstanceFailedErrorCode)
			}
			if len(InstanceDetails) == 0 {
				l.Logger.Errorf("[%s] SearcchInstanceList len = 0 areaId:%d, InstanceId:%d err:%+v", sessionId, item.AreaId, item.InstanceId, err1)
				return nil, errorx.NewDefaultCodeError("资源池没有对应的实例")
			}

			instDetail := InstanceDetails[0]
			if instDetail.Specification == 0 {
				l.Logger.Errorf("[%s] SearcchInstanceList Specification is nil areaId:%d, InstanceId:%d err:%+v", sessionId, item.AreaId, item.InstanceId, err)
				return nil, errorx.NewDefaultCodeError("实例算力池ID为空")
			}
			poolId := instDetail.Specification
			// 查询资源池信息
			instListReq = &instance_types.ListInstancesRequestNew{
				Offset:        0,
				Length:        9999,
				Specification: []int{int(poolId)},
			}

			InstanceDetails2, err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).SearcchInstanceList(int64(item.AreaId), sessionId, instListReq)
			if err != nil {
				l.Logger.Errorf("[%s] SearcchInstanceList areaId:%d, Specification:%d err:%+v", sessionId, item.AreaId, poolId, err)
				return nil, errorx.NewDefaultError(errorx.QueryInstanceFailedErrorCode)
			}
			if len(InstanceDetails2) <= 1 {
				l.Logger.Info("[%s] SearcchInstanceList len <= 1 areaId:%d, InstanceId:%d err:%+v", sessionId, item.AreaId, item.InstanceId, err)
				return nil, errorx.NewDefaultCodeError("实例是算力池最后一个实例，请从资源池中解散")
			}

			// 释放实例回默认算力池 放回默认的算力池
			specificationId := int64(table.InstancePoolDefaultId)
			instReq := &instance_types.UpdateInstanceRequest{
				FlowId:     sessionId,
				InstanceId: int64(item.InstanceId),
				UpdateableInstanceInfo: instance_types.UpdateableInstanceInfo{
					Specification: &specificationId,
				},
			}
			if err = diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).UpdateInstance(int64(item.AreaId), instReq); err != nil {
				l.Logger.Errorf("[%s] diskless.UpdateInstance AreaId[%d] instalnceId[%d] err:%+v", sessionId, item.AreaId, item.InstanceId, err)
				return nil, errorx.NewHandlerError(errorx.ServerErrorCode, "更新无盘实例算力池ID失败")
			}
		}

		// 释放资源池中的实例
		releaseReq := &diskless.ReleasePoolItemRequest{
			FlowId:     item.FlowId,
			Source:     item.Source,
			AreaType:   int32(item.AreaId),
			InstanceId: fmt.Sprintf("%d", item.InstanceId),
		}
		if err = diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).ReleasePoolItem(sessionId, int64(item.AreaId), releaseReq); err != nil {
			l.Logger.Errorf("[%s] diskless.ReleasePoolItem AreaId[%d] instalnceId[%d] err:%+v", sessionId, item.AreaId, item.InstanceId, err)
			return nil, errorx.NewHandlerError(errorx.ServerErrorCode, "释放资源池实例失败")
		}
	}

	return
}
