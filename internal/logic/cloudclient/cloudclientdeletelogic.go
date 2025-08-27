package cloudclient

import (
	"context"
	"fmt"
	"time"

	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/helper/diskless"
	"cdp-admin-service/internal/model/errorx"
	instance_types "cdp-admin-service/internal/proto/instance_service/types"
	proto "cdp-admin-service/internal/proto/location_seat_service"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
)

type CloudClientDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCloudClientDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CloudClientDeleteLogic {
	return &CloudClientDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CloudClientDeleteLogic) CloudClientDelete(req *types.CloudClientDeleteReq) (resp *types.CloudClientDeleteResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)
	resp = new(types.CloudClientDeleteResp)

	// 调用无盘实例的接口 通过IP查询云主机实例信息
	// 不支持批量删除 只支持单个删除
	for _, item := range req.List {

		cloudclientInfo, _, err := table.T_TCdpCloudclientInfoService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d$status:1", item.CloudClientId), nil, nil)
		if err != nil && err != gopublic.ErrNotExist {
			l.Logger.Errorf("[%s] Query T_TCdpCloudclientInfo failed. CloudClientId[%s]  err: %v", sessionId, item.CloudClientId, err)
			return nil, errorx.NewDefaultCodeError("查询客户机信息失败")
		}

		if err == gopublic.ErrNotExist {
			l.Logger.Errorf("[%s] Query T_TCdpCloudclientInfo failed. CloudClientId[%s] err: %v", sessionId, item.CloudClientId, err)
			return nil, errorx.NewHandlerError(errorx.ServerErrorCode, "客户机不存在")
		}

		ipList := []string{cloudclientInfo.HostIp}
		instListReq := &instance_types.ListInstancesRequestNew{
			Offset: 0,
			Length: 9999,
			Ips:    ipList,
		}
		Instancelist2, err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).SearcchInstanceList(cloudclientInfo.AreaId, sessionId, instListReq)
		if err != nil {
			l.Logger.Errorf("[%s] GetInstanceList failed. AreaId[%d]  err:%+v", sessionId, cloudclientInfo.AreaId, err)
			return nil, errorx.NewDefaultCodeError("查询云主机信息失败")
		}

		if len(Instancelist2) > 0 {
			instDetail := Instancelist2[0]
			// AssignStatus = 100: 占用中  PowerStatus = 1: 开机  BusinessStatus = 0: 正常
			if len(instDetail.BootMac) != 0 && instDetail.AssignStatus == 100 && instDetail.PowerStatus == 1 && instDetail.BusinessStatus == 0 {
				l.Logger.Errorf("[%s] GetInstanceList failed. AreaId[%d]  err:%+v", sessionId, cloudclientInfo.AreaId, err)
				return nil, errorx.NewDefaultCodeError("云主机串流中,不可删除")
			}
		}
		deleteSeatReq := &proto.DeleteSeatRequest{
			FlowId: sessionId,
			Id:     int32(cloudclientInfo.DisklessSeatId),
		}
		_, err = diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).DeleteSeat(sessionId, cloudclientInfo.AreaId, deleteSeatReq)
		if err != nil {
			l.Logger.Errorf("[%s] DeleteSeat failed. AreaId[%d]  err:%+v", sessionId, cloudclientInfo.AreaId, err)
			return nil, errorx.NewDefaultCodeError("删除客户机失败")
		}
		l.Logger.Infof("[%s] DeleteSeat success. AreaId[%d]  req:%s", sessionId, cloudclientInfo.AreaId, helper.ToJSON(deleteSeatReq))

		delTime := time.Now().Format("20060102150405")
		newCloudClientInfo, _, err := table.T_TCdpCloudclientInfoService.Update(l.ctx, sessionId, cloudclientInfo.Id, map[string]interface{}{
			"status":       0,
			"name":         fmt.Sprintf("%s-del-%s", cloudclientInfo.Name, delTime),
			"cloudbox_mac": fmt.Sprintf("%s-del-%s", cloudclientInfo.CloudboxMac, delTime),
			"update_by":    helper.GetUserName(l.ctx),
			"update_time":  time.Now(),
		})

		l.Logger.Infof("[%s] Update T_TCdpCloudclientInfo success. CloudClientId[%d]  req:%s, resp:%s", sessionId, item.CloudClientId, helper.ToJSON(req), helper.ToJSON(newCloudClientInfo))
	}

	return
}
