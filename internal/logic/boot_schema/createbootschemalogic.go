package boot_schema

import (
	"context"
	"encoding/json"
	"errors"
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
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
)

type CreateBootSchemaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateBootSchemaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateBootSchemaLogic {
	return &CreateBootSchemaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateBootSchemaLogic) CreateBootSchema(req *types.CreateBootSchemaReq) (resp *types.CreateBootSchemaResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)
	userName := helper.GetUserName(l.ctx)

	err = checkBootSchemaName("", req.Name, req.BizId, l.ctx)
	if err != nil {
		l.Logger.Errorf("[%s] checkBootSchemaName biz_id:%d, name:%s, err:%+v", sessionId, req.BizId, req.Name, err)
		return nil, err
	}

	biz, _, err := table.T_TCdpBizInfoService.Query(l.ctx, sessionId, fmt.Sprintf("biz_id:%d", req.BizId), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] table.T_TCdpBizInfoService.Query biz_id:%d, err:%+v", sessionId, req.BizId, err)
		return nil, errorx.NewDefaultError(errorx.BizNotFoundErrorCode)
	}

	//  无盘编排方案
	schemaConfig, resetSchemaConfig, err := getSchemaConfig(l.ctx, l.svcCtx, int64(biz.AreaId))
	if err != nil {
		l.Logger.Errorf("[%s] getSchemaConfig biz_id:%d, err:%+v", sessionId, req.BizId, err)
		return nil, err
	}
	schemaReq := make(map[string]any)
	err = json.Unmarshal([]byte(schemaConfig), &schemaReq)
	if err != nil {
		l.Logger.Errorf("[%s] json.Unmarshal l.svcCtx.Config.Schema.DefaultConfig, err:%+v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.ServerErrorCode)
	}

	schemaReq["flow_id"] = sessionId
	schemaReq["name"] = req.SchemaInfo.Name
	schemaReq["os_image_id"] = req.SchemaInfo.OsImageId
	schemaReq["os_snapshot_id"] = req.SchemaInfo.OsSnapshotId

	schema, err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).CreateSchema(int64(biz.AreaId), sessionId, schemaReq)
	if err != nil {
		l.Logger.Errorf("[%s] diskless.NewDisklessWebGateway.CreateSchema schemaId:%d, err:%+v", sessionId, schema.SchemeId, err)
		return nil, errorx.NewDefaultError(errorx.CreateSchemaFailedErrorCode)
	}

	l.Logger.Infof("[%s] diskless.NewDisklessWebGateway.CreateSchema req:%+v, schemaId:%d, success", sessionId, schemaReq, schema.SchemeId)

	if resetSchemaConfig != "" {
		resetSchemaReq := diskless.ResetSchemeImageBindReq{}
		err = json.Unmarshal([]byte(resetSchemaConfig), &resetSchemaReq)
		if err != nil {
			l.Logger.Errorf("[%s] json.Unmarshal l.svcCtx.Config.Schema.DefaultResetSchemaConfig, err:%+v", err)
			return nil, errorx.NewDefaultError(errorx.ServerErrorCode)
		}
		resetSchemaReq.FlowId = sessionId
		resetSchemaReq.SchemeId = schema.SchemeId
		for _, image := range resetSchemaReq.Images {
			image.SchemeId = schema.SchemeId
		}
		err = diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).ResetSchemeImageBind(int64(biz.AreaId), sessionId, resetSchemaReq)
		if err != nil {
			l.Logger.Errorf("[%s] ResetSchemeImageBind biz_id:%d, err:%+v", sessionId, req.BizId, err)
			return nil, errorx.NewDefaultError(errorx.ResetSchemeImageBindErrorCode)
		}
		l.Logger.Infof("[%s] ResetSchemeImageBind req:%+v, schemaId:%d, success", sessionId, resetSchemaReq, schema.SchemeId)
	}

	//  命令
	var scriptIds []string
	for _, script := range req.ScriptInfos {
		script, err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).CreateScript(int64(biz.AreaId), sessionId, diskless.CreateScriptReq{
			FlowId: sessionId,
			Name:   script.Name,
			Script: script.Script,
		})
		if err != nil {
			l.Logger.Errorf("[%s] diskless.NewDisklessWebGateway.CreateScript biz_id:%d, err:%+v", sessionId, req.BizId, err)
			return nil, errorx.NewDefaultError(errorx.CreateScriptErrorCode)
		}

		l.Logger.Infof("[%s] diskless.NewDisklessWebGateway.CreateScript biz_id:%d, scriptId:%d, success", sessionId, req.BizId, script.Script.Id)

		_, err = diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).ResetScriptBind(int64(biz.AreaId), sessionId, diskless.ResetScriptBindReq{
			FlowId:    sessionId,
			ScriptId:  script.Script.Id,
			SchemeIds: []int64{schema.SchemeId},
		})
		if err != nil {
			l.Logger.Errorf("[%s] ResetScriptBind biz_id:%d, err:%+v", sessionId, req.BizId, err)
			return nil, errorx.NewDefaultError(errorx.ResetScriptBindErrorCode)
		}

		l.Logger.Infof("[%s] ResetScriptBind biz_id:%d, scriptId:%d, success", sessionId, req.BizId, script.Script.Id)

		scriptIds = append(scriptIds, script.Script.Id)
	}
	//  启动方案
	bootSchema, _, err := table.T_TCdpBootSchemaInfoService.Insert(l.ctx, sessionId, table.TCdpBootSchemaInfo{
		BizId:            req.BizId,
		Name:             req.Name,
		AreaId:           int64(biz.AreaId),
		DisklessSchemaId: schema.SchemeId,
		BootCommandIds:   strings.Join(scriptIds, ","),
		OsImageId:        req.SchemaInfo.OsImageId,
		Status:           1,
		CreateBy:         userName,
		UpdateBy:         userName,
		CreateTime:       time.Now(),
		UpdateTime:       time.Now(),
		ModifyTime:       time.Now(),
	})
	if err != nil {
		l.Logger.Errorf("[%s] table.T_TCdpBootSchemaInfoService.Insert biz_id:%d, err:%+v", sessionId, req.BizId, err)
		return nil, errorx.NewDefaultError(errorx.CreateBootSchemaErrorCode)
	}
	l.Logger.Infof("[%s] table.T_TCdpBootSchemaInfoService.Insert biz_id:%d, success", sessionId, req.BizId)

	return &types.CreateBootSchemaResp{
		Id: bootSchema.Id,
	}, nil
}

