package boot_schema

import (
	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/helper/diskless"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"
	"context"
	"fmt"
	"time"

	"cdp-admin-service/internal/model/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteBootSchemaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteBootSchemaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteBootSchemaLogic {
	return &DeleteBootSchemaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteBootSchemaLogic) DeleteBootSchema(req *types.DeleteBootSchemaReq) (resp *types.CommonNilJson, err error) {

	sessionId := helper.GetSessionId(l.ctx)
	userName := helper.GetUserName(l.ctx)

	if err := l.checkDeleteBootSchema(req.Id); err != nil {
		return nil, err
	}

	bootSchema, _, err := table.T_TCdpBootSchemaInfoService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d", req.Id), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpBootSchemaInfoService.Query id:%d, err:%+v", sessionId, req.Id, err)
		return nil, errorx.NewDefaultError(errorx.BootSchemaNotFoundErrorCode)
	}

	_, err = diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).DeleteScript(int64(bootSchema.AreaId), sessionId, diskless.DeleteScriptReq{
		FlowId: sessionId,
		Ids:    bootSchema.BootCommandIds,
	})
	if err != nil {
		l.Logger.Errorf("[%s] diskless.NewDisklessWebGateway.DeleteScript id:%d, err:%+v", sessionId, req.Id, err)
		return nil, errorx.NewDefaultError(errorx.DeleteScriptFailedErrorCode)
	}

	l.Logger.Infof("[%s] diskless.NewDisklessWebGateway.DeleteScript id:%d, success", sessionId, req.Id)

	_, err = diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).DeleteScheme(int64(bootSchema.AreaId), sessionId, diskless.DeleteSchemeReq{
		FlowId:   sessionId,
		SchemeId: bootSchema.DisklessSchemaId,
	})
	if err != nil {
		l.Logger.Errorf("[%s] diskless.NewDisklessWebGateway.DeleteScheme id:%d, err:%+v", sessionId, req.Id, err)
		return nil, errorx.NewDefaultError(errorx.DeleteSchemeFailedErrorCode)
	}

	l.Logger.Infof("[%s] diskless.NewDisklessWebGateway.DeleteScheme id:%d, success", sessionId, req.Id)

	bootSchemaInfo, _, err := table.T_TCdpBootSchemaInfoService.Update(l.ctx, sessionId, req.Id, map[string]interface{}{
		"status":      0,
		"update_by":   userName,
		"update_time": time.Now(),
	})
	if err != nil {
		l.Logger.Errorf("[%s] table.T_TCdpBootSchemaInfoService.Update id:%d, err:%+v", sessionId, req.Id, err)
		return nil, errorx.NewDefaultError(errorx.UpdateBootSchemaErrorCode)
	}
	l.Logger.Infof("[%s] T_TCdpBootSchemaInfoService.Update id:%d, bootSchemaInfo: %s success", sessionId, req.Id, helper.ToJSON(bootSchemaInfo))
	return
}

func (l *DeleteBootSchemaLogic) checkDeleteBootSchema(id int64) (err error) {
	sessionId := helper.GetSessionId(l.ctx)

	if id == 0 {
		return errorx.NewDefaultError(errorx.ParamErrorCode)
	}

	total, _, _, err := table.T_TCdpCloudclientInfoService.QueryPage(l.ctx, sessionId, fmt.Sprintf("first_boot_schema_id:%d$status:%d", id, table.CloudClientStatusValid), 0, 1, nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpCloudclientInfoService.QueryPage id:%d, err:%+v", sessionId, id, err)
		return errorx.NewDefaultCodeError("查询客户机信息失败")
	}
	if total > 0 {
		return errorx.NewDefaultCodeError("客户机已使用该开机方案，不能删除")
	}
	total, _, _, err = table.T_TCdpCloudclientInfoService.QueryPage(l.ctx, sessionId, fmt.Sprintf("second_boot_schema_id:%d$status:%d", id, table.CloudClientStatusValid), 0, 1, nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpCloudclientInfoService.QueryPage id:%d, err:%+v", sessionId, id, err)
		return errorx.NewDefaultCodeError("查询客户机信息失败")
	}
	if total > 0 {
		return errorx.NewDefaultCodeError("客户机已使用该开机方案，不能删除")
	}

	return nil
}
