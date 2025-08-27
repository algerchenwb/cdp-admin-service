package boot_schema

import (
	"context"
	"fmt"
	"strings"

	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/helper/diskless"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListBootSchemaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListBootSchemaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListBootSchemaLogic {
	return &ListBootSchemaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListBootSchemaLogic) ListBootSchema(req *types.CommonPageRequest) (resp *types.ListBootSchemaResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)
	areaIds := helper.GetAreaIds(l.ctx)
	req.CondList = append(req.CondList, fmt.Sprintf("area_id__in:%s", areaIds))
	req.CondList = append(req.CondList, "status__ex:0")
	qry, err := helper.CheckCommQueryParam(l.ctx, sessionId, req.CondList, req.Orders, req.Sorts, table.TCdpBootSchemaInfo{})
	if err != nil {
		l.Logger.Errorf("[%s] CheckCommQueryParam failed, req:%+v err:%+v", sessionId, req, err)
		return nil, errorx.NewDefaultError(errorx.ParamErrorCode)
	}
	total, bootSchemas, _, err := table.T_TCdpBootSchemaInfoService.QueryPage(l.ctx, sessionId, qry, req.Offset, req.Limit, req.Sorts, req.Orders)
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpBootSchemaInfoService query failed, qry:%s err:%+v", sessionId, qry, err)
		return nil, errorx.NewDefaultError(errorx.BootSchemaNotFoundErrorCode)
	}
	if total == 0 {
		return &types.ListBootSchemaResp{
			Total: 0,
			List:  make([]types.BootSchemaInfo, 0),
		}, nil
	}

	schemaIds := make([]string, 0)
	for _, bootSchema := range bootSchemas {
		schemaIds = append(schemaIds, fmt.Sprintf("%d", bootSchema.DisklessSchemaId))
	}

	areaId := bootSchemas[0].AreaId

	schemas, err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).ListSchemes(areaId, sessionId,
		0, 5000, "", "", strings.Join(schemaIds, ","))
	if err != nil {
		l.Logger.Errorf("[%s] diskless.NewDisklessWebGateway.ListScheme biz_id:%d, err:%+v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.ListSchemesFailedErrorCode)
	}

	var schemaMap = make(map[string]diskless.Scheme)
	for _, schema := range schemas.Schemes {
		schemaMap[fmt.Sprintf("%d", schema.Id)] = schema
	}

	var scriptIds []string
	for _, bootSchema := range bootSchemas {
		scriptIds = append(scriptIds, strings.Split(bootSchema.BootCommandIds, ",")...)
	}

	scripts, err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).ListScript(areaId, sessionId,
		diskless.ListScriptReq{
			FlowId:   sessionId,
			Offset:   0,
			Limit:    5000,
			Orders:   "",
			Sorts:    "",
			CondList: []string{fmt.Sprintf("id__in:%s", strings.Join(scriptIds, ","))},
		})
	if err != nil {
		l.Logger.Errorf("[%s] diskless.NewDisklessWebGateway.ListScripts  err:%+v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.ListScriptsFailedErrorCode)
	}

	l.Logger.Infof("[%s] diskless.NewDisklessWebGateway.ListScripts , scriptIds:%v, success", sessionId, scriptIds)

	var scriptMap = make(map[string]diskless.ScriptInfo)
	for _, script := range scripts.List {
		scriptMap[script.Id] = script
	}

	resp = &types.ListBootSchemaResp{
		Total: int64(total),
		List:  make([]types.BootSchemaInfo, 0),
	}
	for _, bootSchema := range bootSchemas {
		bootSchemaInfo := types.BootSchemaInfo{
			Id:         bootSchema.Id,
			Name:       bootSchema.Name,
			CreateTime: bootSchema.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateTime: bootSchema.UpdateTime.Format("2006-01-02 15:04:05"),
			BizId:      bootSchema.BizId,
			CreateBy:   bootSchema.CreateBy,
			UpdateBy:   bootSchema.UpdateBy,
			ModifyTime: bootSchema.ModifyTime.Format("2006-01-02 15:04:05"),
			Status:     int32(bootSchema.Status),
			AreaId:     int32(bootSchema.AreaId),
			Remark:     bootSchema.Remark,
		}

		schema, ok := schemaMap[fmt.Sprintf("%d", bootSchema.DisklessSchemaId)]
		if ok {
			bootSchemaInfo.SchemaInfo = types.SchemaInfo{
				Id:        schema.Id,
				Name:      schema.Name,
				OsImageId: schema.OsImageId,
			}
		}
		for _, scriptId := range strings.Split(bootSchema.BootCommandIds, ",") {
			script, ok := scriptMap[scriptId]
			if ok {
				bootSchemaInfo.ScriptInfos = append(bootSchemaInfo.ScriptInfos, types.ScriptInfo{
					Id:         script.Id,
					Name:       script.Name,
					Script:     script.Script,
					ScriptType: script.ScriptType,
				})
			}
		}
		resp.List = append(resp.List, bootSchemaInfo)
	}

	return
}
