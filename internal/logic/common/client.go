package common

import (
	"context"
	"fmt"

	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/helper/diskless"
	"cdp-admin-service/internal/model/errorx"
	instance_types "cdp-admin-service/internal/proto/instance_service/types"
	"cdp-admin-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
)

func CheckCloudClientStreaming(ctx context.Context, svcCtx *svc.ServiceContext, cloudClient table.TCdpCloudclientInfo) (err error) {
	sessionId := helper.GetSessionId(ctx)
	// 校验云主机是否串流
	ipList := []string{cloudClient.HostIp}
	instListReq := &instance_types.ListInstancesRequestNew{
		Offset: 0,
		Length: 9999,
		Ips:    ipList,
	}
	Instancelist2, err := diskless.NewDisklessWebGateway(ctx, svcCtx).SearcchInstanceList(cloudClient.AreaId, sessionId, instListReq)
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] GetInstanceList failed. AreaId[%d]  err:%+v", sessionId, cloudClient.AreaId, err)
		return errorx.NewDefaultCodeError("查询云主机信息失败")
	}
	if len(Instancelist2) > 0 {
		instDetail := Instancelist2[0]
		// AssignStatus = 100: 占用中  PowerStatus = 1: 开机  BusinessStatus = 0: 正常
		if len(instDetail.BootMac) != 0 && instDetail.AssignStatus == 100 && instDetail.PowerStatus == 1 && instDetail.BusinessStatus == 0 {
			logx.WithContext(ctx).Errorf("[%s] GetInstanceList failed. AreaId[%d]  err:%+v", sessionId, cloudClient.AreaId, err)
			return errorx.NewDefaultCodeError("云主机串流中,不可更新")
		}
	}
	return nil
}

// // 添加座位
// func ClientAddSeat(ctx context.Context, svcCtx *svc.ServiceContext, clientType int32,
// 	box *table.TCdpCloudboxInfo, addSeatReq *proto.AddSeatRequest) (err error) {
// 	sessionId := helper.GetSessionId(ctx)
// 	if box != nil && box.Mac != "" {
// 		addSeatReq.LocalInstanceMac = box.Mac
// 		addSeatReq.LocalIp = box.Ip
// 		addSeatReq.LocalBootType = box.BootType
// 		addSeatReq.LocalInstanceMac = box.Mac
// 		boxSchema, _, err := table.T_TCdpBootSchemaInfoService.Query(ctx, sessionId, fmt.Sprintf("id:%d$status:%d", box.BootSchemaId, table.TCdpBootSchemaInfoStatusEnable), nil, nil)
// 		if err != nil {
// 			logx.WithContext(ctx).Errorf("[%s] GetBoxSchema failed.  err:%+v", sessionId, err)
// 			return errorx.NewDefaultCodeError("查询云盒信息失败")
// 		}
// 		firstStrategy, _, err := table.T_TCdpResourceStrategyService.Query(ctx, sessionId, fmt.Sprintf("id:%d$status:%d", box.FirstStrategyId, table.ResourceStrategyStatusValid), nil, nil)
// 		if err != nil {
// 			logx.WithContext(ctx).Errorf("[%s] GetFirstStrategy failed.  err:%+v", sessionId, err)
// 			return errorx.NewDefaultCodeError("查询云盒信息失败")
// 		}
// 		addSeatReq.LocalSchemeId = int32(boxSchema.DisklessSchemaId)
// 		addSeatReq.StreamSpecification = int32(firstStrategy.InstPoolId)
// 	}
// 	_, err = diskless.NewDisklessWebGateway(ctx, svcCtx).AddSeat(sessionId, int64(box.AreaId), addSeatReq)
// 	if err != nil {
// 		logx.WithContext(ctx).Errorf("[%s] AddSeat failed. AreaId[%d] addSeatReq[%s] err:%+v", sessionId, box.AreaId, helper.ToJSON(addSeatReq), err)
// 		return errorx.NewDefaultCodeError("添加座位失败")
// 	}

// 	return nil
// }

//  创建1.0座位
// func AddSeat1(ctx context.Context, svcCtx *svc.ServiceContext, bizId int64, areaId int64,
// 	addSeatReq *proto.AddSeatRequest) (disklessSeatId int32, err error) {
// 	sessionId := helper.GetSessionId(ctx)
// 	resp, err := diskless.NewDisklessWebGateway(ctx, svcCtx).AddSeat(sessionId, int64(areaId), addSeatReq)
// 	if err != nil {
// 		logx.WithContext(ctx).Errorf("[%s] AddSeat failed. AreaId[%d] addSeatReq[%s] err:%+v", sessionId, areaId, helper.ToJSON(addSeatReq), err)
// 		return 0, errorx.NewDefaultCodeError("添加座位失败")
// 	}
// 	logx.WithContext(ctx).Infof("[%s] AddSeat success. AreaId[%d] addSeatReq[%s]", sessionId, areaId, helper.ToJSON(addSeatReq))
// 	return resp.Id, nil
// }

// // 1.0 -> 2.0
// func Seat1To2(ctx context.Context, svcCtx *svc.ServiceContext, bizId int64, areaId int64,
// 	seatId int32, seatName string, strategyId int64, box *table.TCdpCloudboxInfo, instanceDetail *disklessType.InstanceDetail, ) (err error) {
// 	sessionId := helper.GetSessionId(ctx)

// }

//  2.0 -> 1.0

//  2.0

func UniqueMac(oldMac, newMac string, ctx context.Context) error {
	sessionId := helper.GetSessionId(ctx)
	if oldMac == newMac {
		return nil
	}
	//  1.0客户机 & 云盒 mac保持全局唯一
	cloudclientInfo, _, err := table.T_TCdpCloudclientInfoService.Query(ctx, sessionId, fmt.Sprintf("client_type:%d$cloudbox_mac:%s$status:%d", table.ClientType1, newMac, table.CloudClientStatusValid), nil, nil)
	if err != nil && err != gopublic.ErrNotExist {
		logx.WithContext(ctx).Errorf("[%s] Query T_TCdpCloudclientInfoService failed. oldMac[%s] newMac[%s] err: %v", sessionId, oldMac, newMac, err)
		return errorx.NewDefaultCodeError("查询云客户机信息失败")
	}
	if cloudclientInfo != nil {
		return errorx.NewDefaultCodeError(fmt.Sprintf("普通客户机mac[%s]已存在", newMac))
	}
	// 云盒 mac保持全局唯一
	cloudboxInfo, _, err := table.T_TCdpCloudboxInfoService.Query(ctx, sessionId, fmt.Sprintf("mac:%s$status:%d", newMac, table.CloudBoxStatusValid), nil, nil)
	if err != nil && err != gopublic.ErrNotExist {
		logx.WithContext(ctx).Errorf("[%s] Query T_TCdpCloudboxInfoService failed. oldMac[%s] newMac[%s] err: %v", sessionId, oldMac, newMac, err)
		return errorx.NewDefaultCodeError("查询云盒信息失败")
	}
	if cloudboxInfo != nil {
		return errorx.NewDefaultCodeError(fmt.Sprintf("云盒mac[%s]已存在", newMac))
	}
	return nil
}
