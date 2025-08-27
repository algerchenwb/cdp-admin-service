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

type CloudClientCancelAdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCloudClientCancelAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CloudClientCancelAdminLogic {
	return &CloudClientCancelAdminLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CloudClientCancelAdminLogic) CloudClientCancelAdmin(req *types.CloudClientCancelAdminReq) (resp *types.CloudClientCancelAdminResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)
	resp = new(types.CloudClientCancelAdminResp)

	cloudclientInfo, _, err := table.T_TCdpCloudclientInfoService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d$status:1", req.CloudClientId), nil, nil)
	if err != nil && err != gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] Query T_TCdpCloudclientInfo failed. CloudClientId[%s]  err: %v", sessionId, req.CloudClientId, err)
		return nil, errorx.NewDefaultCodeError("查询客户机信息失败")
	}

	if err == gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] Query T_TCdpCloudclientInfo failed. CloudClientId[%s] err: %v", sessionId, req.CloudClientId, err)
		return nil, errorx.NewHandlerError(errorx.ServerErrorCode, "客户机不存在")
	}

	if cloudclientInfo.AdminState == uint32(instance_types.BindAdminUser) { // 客户机超管

		var macList = []string{cloudclientInfo.CloudboxMac}
		Instancelist, err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).GetInstanceList(req.AreaId, sessionId, macList)
		if err != nil {
			l.Logger.Errorf("[%s] GetInstanceList failed AreaId[%d] maclist:%s err:%+v", sessionId, req.AreaId, helper.ToJSON(macList), err)
			return nil, errorx.NewHandlerError(errorx.ServerErrorCode, "查询无盘设备信息失败")
		}

		if len(Instancelist) == 0 {
			return nil, errorx.NewHandlerError(errorx.ServerErrorCode, "无盘设备信息不存在")
		}

		InstanceInfo := Instancelist[0]

		// 1.0 为设备Id{"vmid":11} 3 2 27
		instanceId := InstanceInfo.Id
		var configInfo saas.EsportRoomConfigInfo
		if cloudclientInfo.ClientType == table.ClientType2 && cloudclientInfo.ConfigInfo != "" {

			// 调用无盘超管接口, 实例关掉超管
			if err := json.Unmarshal([]byte(cloudclientInfo.ConfigInfo), &configInfo); err != nil {
				l.Logger.Errorf("[%s] Unmarshal ConfigInfo failed. CloudClientId[%s]  err: %v", sessionId, req.CloudClientId, err)
				return nil, errorx.NewDefaultCodeError("解析客户机配置信息失败")
			}
			if configInfo.VmId == 0 {
				l.Logger.Errorf("[%s] Unmarshal ConfigInfo failed.VmId=0 CloudClientId[%s]  err: %v", sessionId, req.CloudClientId, err)
				return nil, errorx.NewDefaultCodeError("解析客户机配置信息失败,没有实例ID")
			}

			instListReq := &instance_types.ListInstancesRequestNew{
				Offset:      0,
				Length:      9999,
				InstanceIds: []int{int(configInfo.VmId)},
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

			if InstanceDetails[0].UserMode == 0 { // 实例开的超管
				l.Logger.Errorf("[%s] insatnce is admin .UserMode == 0  areaId:%d, InstanceId:%d err:%+v", sessionId, req.AreaId, cloudclientInfo.Vmid, err1)
				return nil, errorx.NewDefaultCodeError("客户机串流的主机已经实例管理中开了超管，请在实例管理关闭超管")
			}

			instanceId = int64(configInfo.VmId) // 2.0 客户机 为云主机实例ID

			managerStatus := diskless.MANAGE_STATUS_AVAILABLE
			instanceStatusInfo := instance_types.UpdateInstanceStatusRequest{
				InstanceID: instanceId,
				FlowID:     sessionId,
				UpdateableInstanceStatusInfo: instance_types.UpdateableInstanceStatusInfo{
					ManageStatus: &managerStatus,
				},
			}

			if err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).UpdateInstanceStatus(sessionId, req.AreaId, instanceStatusInfo); err != nil {
				l.Logger.Errorf("[%s] UpdateInstanceStatus err. req:%s, err:%+v", sessionId, helper.ToJSON(instanceStatusInfo), err)
				return nil, errorx.NewHandlerError(errorx.ServerErrorCode, "更新客户机无盘实例状态失败")
			}

		}

		// 调用无盘超管接口, 实例关掉超管
		setAdminReq := &instance_types.SetAdminRequest{
			FlowID:     sessionId,
			AppID:      "diskless-aggregator",
			InstanceID: instanceId,
			UserMode:   instance_types.RegularUser2, // 设置关闭超管
		}

		if err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).SetAdmin(req.AreaId, setAdminReq); err != nil {
			l.Logger.Errorf("[%s] SetAdmin err, areaId[%d]  req:%s err:%+v", sessionId, req.AreaId, helper.ToJSON(setAdminReq), err)
			return nil, errorx.NewHandlerError(errorx.ServerErrorCode, err.Error())
		}

		// 调用制作镜像接口
		if req.OsVersion != "" && req.Name != "" {
			if err := diskless.CreateImageFromAreaInstance(l.ctx,
				l.svcCtx.Config.OutSide.DisklessCloudImageHost,
				sessionId,
				req.ImageId,
				req.Name,
				req.OsVersion,
				req.Remark,
				req.ManagerState,
				req.BizId,
				req.AreaId,
				instanceId,
				int32(req.FlattenFlag)); err != nil {

				l.Logger.Errorf("[%s] CreateImageFromAreaInstance err. req:%s, err:%+v", sessionId, helper.ToJSON(req), err)
				return nil, errorx.NewHandlerError(errorx.ServerErrorCode, "制作镜像失败")
			}
		}

	}

	newCloudClientInfo, _, err := table.T_TCdpCloudclientInfoService.Update(l.ctx, sessionId, cloudclientInfo.Id, map[string]interface{}{
		"admin_state": int(instance_types.RegularUser2),
		"config_info": "",
		"update_by":   helper.GetUserName(l.ctx),
		"update_time": time.Now(),
	})
	if err != nil {
		l.Logger.Errorf("[%s] Update T_TCdpCloudclientInfo failed. CloudClientId[%s]  err: %v", sessionId, req.CloudClientId, err)
		return nil, errorx.NewDefaultCodeError("更新客户机信息失败")
	}
	l.Logger.Infof("[%s] Update T_TCdpCloudclientInfo success. CloudClientId[%d]  req:%s, resp:%s", sessionId, req.CloudClientId, helper.ToJSON(req), helper.ToJSON(newCloudClientInfo))

	return
}
