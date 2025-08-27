package instance

import (
	"context"
	"fmt"
	"strconv"
	"strings"
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

type InstancePoolAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInstancePoolAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InstancePoolAddLogic {
	return &InstancePoolAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InstancePoolAddLogic) InstancePoolAdd(req *types.InstancePoolAddReq) (resp *types.InstancePoolAddResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)
	updateBy := helper.GetUserName(l.ctx)
	resp = &types.InstancePoolAddResp{}

	if req.AreaId == 0 || req.InstIds == "" || req.PoolName == "" {
		return nil, errorx.NewDefaultCodeError("参数错误")
	}

	if !helper.CheckAreaId(l.ctx, fmt.Sprintf("%d", req.AreaId)) {
		return nil, errorx.NewDefaultCodeError("用户无该区域的权限")
	}

	vmIds := strings.Split(req.InstIds, ",")
	if len(vmIds) == 0 {
		return nil, errorx.NewDefaultCodeError("实例列表为空")
	}

	_, _, err = table.T_TCdpInstancePoolService.Query(l.ctx, sessionId, fmt.Sprintf("area_id:%d$inst_pool_name:%s", req.AreaId, req.PoolName), nil, nil)
	if err != nil && err != gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] T_TCdpInstancePoolService query failed  AreaId:%d, err:%+v", sessionId, req.AreaId, err)
		return nil, errorx.NewDefaultCodeError("查询数据失败")
	}

	if err == nil {
		return nil, errorx.NewDefaultCodeError("算力池名称已存在")
	}

	PoolInfos, _, err := table.T_TCdpInstancePoolService.QueryAll(l.ctx, sessionId, fmt.Sprintf("area_id:%d$status:1", req.AreaId), "pool_id", "desc")
	if err != nil && err != gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] T_TCdpInstancePoolService QueryAll failed req.AreaId :%d, err:%+v", sessionId, req.AreaId, err)
		return nil, errorx.NewDefaultCodeError("查询数据失败")
	}

	lastPoolId := int64(0)
	if len(PoolInfos) > 0 {
		lastPoolId = PoolInfos[0].PoolId
	} else {
		lastPoolId = int64(table.InstancePoolDefaultId)
	}

	newPoolId := lastPoolId + 1 // 确认新ID可有
	_, _, err = table.T_TCdpInstancePoolService.Query(l.ctx, sessionId, fmt.Sprintf("pool_id:%d", newPoolId), nil, nil)
	if err != nil && err != gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] T_TCdpInstancePoolService query failed req.AreaId :%d, err:%+v", sessionId, req.AreaId, err)
		return nil, errorx.NewDefaultCodeError("查询数据失败")
	}

	if err == nil {
		newPoolId = newPoolId + 1
	}

	newInstPoolInfo, _, err := table.T_TCdpInstancePoolService.Insert(l.ctx, sessionId, table.TCdpInstancePool{
		PoolId:       newPoolId,
		InstPoolName: req.PoolName,
		BizId:        0,
		AreaId:       int64(req.AreaId),
		Status:       table.InstancePoolStatusValid,
		Remark:       req.Remark,
		CreateBy:     updateBy,
		UpdateBy:     updateBy,
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
		ModifyTime:   time.Now(),
	})
	if err != nil {
		l.Logger.Errorf("[%s] table.T_TCdpInstancePoolService.Insert AreaId:%d, newPoolId:%d, name :%s err:%+v", sessionId, req.AreaId, newPoolId, req.PoolName, err)
		return nil, errorx.NewDefaultCodeError("创建算力池失败")
	}

	l.Logger.Debugf("[%s] T_TCdpInstancePoolService create new success. info:%s", sessionId, helper.ToJSON(newInstPoolInfo))

	//调用无盘的接口
	for _, vmId := range vmIds {

		instanceId, _ := strconv.ParseInt(vmId, 10, 64)
		specificationId := int64(newInstPoolInfo.PoolId)
		instReq := &instance_types.UpdateInstanceRequest{
			FlowId:     sessionId,
			InstanceId: instanceId,
			UpdateableInstanceInfo: instance_types.UpdateableInstanceInfo{
				Specification: &specificationId,
			},
		}
		if err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).UpdateInstance(int64(req.AreaId), instReq); err != nil {
			l.Logger.Errorf("[%s] diskless.UpdateInstance AreaId[%d] err:%+v", sessionId, req.AreaId, err)
			return nil, errorx.NewHandlerError(errorx.ServerErrorCode, "更新无盘实例算力池ID失败")
		}
	}

	resp.PoolId = uint64(newPoolId)
	return
}
