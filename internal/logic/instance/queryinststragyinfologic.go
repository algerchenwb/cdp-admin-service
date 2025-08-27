package instance

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

type QueryInstStragyInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryInstStragyInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryInstStragyInfoLogic {
	return &QueryInstStragyInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryInstStragyInfoLogic) QueryInstStragyInfo(req *types.QueryInstStragyInfoReq) (resp *types.QueryInstStragyInfoResp, err error) {
	sessionId := helper.GetSessionId(l.ctx)
	resp = &types.QueryInstStragyInfoResp{}
	resp.FlowId = req.FlowId

	if req.Mac == "" {
		return nil, errorx.NewDefaultCodeError("参数错误")
	}

	// 实时开机
	qry := fmt.Sprintf("cloudbox_mac:%s$biz_id:%d$status:1", req.Mac, req.BizId)
	cloudClientInfo, _, err := table.T_TCdpCloudclientInfoService.Query(l.ctx, sessionId, qry, nil, nil)
	if err != nil && err != gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] Query T_TCdpBizStrategy failed. flowId[%s]  qry:%s err: %v", sessionId, req.FlowId, qry, err)
		return nil, errorx.NewDefaultCodeError("查询客户机失败")
	}

	if err == nil {
		// mac查询 客户机信息 ,用 实时开机模式
		l.Logger.Infof("[%s] Query T_TCdpCloudclientInfoService success. flowId[%s] BizId[%d] mac[%s] cloudClientInfos: %v", sessionId, req.FlowId, req.BizId, req.Mac, cloudClientInfo)
		strategyIds := make([]string, 0) // 汇总策略ID列表
		if cloudClientInfo.FirstStrategyId != 0 {
			strategyIds = append(strategyIds, fmt.Sprintf("%d", cloudClientInfo.FirstStrategyId))
		}
		if cloudClientInfo.SecondStrategyId != 0 {
			strategyIds = append(strategyIds, fmt.Sprintf("%d", cloudClientInfo.SecondStrategyId))
		}

		qry = fmt.Sprintf("id__in:%s$status:1", strings.Join(strategyIds, ","))
		strategyInfos, _, err1 := table.T_TCdpResourceStrategyService.QueryAll(l.ctx, sessionId, qry, nil, nil)
		if err1 != nil && err1 != gopublic.ErrNotExist {
			l.Logger.Errorf("[%s] Query T_TCdpResourceStrategyService failed. flowId[%s]  qry:%s err: %v", sessionId, req.FlowId, qry, err1)
			return nil, errorx.NewDefaultCodeError("查询客户机算力策略失败")
		}
		if err1 == gopublic.ErrNotExist {
			l.Logger.Errorf("[%s] Query T_TCdpResourceStrategyService ErrNotExist. flowId[%s] qry:%s err: %v", sessionId, req.FlowId, qry, err)
			return nil, errorx.NewDefaultCodeError("客户机算力策略不存在")
		}

		mapStrategyInfo := make(map[int64]table.TCdpResourceStrategy)
		for _, item := range strategyInfos {
			mapStrategyInfo[item.Id] = item
		}

		// 查询启动方案
		bootSchemaIds := make([]string, 0)
		if cloudClientInfo.FirstBootSchemaId != 0 {
			bootSchemaIds = append(bootSchemaIds, fmt.Sprintf("%d", cloudClientInfo.FirstBootSchemaId))
		}
		if cloudClientInfo.SecondBootSchemaId != 0 {
			bootSchemaIds = append(bootSchemaIds, fmt.Sprintf("%d", cloudClientInfo.SecondBootSchemaId))
		}

		qry = fmt.Sprintf("id__in:%s$status:1", strings.Join(bootSchemaIds, ","))
		bootSchemaInfos, _, err1 := table.T_TCdpBootSchemaInfoService.QueryAll(l.ctx, sessionId, qry, nil, nil)
		if err1 != nil && err1 != gopublic.ErrNotExist {
			l.Logger.Errorf("[%s] Query T_TCdpBootSchemaInfoService failed. flowId[%s] qry:%s err: %v", sessionId, req.FlowId, qry, err1)
			return nil, errorx.NewDefaultCodeError("查询客户机启动方案失败")
		}
		if err1 == gopublic.ErrNotExist {
			l.Logger.Errorf("[%s] Query T_TCdpBootSchemaInfoService ErrNotExist. flowId[%s] qry:%s err: %v", sessionId, req.FlowId, qry, err1)
			return nil, errorx.NewDefaultCodeError("客户机启动方案不存在")
		}

		mapBootSchemaInfo := make(map[int64]table.TCdpBootSchemaInfo)
		for _, item := range bootSchemaInfos {
			mapBootSchemaInfo[item.Id] = item
		}

		instPoolId := int64(0)
		stragyName := ""
		if v, ok := mapStrategyInfo[cloudClientInfo.FirstStrategyId]; ok {
			instPoolId = v.InstPoolId
			stragyName = v.Name
		}

		strategy := types.StrategyData{
			BootType:     1, // 实时开机
			StrategyId:   cloudClientInfo.FirstStrategyId,
			StrategyName: stragyName,
			InstPoolId:   instPoolId,
			BootSchemaId: 0,
		}
		if v, ok := mapBootSchemaInfo[cloudClientInfo.FirstBootSchemaId]; ok {
			strategy.BootSchemaId = v.DisklessSchemaId

		}
		resp.StrategyList = append(resp.StrategyList, strategy)

		instPoolId = int64(0)
		stragyName = ""
		if v, ok := mapStrategyInfo[cloudClientInfo.SecondStrategyId]; ok {
			instPoolId = v.InstPoolId
			stragyName = v.Name
		}
		strategy = types.StrategyData{
			BootType:     1, // 实时开机,
			StrategyId:   cloudClientInfo.SecondStrategyId,
			InstPoolId:   instPoolId,
			StrategyName: stragyName,
			BootSchemaId: 0,
		}

		if v, ok := mapBootSchemaInfo[cloudClientInfo.SecondBootSchemaId]; ok {
			strategy.BootSchemaId = v.DisklessSchemaId
		}
		resp.StrategyList = append(resp.StrategyList, strategy)
		return
	}

	// 用预开机模式
	l.Logger.Errorf("[%s] Query T_TCdpBizStrategy ErrNotExist. flowId[%s] qry:%s err: %v", sessionId, req.FlowId, qry, err)

	if req.AreaId == 0 || req.PoolId == 0 || req.BizId == 0 {
		return nil, errorx.NewDefaultCodeError("预开机模式方式：参数错误")
	}

	_, _, err = table.T_TCdpBizInfoService.Query(l.ctx, sessionId, fmt.Sprintf("biz_id:%d$status__ex:0", req.BizId), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] Query T_TCdpBizInfo failed. flowId[%s] BizId[%d]  err: %v", sessionId, req.FlowId, req.BizId, err)
		return nil, errorx.NewDefaultCodeError("查询业务信息失败")
	}

	if err == gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] Query T_TCdpBizInfo ErrNotExist  flowId[%s] BizId[%d] err: %v", sessionId, req.FlowId, req.BizId, err)
		return nil, errorx.NewDefaultCodeError("业务BizId不存在")
	}

	qry = fmt.Sprintf("area_id:%d$inst_pool_id:%d$status:1", req.AreaId, req.PoolId) // todo 筛选预开机

	resourceStrategyInfos, _, err := table.T_TCdpResourceStrategyService.QueryAll(l.ctx, sessionId, qry, nil, nil)
	if err != nil && err != gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] Query T_TCdpResourceStrategyService failed. flowId[%s]  BizId[%d] poolId[%d] err: %v", sessionId, req.FlowId, req.BizId, req.PoolId, err)
		return nil, errorx.NewDefaultCodeError("查询资源策略信息失败")
	}
	if err == gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] Query T_TCdpResourceStrategyService ErrNotExist. flowId[%s] BizId[%d] poolId[%d] err: %v", sessionId, req.FlowId, req.BizId, req.PoolId, err)
		return nil, errorx.NewDefaultCodeError("该租户的算力池没有配置策略")
	}

	resStrategyIds := make([]string, 0) // 策略ID列表
	mapStrategyInfo := make(map[int64]table.TCdpResourceStrategy)
	for _, item := range resourceStrategyInfos {
		resStrategyIds = append(resStrategyIds, fmt.Sprintf("%d", item.Id))
		mapStrategyInfo[item.Id] = item
	}
	if len(resStrategyIds) == 0 {
		l.Logger.Errorf("[%s] Query T_TCdpResourceStrategyService failed. flowId[%s] BizId[%d] poolId[%d] err: %v", sessionId, req.FlowId, req.BizId, req.PoolId, err)
		return nil, errorx.NewDefaultCodeError("该租户没有配置对应的算力池ID")
	}

	qry = fmt.Sprintf("biz_id:%d$status:1$inst_strategy_id__in:%s", req.BizId, strings.Join(resStrategyIds, ","))
	bizStrategyInfos, _, err := table.T_TCdpBizStrategyService.QueryAll(l.ctx, sessionId, qry, nil, nil)
	if err != nil && err != gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] Query T_TCdpBizStrategyService failed. flowId[%s]  BizId[%d] resStrategyIds:%s  err: %v", sessionId, req.FlowId, req.BizId, helper.ToJSON(resStrategyIds), err)
		return nil, errorx.NewDefaultCodeError("查询租户策略信息失败")
	}
	if err == gopublic.ErrNotExist || len(bizStrategyInfos) == 0 {
		l.Logger.Errorf("[%s] Query T_TCdpBizStrategy ErrNotExist. flowId[%s] qry:%s err: %v", sessionId, req.FlowId, qry, err)
		return nil, errorx.NewDefaultCodeError("该租户没有配置对应的策略")
	}

	bizStrategyInfo := bizStrategyInfos[0]
	strategyInfo, ok := mapStrategyInfo[bizStrategyInfo.InstStrategyId]
	if !ok {
		l.Logger.Errorf("[%s] Query T_TCdpBizStrategy failed. flowId[%s] BizId[%d] err: %v", sessionId, req.FlowId, req.BizId, err)
		return nil, errorx.NewDefaultCodeError("该租户没有配置对应的算力策略")
	}
	l.Logger.Infof("[%s] get T_TCdpBizStrategy  . BizId[%d] strategyInfo:%s", sessionId, req.BizId, helper.ToJSON(strategyInfo))

	if strategyInfo.BootType == table.ResourceStrategyBootTypePre {

		// 预开机模式
		strategy := types.StrategyData{
			BootType:     int64(strategyInfo.BootType),
			StrategyId:   strategyInfo.Id,
			InstPoolId:   strategyInfo.InstPoolId,
			StrategyName: strategyInfo.Name,
			BootSchemaId: 0, // 预开机模式 不需要启动方案
		}
		resp.StrategyList = append(resp.StrategyList, strategy)
		return resp, nil

	}
	return
}
