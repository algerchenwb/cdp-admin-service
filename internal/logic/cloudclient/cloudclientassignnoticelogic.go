package cloudclient

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
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

type CloudClientAssignNoticeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCloudClientAssignNoticeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CloudClientAssignNoticeLogic {
	return &CloudClientAssignNoticeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CloudClientAssignNoticeLogic) CloudClientAssignNotice(req *types.CloudClientAssignNoticeReq) (resp *types.CloudClientAssignNoticeResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)
	//updateBy := helper.GetUserName(l.ctx)

	resp = new(types.CloudClientAssignNoticeResp)

	// 通知灼光平台的esport-diskless-aggregator 服务
	go func() {
		req1 := diskless.CloudClientAssignNoticeReq{
			MAC:           req.MAC,
			BizId:         req.BizId,
			AreaId:        req.AreaId,
			VmId:          req.VmId,
			FlowID:        req.FlowID,
			CloudHostName: req.CloudHostName,
		}

		diskless.CloudClientAssignNotice(l.ctx, l.svcCtx.Config, sessionId, req1)

	}()

	mac := strings.ToLower(req.MAC)
	qry := fmt.Sprintf("area_id:%d$cloudbox_mac:%s$status:1", req.AreaId, mac)
	if len(mac) != 17 {

		// 如果mac 不对就尝试用主机名查询
		l.Logger.Infof("[%s] FormatMacAddress failed. req:%s err: %v", sessionId, helper.ToJSON(req), err)
		if req.CloudHostName == "" {
			return nil, errorx.NewDefaultCodeError("主机名不能为空")
		}
		qry = fmt.Sprintf("area_id:%d$cloudbox_mac:%s$status:1", req.AreaId, mac)
	}

	cloudclientInfo, _, err := table.T_TCdpCloudclientInfoService.Query(l.ctx, sessionId, qry, nil, nil)
	if err != nil && err != gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] Query T_TCdpCloudclientInfo failed. qry:%s  err: %v", sessionId, qry, err)
		return nil, errorx.NewDefaultCodeError("查询客户机信息失败")
	}

	if err == gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] Query T_TCdpCloudclientInfo failed. req.MAC[%s] call diskless-aggregator service err: %v", sessionId, req.MAC, err)
		return nil, errorx.NewDefaultCodeError("客户机不存在")
	}

	if cloudclientInfo.BizId != req.BizId {
		l.Logger.Errorf("[%s] Query T_TCdpCloudclientInfo failed. req.MAC[%s] cloudclientInfo.BizId[%d] != req.BizId[%d]", sessionId, req.MAC, cloudclientInfo.BizId, req.BizId)
		return nil, errorx.NewDefaultCodeError("bizId不匹配")
	}

	updateData := map[string]any{}
	updateData["update_by"] = "instance-event-processor"
	updateData["update_time"] = time.Now()

	switch req.EventType {
	case types.InstanceEventOpen:
		// 实例开机，分配成功
		updateData["flow_id"] = req.FlowID
		updateData["vmid"] = req.VmId
	case types.InstanceEventClose:
		// 实例关机，释放
		updateData["flow_id"] = ""
		updateData["vmid"] = 0

	}

	if req.EventType == types.InstanceEventOpen || req.EventType == types.InstanceEventClose {
		newCloudClientInfo, _, err1 := table.T_TCdpCloudclientInfoService.Update(l.ctx, sessionId, cloudclientInfo.Id, updateData)
		if err1 != nil {
			l.Logger.Errorf("[%s] Update T_TCdpCloudclientInfo failed. req:%s  err: %v", sessionId, req, err1)
			return nil, errorx.NewDefaultCodeError("更新客户机信息失败")
		}

		l.Logger.Infof("[%s] Update T_TCdpCloudclientInfo success. req:%s  newCloudClientInfo: %v", sessionId, helper.ToJSON(req), newCloudClientInfo)

	}

	if req.EventType == types.InstanceEventClose {
		// 关机事件，不处理超管逻辑
		return
	}
	if cloudclientInfo.AdminState != uint32(instance_types.BindAdminUser) {
		l.Logger.Debugf("[%s] CloudClientAssignNotice CloudClient not BindAdminUser. req:%s", sessionId, helper.ToJSON(req))
		return
	}

	managerStatus := diskless.MANAGE_STATUS_DEBUG
	instanceStatusInfo := instance_types.UpdateInstanceStatusRequest{
		InstanceID: int64(req.VmId),
		FlowID:     sessionId,
		UpdateableInstanceStatusInfo: instance_types.UpdateableInstanceStatusInfo{
			ManageStatus: &managerStatus,
		},
	}

	if err = diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).UpdateInstanceStatus(sessionId, req.AreaId, instanceStatusInfo); err != nil {
		l.Logger.Errorf("[%s] UpdateInstanceStatus err. mac[%s] req:%s, err:%+v", sessionId, req.MAC, helper.ToJSON(instanceStatusInfo), err)
		return nil, errorx.NewDefaultCodeError("更新客户机无盘实例状态失败")
	}
	// 调用无盘超管接口
	setAdminReq := &instance_types.SetAdminRequest{
		FlowID:     req.FlowID,
		AppID:      "diskless-aggregator",
		InstanceID: req.VmId,
		UserMode:   instance_types.BindAdminUser, // 设置为 被调度绑定的超管模式
	}

	if err = diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).SetAdmin(req.AreaId, setAdminReq); err != nil {
		l.Logger.Errorf("[%s] diskless SetAdmin error, req:%s, err:%+v", helper.ToJSON(setAdminReq), err)
		return nil, errorx.NewDefaultCodeError("设置客户机超管失败")
	}

	// 云客机中记录超管的状态
	configInfo, err := json.Marshal(saas.EsportRoomConfigInfo{
		VmId: int(req.VmId),
	})
	if err != nil {
		l.Logger.Errorf("[%s] json.Marshal, req:%s, err:%+v", helper.ToJSON(setAdminReq), err)
		return nil, errorx.NewDefaultCodeError("序列化失败")
	}

	updateData["config_info"] = string(configInfo)
	newCloudClientInfo, _, err := table.T_TCdpCloudclientInfoService.Update(l.ctx, sessionId, cloudclientInfo.Id, updateData)
	if err != nil {
		l.Logger.Errorf("[%s] Update T_TCdpCloudclientInfo failed. req:%s  err: %v", sessionId, helper.ToJSON(req), err)
		return nil, errorx.NewDefaultCodeError("更新客户机信息失败")
	}

	l.Logger.Infof("[%s] Update T_TCdpCloudclientInfo success. req:%s  newCloudClientInfo: %v", sessionId, helper.ToJSON(req), newCloudClientInfo)

	return
}
