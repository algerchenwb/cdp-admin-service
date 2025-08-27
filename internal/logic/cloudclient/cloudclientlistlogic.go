package cloudclient

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
	"gitlab.vrviu.com/cloud_esport_backend/saas_user_server/utils"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
)

type CloudClientListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCloudClientListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CloudClientListLogic {
	return &CloudClientListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CloudClientListLogic) CloudClientList(req *types.CloudClientListReq) (resp *types.CloudClientListResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)
	resp = new(types.CloudClientListResp)

	// 处理mac地址
	req.CondList = append(req.CondList, "status:1")
	qry, err := helper.CheckCommQueryParam(l.ctx, sessionId, req.CondList, req.Orders, req.Sorts, table.TCdpCloudclientInfo{})
	if err != nil {
		l.Logger.Errorf("[%s] CheckCommQueryParam failed, req:%+v err:%+v", sessionId, req, err)
		return nil, errorx.NewDefaultError(errorx.ParamErrorCode)
	}
	total, list, _, err := table.T_TCdpCloudclientInfoService.QueryPage(l.ctx, sessionId, qry, int(req.Offset), int(req.Limit), req.Sorts, req.Orders)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("[%s] T_TCdpCloudclientInfoService QueryPage failed qry:%s err:%+v", sessionId, qry, err)
		return nil, errorx.NewDefaultCodeError("查询云盒列表失败")
	}

	StrategyIds := make([]string, 0)
	BootSchemaIds := make([]string, 0)
	cloudMacList := make([]string, 0)
	for _, item := range list {
		if item.FirstStrategyId != 0 {
			StrategyId := fmt.Sprintf("%d", item.FirstStrategyId)
			if !gopublic.StringInArray(StrategyId, StrategyIds) {
				StrategyIds = append(StrategyIds, StrategyId)
			}
		}
		if item.SecondStrategyId != 0 {
			StrategyId := fmt.Sprintf("%d", item.SecondStrategyId)
			if !gopublic.StringInArray(StrategyId, StrategyIds) {
				StrategyIds = append(StrategyIds, StrategyId)
			}
		}

		if item.FirstBootSchemaId != 0 {
			BootSchemaId := fmt.Sprintf("%d", item.FirstBootSchemaId)
			if !gopublic.StringInArray(BootSchemaId, BootSchemaIds) {
				BootSchemaIds = append(BootSchemaIds, BootSchemaId)
			}
		}

		if item.SecondBootSchemaId != 0 {
			BootSchemaId := fmt.Sprintf("%d", item.SecondBootSchemaId)
			if !gopublic.StringInArray(BootSchemaId, BootSchemaIds) {
				BootSchemaIds = append(BootSchemaIds, BootSchemaId)
			}
		}

		if item.CloudboxMac != "" {
			cloudMacList = append(cloudMacList, item.CloudboxMac)
		}

	}

	l.Logger.Debugf("[%s] CloudClientList StrategyIds:%s BootSchemaIds:%s", sessionId, helper.ToJSON(StrategyIds), helper.ToJSON(BootSchemaIds))

	// 查询出策略的名字
	mapStrategyInfo := make(map[int64]table.TCdpResourceStrategy)
	if len(StrategyIds) > 0 {
		qry := fmt.Sprintf("id__in:%s", strings.Join(StrategyIds, ","))
		strategyInfs, _, err1 := table.T_TCdpResourceStrategyService.QueryAll(l.ctx, sessionId, qry, nil, nil)
		if err1 != nil && err1 != gopublic.ErrNotExist {
			logx.WithContext(l.ctx).Errorf("[%s] T_TCdpResourceStrategyService QueryAll failed StrategyIds:%s err:%+v", sessionId, helper.ToJSON(StrategyIds), err1)
			return nil, errorx.NewDefaultCodeError("查询策略列表失败")
		}

		for _, item := range strategyInfs {
			mapStrategyInfo[item.Id] = item
		}
	}

	// 查询出策略的名字
	mapBootSchemaInfo := make(map[int64]table.TCdpBootSchemaInfo)
	if len(BootSchemaIds) > 0 {
		bootSchemaInfos, _, err1 := table.T_TCdpBootSchemaInfoService.QueryAll(l.ctx, sessionId, fmt.Sprintf("id__in:%s", strings.Join(BootSchemaIds, ",")), nil, nil)
		if err1 != nil && err1 != gopublic.ErrNotExist {
			logx.WithContext(l.ctx).Errorf("[%s] T_TCdpBootSchemaInfoService QueryAll failed BootSchemaIds:%s err:%+v", sessionId, helper.ToJSON(BootSchemaIds), err1)
			return nil, errorx.NewDefaultCodeError("查询策略列表失败")
		}

		for _, item := range bootSchemaInfos {
			mapBootSchemaInfo[item.Id] = item
		}
	}

	var vmIds []int
	for _, info := range list {
		cloudClient := types.CloudClientInfo{
			CloudClientId:      int64(info.Id),
			CloudClientName:    info.Name,
			CloudHostName:      info.Name,
			FirstStrategyId:    info.FirstStrategyId,
			SecondStrategyId:   info.SecondStrategyId,
			FirstBootSchemaId:  info.FirstBootSchemaId,
			SecondBootSchemaId: info.SecondBootSchemaId,
			StreamState:        0, // 有mac址址且 实例是占用状态
			CloudHostIP:        info.HostIp,
			ManageStatusDesc:   "普通",
			ClientType:         int(info.ClientType),
			CloudBoxMAC:        info.CloudboxMac,
			InstanceId:         info.Vmid,
			UserMode:           int(info.AdminState),
		}

		if StrategyInfo, ok := mapStrategyInfo[info.FirstStrategyId]; ok {
			cloudClient.FirstStrategyName = StrategyInfo.Name
		}

		if StrategyInfo, ok := mapStrategyInfo[info.SecondStrategyId]; ok {
			cloudClient.SecondStrategyName = StrategyInfo.Name
		}
		if BootSchemaInfo, ok := mapBootSchemaInfo[info.FirstBootSchemaId]; ok {
			cloudClient.FirstBootSchemaName = BootSchemaInfo.Name
		}
		if BootSchemaInfo, ok := mapBootSchemaInfo[info.SecondBootSchemaId]; ok {
			cloudClient.SecondBootSchemaName = BootSchemaInfo.Name
		}

		if info.ClientType == table.ClientType1 {
			cloudClient.CloudHostMAC = info.CloudboxMac
		}

		// 创建客户机后，AdminState 被默认成了0，客户机的超管状态是3 非超管是4，其他的状态不做处理
		if info.AdminState == uint32(instance_types.RegularUser2) || info.AdminState == uint32(instance_types.BindAdminUser) {
			if desc, ok := types.ManageStatusDescMap[int(info.AdminState)]; ok {
				cloudClient.ManageStatusDesc = desc
			}
		}

		resp.List = append(resp.List, cloudClient)
		if info.Vmid != 0 {
			vmIds = append(vmIds, int(info.Vmid))
		}

	}
	resp.Total = int64(total)

	l.Logger.Debugf("[%s] CloudClientList Total[%d] List:%s", sessionId, resp.Total, helper.ToJSON(resp.List))

	// 调用无盘实例的接口,通过 mac地址查询云盒信息
	Instancelist, err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).GetInstanceList(req.AreaId, sessionId, cloudMacList)
	if err != nil {
		l.Logger.Errorf("[%s] GetInstanceList failed. AreaId[%d] macList:%s, err:%+v", sessionId, req.AreaId, helper.ToJSON(cloudMacList), err)
		//return nil, types.NewDefaultError(err.Error())
	}

	instanceMap := make(map[string]instance_types.InstanceDetail)
	for _, instDetail := range Instancelist {
		instanceMap[instDetail.BootMac] = instDetail
	}
	var instanceIds []int64

	for index, item := range resp.List {
		instDetail, ok := instanceMap[item.CloudBoxMAC]
		if ok {
			resp.List[index].BootTime = instDetail.BootTime.Format(utils.TimeFormat)
			resp.List[index].CloudBoxIP = instDetail.NetInfo.Ip
			resp.List[index].HostId = instDetail.HostId
			//resp.List[index].UserMode = int(instDetail.UserMode)
			resp.List[index].PowerStatusDesc = instDetail.PowerStatusDesc
			resp.List[index].KeepType = instDetail.KeepType
			resp.List[index].CloudBoxIP = instDetail.NetInfo.Ip
			resp.List[index].CloudBoxMAC = instDetail.BootMac
			instanceIds = append(instanceIds, instDetail.Id)
		}
	}

	// 调用无盘实例的接口 通过IP查询云主机实例信息
	var Instancelist2 []instance_types.InstanceDetail
	if len(vmIds) > 0 {
		instListReq := &instance_types.ListInstancesRequestNew{
			Offset:      0,
			Length:      9999,
			InstanceIds: vmIds,
		}
		Instancelist2, err = diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).SearcchInstanceList(req.AreaId, sessionId, instListReq)
		if err != nil {
			l.Logger.Errorf("[%s] GetInstanceList failed. AreaId[%d]  err:%+v", sessionId, req.AreaId, err)
			//return nil, types.NewDefaultError(err.Error())
		}
	}

	mapInstnceInfo := make(map[int64]instance_types.InstanceDetail)
	SchemeIdList := make([]string, 0)
	for _, instDetail := range Instancelist2 {
		mapInstnceInfo[instDetail.Id] = instDetail
		SchemeIdList = append(SchemeIdList, fmt.Sprintf("%d", instDetail.SchemeId))
		instanceIds = append(instanceIds, instDetail.Id)
	}

	qry = fmt.Sprintf("area_id:%d$diskless_schema_id__in:%s", req.AreaId, strings.Join(SchemeIdList, ","))
	bootSchemaInfos, _, err := table.T_TCdpBootSchemaInfoService.QueryAll(l.ctx, sessionId, qry, nil, nil)
	if err != nil && err != gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] T_TCdpBootSchemaInfoService QueryAll failed AreaId[%d] qry:%s err:%+v", sessionId, req.AreaId, qry, err)
		//return nil, errorx.NewDefaultCodeError("查询启动方案失败"),
	}

	mBootSchemaInfo := make(map[int64]table.TCdpBootSchemaInfo)
	for _, item := range bootSchemaInfos {
		mBootSchemaInfo[item.DisklessSchemaId] = item
	}

	mapBootSession := make(map[int64]diskless.BootSessionDetail)
	if len(instanceIds) != 0 {
		bootSessions, err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).QueryBootSession(sessionId, req.AreaId, diskless.QueryBootSessionRequest{
			FlowId:     sessionId,
			InstanceId: instanceIds,
		})
		if err != nil {
			l.Logger.Errorf("[%s] QueryBootSession failed. AreaId[%d] instanceIds:%s, err:%+v", sessionId, req.AreaId, helper.ToJSON(instanceIds), err)
			return nil, errorx.NewDefaultCodeError("查询开机进度失败")
		}

		for _, item := range bootSessions.List {
			mapBootSession[item.InstanceId] = item
		}
	}

	for index, item := range resp.List {
		if instDetail, ok := mapInstnceInfo[item.InstanceId]; ok {
			resp.List[index].BootTime = instDetail.BootTime.Format(utils.TimeFormat)
			resp.List[index].BootTime = instDetail.BootTime.Format(utils.TimeFormat)
			resp.List[index].PowerStatusDesc = instDetail.PowerStatusDesc
			resp.List[index].InstanceId = instDetail.Id
			resp.List[index].KeepType = instDetail.KeepType
			resp.List[index].ImageUse = instDetail.OsImage
			resp.List[index].DataImage = instDetail.DataImage
			resp.List[index].CloudHostMAC = instDetail.BootMac

			if bootSchemaInfo, ok := mBootSchemaInfo[instDetail.SchemeId]; ok && resp.List[index].ClientType == table.ClientType2 {
				resp.List[index].StreamBootSchemaName = bootSchemaInfo.Name
			}
			if bootSession, ok := mapBootSession[instDetail.Id]; ok {
				resp.List[index].Process = bootSession.Process
			}

			// AssignStatus = 100: 占用中  PowerStatus = 1: 开机  BusinessStatus = 0: 正常
			if len(instDetail.BootMac) != 0 && instDetail.AssignStatus == 100 && instDetail.PowerStatus == 1 && instDetail.BusinessStatus == 0 && resp.List[index].ClientType == table.ClientType2 {
				resp.List[index].StreamState = 5 // 串流中
			}
		}
		// 1.0客户机 无串流 无云盒
		if resp.List[index].ClientType == table.ClientType1 {
			resp.List[index].CloudBoxIP = ""
			resp.List[index].CloudBoxMAC = ""
			resp.List[index].StreamBootSchemaName = ""
			resp.List[index].StreamState = 0
			resp.List[index].ImageUse = ""
		}
	}

	return
}
