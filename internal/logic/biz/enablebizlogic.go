package biz

import (
	"context"
	"fmt"
	"time"

	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EnableBizLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEnableBizLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EnableBizLogic {
	return &EnableBizLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EnableBizLogic) EnableBiz(req *types.EnableBizReq) (resp *types.EnableBizResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)
	userName := helper.GetUserName(l.ctx)
	resp = &types.EnableBizResp{}

	if req.BizId == 0 || req.Platform == 0 {
		return nil, errorx.NewDefaultError(errorx.ParamErrorCode)
	}

	biz, _, err := table.T_TCdpBizInfoService.Query(l.ctx, sessionId, fmt.Sprintf("biz_id:%d", req.BizId), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpBizInfoService query biz[%d] info failed, err: %v", sessionId, req.BizId, err)
		return nil, errorx.NewDefaultError(errorx.QueryBizFailedErrorCode)
	}

	switch biz.Status {
	case table.BizStatusDeleted:
		l.Logger.Errorf("[%s]  biz status is deleted, status: %d", sessionId, biz.Status)
		return nil, errorx.NewDefaultError(errorx.BizAlreadyDeletedErrorCode)

	case table.BizStatusWaitWorking:
		err = l.bizConfigCompleted(biz)
		if err != nil {
			l.Logger.Errorf("[%s] biz config completed failed, bizId[%d] err: %v", sessionId, req.BizId, err)
			return nil, err
		}

		l.Logger.Infof("[%s] biz config completed, bizId[%d] status: %d", sessionId, req.BizId, table.BizStatusWorking)

		var status int32
		if req.Platform == table.PlatformConstruction {
			status = table.BizStatusOnline
		} else if req.Platform == table.PlatformPower {
			status = table.BizStatusWorking
		}

		if _, _, err = table.T_TCdpBizInfoService.Update(l.ctx, sessionId, biz.Id, map[string]interface{}{
			"status":      status,
			"update_by":   userName,
			"update_time": time.Now(),
		}); err != nil {
			l.Logger.Errorf("[%s] T_TCdpBizInfoService Update err. id[%d] err:%+v", sessionId, req.BizId, err)
			return nil, errorx.NewDefaultError(errorx.EnableBizFailedErrorCode)
		}

	case table.BizStatusWorking:

		if req.Platform != table.PlatformConstruction {
			return nil, errorx.NewDefaultError(errorx.NotPlatformConstructionAuitRefusedErrorCode)
		}

		err = l.bizWorkingEnable(biz)
		if err != nil {
			return nil, errorx.NewDefaultError(errorx.EnableBizFailedErrorCode)
		}

	case table.BizStatusOnline:

		l.Logger.Errorf("[%s] biz [%d] status is online, status: %d", req.BizId, biz.Status)
		return nil, errorx.NewDefaultError(errorx.EnableBizFailedErrorCode)

	default:

		l.Logger.Errorf("[%s] biz[%s] status is not working, status: %d", biz.Status)
		return nil, errorx.NewDefaultError(errorx.EnableBizFailedErrorCode)
	}

	return
}

