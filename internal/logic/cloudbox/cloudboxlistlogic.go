package cloudbox

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
	"gitlab.vrviu.com/cloud_esport_backend/saas_user_server/utils"
)

type CloudBoxListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCloudBoxListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CloudBoxListLogic {
	return &CloudBoxListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CloudBoxListLogic) CloudBoxList(req *types.GetCloudBoxListReq) (resp *types.CloudBoxListResp, err error) {
	sessionId := helper.GetSessionId(l.ctx)
	resp = new(types.CloudBoxListResp)

	if !helper.CheckAreaId(l.ctx, fmt.Sprintf("%d", req.AreaId)) {
		return nil, errorx.NewDefaultCodeError("无权限查看该区域云盒")
	}

	// 处理mac地址
	req.CondList = append(req.CondList, "status:1")

	if req.CloudClientName != "" {
		cloudClientInfos, _, err := table.T_TCdpCloudclientInfoService.QueryAll(l.ctx, sessionId, fmt.Sprintf("biz_id:%d$name__contains:%s$status:%d", req.BizId, req.CloudClientName, table.CloudClientStatusValid), nil, nil)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("[%s] T_TCdpCloudclientInfoService QueryAll failed req:%s err:%+v", sessionId, helper.ToJSON(req), err)
			return nil, errorx.NewDefaultCodeError("查询客户机列表失败")
		}
		var macList []string
		for _, cloudClientInfo := range cloudClientInfos {
			if cloudClientInfo.CloudboxMac != "" {
				macList = append(macList, cloudClientInfo.CloudboxMac)
			}
		}
		if len(macList) == 0 {
			return resp, nil
		}
		req.CondList = append(req.CondList, fmt.Sprintf("mac__in:%s", strings.Join(macList, ",")))
	}
	if req.FirstStrategyName != "" {
		strategyInfos, _, err := table.T_TCdpResourceStrategyService.QueryAll(l.ctx, sessionId, fmt.Sprintf("name__contains:%s$status:%d", req.FirstStrategyName, table.ResourceStrategyStatusValid), nil, nil)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("[%s] T_TCdpResourceStrategyService QueryAll failed req:%s err:%+v", sessionId, helper.ToJSON(req), err)
			return nil, errorx.NewDefaultCodeError("查询主算力策略列表失败")
		}
		var strategyIdList []string
		for _, strategyInfo := range strategyInfos {
			strategyIdList = append(strategyIdList, fmt.Sprintf("%d", strategyInfo.Id))
		}
		if len(strategyIdList) == 0 {
			return resp, nil
		}
		req.CondList = append(req.CondList, fmt.Sprintf("first_strategy_id__in:%s", strings.Join(strategyIdList, ",")))
	}
	if req.SecondStrategyName != "" {
		strategyInfos, _, err := table.T_TCdpResourceStrategyService.QueryAll(l.ctx, sessionId, fmt.Sprintf("name__contains:%s$status:%d", req.SecondStrategyName, table.ResourceStrategyStatusValid), nil, nil)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("[%s] T_TCdpResourceStrategyService QueryAll failed req:%s err:%+v", sessionId, helper.ToJSON(req), err)
			return nil, errorx.NewDefaultCodeError("查询从算力策略列表失败")
		}
		var strategyIdList []string
		for _, strategyInfo := range strategyInfos {
			strategyIdList = append(strategyIdList, fmt.Sprintf("%d", strategyInfo.Id))
		}
		if len(strategyIdList) == 0 {
			return resp, nil
		}
		req.CondList = append(req.CondList, fmt.Sprintf("second_strategy_id__in:%s", strings.Join(strategyIdList, ",")))
	}

	qry, err := helper.CheckCommQueryParam(l.ctx, sessionId, req.CondList, req.Orders, req.Sorts, table.TCdpCloudboxInfo{})
	if err != nil {
		l.Logger.Errorf("[%s] CheckCommQueryParam failed, req:%+v err:%+v", sessionId, req, err)
		return nil, errorx.NewDefaultError(errorx.ParamErrorCode)
	}

	total, list, _, err := table.T_TCdpCloudboxInfoService.QueryPage(l.ctx, sessionId, qry, int(req.Offset), int(req.Limit), req.Sorts, req.Orders)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("[%s] T_TCdpCloudboxInfoService QueryPage failed req:%s err:%+v", sessionId, helper.ToJSON(req), err)
		return nil, errorx.NewDefaultCodeError("查询云盒列表失败")
	}

	var macList, strategyIdList []string
	for _, info := range list {
		cloudClientInfo := types.CloudBoxInfo{
			CloudBoxId:       int64(info.Id),
			Name:             info.Name,
			MAC:              info.Mac,
			FirstStrategyId:  int64(info.FirstStrategyId),
			SecondStrategyId: int64(info.SecondStrategyId),
			BootSchemaId:     int64(info.BootSchemaId),
			StartMode:        int(info.BootType),
		}
		resp.List = append(resp.List, cloudClientInfo)
		macList = append(macList, cloudClientInfo.MAC)
		strategyIdList = append(strategyIdList, fmt.Sprintf("%d", info.FirstStrategyId), fmt.Sprintf("%d", info.SecondStrategyId))
	}
	resp.Total = int64(total)

	if len(macList) == 0 {
		return resp, nil
	}
	// 获取客户机列表
	cloudClientInfos, _, err := table.T_TCdpCloudclientInfoService.QueryAll(l.ctx, sessionId, fmt.Sprintf("cloudbox_mac__in:%s", strings.Join(macList, ",")), nil, nil)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("[%s] T_TCdpCloudclientInfoService QueryAll failed req:%s err:%+v", sessionId, helper.ToJSON(req), err)
		return nil, errorx.NewDefaultCodeError("查询客户机列表失败")
	}
	var cloudClientMap = make(map[string]table.TCdpCloudclientInfo)
	for _, cloudClientInfo := range cloudClientInfos {
		cloudClientMap[cloudClientInfo.CloudboxMac] = cloudClientInfo
	}

	// 获取主从算力策略列表
	strategyInfos, _, err := table.T_TCdpResourceStrategyService.QueryAll(l.ctx, sessionId, fmt.Sprintf("id__in:%s", strings.Join(strategyIdList, ",")), nil, nil)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("[%s] T_TCdpResourceStrategyService QueryAll failed req:%s err:%+v", sessionId, helper.ToJSON(req), err)
		return nil, errorx.NewDefaultCodeError("查询主从算力策略列表失败")
	}
	var strategyMap = make(map[string]table.TCdpResourceStrategy)
	for _, strategyInfo := range strategyInfos {
		strategyMap[fmt.Sprintf("%d", strategyInfo.Id)] = strategyInfo
	}
	for index, item := range resp.List {
		if strategyInfo, ok := strategyMap[fmt.Sprintf("%d", item.FirstStrategyId)]; ok {
			resp.List[index].FirstStrategyName = strategyInfo.Name
		}
		if strategyInfo, ok := strategyMap[fmt.Sprintf("%d", item.SecondStrategyId)]; ok {
			resp.List[index].SecondStrategyName = strategyInfo.Name
		}
		if cloudClientInfo, ok := cloudClientMap[item.MAC]; ok {
			resp.List[index].CloudClientName = cloudClientInfo.Name
			resp.List[index].CloudClientId = int64(cloudClientInfo.Id)
		}
	}

	// 调用无盘实例的接口
	Instancelist, err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).GetInstanceList(req.AreaId, sessionId, macList)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("[%s] GetInstanceList failed macList:%+v err:%+v", sessionId, macList, err)
		//return nil, types.NewDefaultError(err.Error())
	}

	for _, instDetail := range Instancelist {
		for index, item := range resp.List {

			if instDetail.BootMac == item.MAC {
				resp.List[index].BootTime = instDetail.BootTime.Format(utils.TimeFormat)
				resp.List[index].NetCardSpeed = instDetail.HostInfo.Net
				resp.List[index].Ip = instDetail.NetInfo.Ip
				resp.List[index].StartMode = instDetail.BootType
				resp.List[index].ImageUse = instDetail.OsImage
				resp.List[index].ConfigId = fmt.Sprintf("%d", instDetail.SchemeId)
				resp.List[index].PowerStatusDesc = instDetail.PowerStatusDesc
				resp.List[index].ManageStatusDesc = instDetail.ManageStatusDesc
				resp.List[index].RunningStatusDesc = instDetail.RunningStatusDesc
				resp.List[index].InstanceId = instDetail.Id
				resp.List[index].UserMode = int(instDetail.UserMode)
				resp.List[index].HostId = instDetail.HostId
				resp.List[index].KeepType = instDetail.KeepType
			}
		}
	}

	return
}
