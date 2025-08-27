package instance

import (
	"context"
	"fmt"
	"time"

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

type InstancePoolReleaseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInstancePoolReleaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InstancePoolReleaseLogic {
	return &InstancePoolReleaseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InstancePoolReleaseLogic) InstancePoolRelease(req *types.InstancePoolReleaseReq) (resp *types.InstancePoolReleaseResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)
	updateBy := helper.GetUserName(l.ctx)
	resp = &types.InstancePoolReleaseResp{}

	if req.Id == 0 || req.AreaId == 0 {
		return nil, errorx.NewDefaultError(errorx.ParamErrorCode)
	}

	if !helper.CheckAreaId(l.ctx, fmt.Sprintf("%d", req.AreaId)) {
		return nil, errorx.NewDefaultCodeError("用户无该区域的权限")
	}

	qry := fmt.Sprintf("id:%d$status:%d", req.Id, table.InstancePoolStatusValid)
	poolInfo, _, err := table.T_TCdpInstancePoolService.Query(l.ctx, sessionId, qry, nil, nil)
	if err != nil && err != gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] Table Query err. Id[%d] qry:%s err:%+v", sessionId, req.Id, qry, err)
		return nil, errorx.NewDefaultError(errorx.DalGetErrorCode)
	}

	if err == gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] Table Query err. ErrNotExist Id[%d]  qry:%s, err:%+v", sessionId, req.Id, qry, err)
		return nil, errorx.NewDefaultError(errorx.InstancePoolNotFoundErrorCode)
	}

	if poolInfo.PoolId == table.InstancePoolDefaultId {
		return nil, errorx.NewDefaultError(errorx.InstancePoolDefaultIdErrorCode)
	}

	//  检测是否还有算力策略绑定
	qry = fmt.Sprintf("area_id:%d$inst_pool_id:%d$status:1", req.AreaId, poolInfo.PoolId)
	strategyInfos, _, err := table.T_TCdpResourceStrategyService.QueryAll(l.ctx, sessionId, qry, nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpResourceStrategyService QueryAll err. poolInfo.PoolId[%d] err:%+v", sessionId, poolInfo.PoolId, err)
		return nil, errorx.NewDefaultError(errorx.DalGetErrorCode)
	}
	if len(strategyInfos) > 0 {
		name := make([]string, len(strategyInfos))
		for _, v := range strategyInfos {
			name = append(name, v.Name)
		}
		return nil, errorx.NewDefaultCodeError(fmt.Sprintf("该算力池下有%d个策略绑定，无法解散，策略名称：%v", len(strategyInfos), name))
	}
	// 调用无盘的接口
	instListReq := &instance_types.ListInstancesRequestNew{
		Offset:        0,
		Length:        9999,
		Specification: []int{int(poolInfo.PoolId)},
	}
	InstanceDetails, err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).SearcchInstanceList(int64(req.AreaId), sessionId, instListReq)
	if err != nil {
		l.Logger.Errorf("[%s] SearcchInstanceList areaId:%d, PoolId:%d err:%+v", sessionId, req.AreaId, poolInfo.PoolId, err)
		return nil, errorx.NewDefaultError(errorx.QueryInstanceFailedErrorCode)
	}

	for _, v := range InstanceDetails {

		// 放回默认的算力池
		specificationId := int64(table.InstancePoolDefaultId)
		instReq := &instance_types.UpdateInstanceRequest{
			FlowId:     sessionId,
			InstanceId: v.Id,
			UpdateableInstanceInfo: instance_types.UpdateableInstanceInfo{
				Specification: &specificationId,
			},
		}
		if err = diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).UpdateInstance(int64(req.AreaId), instReq); err != nil {
			l.Logger.Errorf("[%s] diskless.UpdateInstance AreaId[%d] instalnceId[%d] err:%+v", sessionId, req.AreaId, v.Id, err)
			return nil, errorx.NewHandlerError(errorx.ServerErrorCode, "更新无盘实例算力池ID失败")
		}
	}

	newPool, _, err := table.T_TCdpInstancePoolService.Update(l.ctx, sessionId, req.Id, map[string]interface{}{
		"status":         table.InstancePoolStatusInvalid,
		"inst_pool_name": fmt.Sprintf("%s-del-%s", poolInfo.InstPoolName, time.Now().Format("20060102150405")),
		"update_by":      updateBy,
	})
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpInstancePoolService Update err. Id[%d] err:%+v", sessionId, req.Id, err)
		return
	}

	l.Logger.Debugf("[%s] InstancePoolRelease success. Id[%d] newPool:%s", sessionId, req.Id, helper.ToJSON(newPool))
	return
}