func (l *EnableBizLogic) bizWaitWorkingEnable(biz *table.TCdpBizInfo) error {

	userName := helper.GetUserName(l.ctx)
	sessionId := helper.GetSessionId(l.ctx)

	if biz.VlanId == 0 {
		return errorx.NewDefaultError(errorx.VlanNotConfigErrorCode)
	}

	total, _, _, err := table.T_TCdpBizStrategyService.QueryPage(l.ctx, sessionId, fmt.Sprintf("biz_id:%d$status:1", biz.BizId), 0, 1, nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] query biz strategy failed, bizId[%d] err: %v", sessionId, biz.BizId, err)
		return errorx.NewDefaultError(errorx.StrategyNotConfigErrorCode)
	}
	if total == 0 {
		l.Logger.Errorf("[%s] biz not config  strategy failed, bizId[%d]", sessionId, biz.BizId)
		return errorx.NewDefaultError(errorx.BizStrategyNotConfigErrorCode)
	}

	// todo 后面可以需要校验
	// if biz.Serverinfo == "" || biz.Serverinfo == "{}" {
	// 	return errorx.NewDefaultError(errorx.ServerNotConfigErrorCode)
	// } else {
	// 	var serverInfo map[string]interface{} = make(map[string]interface{})
	// 	err := json.Unmarshal([]byte(biz.Serverinfo), &serverInfo)
	// 	if err != nil {
	// 		l.Logger.Errorf(fmt.Sprintf("unmarshal server info failed, err: %v", err))
	// 		return errorx.NewDefaultError(errorx.ServerInvalidErrorCode)
	// 	}

	// }

	total, _, _, err = table.T_TCdpBootSchemaInfoService.QueryPage(l.ctx, sessionId, fmt.Sprintf("biz_id:%d", biz.BizId), 0, 1, nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpBootSchemaInfoService query failed, bizId[%d] err: %v", sessionId, biz.BizId, err)
		return errorx.NewDefaultError(errorx.BootSchemaNotFoundErrorCode)
	}
	if total == 0 {
		return errorx.NewDefaultError(errorx.BootSchemaNotFoundErrorCode)
	}

	newBizInfo, _, err := table.T_TCdpBizInfoService.Update(l.ctx, sessionId, biz.Id, map[string]interface{}{
		"status":      table.BizStatusWorking,
		"update_by":   userName,
		"update_time": time.Now(),
	})
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpBizInfoService, update biz info failed, bizId[%d] err: %v", sessionId, biz.Id, err)
		return errorx.NewDefaultError(errorx.EnableBizFailedErrorCode)
	}

	l.Logger.Infof("[%s] bizWaitWorkingEnable T_TCdpBizInfoService Update success. bizId[%d] newBizInfo:%s", sessionId, biz.BizId, helper.ToJSON(newBizInfo))

	return nil

}

func (l *EnableBizLogic) bizWorkingEnable(biz *table.TCdpBizInfo) error {

	userName := helper.GetUserName(l.ctx)
	sessionId := helper.GetSessionId(l.ctx)

	newBizInfo, _, err := table.T_TCdpBizInfoService.Update(l.ctx, sessionId, biz.Id, map[string]interface{}{
		"status":      table.BizStatusOnline,
		"update_by":   userName,
		"update_time": time.Now(),
	})
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpBizInfoService update biz info failed, bizId[%s] err: %v", sessionId, biz.BizId, err)
		return errorx.NewDefaultError(errorx.EnableBizFailedErrorCode)
	}

	l.Logger.Infof("[%s] bizWorkingEnable T_TCdpBizInfoService Update success. bizId[%d] newBizInfo:%s", sessionId, biz.BizId, helper.ToJSON(newBizInfo))

	return nil
}

// 租户是否配置完成
func (l *EnableBizLogic) bizConfigCompleted(biz *table.TCdpBizInfo) error {

	sessionId := helper.GetSessionId(l.ctx)

	if biz.VlanId == 0 {
		return errorx.NewDefaultError(errorx.VlanNotConfigErrorCode)
	}

	total, _, _, err := table.T_TCdpBizStrategyService.QueryPage(l.ctx, sessionId, fmt.Sprintf("biz_id:%d$status:1", biz.BizId), 0, 1, nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpBizStrategyService query failed, bizId[%d] err: %v", sessionId, biz.BizId, err)
		return errorx.NewDefaultError(errorx.StrategyNotConfigErrorCode)
	}
	if total == 0 {
		l.Logger.Errorf("[%s] biz not config  strategy failed, bizId[%d]", sessionId, biz.BizId)
		return errorx.NewDefaultError(errorx.BizStrategyNotConfigErrorCode)
	}

	// todo 后面可以需要校验
	// if biz.Serverinfo == "" || biz.Serverinfo == "{}" {
	// 	return errorx.NewDefaultError(errorx.ServerNotConfigErrorCode)
	// } else {
	// 	var serverInfo map[string]interface{} = make(map[string]interface{})
	// 	err := json.Unmarshal([]byte(biz.Serverinfo), &serverInfo)
	// 	if err != nil {
	// 		l.Logger.Errorf(fmt.Sprintf("unmarshal server info failed, err: %v", err))
	// 		return errorx.NewDefaultError(errorx.ServerInvalidErrorCode)
	// 	}

	// }

	return nil
}
