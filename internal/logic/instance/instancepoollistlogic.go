package instance

import (
	"context"
	"fmt"
	"strings"

	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/helper/diskless"
	"cdp-admin-service/internal/model/errorx"
	instance_types "cdp-admin-service/internal/proto/instance_service/types"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
)

type InstancePoolListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInstancePoolListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InstancePoolListLogic {
	return &InstancePoolListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InstancePoolListLogic) InstancePoolList(req *types.CommonPageRequest) (resp *types.InstancePoolListResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)

	req.CondList = append(req.CondList, fmt.Sprintf("status:%d", table.InstancePoolStatusValid))
	qry, err := helper.CheckCommQueryParam(l.ctx, sessionId, req.CondList, req.Orders, req.Sorts, table.TCdpInstancePool{})
	if err != nil {
		l.Logger.Errorf("[%s] CheckCommQueryParam failed, req:%+v err:%+v", sessionId, req, err)
		return nil, errorx.NewDefaultError(errorx.ParamErrorCode)
	}
	total, list, _, err := table.T_TCdpInstancePoolService.QueryPage(l.ctx, sessionId, qry, req.Offset, req.Limit, req.Sorts, req.Orders)
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpInstancePoolService querypage qry :%s err: %v", sessionId, qry, err)
		return nil, errorx.NewDefaultError(errorx.ListStrategyErrorCode)
	}

	resp = &types.InstancePoolListResp{
		Total: total,
	}

	if len(list) == 0 {
		return
	}

	var areaId int64 = 0
	poolIds := make([]string, 0)
	poolIdList := make([]int, 0)
	for _, v := range list {
		areaId = v.AreaId
		poolId := fmt.Sprintf("%d", v.PoolId)
		if v.Id != 0 && !gopublic.StringInArray(poolId, poolIds) {
			poolIds = append(poolIds, poolId)
			poolIdList = append(poolIdList, int(v.PoolId))
		}
	}

	bizStrategyInstCountMap := make(map[int64]int64) // kye : poolID value : 实例数量
	strategyNameMap := make(map[int64][]string)
	if len(poolIds) != 0 {
		bizStrategyInstCountMap, strategyNameMap = l.GetPoolInstCount(areaId, poolIds)
	}

	// 优化，批量查询无盘实例接口
	instListReq := &instance_types.ListInstancesRequestNew{
		Offset:        0,
		Length:        9999,
		Specification: poolIdList,
	}
	InstanceDetails, err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).SearcchInstanceList(areaId, sessionId, instListReq)
	if err != nil {
		l.Logger.Errorf("[%s] GetInstanceListBySpecification areaId:%d, instListReq:%d err:%+v", sessionId, areaId, helper.ToJSON(instListReq), err)
		return nil, errorx.NewDefaultCodeError("查询实例列表失败")
	}

	mapPoolInstanceDetail := make(map[int64][]instance_types.InstanceDetail) // key : poolID value : 实例列表
	for _, inst := range InstanceDetails {
		mapPoolInstanceDetail[inst.Specification] = append(mapPoolInstanceDetail[inst.Specification], inst)
	}

	for _, v := range list {

		if value, ok := mapPoolInstanceDetail[v.PoolId]; ok {
			InstanceDetails = value
		} else {
			InstanceDetails = []instance_types.InstanceDetail{}
		}

		instPoolInfo := types.InstancePoolInfo{
			Id:               uint64(v.Id),
			PoolId:           uint64(v.PoolId),
			PoolName:         v.InstPoolName,
			Remark:           v.Remark,
			Status:           int(v.Status),
			CreateTime:       v.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateTime:       v.UpdateTime.Format("2006-01-02 15:04:05"),
			StrategyNameList: []string{},
		}

		if strategyNameList, ok := strategyNameMap[int64(v.Id)]; ok {
			instPoolInfo.StrategyNameList = strategyNameList
		}

		instPoolInfo.InstCount = uint64(len(InstanceDetails))

		// 占用数量，要确认策略是否已经绑定租户
		instPoolInfo.OccupyCount = 0
		if instCount, ok := bizStrategyInstCountMap[int64(v.PoolId)]; ok {
			instPoolInfo.OccupyCount = uint64(instCount)
		}

		idleCount := instPoolInfo.InstCount - instPoolInfo.OccupyCount
		if instPoolInfo.InstCount < instPoolInfo.OccupyCount {
			idleCount = 0
		}

		instPoolInfo.IdleCount = idleCount
		var specNames []string
		for _, detail := range InstanceDetails {
			if detail.HostInfo.Gpu != "" && !gopublic.StringInArray(detail.HostInfo.Gpu, specNames) {
				specNames = append(specNames, detail.HostInfo.Gpu)
			}
		}

		instPoolInfo.SpecName = strings.Join(specNames, ",")
		resp.List = append(resp.List, instPoolInfo)
	}

	return
}
func (l *InstancePoolListLogic) GetPoolInstCount(areaId int64, poolIds []string) (bizStrategyInstCountMap map[int64]int64, mapStrateyNameList map[int64][]string) {

	sessionId := helper.GetSessionId(l.ctx)
	bizStrategyInstCountMap = make(map[int64]int64)
	mapStrateyNameList = make(map[int64][]string)

	// 查询策略信息，拿到策略ID
	qry := fmt.Sprintf("area_id:%d$inst_pool_id__in:%s$status:1", areaId, strings.Join(poolIds, ","))
	strategyInfos, _, err := table.T_TCdpResourceStrategyService.QueryAll(l.ctx, sessionId, qry, nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] GetPoolInstCount T_TCdpResourceStrategyService QueryAll qry :%s err: %v", sessionId, qry, err)
		return
	}

	// 查询这个策略是否绑定了租户
	// 优化，批量查询策略信息
	strategyIds := make([]string, 0)
	for _, v := range strategyInfos {
		if v.Id != 0 && !gopublic.StringInArray(fmt.Sprintf("%d", v.Id), strategyIds) {
			strategyIds = append(strategyIds, fmt.Sprintf("%d", v.Id))
		}
	}

	qry = fmt.Sprintf("area_id:%d$inst_strategy_id__in:%s$status:1", areaId, strings.Join(strategyIds, ","))
	bizStrategyInfos, _, err := table.T_TCdpBizStrategyService.QueryAll(l.ctx, sessionId, qry, nil, nil)
	if err != nil && err != gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] GetPoolInstCount T_TCdpBizStrategyService QueryAll qry:%s err: %v", sessionId, qry, err)
		return
	}
	if err == gopublic.ErrNotExist {
		l.Logger.Infof("[%s] GetPoolInstCount T_TCdpBizStrategyService QueryAll ErrNotExist qry:%s", sessionId, qry)
	}

	// 有绑定租户的策略
	mapStrateyBindBiz := make(map[int64]struct{})
	for _, vv := range bizStrategyInfos {
		mapStrateyBindBiz[vv.InstStrategyId] = struct{}{}
	}

	for _, v := range strategyInfos {

		if value, ok := mapStrateyNameList[v.InstPoolId]; ok {
			mapStrateyNameList[v.InstPoolId] = append(value, v.Name)
		}

		if _, ok := mapStrateyBindBiz[v.Id]; !ok {
			continue
		}

		bizStrategyInstCountMap[v.InstPoolId] = bizStrategyInstCountMap[v.InstPoolId] + int64(v.TotalInstances)

	}

	l.Logger.Debugf("[%s]  GetPoolInstCount bizStrategyInstCountMap :%s mapStrateyNameList:%s",
		sessionId, helper.ToJSON(bizStrategyInstCountMap), helper.ToJSON(mapStrateyNameList))

	return
}
