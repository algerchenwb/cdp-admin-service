package cloudclient

import (
	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/logic/common"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"
	"context"
	"fmt"
	"time"

	"cdp-admin-service/internal/helper/diskless"
	proto "cdp-admin-service/internal/proto/location_seat_service"

	"github.com/zeromicro/go-zero/core/logx"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
)

type CloudClientBatchUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCloudClientBatchUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CloudClientBatchUpdateLogic {
	return &CloudClientBatchUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CloudClientBatchUpdateLogic) CloudClientBatchUpdate(req *types.CloudClientBatchUpdateReq) (resp *types.CloudClientBatchUpdateResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)
	userName := helper.GetUserName(l.ctx)

	if !helper.CheckAreaId(l.ctx, fmt.Sprintf("%d", req.AreaId)) {
		return nil, errorx.NewDefaultCodeError("当前用户没有权限操作该区域")
	}

	resp = &types.CloudClientBatchUpdateResp{}

	cloudClients, _, err := table.T_TCdpCloudclientInfoService.QueryAll(l.ctx, sessionId, fmt.Sprintf("id__in:%s$status:%d", helper.SliceToString(req.CloudClientIds), table.CloudClientStatusValid), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] Query T_TCdpCloudclientInfo failed. CloudClientIds[%v]  err: %v", sessionId, req.CloudClientIds, err)
		return nil, errorx.NewDefaultCodeError("查询客户机信息失败")
	}

	var updateInfo map[string]any = make(map[string]any)
	var bootSchema *table.TCdpBootSchemaInfo

	if req.BootSchemaId != 0 {

		updateInfo["first_boot_schema_id"] = req.BootSchemaId
		updateInfo["update_by"] = userName
		updateInfo["update_time"] = time.Now()

		bootSchema, _, err = table.T_TCdpBootSchemaInfoService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d$biz_id:%d$status:%d", req.BootSchemaId, req.BizId, table.TCdpBootSchemaInfoStatusEnable), nil, nil)
		if err != nil {
			l.Logger.Errorf("[%s] Query T_TCdpBootSchemaInfo failed. BootSchemaId[%d]  err: %v", sessionId, req.BootSchemaId, err)
			return nil, errorx.NewDefaultCodeError("查询启动方案信息失败")
		}
	}

	for _, cloudClient := range cloudClients {

		err = common.CheckCloudClientStreaming(l.ctx, l.svcCtx, cloudClient)
		if err != nil {
			resp.FailedItem = append(resp.FailedItem, types.CloudClientBatchUpdateItem{
				CloudClientId: int64(cloudClient.Id),
				ErrorMsg:      err.Error(),
			})
			continue
		}

		if req.BootSchemaId != cloudClient.FirstBootSchemaId {
			updateSeatReq := &proto.UpdateSeatRequest{
				FlowId: sessionId,
				Id:     int32(cloudClient.DisklessSeatId),
			}
			if cloudClient.ClientType == table.ClientType1 {
				updateSeatReq.LocalSchemeId = int32(bootSchema.DisklessSchemaId)
			} else if cloudClient.ClientType == table.ClientType2 {
				updateSeatReq.StreamSchemeId = int32(bootSchema.DisklessSchemaId)
			}
			_, err = diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).UpdateSeat(sessionId, int64(cloudClient.AreaId), updateSeatReq)
			if err != nil {
				l.Logger.Errorf("[%s] UpdateSeat failed. CloudClientId[%d] err: %v", sessionId, cloudClient.Id, err)
				resp.FailedItem = append(resp.FailedItem, types.CloudClientBatchUpdateItem{
					CloudClientId: int64(cloudClient.Id),
					ErrorMsg:      "更新客户机失败",
				})
				continue
			}
		}
		_, _, err = table.T_TCdpCloudclientInfoService.Update(l.ctx, sessionId, cloudClient.Id, updateInfo)
		if err != nil {
			l.Logger.Errorf("[%s] Update T_TCdpCloudclientInfo failed. CloudClientId[%d] updateInfo[%s]  err: %v", sessionId, cloudClient.Id, gopublic.ToJSON(updateInfo), err)
			resp.FailedItem = append(resp.FailedItem, types.CloudClientBatchUpdateItem{
				CloudClientId: int64(cloudClient.Id),
				ErrorMsg:      "更新客户机失败",
			})
			continue
		}
		resp.SuccessItem = append(resp.SuccessItem, types.CloudClientBatchUpdateItem{
			CloudClientId: int64(cloudClient.Id),
		})
	}

	return
}
