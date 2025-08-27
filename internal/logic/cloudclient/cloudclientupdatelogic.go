package cloudclient

import (
	"context"
	"fmt"
	"time"

	"cdp-admin-service/internal/helper"
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

type CloudClientUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCloudClientUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CloudClientUpdateLogic {
	return &CloudClientUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CloudClientUpdateLogic) CloudClientUpdate(req *types.CloudClientUpdateReq) (resp *types.CloudClientUpdateResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)

	resp = new(types.CloudClientUpdateResp)

	if req.AreaId == 0 || req.BizId == 0 {
		return nil, errorx.NewDefaultCodeError("参数错误")
	}

	if err := l.checkClient(req); err != nil {
		l.Logger.Errorf("[%s] checkClient err. req:%s, err: %v", sessionId, gopublic.ToJSON(req), err)
		return nil, err
	}

	if len(req.CloudHostName) > 15 {
		return nil, errorx.NewDefaultCodeError("云主机名不能为空且不能超过15个字符")
	}

	// 客户机
	cloudclientInfo, _, err := table.T_TCdpCloudclientInfoService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d$status:1", req.CloudClientId), nil, nil)
	if err != nil && err != gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] Query T_TCdpCloudclientInfo failed. CloudClientId[%d]  err: %v", sessionId, req.CloudClientId, err)
		return nil, errorx.NewDefaultCodeError("查询客户机信息失败")
	}

	if err == gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] Query T_TCdpCloudclientInfo failed. CloudClientId[%d] err: %v", sessionId, req.CloudClientId, err)
		return nil, errorx.NewHandlerError(errorx.ServerErrorCode, "客户机不存在")
	}

	//    2.0切换1.0 新增1.0客户机校验全局ma
	if cloudclientInfo.ClientType == table.ClientType2 && req.ClientType == table.ClientType1 {
		err = common.UniqueMac(cloudclientInfo.CloudboxMac, req.CloudBoxMAC, l.ctx)
		if err != nil {
			return nil, err
		}
	}

	// 云盒
	var box *table.TCdpCloudboxInfo = new(table.TCdpCloudboxInfo)
	if req.ClientType == table.ClientType2 && req.CloudBoxMAC != "" {
		box, _, err = table.T_TCdpCloudboxInfoService.Query(l.ctx, sessionId, fmt.Sprintf("mac:%s$status:%d", req.CloudBoxMAC, table.CloudBoxStatusValid), nil, nil)
		if err != nil {
			l.Logger.Errorf("[%s] Query T_TCdpCloudboxInfo failed. CloudBoxMAC[%s] err: %v", sessionId, req.CloudBoxMAC, err)
			return nil, errorx.NewDefaultCodeError("查询云盒信息失败")
		}
		// 切换云盒需要判断云盒是否被占用
		if box.Mac != cloudclientInfo.CloudboxMac {
			_cloudclientInfo, _, err := table.T_TCdpCloudclientInfoService.Query(l.ctx, sessionId, fmt.Sprintf("cloudbox_mac:%s$status:1", box.Mac), nil, nil)
			if err != nil && err != gopublic.ErrNotExist {
				l.Logger.Errorf("[%s] Query T_TCdpCloudclientInfo failed. CloudBoxMAC[%s] err: %v", sessionId, box.Mac, err)
				return nil, errorx.NewDefaultCodeError("查询客户机信息失败")
			}
			if _cloudclientInfo != nil && _cloudclientInfo.Id != cloudclientInfo.Id {
				l.Logger.Errorf("[%s] CloudClientUpdate failed. item[%s], cloudclientInfo[%s]", sessionId, gopublic.ToJSON(req), gopublic.ToJSON(cloudclientInfo))
				return nil, errorx.NewDefaultCodeError(fmt.Sprintf("当前云盒已绑定客户机【%s】", _cloudclientInfo.Name))
			}
		}
		if box.FirstStrategyId == 0 {
			return nil, errorx.NewDefaultCodeError("2.0客户机,绑定云盒主算力策略不能为空")
		}
	}

	// 实例
	// 校验云主机是否串流
	ipList := []string{cloudclientInfo.HostIp}
	instListReq := &instance_types.ListInstancesRequestNew{
		Offset: 0,
		Length: 9999,
		Ips:    ipList,
	}
	Instancelist2, err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).SearcchInstanceList(req.AreaId, sessionId, instListReq)
	if err != nil {
		l.Logger.Errorf("[%s] GetInstanceList failed. AreaId[%d]  err:%+v", sessionId, req.AreaId, err)
		return nil, errorx.NewDefaultCodeError("查询云主机信息失败")
	}
	var instDetail instance_types.InstanceDetail
	if len(Instancelist2) > 0 {
		instDetail = Instancelist2[0]
		// AssignStatus = 100: 占用中  PowerStatus = 1: 开机  BusinessStatus = 0: 正常
		if len(instDetail.BootMac) != 0 && instDetail.AssignStatus == 100 && instDetail.PowerStatus == 1 && instDetail.BusinessStatus == 0 {
			l.Logger.Errorf("[%s] GetInstanceList failed. AreaId[%d]  err:%+v", sessionId, req.AreaId, err)
			return nil, errorx.NewDefaultCodeError("云主机串流中,不可更新")
		}
	}

	// 校验客户机名称是否存在
	if req.CloudHostName != cloudclientInfo.Name {
		ccInfo, _, err1 := table.T_TCdpCloudclientInfoService.Query(l.ctx, sessionId, fmt.Sprintf("name:%s$area_id:%d$status:1$id__ne:%d", req.CloudHostName, req.AreaId, cloudclientInfo.Id), nil, nil)
		if err1 != nil && err1 != gopublic.ErrNotExist {
			l.Logger.Errorf("[%s] Query T_TCdpCloudclientInfo failed. CloudHostName[%s]  err: %v", sessionId, req.CloudHostName, err1)
			return nil, errorx.NewDefaultCodeError("查询客户机信息失败")
		}

		if err1 == nil && ccInfo != nil && len(ccInfo.Name) > 0 {
			l.Logger.Errorf("[%s] 客户机名称已存在 CloudHostName [%s] err: %v", sessionId, req.CloudHostName, err1)
			return nil, errorx.NewHandlerError(errorx.ServerErrorCode, "客户机名称已存在")
		}
	}

	err = l.updateSeat(sessionId, req.AreaId, req.BizId, req, cloudclientInfo, box)
	if err != nil {
		l.Logger.Errorf("[%s] CloudClientUpdate failed. item[%s], cloudclientInfo[%s]", sessionId, gopublic.ToJSON(req), gopublic.ToJSON(cloudclientInfo))
		return nil, err
	}
	l.Logger.Infof("[%s] 无盘CloudClientUpdate success. item[%s], cloudclientInfo[%s]", sessionId, gopublic.ToJSON(req), gopublic.ToJSON(cloudclientInfo))

	err = l.updateClientInfo(sessionId, req, box)
	if err != nil {
		l.Logger.Errorf("[%s] CloudClientUpdate failed. item[%s], cloudclientInfo[%s]", sessionId, gopublic.ToJSON(req), gopublic.ToJSON(cloudclientInfo))
		return nil, err
	}
	l.Logger.Infof("[%s] Update T_TCdpCloudclientInfo success. CloudClientId[%d], updateInfo[%s]", sessionId, req.CloudClientId, gopublic.ToJSON(req))

	return
}

