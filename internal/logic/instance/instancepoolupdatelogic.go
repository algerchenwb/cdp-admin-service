package instance

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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

type InstancePoolUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInstancePoolUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InstancePoolUpdateLogic {
	return &InstancePoolUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InstancePoolUpdateLogic) InstancePoolUpdate(req *types.InstancePoolUpdateReq) (resp *types.InstancePoolUpdateResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)
	updateBy := helper.GetUserName(l.ctx)
	resp = &types.InstancePoolUpdateResp{}

	if req.Id == 0 || req.PoolName == "" || req.AreaId == 0 {
		return nil, errorx.NewDefaultError(errorx.ParamErrorCode)
	}

	if !helper.CheckAreaId(l.ctx, fmt.Sprintf("%d", req.AreaId)) {
		return nil, errorx.NewDefaultCodeError("用户无该区域的权限")
	}

	InstIds := strings.Split(req.InstIds, ",")
	if len(InstIds) == 0 {
		return nil, errorx.NewDefaultCodeError("实例ID不能为空")
	}

	qry := fmt.Sprintf("id__ex:%d$area_id:%d$inst_pool_name:%s", req.Id, req.AreaId, req.PoolName)
	_, _, err = table.T_TCdpInstancePoolService.Query(l.ctx, sessionId, qry, nil, nil)
	if err != nil && err != gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] T_TCdpInstancePoolService query failed, AreaId:%d,qry:%s err:%+v", sessionId, req.AreaId, qry, err)
		return nil, errorx.NewDefaultCodeError("查询数据失败")
	}
	if err == nil {
		return nil, errorx.NewDefaultCodeError("算力池名称已存在")
	}

	poolInfo, _, err := table.T_TCdpInstancePoolService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d$status:%d", req.Id, table.InstancePoolStatusValid), nil, nil)
	if err != nil && err != gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] T_TCdpInstancePoolService Query err. id[%d] err:%+v", sessionId, req.Id, err)
		return nil, errorx.NewDefaultError(errorx.DalGetErrorCode)
	}

	if err == gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] T_TCdpInstancePoolService Query err. ErrNotExist Id[%d] err:%+v", sessionId, req.Id, err)
		return nil, errorx.NewDefaultError(errorx.InstancePoolNotFoundErrorCode)
	}

	// 查询原有资源池的实例
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

	var InInstIds []string
	for _, v := range InstanceDetails {
		instanceId := fmt.Sprintf("%d", v.Id)

		InInstIds = append(InInstIds, instanceId)

		if !gopublic.StringInArray(instanceId, InstIds) {
			// 不在算力池里，则要放回默认的算力池
			specificationId := int64(table.InstancePoolDefaultId)
			instReq := &instance_types.UpdateInstanceRequest{
				FlowId:     sessionId,
				InstanceId: v.Id,
				UpdateableInstanceInfo: instance_types.UpdateableInstanceInfo{
					Specification: &specificationId,
				},
			}

			if err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).UpdateInstance(int64(req.AreaId), instReq); err != nil {
				l.Logger.Errorf("[%s] diskless.UpdateInstance AreaId[%d] instalnceId[%d] err:%+v", sessionId, req.AreaId, v.Id, err)
				//return nil, errorx.NewHandlerError(errorx.ServerErrorCode, "更新无盘实例算力池ID失败")
			}
		}
	}

	for _, v := range InstIds {

		if !gopublic.StringInArray(v, InInstIds) {
			// 新添加到算力池的实例
			instanceId, _ := strconv.ParseInt(v, 10, 64)
			specificationId := int64(poolInfo.PoolId)
			instReq := &instance_types.UpdateInstanceRequest{
				FlowId:     sessionId,
				InstanceId: instanceId,
				UpdateableInstanceInfo: instance_types.UpdateableInstanceInfo{
					Specification: &specificationId,
				},
			}

			if err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).UpdateInstance(int64(req.AreaId), instReq); err != nil {
				l.Logger.Errorf("[%s] diskless.UpdateInstance AreaId[%d] instalnceId[%d] err:%+v", sessionId, req.AreaId, specificationId, err)
				//return nil, errorx.NewHandlerError(errorx.ServerErrorCode, "更新无盘实例算力池ID失败")
			}
		}
	}

	if req.PoolName != poolInfo.InstPoolName || req.Remark != poolInfo.Remark {
		newPool, _, err := table.T_TCdpInstancePoolService.Update(l.ctx, sessionId, req.Id, map[string]interface{}{
			"inst_pool_name": req.PoolName,
			"remark":         req.Remark,
			"update_by":      updateBy,
		})
		if err != nil {
			l.Logger.Errorf("[%s] Table Update err. Id[%d] err:%+v", sessionId, req.Id, err)
			return nil, errorx.NewDefaultCodeError("更新算力信息池失败")
		}
		l.Logger.Debugf("[%s] T_TCdpInstancePoolService Update success. Id[%d] newPool:%s", sessionId, req.Id, helper.ToJSON(newPool))
	}

	return
}
