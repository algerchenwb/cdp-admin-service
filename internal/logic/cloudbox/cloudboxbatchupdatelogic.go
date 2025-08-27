package cloudbox

import (
	"context"
	"fmt"
	"strings"
	"time"

	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/helper/diskless"
	"cdp-admin-service/internal/logic/common"
	"cdp-admin-service/internal/model/errorx"
	proto "cdp-admin-service/internal/proto/location_seat_service"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CloudBoxBatchUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCloudBoxBatchUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CloudBoxBatchUpdateLogic {
	return &CloudBoxBatchUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CloudBoxBatchUpdateLogic) CloudBoxBatchUpdate(req *types.CloudBoxBatchUpdateReq) (resp *types.CloudBoxBatchUpdateResp, err error) {
	sessionId := helper.GetSessionId(l.ctx)
	userName := helper.GetUserName(l.ctx)

	if !helper.CheckAreaId(l.ctx, fmt.Sprintf("%d", req.AreaId)) {
		return nil, errorx.NewDefaultCodeError("无权限操作该节点数据")
	}
	if req.FirstStrategyId == req.SecondStrategyId {
		return nil, errorx.NewDefaultCodeError("主算力策略不能与从算力策略相同")
	}

	resp = &types.CloudBoxBatchUpdateResp{}

	boxIds := make([]string, 0)
	for _, item := range req.CloudBoxIds {
		boxIds = append(boxIds, fmt.Sprintf("%d", item))
	}
	cloudBoxInfos, _, err := table.T_TCdpCloudboxInfoService.QueryAll(l.ctx, sessionId, fmt.Sprintf("id__in:%s$status:%d$biz_id:%d", strings.Join(boxIds, ","), table.CloudBoxStatusValid, req.BizId), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpCloudboxInfoService QueryAll failed. cloudBoxIds[%s] err: %v", sessionId, strings.Join(boxIds, ","), err)
		return nil, errorx.NewDefaultCodeError(fmt.Sprintf("获取云盒信息失败: %s", err.Error()))
	}

	macs := make([]string, 0)
	for _, cloudBoxInfo := range cloudBoxInfos {
		macs = append(macs, cloudBoxInfo.Mac)
	}
	cloudClientInfos, _, err := table.T_TCdpCloudclientInfoService.QueryAll(l.ctx, sessionId, fmt.Sprintf("cloudbox_mac__in:%s$status:%d$biz_id:%d", strings.Join(macs, ","), table.CloudClientStatusValid, req.BizId), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpCloudclientInfoService QueryAll failed. cloudBoxIds[%s] err: %v", sessionId, strings.Join(boxIds, ","), err)
		return nil, errorx.NewDefaultCodeError(fmt.Sprintf("获取客户机信息失败: %s", err.Error()))
	}

	cloudClientMap := make(map[string]table.TCdpCloudclientInfo)
	for _, cloudClientInfo := range cloudClientInfos {
		cloudClientMap[cloudClientInfo.CloudboxMac] = cloudClientInfo
	}

	for _, cloudBoxInfo := range cloudBoxInfos {

		cloudClintInfo, ok := cloudClientMap[cloudBoxInfo.Mac]
		if ok {
			err = common.CheckCloudClientStreaming(l.ctx, l.svcCtx, cloudClintInfo)
			if err != nil {
				l.Logger.Errorf("[%s] CheckCloudClientStreaming failed. cloudClientId[%d] err: %v", sessionId, cloudClintInfo.Id, err)
				resp.FailedItem = append(resp.FailedItem, types.CloudBoxBatchUpdateItem{
					CloudBoxId: int64(cloudBoxInfo.Id),
					ErrorMsg:   err.Error(),
				})
				continue
			}
		}
		if err := checkStrategy(cloudBoxInfo, req.FirstStrategyId, req.SecondStrategyId); err != nil {
			l.Logger.Errorf("[%s] checkStrategy failed. cloudBoxId[%d] err: %v", sessionId, cloudBoxInfo.Id, err)
			resp.FailedItem = append(resp.FailedItem, types.CloudBoxBatchUpdateItem{
				CloudBoxId: int64(cloudBoxInfo.Id),
				ErrorMsg:   err.Error(),
			})
			continue
		}

		updateInfo := map[string]interface{}{
			"update_by":   userName,
			"update_time": time.Now(),
		}

		if req.FirstStrategyId != 0 && cloudBoxInfo.FirstStrategyId != req.FirstStrategyId {
			updateInfo["first_strategy_id"] = req.FirstStrategyId

			if cloudClintInfo.DisklessSeatId != 0 {
				firstStrategy, _, err := table.T_TCdpResourceStrategyService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d$status:%d", req.FirstStrategyId, table.ResourceStrategyStatusValid), nil, nil)
				if err != nil {
					l.Logger.Errorf("[%s] T_TCdpResourceStrategyService Query failed. strategyId[%d] err: %v", sessionId, req.FirstStrategyId, err)
					resp.FailedItem = append(resp.FailedItem, types.CloudBoxBatchUpdateItem{
						CloudBoxId: int64(cloudBoxInfo.Id),
						ErrorMsg:   err.Error(),
					})
					continue
				}

				_, err = diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).UpdateSeat(sessionId, int64(cloudBoxInfo.AreaId), &proto.UpdateSeatRequest{
					FlowId:              sessionId,
					Id:                  int32(cloudClintInfo.DisklessSeatId),
					StreamSpecification: int32(firstStrategy.InstPoolId),
				})
				if err != nil {
					l.Logger.Errorf("[%s] UpdateSeat failed. cloudClientId[%d] err: %v", sessionId, cloudClintInfo.Id, err)
					resp.FailedItem = append(resp.FailedItem, types.CloudBoxBatchUpdateItem{
						CloudBoxId: int64(cloudBoxInfo.Id),
						ErrorMsg:   err.Error(),
					})
					continue
				}
			}
		}

		if req.SecondStrategyId != 0 && cloudBoxInfo.SecondStrategyId != req.SecondStrategyId {
			updateInfo["second_strategy_id"] = req.SecondStrategyId
		}

		_, _, err = table.T_TCdpCloudboxInfoService.Update(l.ctx, sessionId, cloudBoxInfo.Id, updateInfo)
		if err != nil {
			l.Logger.Errorf("[%s] T_TCdpCloudboxInfoService Update failed. cloudBoxId[%d]updateInfo[%s] err: %v", sessionId, cloudBoxInfo.Id, helper.ToJSON(updateInfo), err)
			resp.FailedItem = append(resp.FailedItem, types.CloudBoxBatchUpdateItem{
				CloudBoxId: int64(cloudBoxInfo.Id),
				ErrorMsg:   err.Error(),
			})
			continue
		}

		// 更新客户机信息
		if cloudClintInfo.Id != 0 {
			_, _, err = table.T_TCdpCloudclientInfoService.Update(l.ctx, sessionId, cloudClintInfo.Id, updateInfo)
			if err != nil {
				l.Logger.Errorf("[%s] T_TCdpCloudclientInfoService Update failed. cloudClientId[%d]updateInfo[%s] err: %v", sessionId, cloudClintInfo.Id, helper.ToJSON(updateInfo), err)
				resp.FailedItem = append(resp.FailedItem, types.CloudBoxBatchUpdateItem{
					CloudBoxId: int64(cloudBoxInfo.Id),
					ErrorMsg:   err.Error(),
				})
			}
		}

		resp.SuccessItem = append(resp.SuccessItem, types.CloudBoxBatchUpdateItem{
			CloudBoxId: int64(cloudBoxInfo.Id),
		})
	}

	return
}

func checkStrategy(cloudBoxInfo table.TCdpCloudboxInfo, firstStrategyId int64, secondStrategyId int64) (err error) {

	if firstStrategyId != 0 {
		cloudBoxInfo.FirstStrategyId = firstStrategyId
	}
	if secondStrategyId != 0 {
		cloudBoxInfo.SecondStrategyId = secondStrategyId
	}

	if cloudBoxInfo.FirstStrategyId == cloudBoxInfo.SecondStrategyId {
		return errorx.NewDefaultCodeError("主算力策略不能与从算力策略相同")
	}

	return nil
}
