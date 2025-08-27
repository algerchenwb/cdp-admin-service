package cloudclient

import (
	"context"

	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CloudClientBindBoxLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCloudClientBindBoxLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CloudClientBindBoxLogic {
	return &CloudClientBindBoxLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 废弃 todo 删除掉
func (l *CloudClientBindBoxLogic) CloudClientBindBox(req *types.ClientBindBoxReq) (resp *types.ClientBindBoxResp, err error) {
	// sessionId := helper.GetSessionId(l.ctx)
	// userName := helper.GetUserName(l.ctx)

	// biz := cdp_cache.GetBizCache(l.ctx, sessionId, req.BizId)
	// if biz == nil {
	// 	l.Logger.Errorf("[%s] biz not found. bizId[%d]", sessionId, req.BizId)
	// 	return nil, errorx.NewDefaultError(errorx.BizNotFoundErrorCode)
	// }

	// client, _, err := table.T_TCdpCloudclientInfoService.Query(l.ctx, sessionId, fmt.Sprintf("name:%s$status:%d", req.ClientName, table.CloudClientStatusValid), nil, nil)
	// if err != nil {
	// 	l.Logger.Errorf("[%s] query client failed. clientName[%s] err:%+v", sessionId, req.ClientName, err)
	// 	return nil, errorx.NewDefaultError(errorx.QueryClientFailedErrorCode)
	// }

	// if client.CloudboxMac != "" {
	// 	l.Logger.Infof("[%s] client already bind box. clientId[%d] boxMac[%s]", sessionId, client.Id, client.CloudboxMac)
	// 	return nil, errorx.NewDefaultCodeError("客户机已绑定云盒")
	// }

	// box, _, err := table.T_TCdpCloudboxInfoService.Query(l.ctx, sessionId, fmt.Sprintf("mac:%s$status:%d", req.CloudBoxMAC, table.CloudBoxStatusValid), nil, nil)
	// if err != nil {
	// 	l.Logger.Errorf("[%s] query box failed. cloudBoxMAC[%s] err:%+v", sessionId, req.CloudBoxMAC, err)
	// 	return nil, errorx.NewDefaultError(errorx.QueryBoxFailedErrorCode)
	// }
	// err = common.ClientBindBoxDisklessClientBindBox( l.svcCtx, l.ctx,  *client, box, box.FirstStrategyId, box.BootSchema, instanceDetail)
	// if err != nil {
	// 	l.Logger.Errorf("[%s] ClientBindBox failed. clientId[%d] boxMac[%s] err:%+v", sessionId, client.Id, box.Mac, err)
	// 	return nil, errorx.NewDefaultError(errorx.UpdateClientFailedErrorCode)
	// }
	// l.Logger.Infof("[%s] 无盘 bind box. clientId[%d] boxMac[%s] success", sessionId, client.Id, box.Mac)
	// updateInfo := map[string]interface{}{
	// 	"cloudbox_mac":       box.Mac,
	// 	"first_strategy_id":  box.FirstStrategyId,
	// 	"second_strategy_id": box.SecondStrategyId,
	// 	"client_type":        table.ClientType2,
	// 	"update_by":          userName,
	// 	"update_time":        time.Now(),
	// }
	// _, _, err = table.T_TCdpCloudclientInfoService.Update(l.ctx, sessionId, client.Id, updateInfo)
	// if err != nil {
	// 	l.Logger.Errorf("[%s] update client failed. clientId[%d] updateInfo[%+v] err:%+v", sessionId, client.Id, updateInfo, err)
	// 	return nil, errorx.NewDefaultError(errorx.UpdateClientFailedErrorCode)
	// }
	// l.Logger.Infof("[%s] 平台 bind box. clientId[%d] boxMac[%s] success", sessionId, client.Id, box.Mac)

	return
}
