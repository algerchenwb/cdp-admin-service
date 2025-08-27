package cloudclient

import (
	"context"
	"fmt"
	"time"

	"cdp-admin-service/internal/helper"
	"cdp-admin-service/internal/helper/cdp_cache"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/helper/diskless"
	"cdp-admin-service/internal/logic/common"
	"cdp-admin-service/internal/model/errorx"
	proto "cdp-admin-service/internal/proto/location_seat_service"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
)

type CloudClientAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCloudClientAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CloudClientAddLogic {
	return &CloudClientAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CloudClientAddLogic) CloudClientAdd(req *types.CloudClientAddReq) (resp *types.CloudClientAddResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)
	updateBy := helper.GetUserName(l.ctx)

	if req.CloudHostName == "" || len(req.CloudHostName) > 15 {
		return nil, errorx.NewDefaultCodeError("云主机名不能为空且不能超过15个字符")
	}
	if req.ClientType == table.ClientType1 && req.CloudBoxMAC == "" {
		return nil, errorx.NewDefaultCodeError("普通客户机类型MAC不能为空")
	}

	resp = new(types.CloudClientAddResp)

	qry := fmt.Sprintf("area_id:%d$cloudbox_mac:%s$status:1", req.AreaId, req.CloudBoxMAC)
	cloudclientInfo, _, err := table.T_TCdpCloudclientInfoService.Query(l.ctx, sessionId, qry, nil, nil)
	if err != nil && err != gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] Query T_TCdpCloudclientInfo failed. CloudHostName[%s] cloudboxMac[%s]  err: %v", sessionId, req.CloudHostName, req.CloudHostIP, err)
		return nil, errorx.NewDefaultCodeError("查询客户机信息失败")
	}

	if err == nil && cloudclientInfo != nil && len(cloudclientInfo.CloudboxMac) > 0 {
		if cloudclientInfo.ClientType == table.ClientType1 {
			return nil, errorx.NewDefaultCodeError(fmt.Sprintf("普通客户机[%s]已经被绑定", req.CloudBoxMAC))
		}
		if cloudclientInfo.ClientType == table.ClientType2 {
			return nil, errorx.NewDefaultCodeError(fmt.Sprintf("云盒[%s]已经被绑定", req.CloudBoxMAC))
		}
	}

	qry = fmt.Sprintf("name:%s$status:1", req.CloudHostName)
	cloudclientInfo, _, err = table.T_TCdpCloudclientInfoService.Query(l.ctx, sessionId, qry, nil, nil)
	if err != nil && err != gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] Query T_TCdpCloudclientInfo failed. CloudHostName[%s] cloudboxMac[%s]  err: %v", sessionId, req.CloudHostName, req.CloudHostIP, err)
		return nil, errorx.NewDefaultCodeError("查询云客户机信息失败")
	}

	if err == nil && cloudclientInfo != nil && len(cloudclientInfo.Name) > 0 {
		l.Logger.Errorf("[%s] CloudClientAdd BizId[%d] CloudHostName[%s] CloudHostName[%s] already exist", sessionId, req.BizId, req.CloudHostName, req.CloudHostIP)
		return nil, errorx.NewDefaultCodeError(fmt.Sprintf("客户机[%s]信息已存在", req.CloudHostName))
	}

	//  校验租户客户机数量
	biz := cdp_cache.GetBizCache(l.ctx, sessionId, req.BizId)
	if biz == nil {
		return nil, errorx.NewDefaultCodeError("租户不存在")
	}
	bizClientCount, _, _, err := table.T_TCdpCloudclientInfoService.QueryPage(l.ctx, sessionId, fmt.Sprintf("biz_id:%d$status:%d", req.BizId, table.CloudClientStatusValid), 0, 1, nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] Query T_TCdpCloudclientInfo failed. err: %v", sessionId, err)
		return nil, errorx.NewDefaultCodeError("查询云客户机信息失败")
	}
	if bizClientCount+1 > int(biz.ClientNumLimit) {
		l.Logger.Errorf("[%s] CloudClientAdd BizId[%d]  bizClientCount[%d] bizClientNumLimit[%d]", sessionId, req.BizId, bizClientCount, biz.ClientNumLimit)
		return nil, errorx.NewDefaultCodeError("租户客户机数量超过限制")
	}

	firstBootSchema := new(table.TCdpBootSchemaInfo)
	if req.FirstBootSchemaId > 0 {
		firstBootSchema, _, err = table.T_TCdpBootSchemaInfoService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d$status:%d", req.FirstBootSchemaId, table.TCdpBootSchemaInfoStatusEnable), nil, nil)
		if err != nil {
			l.Logger.Errorf("[%s] Query T_TCdpBootSchemaInfo failed. CloudHostName[%s] cloudboxMac[%s]  err: %v", sessionId, req.CloudHostName, req.CloudHostIP, err)
			return nil, errorx.NewDefaultCodeError("查询启动方案失败")
		}
	}

	if req.ClientType == table.ClientType1 {
		err = common.UniqueMac("", req.CloudBoxMAC, l.ctx)
		if err != nil {
			return nil, err
		}
	}

	box := &table.TCdpCloudboxInfo{}
	firstStrategy := &table.TCdpResourceStrategy{}
	if req.ClientType == table.ClientType2 {
		box, _, err = table.T_TCdpCloudboxInfoService.Query(l.ctx, sessionId, fmt.Sprintf("mac:%s$status:%d", req.CloudBoxMAC, table.CloudBoxStatusValid), nil, nil)
		if err != nil {
			l.Logger.Errorf("[%s] Query T_TCdpCloudboxInfo failed. CloudHostName[%s] cloudboxMac[%s]  err: %v", sessionId, req.CloudHostName, req.CloudHostIP, err)
			return nil, errorx.NewDefaultCodeError("查询云盒信息失败")
		}

		firstStrategy, _, err = table.T_TCdpResourceStrategyService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d$status:%d", box.FirstStrategyId, table.ResourceStrategyStatusValid), nil, nil)
		if err != nil {
			l.Logger.Errorf("[%s] Query T_TCdpResourceStrategy failed. CloudHostName[%s] cloudboxMac[%s]  err: %v", sessionId, req.CloudHostName, req.CloudHostIP, err)
			return nil, errorx.NewDefaultCodeError("查询策略信息失败")
		}

	}

	seatType := proto.SeatType_SeatTypeDisklessUnknown
	addSetReq := &proto.AddSeatRequest{}

	switch req.ClientType {
	case table.ClientType1:
		seatType = proto.SeatType_SeatTypeDisklessComputer
		addSetReq = &proto.AddSeatRequest{
			FlowId:           sessionId,
			Name:             req.CloudHostName,
			Type:             int32(seatType),
			LocationBizId:    int32(req.BizId),
			LocalInstanceMac: req.CloudBoxMAC,
			LocalBootType:    int32(proto.BootType_BOOTTYPE_DISKLESS_UPGRADE),
			LocalIp:          req.CloudHostIP,
			LocalSchemeId:    int32(req.FirstBootSchemaId),
			ManagerState:     int32(proto.SeatManagerState_SeatManagerStateEnable),
		}
	case table.ClientType2:
		seatType = proto.SeatType_SeatTypeBoxStreamCloud
		addSetReq = &proto.AddSeatRequest{
			FlowId:              sessionId,
			Name:                req.CloudHostName,
			Type:                int32(seatType),
			LocationBizId:       int32(req.BizId),
			LocalInstanceMac:    box.Mac,
			StreamIp:            req.CloudHostIP,
			StreamSchemeId:      int32(firstBootSchema.DisklessSchemaId),
			StreamSpecification: int32(firstStrategy.InstPoolId),
			ManagerState:        int32(proto.SeatManagerState_SeatManagerStateEnable),
		}
	}

	disklessSet, err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).AddSeat(sessionId, req.AreaId, addSetReq)
	if err != nil {
		l.Logger.Errorf("[%s] AddSeat failed. CloudHostName[%s] cloudboxMac[%s]  err: %v", sessionId, req.CloudHostName, req.CloudHostIP, err)
		return nil, errorx.NewDefaultCodeError("添加云客户机信息失败")
	}

	newCloudClientInfo, _, err := table.T_TCdpCloudclientInfoService.Insert(l.ctx, sessionId, &table.TCdpCloudclientInfo{
		Name:               req.CloudHostName,
		BizId:              req.BizId,
		AreaId:             req.AreaId,
		FirstStrategyId:    box.FirstStrategyId,
		SecondStrategyId:   box.SecondStrategyId,
		FirstBootSchemaId:  req.FirstBootSchemaId,
		SecondBootSchemaId: req.SecondBootSchemaId,
		CloudboxMac:        req.CloudBoxMAC,
		HostIp:             req.CloudHostIP,
		ConfigInfo:         "",
		AdminState:         0,
		Status:             1,
		Remark:             "",
		CreateBy:           updateBy,
		UpdateBy:           updateBy,
		UpdateTime:         time.Now(),
		CreateTime:         time.Now(),
		ModifyTime:         time.Now(),
		ClientType:         int32(req.ClientType),
		DisklessSeatId:     int64(disklessSet.Id),
	})
	if err != nil {
		l.Logger.Errorf("[%s] Insert T_TCdpCloudclientInfo failed. CloudHostName[%s] CloudHostIP[%s]  err: %v", sessionId, req.CloudHostName, req.CloudHostIP, err)
		return nil, errorx.NewDefaultCodeError("添加云客户机信息失败")
	}

	l.Logger.Infof("[%s] CloudClientAdd success. CloudHostName[%s] CloudHostIP[%s] newCloudClientInfo:%s", sessionId, req.CloudHostName, req.CloudHostIP, helper.ToJSON(newCloudClientInfo))

	return
}
