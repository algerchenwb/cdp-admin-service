package instance

import (
	"context"
	"fmt"
	"time"

	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/helper/diskless"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/proto/instance_scheduler"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
)

type StrategyAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStrategyAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StrategyAddLogic {
	return &StrategyAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StrategyAddLogic) StrategyAdd(req *types.StrategyAddReq) (resp *types.StrategyAddResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)
	updateBy := helper.GetUserName(l.ctx)

	if req.AreaId == 0 || req.InstPoolId == 0 {
		return nil, errorx.NewDefaultError(errorx.ParamErrorCode)
	}

	if !helper.CheckAreaId(l.ctx, fmt.Sprintf("%d", req.AreaId)) {
		return nil, errorx.NewDefaultCodeError("用户无该区域的权限")
	}
	// TODO 添加状态判断
	_, _, err = table.T_TCdpAreaInfoService.Query(l.ctx, sessionId, fmt.Sprintf("area_id:%d", req.AreaId), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpAreaInfoService Query err. areaId[%d] err:%+v", sessionId, req.AreaId, err)
		return nil, errorx.NewDefaultError(errorx.ParamErrorCode)
	}
	if err == gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] T_TCdpAreaInfoService Query err. ErrNotExist areaId[%d] err:%+v", sessionId, req.AreaId, err)
		return nil, errorx.NewDefaultError(errorx.AreaNotFoundErrorCode)
	}

	strategyInfo, _, err := table.T_TCdpResourceStrategyService.Query(l.ctx, sessionId, fmt.Sprintf("area_id:%d$name:%s$status:%d", req.AreaId, req.Name, table.ResourceStrategyStatusValid), nil, nil)
	if err != nil && err != gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] T_TCdpResourceStrategyService Query err. AreaId[%d] name[%s] err:%+v", sessionId, req.AreaId, req.Name, err)
		return nil, errorx.NewDefaultError(errorx.DalGetErrorCode)
	}

	if err == nil && strategyInfo != nil && strategyInfo.Id != 0 {
		l.Logger.Errorf("[%s] T_TCdpResourceStrategyService Query Name Exist name[%d] strategyInfo:%s", sessionId, req.Name, helper.ToJSON(strategyInfo))
		return nil, errorx.NewDefaultError(errorx.StrategyNameErrorCode)
	}

	// 检测实例池Id是否存在
	_, _, err = table.T_TCdpInstancePoolService.Query(l.ctx, sessionId, fmt.Sprintf("pool_id:%d$area_id:%d$status:%d", req.InstPoolId, req.AreaId, table.InstancePoolStatusValid), nil, nil)
	if err != nil && err != gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] T_TCdpInstancePoolService Query err. InstPoolId[%d] err:%+v", sessionId, req.InstPoolId, err)
		return nil, errorx.NewDefaultError(errorx.DalGetErrorCode)
	}
	if err == gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] T_TCdpInstancePoolService Query err. ErrNotExist InstPoolId[%d] err:%+v", sessionId, req.InstPoolId, err)
		return nil, errorx.NewDefaultError(errorx.InstancePoolNotFoundErrorCode)
	}

	newStrategy, _, err := table.T_TCdpResourceStrategyService.Insert(l.ctx, sessionId, &table.TCdpResourceStrategy{
		Name:            req.Name,
		ApplicableLever: 0,
		AreaId:          int64(req.AreaId),
		SpecId:          0,
		OuterSpecId:     0,
		SpecName:        "",
		InstPoolId:      int64(req.InstPoolId),
		TotalInstances:  int32(req.TotalInstances),
		BootType:        int32(req.BootType),
		Remark:          req.Remark,
		Status:          table.ResourceStrategyStatusValid,
		CreateBy:        updateBy,
		UpdateBy:        updateBy,
		CreateTime:      time.Now(),
		UpdateTime:      time.Now(),
		ModifyTime:      time.Now(),
	})

	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpResourceStrategyService Insert err. AreaId[%d] err:%+v", sessionId, req.AreaId, err)
		return nil, errorx.NewDefaultError(errorx.DalInsertErrorCode)
	}

	// 无盘接口 todo 优先调用无盘接口
	newResReq := &instance_scheduler.NewResourceRequest{
		FlowId:   sessionId,
		AreaType: int32(req.AreaId),
		ResourceConfig: &instance_scheduler.ResourceConfig{
			ResourceId:    newStrategy.Id,
			AreaType:      int32(req.AreaId),
			Type:          0,
			Name:          req.Name,
			Specification: int64(req.InstPoolId),
			Vlan:          int32(req.VlanId),
			Mode:          int64(req.BootType),
			Capacity:      int32(req.TotalInstances),
			Buffer:        0,
			Init:          int32(req.PreBootCount),
			Concurrent:    0,
			Priority:      0,
			Preemptable:   0,
			AssignConfig:  req.PreBootSchemaIdInfo,
			State:         1,
		},
	}
	if err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).NewResource(sessionId, int64(req.AreaId), newResReq); err != nil {
		l.Logger.Errorf("[%s] diskless.NewResource AreaId[%d] newResReq:%s err:%+v", sessionId, req.AreaId, helper.ToJSON(newResReq), err)
		return nil, errorx.NewDefaultCodeError("创建无盘策略失败")
	}

	resp = &types.StrategyAddResp{
		Id:                  uint64(newStrategy.Id),
		Name:                newStrategy.Name,
		InstPoolId:          req.InstPoolId,
		TotalInstances:      req.TotalInstances,
		VlanId:              req.VlanId,
		BootType:            req.BootType,
		PreBootSchemaIdInfo: req.PreBootSchemaIdInfo,
		PreBootCount:        req.PreBootCount,
		Remark:              newStrategy.Remark,
		Status:              int(newStrategy.Status),
		CreateBy:            newStrategy.CreateBy,
		UpdateBy:            newStrategy.UpdateBy,
		CreateTime:          newStrategy.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateTime:          newStrategy.UpdateTime.Format("2006-01-02 15:04:05"),
		ModifyTime:          newStrategy.ModifyTime.Format("2006-01-02 15:04:05"),
	}

	return
}
