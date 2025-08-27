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
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
)

type InstanceBindPoolLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInstanceBindPoolLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InstanceBindPoolLogic {
	return &InstanceBindPoolLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InstanceBindPoolLogic) InstanceBindPool(req *types.InstanceBindPoolReq) (resp *types.InstanceBindPoolResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)
	resp = &types.InstanceBindPoolResp{}

	if req.PoolId == 0 || req.AreaId == 0 || len(req.InstanceIds) == 0 {
		return nil, errorx.NewDefaultError(errorx.ParamErrorCode)
	}
	pool, _, err := table.T_TCdpInstancePoolService.Query(l.ctx, sessionId, fmt.Sprintf("pool_id:%d$area_id:%d$status:%d", req.PoolId, req.AreaId, table.InstancePoolStatusValid), nil, nil)
	if err != nil && err != gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] T_TCdpInstancePoolService Query err. id[%d] err:%+v", sessionId, req.PoolId, err)
		return nil, errorx.NewDefaultCodeError("查询算力池失败")
	}

	if !helper.CheckAreaId(l.ctx, fmt.Sprintf("%d", pool.AreaId)) {
		return nil, errorx.NewDefaultCodeError("用户无该区域的权限")
	}

	// 查询原有资源池的实例
	// 调用无盘的接口
	instListReq := &instance_types.ListInstancesRequestNew{
		Offset:        0,
		Length:        9999,
		Specification: []int{int(pool.PoolId)},
	}
	InstanceDetails, err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).SearcchInstanceList(int64(pool.AreaId), sessionId, instListReq)
	if err != nil {
		l.Logger.Errorf("[%s] SearcchInstanceList areaId:%d, PoolId:%d err:%+v", sessionId, pool.AreaId, pool.PoolId, err)
		return nil, errorx.NewDefaultError(errorx.QueryInstanceFailedErrorCode)
	}
	existInstIds := make(map[string]struct{})
	for _, v := range InstanceDetails {
		existInstIds[fmt.Sprintf("%d", v.Id)] = struct{}{}
	}

	for _, v := range req.InstanceIds {
		if _, ok := existInstIds[fmt.Sprintf("%d", v)]; ok {
			continue
		}
		// 新添加到算力池的实例
		specificationId := int64(pool.PoolId)
		instReq := &instance_types.UpdateInstanceRequest{
			FlowId:     sessionId,
			InstanceId: int64(v),
			UpdateableInstanceInfo: instance_types.UpdateableInstanceInfo{
				Specification: &specificationId,
			},
		}

		if err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).UpdateInstance(int64(pool.AreaId), instReq); err != nil {
			l.Logger.Errorf("[%s] diskless.UpdateInstance AreaId[%d] instalnceId[%d] err:%+v", sessionId, pool.AreaId, v, err)
			resp.FailedInstanceIds = append(resp.FailedInstanceIds, v)
			continue
		}
		resp.SuccessInstanceIds = append(resp.SuccessInstanceIds, v)
	}
	return
}
