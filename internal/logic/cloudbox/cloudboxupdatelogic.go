package cloudbox

import (
	"context"
	"fmt"
	"time"

	"cdp-admin-service/internal/helper"
	cdp_cache "cdp-admin-service/internal/helper/cdp_cache"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/helper/diskless"
	"cdp-admin-service/internal/logic/common"
	"cdp-admin-service/internal/model/errorx"
	instance_types "cdp-admin-service/internal/proto/instance_service/types"
	proto "cdp-admin-service/internal/proto/location_seat_service"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
)

type CloudBoxUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCloudBoxUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CloudBoxUpdateLogic {
	return &CloudBoxUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CloudBoxUpdateLogic) CloudBoxUpdate(req *types.CloudBoxUpdateReq) (resp *types.CloudBoxUpdateResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)
	updateBy := helper.GetUserName(l.ctx)

	resp = new(types.CloudBoxUpdateResp)

	cloudBoxInfo, _, err := table.T_TCdpCloudboxInfoService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d$status:1", req.CloudBoxId), nil, nil)
	if err != nil && err != gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] Query T_TCdpCloudboxInfo failed. cloudbox_id:%d err: %v", sessionId, req.CloudBoxId, err)
		return nil, errorx.NewDefaultCodeError("查询云盒信息失败")
	}
	if err == gopublic.ErrNotExist {
		return nil, errorx.NewHandlerError(errorx.ServerErrorCode, "云盒信息不存在")
	}
	if err := checkStrategy(*cloudBoxInfo, req.FirstStrategyId, req.SecondStrategyId); err != nil {
		return nil, errorx.NewDefaultCodeError(err.Error())
	}

	biz := cdp_cache.GetBizCache(l.ctx, sessionId, req.BizId)
	if biz == nil {
		return nil, errorx.NewDefaultCodeError("当前租户不存在")
	}
	if biz.BoxVlanId == 0 {
		return nil, errorx.NewDefaultCodeError("当前租户未绑定云盒vlan")
	}

	if req.Ip != cloudBoxInfo.Ip {
		vlan, err := common.LoadVlanIP(l.ctx, l.svcCtx, int64(biz.VlanId), int64(biz.AreaId), diskless.IpTypeBox)
		if err != nil {
			l.Logger.Errorf("[%s] LoadVlanIP failed. biz_id:%d err: %v", sessionId, req.BizId, err)
			return nil, errorx.NewDefaultCodeError("查询vlan信息失败")
		}
		if !vlan.FreeIp(req.Ip) {
			return nil, errorx.NewDefaultCodeError("IP不存在或已被使用")
		}
	}

	if req.FirstStrategyId == 0 {
		return nil, errorx.NewDefaultCodeError("主算力策略不能为空")
	}
	_, _, err = table.T_TCdpBizStrategyService.Query(l.ctx, sessionId, fmt.Sprintf("biz_id:%d$inst_strategy_id:%d$status:%d", req.BizId, req.FirstStrategyId, 1), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] Query T_TCdpBizStrategyService failed. biz_id:%d inst_strategy_id:%d err: %v", sessionId, req.BizId, req.FirstStrategyId, err)
		return nil, errorx.NewDefaultCodeError("当前租户不存在该主算力策略")
	}
	if req.SecondStrategyId != 0 {
		_, _, err := table.T_TCdpBizStrategyService.Query(l.ctx, sessionId, fmt.Sprintf("biz_id:%d$inst_strategy_id:%d$status:%d", req.BizId, req.SecondStrategyId, 1), nil, nil)
		if err != nil {
			l.Logger.Errorf("[%s] Query T_TCdpBizStrategyService failed. biz_id:%d inst_strategy_id:%d err: %v", sessionId, req.BizId, req.SecondStrategyId, err)
			return nil, errorx.NewDefaultCodeError("当前租户不存在该从算力策略")
		}
	}

	var bootSchemaInfo *table.TCdpBootSchemaInfo = new(table.TCdpBootSchemaInfo)
	if req.BootSchemaId != 0 {
		bootSchemaInfo, _, err = table.T_TCdpBootSchemaInfoService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d$status:%d", req.BootSchemaId, table.TCdpBootSchemaInfoStatusEnable), nil, nil)
		if err != nil {
			l.Logger.Errorf("[%s] Query T_TCdpBootSchemaInfo failed. id:%d err: %v", sessionId, req.BootSchemaId, err)
			return nil, errorx.NewDefaultCodeError("查询启动方案信息失败")
		}
	}

	// 调用无盘实例的接口
	instReq := &instance_types.UpdateInstanceRequest{
		FlowId:     sessionId,
		InstanceId: req.InstanceId,
		UpdateableInstanceInfo: instance_types.UpdateableInstanceInfo{
			SchemeId: &bootSchemaInfo.DisklessSchemaId,
			BootType: &req.StartMode,
			NetInfo: &instance_types.NetInfo{
				Ip:       req.Ip,
				Hostname: cloudBoxInfo.Name,
			},
			DefaultConfig: &instance_types.DefaultConfig{
				NetInfo: instance_types.NetInfo{
					Ip:       req.Ip,
					Hostname: cloudBoxInfo.Name,
				},
			},
		},
	}

	if err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).UpdateInstance(req.AreaId, instReq); err != nil {
		l.Logger.Errorf("[%s] diskless.UpdateInstance AreaId[%d] err:%+v", sessionId, req.AreaId, err)
		return nil, errorx.NewHandlerError(errorx.ServerErrorCode, "更新无盘实例失败")
	}

	if cloudBoxInfo.FirstStrategyId != req.FirstStrategyId || cloudBoxInfo.SecondStrategyId != req.SecondStrategyId {
		cloudClientInfo, _, err := table.T_TCdpCloudclientInfoService.Query(l.ctx, sessionId, fmt.Sprintf("cloudbox_mac:%s$status:%d", cloudBoxInfo.Mac, table.CloudClientStatusValid), nil, nil)
		if err != nil && err != gopublic.ErrNotExist {
			l.Logger.Errorf("[%s] Query T_TCdpCloudClientService failed. cloudbox_mac:%s err: %v", sessionId, cloudBoxInfo.Mac, err)
			return nil, errorx.NewDefaultCodeError("查询云盒信息失败")
		}

		if cloudClientInfo != nil {
			// 通知无盘规格调整
			strategy, _, err := table.T_TCdpResourceStrategyService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d$status:%d", req.FirstStrategyId, table.ResourceStrategyStatusValid), nil, nil)
			if err != nil {
				l.Logger.Errorf("[%s] Query T_TCdpResourceStrategyService failed. id:%d err: %v", sessionId, req.FirstStrategyId, err)
				return nil, errorx.NewDefaultCodeError("查询算力策略失败")
			}
			updateSeatReq := &proto.UpdateSeatRequest{
				FlowId:              sessionId,
				Id:                  int32(cloudClientInfo.DisklessSeatId),
				StreamSpecification: int32(strategy.InstPoolId),
			}
			_, err = diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).UpdateSeat(sessionId, cloudClientInfo.AreaId, updateSeatReq)
			if err != nil {
				l.Logger.Errorf("[%s] UpdateSeat failed. areaId[%d] updateSeatReq:%s err: %v", sessionId, cloudClientInfo.AreaId, helper.ToJSON(updateSeatReq), err)
				return nil, errorx.NewDefaultCodeError("更新无盘规格失败")
			}
			updateInfo := map[string]interface{}{
				"first_strategy_id":  req.FirstStrategyId,
				"second_strategy_id": req.SecondStrategyId,
				"update_time":        time.Now(),
				"update_by":          updateBy,
			}
			_, _, err = table.T_TCdpCloudclientInfoService.Update(l.ctx, sessionId, cloudClientInfo.Id, updateInfo)
			if err != nil {
				l.Logger.Errorf("[%s] Update T_TCdpCloudclientInfo failed. cloudclient_id:%d update_info:%s err: %v", sessionId, cloudClientInfo.Id, helper.ToJSON(updateInfo), err)
				return nil, errorx.NewDefaultCodeError("更新客户机信息失败")
			}
		}
	}

	newCloudBoxInfo, _, err := table.T_TCdpCloudboxInfoService.Update(l.ctx, sessionId, cloudBoxInfo.Id, map[string]any{
		"ip":                 req.Ip,
		"update_time":        time.Now(),
		"update_by":          updateBy,
		"boot_schema_id":     req.BootSchemaId,
		"boot_type":          req.StartMode,
		"first_strategy_id":  req.FirstStrategyId,
		"second_strategy_id": req.SecondStrategyId,
	})
	if err != nil {
		l.Logger.Errorf("[%s] Update T_TCdpCloudboxInfo failed. cloudbox_id:%d err: %v", sessionId, req.CloudBoxId, err)
		return nil, errorx.NewDefaultCodeError("更新云盒信息失败")
	}

	l.Logger.Infof("[%s] Update T_TCdpCloudboxInfo req.cloudbox_id:%d newCloudBoxInfo:%s", sessionId, req.CloudBoxId, helper.ToJSON(newCloudBoxInfo))

	return
}
