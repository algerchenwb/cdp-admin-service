package boot_schema

import (
	"context"
	"fmt"
	"strings"
	"time"

	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/helper/diskless"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"cdp-admin-service/internal/model/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateBootSchemaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateBootSchemaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBootSchemaLogic {
	return &UpdateBootSchemaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateBootSchemaLogic) UpdateBootSchema(req *types.UpdateBootSchemaReq) (resp *types.CommonNilJson, err error) {

	sessionId := helper.GetSessionId(l.ctx)
	userName := helper.GetUserName(l.ctx)

	bootSchema, _, err := table.T_TCdpBootSchemaInfoService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d", req.Id), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] query boot schema failed, err:%+v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.BootSchemaNotFoundErrorCode)
	}

	err = checkBootSchemaName(bootSchema.Name, req.Name, bootSchema.BizId, l.ctx)
	if err != nil {
		l.Logger.Errorf("[%s] checkBootSchemaName biz_id:%d, name:%s, err:%+v", sessionId, bootSchema.BizId, req.Name, err)
		return nil, err
	}

	var destScriptIds map[string]struct{} = make(map[string]struct{})

	for _, scriptInfo := range req.ScriptInfos {
		destScriptIds[scriptInfo.Id] = struct{}{}
	}

	var deleteScriptIds []string

	for _, scriptId := range strings.Split(bootSchema.BootCommandIds, ",") {
		if _, ok := destScriptIds[scriptId]; !ok {
			deleteScriptIds = append(deleteScriptIds, scriptId)
		}
	}

	for _, scriptId := range deleteScriptIds {

		_, err = diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).DeleteScript(int64(bootSchema.AreaId), sessionId, diskless.DeleteScriptReq{
			FlowId: sessionId,
			Ids:    scriptId,
		})
		if err != nil {
			l.Logger.Errorf("[%s] diskless.NewDisklessWebGateway.DeleteScript id:%d, err:%+v", sessionId, req.Id, err)
			return nil, errorx.NewDefaultError(errorx.DeleteScriptFailedErrorCode)
		}

		l.Logger.Infof("[%s] diskless.NewDisklessWebGateway.DeleteScript id:%d, success", sessionId, req.Id)
	}

	schemaReq := make(map[string]any)
	schemaReq["flow_id"] = sessionId
	schemaReq["name"] = req.SchemaInfo.Name
	schemaReq["os_image_id"] = req.SchemaInfo.OsImageId
	schemaReq["os_snapshot_id"] = req.SchemaInfo.OsSnapshotId
	schemaReq["scheme_id"] = bootSchema.DisklessSchemaId

	schemaId, err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).UpdateSchema(int64(bootSchema.AreaId), sessionId, schemaReq)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("[%s] UpdateSchema schema:%+v, err:%+v", sessionId, schemaId, err)
		return nil, errorx.NewDefaultError(errorx.UpdateSchemaErrorCode)
	}

	l.Logger.Infof("[%s] diskless.UpdateSchema success req:%+v, schemaId:%d ", sessionId, schemaReq, schemaId.SchemeId)

	var scriptIds []string
	for _, scriptInfo := range req.ScriptInfos {

		if scriptInfo.Id == "" {

			script, err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).CreateScript(bootSchema.AreaId, sessionId, diskless.CreateScriptReq{
				FlowId:     sessionId,
				Name:       scriptInfo.Name,
				Script:     scriptInfo.Script,
				ScriptType: scriptInfo.ScriptType,
			})
			if err != nil {
				l.Logger.Errorf("[%s] diskless.NewDisklessWebGateway.CreateScript id:%d, err:%+v", sessionId, req.Id, err)
				return nil, errorx.NewDefaultError(errorx.CreateScriptErrorCode)
			}

			l.Logger.Infof("[%s] diskless.NewDisklessWebGateway.CreateScript id:%d, scriptId:%d, success", sessionId, req.Id, script.Script.Id)

			_, err = diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).ResetScriptBind(bootSchema.AreaId, sessionId, diskless.ResetScriptBindReq{
				FlowId:    sessionId,
				SchemeIds: []int64{bootSchema.DisklessSchemaId},
				ScriptId:  script.Script.Id,
			})
			if err != nil {
				l.Logger.Errorf("[%s] diskless.NewDisklessWebGateway.ResetScriptBind id:%d, err:%+v", sessionId, req.Id, err)
				return nil, errorx.NewDefaultError(errorx.ResetScriptBindErrorCode)
			}

			l.Logger.Infof("[%s] diskless.NewDisklessWebGateway.ResetScriptBind id:%d, scriptId:%d, success", sessionId, req.Id, script.Script.Id)

			scriptIds = append(scriptIds, script.Script.Id)

		} else {

			_, err = diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).UpdateScript(bootSchema.AreaId, sessionId, diskless.UpdateScriptReq{
				FlowId:     sessionId,
				Id:         scriptInfo.Id,
				Name:       scriptInfo.Name,
				Script:     scriptInfo.Script,
				ScriptType: scriptInfo.ScriptType,
			})
			if err != nil {
				l.Logger.Errorf("[%s] diskless.NewDisklessWebGateway.UpdateScript id:%d, err:%+v", sessionId, req.Id, err)
				return nil, errorx.NewDefaultError(errorx.UpdateScriptErrorCode)
			}

			l.Logger.Infof("[%s] diskless.NewDisklessWebGateway.UpdateScript id:%d, scriptId:%d, success", sessionId, req.Id, scriptInfo.Id)

			_, err = diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).ResetScriptBind(bootSchema.AreaId, sessionId, diskless.ResetScriptBindReq{
				FlowId:    sessionId,
				SchemeIds: []int64{bootSchema.DisklessSchemaId},
				ScriptId:  scriptInfo.Id,
			})
			if err != nil {
				l.Logger.Errorf("[%s] diskless.NewDisklessWebGateway.ResetScriptBind id:%d, err:%+v", sessionId, req.Id, err)
				return nil, errorx.NewDefaultError(errorx.ResetScriptBindErrorCode)
			}

			l.Logger.Infof("[%s] diskless.NewDisklessWebGateway.ResetScriptBind id:%d, scriptId:%d, success", sessionId, req.Id, scriptInfo.Id)

			scriptIds = append(scriptIds, scriptInfo.Id)
		}
	}

	_, _, err = table.T_TCdpBootSchemaInfoService.Update(l.ctx, sessionId, req.Id, map[string]interface{}{
		"name":             req.Name,
		"update_time":      time.Now(),
		"update_by":        userName,
		"os_image_id":      req.SchemaInfo.OsImageId,
		"boot_command_ids": strings.Join(scriptIds, ","),
	})
	if err != nil {
		l.Logger.Errorf("[%s] update boot schema failed, err:%+v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.UpdateBootSchemaErrorCode)
	}
	return &types.CommonNilJson{}, nil
}
