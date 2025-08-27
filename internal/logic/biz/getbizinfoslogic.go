package biz

import (
	"context"
	"fmt"
	"strings"

	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
)

type GetBizInfosLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBizInfosLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBizInfosLogic {
	return &GetBizInfosLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBizInfosLogic) GetBizInfos(req *types.CommonPageRequest) (resp *types.BizInfosResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)
	req.CondList = append(req.CondList, "status__ex:0")
	req.CondList = append(req.CondList, fmt.Sprintf("area_id__in:%s", helper.GetAreaIds(l.ctx)))

	qry, err := helper.CheckCommQueryParam(l.ctx, sessionId, req.CondList, req.Orders, req.Sorts, table.TCdpBizInfo{})
	if err != nil {
		l.Logger.Errorf("[%s] CheckCommQueryParam failed, req:%+v err:%+v", sessionId, req, err)
		return nil, errorx.NewDefaultError(errorx.ParamErrorCode)
	}
	total, bizInfos, _, err := table.T_TCdpBizInfoService.QueryPage(l.ctx, sessionId, qry, req.Offset, req.Limit, req.Sorts, req.Orders)
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpBizInfoService query biz infos failed, qry:%s err: %v", sessionId, qry, err)
		return nil, errorx.NewDefaultError(errorx.QueryBizFailedErrorCode)
	}

	resp = &types.BizInfosResp{
		Total: int64(total),
	}

	areaIds := make([]string, 0)
	bizIds := make([]string, 0)

	for _, bizInfo := range bizInfos {
		areaIds = append(areaIds, fmt.Sprintf("%d", bizInfo.AreaId))
		bizIds = append(bizIds, fmt.Sprintf("%d", bizInfo.BizId))
	}

	clientInfos, _, err := table.T_TCdpCloudclientInfoService.QueryAll(l.ctx, sessionId, fmt.Sprintf("biz_id__in:%s$status:%d", strings.Join(bizIds, ","), table.CloudClientStatusValid), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpClientInfoService query client infos failed,bizIds:%+v err: %v", sessionId, bizIds, err)
		return nil, errorx.NewDefaultError(errorx.QueryBizFailedErrorCode)
	}

	clientMap := make(map[int64]int64)
	for _, clientInfo := range clientInfos {
		clientMap[clientInfo.BizId]++
	}

	areaInfos, _, err := table.T_TCdpAreaInfoService.QueryAll(l.ctx, sessionId, fmt.Sprintf("area_id__in:%s$status:%d", strings.Join(areaIds, ","), table.AreaStatusEnable), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpAreaInfoService query area infos failed,areaIds:%+v err: %v", sessionId, areaIds, err)
		return nil, errorx.NewDefaultError(errorx.QueryBizFailedErrorCode)
	}

	areaMap := make(map[int32]string)
	for _, areaInfo := range areaInfos {
		areaMap[areaInfo.AreaId] = areaInfo.Name
	}

	strategyInfos, _, err := table.T_TCdpBizStrategyService.QueryAll(l.ctx, sessionId, fmt.Sprintf("biz_id__in:%s$status__ex:%d", strings.Join(bizIds, ","), table.BizStrategyStatusInvalid), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpBizStrategyService query strategy infos failed, err: %v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.QueryBizFailedErrorCode)
	}

	var instStrategyIds []string
	for _, strategy := range strategyInfos {
		if !gopublic.StringInArray(fmt.Sprintf("%d", strategy.InstStrategyId), instStrategyIds) {
			instStrategyIds = append(instStrategyIds, fmt.Sprintf("%d", strategy.InstStrategyId))
		}
	}

	instStrategyInfos, _, err := table.T_TCdpResourceStrategyService.QueryAll(l.ctx, sessionId, fmt.Sprintf("id__in:%s", strings.Join(instStrategyIds, ",")), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpResourceStrategyService query inst strategy infos failed, instStrategyIds:%+v err: %v", sessionId, instStrategyIds, err)
		return nil, errorx.NewDefaultError(errorx.QueryStrategyFailedErrorCode)
	}

	instStrategyMap := make(map[int64]table.TCdpResourceStrategy)
	instPoolIds := make([]string, 0)
	for _, instStrategy := range instStrategyInfos {
		instStrategyMap[instStrategy.Id] = instStrategy
		instPoolIds = append(instPoolIds, fmt.Sprintf("%d", instStrategy.InstPoolId))
	}

	instPoolInfos, _, err := table.T_TCdpInstancePoolService.QueryAll(l.ctx, sessionId, fmt.Sprintf("pool_id__in:%s$status:%d", strings.Join(instPoolIds, ","), table.InstancePoolStatusValid), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpInstancePoolService query inst pool infos failed,instPoolIds:%+v err: %v", sessionId, instPoolIds, err)
		return nil, errorx.NewDefaultError(errorx.QueryBizFailedErrorCode)
	}

	instPoolMap := make(map[string]table.TCdpInstancePool)
	for _, instPool := range instPoolInfos {
		instPoolMap[fmt.Sprintf("%d-%d", instPool.AreaId, instPool.PoolId)] = instPool
	}
	strategyMap := make(map[int64][]types.BizStrategyInfo)
	for _, strategyInfo := range strategyInfos {
		strategyMap[strategyInfo.BizId] = append(strategyMap[strategyInfo.BizId], types.BizStrategyInfo{
			Id:               strategyInfo.Id,
			InstStrategyId:   strategyInfo.InstStrategyId,
			InstStrategyName: instStrategyMap[strategyInfo.InstStrategyId].Name,
			SpecId:           instPoolMap[fmt.Sprintf("%d-%d", instStrategyMap[strategyInfo.InstStrategyId].AreaId, instStrategyMap[strategyInfo.InstStrategyId].InstPoolId)].PoolId,
			SpecName:         instPoolMap[fmt.Sprintf("%d-%d", instStrategyMap[strategyInfo.InstStrategyId].AreaId, instStrategyMap[strategyInfo.InstStrategyId].InstPoolId)].InstPoolName,
			TotalInstances:   int64(instStrategyMap[strategyInfo.InstStrategyId].TotalInstances),
			Remark:           instStrategyMap[strategyInfo.InstStrategyId].Remark,
			OuterSpecId:      instStrategyMap[strategyInfo.InstStrategyId].OuterSpecId,
			BootType:         int64(instStrategyMap[strategyInfo.InstStrategyId].BootType),
		})
	}
	for _, bizInfo := range bizInfos {
		biz := types.BizInfo{
			BizId:          bizInfo.BizId,
			BizName:        bizInfo.BizName,
			Status:         int64(bizInfo.Status),
			CreateTime:     bizInfo.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateTime:     bizInfo.UpdateTime.Format("2006-01-02 15:04:05"),
			Remark:         bizInfo.Remark,
			AreaId:         int64(bizInfo.AreaId),
			AreaName:       areaMap[bizInfo.AreaId],
			RegionId:       int64(bizInfo.RegionId),
			ContactPerson:  bizInfo.ContactPerson,
			Mobile:         bizInfo.Mobile,
			CreateBy:       bizInfo.CreateBy,
			VlanId:         int64(bizInfo.VlanId),
			BoxVlanId:      int64(bizInfo.BoxVlanId),
			ServerInfo:     bizInfo.Serverinfo,
			ClientNumLimit: int64(bizInfo.ClientNumLimit),
			AuthorizedClientNum: clientMap[bizInfo.BizId],
		}
		biz.StrategyList = strategyMap[bizInfo.BizId]
		resp.List = append(resp.List, biz)
	}
	return
}