func (l *CloudClientUpdateLogic) updateSeat(sessionId string, areaId int64, bizId int64, req *types.CloudClientUpdateReq,
	client *table.TCdpCloudclientInfo, box *table.TCdpCloudboxInfo) error {

	updateSetReq := &proto.UpdateSeatRequest{
		FlowId: sessionId,
		Id:     int32(client.DisklessSeatId),
		Name:   req.CloudHostName,

		LocationBizId: int32(bizId),
		ManagerState:  int32(proto.SeatManagerState_SeatManagerStateEnable),
	}
	firstBootSchema, _, err := table.T_TCdpBootSchemaInfoService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d$status:%d", req.FirstBootSchemaId, table.TCdpBootSchemaInfoStatusEnable), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] Query T_TCdpBootSchemaInfo failed. BootSchemaId[%d] err: %v", sessionId, req.FirstBootSchemaId, err)
		return errorx.NewDefaultCodeError("查询客户机启动方案信息失败")
	}
	switch req.ClientType {

	case table.ClientType1:
		if req.FirstBootSchemaId == 0 {
			l.Logger.Errorf("[%s] CloudClientUpdate failed. item[%s], cloudclientInfo[%s]", sessionId, gopublic.ToJSON(req), gopublic.ToJSON(client))
			return errorx.NewDefaultCodeError("客户机类型,启动方案ID不能为空")
		}

		updateSetReq.Type = int32(proto.SeatType_SeatTypeDisklessComputer)
		updateSetReq.LocalBootType = int32(proto.BootType_BOOTTYPE_DISKLESS_UPGRADE)
		updateSetReq.LocalInstanceMac = req.CloudBoxMAC
		updateSetReq.LocalIp = req.CloudHostIP
		updateSetReq.LocalSchemeId = int32(firstBootSchema.DisklessSchemaId)
		updateSetReq.StreamSpecification = diskless.DisklessZero
		updateSetReq.StreamIp = diskless.DisklessEmptyString
		updateSetReq.StreamSchemeId = diskless.DisklessZero
		updateSetReq.StreamInstanceMac = diskless.DisklessEmptyString

	case table.ClientType2:

		updateSetReq.StreamIp = req.CloudHostIP
		updateSetReq.Type = int32(proto.SeatType_SeatTypeBoxStreamCloud)
		updateSetReq.StreamSchemeId = int32(firstBootSchema.DisklessSchemaId)

		if box != nil && box.FirstStrategyId != 0 {
			boxStrategy, _, err := table.T_TCdpResourceStrategyService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d$status:%d", box.FirstStrategyId, table.ResourceStrategyStatusValid), nil, nil)
			if err != nil {
				l.Logger.Errorf("[%s] Query T_TCdpResourceStrategy failed. StrategyId[%d] err: %v", sessionId, box.FirstStrategyId, err)
				return errorx.NewDefaultCodeError("查询云盒策略信息失败")
			}
			updateSetReq.LocalInstanceMac = box.Mac
			updateSetReq.StreamSpecification = int32(boxStrategy.InstPoolId)

		} else {
			updateSetReq.LocalBootType = diskless.DisklessZero
			updateSetReq.LocalInstanceMac = diskless.DisklessEmptyString
			updateSetReq.LocalIp = diskless.DisklessEmptyString
			updateSetReq.LocalSchemeId = diskless.DisklessZero
			updateSetReq.StreamSpecification = diskless.DisklessZero
		}
	}

	_, err = diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).UpdateSeat(sessionId, areaId, updateSetReq)
	if err != nil {
		l.Logger.Errorf("[%s] UpdateSeat failed. err: %v", sessionId, err)
		return errorx.NewDefaultCodeError(fmt.Sprintf("更新无盘座位信息失败:%s", err.Error()))
	}
	return nil
}