func checkBootSchemaName(oldName string, newName string, bizId int64, ctx context.Context) error {
	sessionId := helper.GetSessionId(ctx)
	if oldName != newName {
		bootSchema, _, err := table.T_TCdpBootSchemaInfoService.Query(ctx, sessionId, fmt.Sprintf("biz_id:%d$name:%s$status:%d", bizId, newName, table.TCdpBootSchemaInfoStatusEnable), nil, nil)
		if err != nil && !errors.Is(err, gopublic.ErrNotExist) {
			logx.WithContext(ctx).Errorf("[%s] table.T_TCdpBootSchemaInfoService.Query biz_id:%d, name:%s err:%+v", sessionId, bizId, newName, err)
			return errorx.NewDefaultError(errorx.BootSchemaQueryFailedErrorCode)
		}
		if bootSchema != nil {
			return errorx.NewDefaultError(errorx.BootSchemaNameExistErrorCode)
		}
	}
	return nil
}

func getSchemaConfig(ctx context.Context, svcCtx *svc.ServiceContext, areaId int64) (string, string, error) {
	sessionId := helper.GetSessionId(ctx)
	areaInfo, _, err := table.T_TCdpAreaInfoService.Query(ctx, sessionId, fmt.Sprintf("area_id:%d$status:%d", areaId, table.AreaStatusEnable), nil, nil)
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] T_TCdpAreaInfoService Query err. AreaId[%d] err:%+v", sessionId, areaId, err)
		return "", "", errorx.NewDefaultError(errorx.AreaNotFound)
	}
	if areaInfo.SchemaConfig != "" {
		return areaInfo.SchemaConfig, areaInfo.ResetSchemaConfig, nil
	}
	return svcCtx.Config.Schema.DefaultSchemaConfig, svcCtx.Config.Schema.DefaultResetSchemaConfig, nil
}
