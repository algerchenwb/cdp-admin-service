package instance

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

type InstanceStatisticLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInstanceStatisticLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InstanceStatisticLogic {
	return &InstanceStatisticLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InstanceStatisticLogic) InstanceStatistic(req *types.InstanceStatisticReq) (resp *types.InstanceStatisticResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)

	bizStrategys, _, err := table.T_TCdpBizStrategyService.QueryAll(l.ctx, sessionId, fmt.Sprintf("area_id:%d$status:%d", req.AreaId, table.BizStrategyStatusValid), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] InstanceStatistic err: %v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.QueryStrategyFailedErrorCode)
	}
	var strategyIds []string
	for _, strategy := range bizStrategys {
		strategyIds = append(strategyIds, fmt.Sprintf("%d", strategy.InstStrategyId))
	}
	instStrategys, _, err := table.T_TCdpResourceStrategyService.QueryAll(l.ctx, sessionId, fmt.Sprintf("id__in:%s", strings.Join(strategyIds, ",")), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] InstanceStatistic err: %v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.QueryStrategyFailedErrorCode)
	}
	strategyMap := make(map[string]table.TCdpResourceStrategy)
	for _, s := range instStrategys {
		strategyMap[fmt.Sprintf("%d", s.Id)] = s
	}
	usedCount, inValidCount := 0, 0
	var usedStrategyMap = make(map[string]struct{})
	for _, bizStrategy := range bizStrategys {
		strategyId := fmt.Sprintf("%d", bizStrategy.InstStrategyId)
		strategy, ok := strategyMap[strategyId]
		if !ok {
			continue
		}
		if _, ok := usedStrategyMap[strategyId]; ok {
			continue
		}
		usedCount += int(strategy.TotalInstances)
		usedStrategyMap[strategyId] = struct{}{}
	}

	areas, err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).SearchInstances(req.AreaId, sessionId, &diskless.SearchInstancesReq{
		Offset:      0,
		Length:      10000,
		DeviceTypes: []int64{0},
	})
	if err != nil {
		l.Logger.Errorf("[%s] InstanceStatistic err: %v", sessionId, err)
		return nil, errorx.NewDefaultCodeError("查询无盘实例失败")
	}
	for _, instance := range areas.Instances {
		if instance.Specification == l.svcCtx.Config.Instance.InvaildVersion {
			inValidCount++
		}
	}
	if usedCount > len(areas.Instances) {
		usedCount = len(areas.Instances)
	}

	return &types.InstanceStatisticResp{
		TotalInstances:   uint64(len(areas.Instances)),
		UsedInstances:    uint64(usedCount),
		FreeInstances:    uint64(len(areas.Instances) - usedCount),
		InValidInstances: uint64(inValidCount),
	}, nil
}