func (l *CloudClientUpdateLogic) updateClientInfo(sessionId string, req *types.CloudClientUpdateReq, box *table.TCdpCloudboxInfo) error {

	updateBy := helper.GetUserName(l.ctx)
	updateInfo := map[string]any{
		"name":                  req.CloudHostName,
		"host_ip":               req.CloudHostIP,
		"client_type":           req.ClientType,
		"first_boot_schema_id":  req.FirstBootSchemaId,
		"second_boot_schema_id": req.SecondBootSchemaId,
		"cloudbox_mac":          "",
		"first_strategy_id":     0,
		"second_strategy_id":    0,
		"update_by":             updateBy,
		"update_time":           time.Now(),
	}
	if req.ClientType == table.ClientType2 && box != nil && box.Mac != "" {
		updateInfo["cloudbox_mac"] = box.Mac
		updateInfo["first_strategy_id"] = box.FirstStrategyId
		updateInfo["second_strategy_id"] = box.SecondStrategyId
	}
	if req.ClientType == table.ClientType1 {
		updateInfo["cloudbox_mac"] = req.CloudBoxMAC
	}

	_, _, err := table.T_TCdpCloudclientInfoService.Update(l.ctx, sessionId, req.CloudClientId, updateInfo)
	if err != nil {
		l.Logger.Errorf("[%s] Update T_TCdpCloudclientInfo failed. CloudClientId[%d], updateInfo[%s]  err: %v", sessionId, req.CloudClientId, gopublic.ToJSON(updateInfo), err)
		return errorx.NewDefaultCodeError("更新客户机信息失败")
	}
	return nil
}

func (l *CloudClientUpdateLogic) checkClient(req *types.CloudClientUpdateReq) error {
	switch req.ClientType {
	case table.ClientType1:
		if req.CloudBoxMAC == "" {
			return errorx.NewDefaultCodeError("1.0客户机,MAC不能为空")
		}
		if req.CloudHostName == "" {
			return errorx.NewDefaultCodeError("1.0客户机,名称不能为空")
		}
		if req.CloudHostIP == "" {
			return errorx.NewDefaultCodeError("1.0客户机,IP不能为空")
		}
		if req.FirstBootSchemaId == 0 {
			return errorx.NewDefaultCodeError("1.0客户机,启动方案不能为空")
		}
	case table.ClientType2:
		if req.CloudHostIP == "" {
			return errorx.NewDefaultCodeError("2.0客户机,IP不能为空")
		}
		if req.CloudHostName == "" {
			return errorx.NewDefaultCodeError("2.0客户机,名称不能为空")
		}
		if req.FirstBootSchemaId == 0 {
			return errorx.NewDefaultCodeError("2.0客户机,主启动方案不能为空")
		}

	}
	return nil
}
