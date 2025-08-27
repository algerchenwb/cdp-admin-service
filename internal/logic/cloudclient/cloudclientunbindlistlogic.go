package cloudclient

import (
	"context"
	"fmt"

	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CloudClientUnbindListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCloudClientUnbindListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CloudClientUnbindListLogic {
	return &CloudClientUnbindListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CloudClientUnbindListLogic) CloudClientUnbindList(req *types.CloudClientUnbindListReq) (resp *types.CloudClientUnbindListResp, err error) {
	sessionId := helper.GetSessionId(l.ctx)
	resp = &types.CloudClientUnbindListResp{
		List: make([]types.CloudClientInfo, 0),
	}
	// 1.0客户机
	clients1, _, err := table.T_TCdpCloudclientInfoService.QueryAll(l.ctx, sessionId, fmt.Sprintf("biz_id:%d$area_id:%d$status:%d$client_type:%d", req.BizId, req.AreaId, table.CloudClientStatusValid, table.ClientType1), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] Query T_TCdpCloudclientInfo failed. bizId[%d] areaId[%d] err: %v", sessionId, req.BizId, req.AreaId, err)
		return nil, errorx.NewDefaultCodeError("查询客户机信息失败")
	}
	// 2.0客户机
	clients2, _, err := table.T_TCdpCloudclientInfoService.QueryAll(l.ctx, sessionId, fmt.Sprintf("biz_id:%d$area_id:%d$status:%d$client_type:%d$cloudbox_mac:%s", req.BizId, req.AreaId, table.CloudClientStatusValid, table.ClientType2, ""), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] Query T_TCdpCloudclientInfo failed. bizId[%d] areaId[%d] err: %v", sessionId, req.BizId, req.AreaId, err)
		return nil, errorx.NewDefaultCodeError("查询客户机信息失败")
	}
	for _, client := range clients1 {
		resp.List = append(resp.List, types.CloudClientInfo{
			CloudClientId:      int64(client.Id),
			CloudClientName:    client.Name,
			FirstStrategyId:    client.FirstStrategyId,
			FirstBootSchemaId:  client.FirstBootSchemaId,
			SecondStrategyId:   client.SecondStrategyId,
			SecondBootSchemaId: client.SecondBootSchemaId,
		})
	}
	for _, client := range clients2 {
		resp.List = append(resp.List, types.CloudClientInfo{
			CloudClientId:     int64(client.Id),
			CloudClientName:   client.Name,
			FirstStrategyId:   client.FirstStrategyId,
			FirstBootSchemaId: client.FirstBootSchemaId,
			SecondStrategyId:  client.SecondStrategyId,
		})
	}
	return
}
