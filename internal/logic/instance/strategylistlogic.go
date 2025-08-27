package instance

import (
	"context"
	"fmt"
	"strings"

	"cdp-admin-service/internal/helper"
	"cdp-admin-service/internal/helper/cdp_cache"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/helper/diskless"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/proto/instance_scheduler"
	instance_types "cdp-admin-service/internal/proto/instance_service/types"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
)

type StrategyListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStrategyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StrategyListLogic {
	return &StrategyListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StrategyListLogic) StrategyList(req *types.CommonPageRequest) (resp *types.StrategyListResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)
	areaIds := helper.GetAreaIds(l.ctx)

	req.CondList = append(req.CondList, fmt.Sprintf("area_id__in:%s", areaIds))
	req.CondList = append(req.CondList, fmt.Sprintf("status:%d", table.ResourceStrategyStatusValid))
	qry, err := helper.CheckCommQueryParam(l.ctx, sessionId, req.CondList, req.Orders, req.Sorts, table.TCdpResourceStrategy{})
	if err != nil {
		l.Logger.Errorf("[%s] CheckCommQueryParam failed, req:%+v err:%+v", sessionId, req, err)
		return nil, errorx.NewDefaultError(errorx.ParamErrorCode)
	}
	total, list, _, err := table.T_TCdpResourceStrategyService.QueryPage(l.ctx, sessionId, qry, req.Offset, req.Limit, req.Sorts, req.Orders)
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpResourceStrategyService err  qry:%s,  err: %v", sessionId, qry, err)
		return nil, errorx.NewDefaultError(errorx.ListStrategyErrorCode)
	}

	resp = &types.StrategyListResp{
		Total: total,
	}
	if len(list) == 0 {
		return resp, nil
	}

	var areaId int64 = 0
	resourceIds := make([]string, 0)
	poolIds := make([]int, 0)
	for _, v := range list {
		areaId = v.AreaId
		resourceIds = append(resourceIds, fmt.Sprintf("%d", v.Id))
		if v.InstPoolId != 0 {
			poolIds = append(poolIds, int(v.InstPoolId))
		}
	}

	// 查询无盘资源，策略
	listResReq := &instance_scheduler.SearchResourceRequest{
		FlowId:     sessionId,
		AreaType:   int32(areaId),
		Offset:     0,
		Length:     9999,
		Conditions: []string{fmt.Sprintf("resource_id__in:%s", strings.Join(resourceIds, ","))},
	}
	resList, err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).SearchResource(sessionId, areaId, listResReq)
	if err != nil {
		l.Logger.Errorf("[%s] diskless.SearchResource AreaId[%d] listResReq:%s err:%+v", sessionId, areaId, helper.ToJSON(listResReq), err)
		return nil, errorx.NewDefaultCodeError("查询无盘策略失败")
	}

	var mapRes = make(map[int64]*instance_scheduler.ResourceConfig, 0)
	for _, res := range resList {
		mapRes[res.ResourceId] = res
	}

	qry = fmt.Sprintf("inst_strategy_id__in:%s$status:%d", strings.Join(resourceIds, ","), table.BizStrategyStatusValid)
	bizStraegys, _, err := table.T_TCdpBizStrategyService.QueryAll(l.ctx, sessionId, qry, nil, nil)
	if err != nil && err != gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] T_TCdpBizStrategyService Query err. qry:%s err:%+v", sessionId, qry, err)
		return nil, errorx.NewDefaultError(errorx.ListStrategyErrorCode)
	}

	// 查询合约名称
	mapBizName := make(map[int64][]string)
	for _, bizStraegy := range bizStraegys {

		bizName := cdp_cache.GetBizName(l.ctx, sessionId, bizStraegy.BizId)
		if bizName != "" {
			mapBizName[bizStraegy.InstStrategyId] = append(mapBizName[bizStraegy.InstStrategyId], bizName)
		}
	}

	// 查询 资源池中实例的规格 poolIds

	// 调用无盘的接口
	instListReq := &instance_types.ListInstancesRequestNew{
		Offset:        0,
		Length:        9999,
		Specification: poolIds,
	}
	InstanceDetails, err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).SearcchInstanceList(areaId, sessionId, instListReq)
	if err != nil {
		l.Logger.Errorf("[%s] SearcchInstanceList areaId:%d, instListReq:%s err:%+v", sessionId, areaId, helper.ToJSON(instListReq), err)
		return nil, errorx.NewDefaultError(errorx.QueryInstanceFailedErrorCode)
	}

	var mapGPUs = make(map[int64][]string, 0)
	for _, inst := range InstanceDetails {
		if inst.HostInfo.Gpu != "" && !gopublic.StringInArray(inst.HostInfo.Gpu, mapGPUs[inst.Specification]) {
			mapGPUs[inst.Specification] = append(mapGPUs[inst.Specification], inst.HostInfo.Gpu)
		}
	}

	l.Logger.Infof("[%s] all strategy  mapGPUs:%s", sessionId, helper.ToJSON(mapGPUs))

	for _, v := range list {
		resourceInfo := new(instance_scheduler.ResourceConfig)
		if res, ok := mapRes[v.Id]; ok {
			resourceInfo = res
		} else {
			l.Logger.Errorf("[%s] SearchResource no this resourceId[%d] ResourceStrategy :%s", sessionId, v.Id, helper.ToJSON(v))
		}

		bizNameList := make([]string, 0)
		if _, ok := mapBizName[v.Id]; ok {
			bizNameList = mapBizName[v.Id]
		}

		resp.List = append(resp.List, types.StrategyInfo{
			Id:                  uint64(v.Id),
			Name:                v.Name,
			ApplicableLever:     uint64(v.ApplicableLever),
			OuterSpecId:         uint64(v.OuterSpecId),
			SpecId:              uint64(v.SpecId),
			InstPoolId:          uint64(v.InstPoolId),
			SpecName:            strings.Join(mapGPUs[v.InstPoolId], ","),
			TotalInstances:      uint64(v.TotalInstances),
			VlanId:              uint64(resourceInfo.Vlan),
			BootType:            int(resourceInfo.Mode),
			PreBootSchemaIdInfo: resourceInfo.AssignConfig,
			PreBootCount:        uint64(resourceInfo.Init),
			Remark:              v.Remark,
			Status:              int(v.Status),
			CreateBy:            v.CreateBy,
			UpdateBy:            v.UpdateBy,
			CreateTime:          v.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateTime:          v.UpdateTime.Format("2006-01-02 15:04:05"),
			ModifyTime:          v.ModifyTime.Format("2006-01-02 15:04:05"),
			BizNameList:         bizNameList,
		})

	}

	return
}
