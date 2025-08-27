package cloudclient

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/helper/diskless"
	"cdp-admin-service/internal/helper/saas"
	"cdp-admin-service/internal/model/errorx"
	instance_types "cdp-admin-service/internal/proto/instance_service/types"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
)

type CloudClientSetAdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCloudClientSetAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CloudClientSetAdminLogic {
	return &CloudClientSetAdminLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CloudClientSetAdminLogic) CloudClientSetAdmin(req *types.CloudClientSetAdminReq) (resp *types.CloudClientSetAdminResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)
	resp = new(types.CloudClientSetAdminResp)

	cloudclientInfo, _, err := table.T_TCdpCloudclientInfoService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d$status:1", req.CloudClientId), nil, nil)
	if err != nil && err != gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] Query T_TCdpCloudclientInfo failed. CloudClientId[%s]  err: %v", sessionId, req.CloudClientId, err)
		return nil, errorx.NewDefaultCodeError("查询客户机信息失败")
	}

	if err == gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] Query T_TCdpCloudclientInfo failed. CloudClientId[%s] err: %v", sessionId, req.CloudClientId, err)
		return nil, errorx.NewHandlerError(errorx.ServerErrorCode, "客户机不存在")
	}

	// 检测云盒是否已经超管
	var macList = []string{cloudclientInfo.CloudboxMac}
	Instancelist, err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).GetInstanceList(req.AreaId, sessionId, macList)
	if err != nil {
		l.Logger.Errorf("[%s] GetInstanceList failed AreaId[%d] maclist:%s err:%+v", sessionId, req.AreaId, helper.ToJSON(macList), err)
		return nil, errorx.NewHandlerError(errorx.ServerErrorCode, "查询无盘实例失败")
	}

	if len(Instancelist) == 0 {
		return nil, errorx.NewHandlerError(errorx.ServerErrorCode, "云盒实例不存在")
	}
	InstanceInfo := Instancelist[0]
	if InstanceInfo.UserMode == instance_types.AdminUser || InstanceInfo.UserMode == instance_types.BindAdminUser {
		return nil, errorx.NewHandlerError(errorx.ServerErrorCode, "该设备已是超管状态")
	}

	if cloudclientInfo.Vmid != 0 {

		instListReq := &instance_types.ListInstancesRequestNew{
			Offset:      0,
			Length:      9999,
			InstanceIds: []int{int(cloudclientInfo.Vmid)},
		}
		InstanceDetails, err1 := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).
			SearcchInstanceList(req.AreaId, sessionId, instListReq)
		if err1 != nil {
			l.Logger.Errorf("[%s] SearcchInstanceList areaId:%d, InstanceId:%d err:%+v", sessionId, req.AreaId, cloudclientInfo.Vmid, err1)
			return nil, errorx.NewDefaultCodeError("查询实例失败")
		}
		if len(InstanceDetails) == 0 {
			l.Logger.Errorf("[%s] SearcchInstanceList len = 0 areaId:%d, InstanceId:%d err:%+v", sessionId, req.AreaId, cloudclientInfo.Vmid, err1)
			return nil, errorx.NewDefaultCodeError("没有查询无盘的实例")
		}

		if InstanceDetails[0].UserMode == 0 { // 实例已经开超管
			l.Logger.Errorf("[%s] insatnce is admin .UserMode == 0  areaId:%d, InstanceId:%d err:%+v", sessionId, req.AreaId, cloudclientInfo.Vmid, err1)
			return nil, errorx.NewDefaultCodeError("客户机串流的主机已经在实例管理开了超管状态")
		}
	}

	if cloudclientInfo.ClientType == table.ClientType1 {

		// 1.0 客户机直接调用设备设置超管
		// 调用无盘超管接口
		setAdminReq := &instance_types.SetAdminRequest{
			FlowID:     sessionId,
			AppID:      "diskless-aggregator",
			InstanceID: InstanceInfo.Id,
			UserMode:   instance_types.BindAdminUser, // 设置被调度绑定的超管模式
		}
		if err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).SetAdmin(req.AreaId, setAdminReq); err != nil {
			l.Logger.Errorf("[%s] SetAdmin failed  req:%s,  err:%+v", sessionId, helper.ToJSON(setAdminReq), err)
			return nil, errorx.NewHandlerError(errorx.ServerErrorCode, "设置超管失败")
		}

		var configInfo string = ""
		if _, _, err := table.T_TCdpCloudclientInfoService.Update(l.ctx, sessionId, cloudclientInfo.Id, map[string]interface{}{
			"admin_state": int(instance_types.BindAdminUser),
			"config_info": configInfo,
			"update_by":   helper.GetUserName(l.ctx),
			"update_time": time.Now(),
		}); err != nil {
			l.Logger.Errorf("[%s] Update T_TCdpCloudclientInfo failed. CloudClientId[%s]  req:%s, err:%+v", sessionId, req.CloudClientId, helper.ToJSON(req), err)
			return nil, errorx.NewDefaultCodeError("更新客户机信息失败")
		}

		l.Logger.Infof("[%s] SetAdmin 1.0 client success. cloudclientInfo[%s] ", sessionId, cloudclientInfo.CloudboxMac)

		return
	}

	// 2.0 客户机直接调用设备设置超管
	instListReq := &instance_types.ListInstancesRequestNew{
		Offset: 0,
		Length: 9999,
		Ips:    []string{cloudclientInfo.HostIp},
	}
	Instancelist2, err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).SearcchInstanceList(req.AreaId, sessionId, instListReq)
	if err != nil {
		l.Logger.Errorf("[%s] GetInstanceList failed. AreaId[%d]  err:%+v", sessionId, req.AreaId, err)
	}

	// 是否为串流
	if len(Instancelist2) > 0 && Instancelist2[0].Id != 0 { // 串流中 云主机的mac 地址
		instanceId := Instancelist2[0].Id

		// 云客机中记录超管的状态
		configInfo, err := json.Marshal(saas.EsportRoomConfigInfo{
			VmId: int(instanceId),
		})
		if err != nil {
			return nil, errorx.NewHandlerError(errorx.ServerErrorCode, "序列化失败")
		}

		newCloudClientInfo, _, err := table.T_TCdpCloudclientInfoService.Update(l.ctx, sessionId, cloudclientInfo.Id, map[string]interface{}{
			"admin_state": int(instance_types.BindAdminUser),
			"config_info": string(configInfo),
			"update_by":   helper.GetUserName(l.ctx),
			"update_time": time.Now(),
		})

		if err != nil {
			l.Logger.Errorf("[%s] Update T_TCdpCloudclientInfo failed. CloudClientId[%s]  req:%s, err:%+v", sessionId, req.CloudClientId, helper.ToJSON(req), err)
			return nil, errorx.NewDefaultCodeError("更新客户机信息失败")
		}

		l.Logger.Infof("[%s] Update T_TCdpCloudclientInfo success. CloudClientId[%d]  req:%s, resp:%s", sessionId, req.CloudClientId, helper.ToJSON(req), helper.ToJSON(newCloudClientInfo))

		managerStatus := diskless.MANAGE_STATUS_DEBUG
		instanceStatusInfo := instance_types.UpdateInstanceStatusRequest{
			InstanceID: int64(instanceId),
			FlowID:     sessionId,
			UpdateableInstanceStatusInfo: instance_types.UpdateableInstanceStatusInfo{
				ManageStatus: &managerStatus,
			},
		}

		if err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).UpdateInstanceStatus(sessionId, req.AreaId, instanceStatusInfo); err != nil {
			l.Logger.Errorf("[%s] UpdateInstanceStatus failed. req:%s, err:%+v", sessionId, helper.ToJSON(instanceStatusInfo), err)
			return nil, errorx.NewHandlerError(errorx.ServerErrorCode, "更新无盘实例状态失败")
		}

		// 调用无盘超管接口
		setAdminReq := &instance_types.SetAdminRequest{
			FlowID:     sessionId,
			AppID:      "diskless-aggregator",
			InstanceID: instanceId,
			UserMode:   instance_types.BindAdminUser, // 设置被调度绑定的超管模式
		}
		if err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).SetAdmin(req.AreaId, setAdminReq); err != nil {
			l.Logger.Errorf("[%s] SetAdmin failed  req:%s,  err:%+v", sessionId, helper.ToJSON(setAdminReq), err)
			return nil, errorx.NewHandlerError(errorx.ServerErrorCode, "设置超管失败")
		}

	} else {

		// 只在云客机中记录超管的状态 等待分配成功的回调
		var configInfo string = ""
		newCloudClientInfo, _, err := table.T_TCdpCloudclientInfoService.Update(l.ctx, sessionId, cloudclientInfo.Id, map[string]interface{}{
			"admin_state": int(instance_types.BindAdminUser),
			"config_info": configInfo,
			"update_by":   helper.GetUserName(l.ctx),
			"update_time": time.Now(),
		})

		if err != nil {
			l.Logger.Errorf("[%s] Update T_TCdpCloudclientInfo failed. CloudClientId[%s]  req:%s, err:%+v", sessionId, req.CloudClientId, helper.ToJSON(req), err)
			return nil, errorx.NewDefaultCodeError("更新客户机信息失败")
		}
		l.Logger.Infof("[%s] Update T_TCdpCloudclientInfo success. CloudClientId[%d]  req:%s, resp:%s", sessionId, req.CloudClientId, helper.ToJSON(req), helper.ToJSON(newCloudClientInfo))

	}

	return
}
