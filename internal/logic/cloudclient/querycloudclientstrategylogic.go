package cloudclient

import (
	"context"
	"fmt"

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

type QueryCloudClientStrategyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryCloudClientStrategyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryCloudClientStrategyLogic {
	return &QueryCloudClientStrategyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryCloudClientStrategyLogic) QueryCloudClientStrategy(req *types.QueryCloudClientStrategyReq) (resp *types.QueryCloudClientStrategyResp, err error) {
	sessionId := helper.GetSessionId(l.ctx)
	resp = &types.QueryCloudClientStrategyResp{}

	if req.BizId == 0 || req.Mac == "" {
		return nil, errorx.NewDefaultCodeError("参数错误")
	}

	qry := fmt.Sprintf("biz_id:%d$status__ex:0", req.BizId)
	bizInfo, _, err := table.T_TCdpBizInfoService.Query(l.ctx, sessionId, qry, nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] Query T_TCdpBizInfo failed. flowId[%s] BizId[%d] qry:%s err: %v", sessionId, req.FlowId, req.BizId, qry, err)
		return nil, errorx.NewDefaultCodeError("查询业务信息失败")
	}
	if err == gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] Query T_TCdpBizInfo ErrNotExist. flowId[%s] BizId[%d] qry:%s err: %v", sessionId, req.FlowId, req.BizId, qry, err)
		return nil, errorx.NewDefaultCodeError("BizId不存在")
	}

	qry = fmt.Sprintf("cloudbox_mac:%s$biz_id:%d$status:1", req.Mac, req.BizId)
	cloudClientInfo, _, err := table.T_TCdpCloudclientInfoService.Query(l.ctx, sessionId, qry, nil, nil)
	if err != nil && err != gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] Query T_TCdpCloudclientInfoService failed. flowId[%s]  qry:%s err: %v", sessionId, req.FlowId, qry, err)
		return nil, errorx.NewDefaultCodeError("查询客户机失败")
	}
	if err == gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] Query T_TCdpCloudclientInfoService ErrNotExist. flowId[%s] qry:%s err: %v", sessionId, req.FlowId, qry, err)
		return nil, errorx.NewDefaultCodeError("查询客户机不存在")
	}

	l.Logger.Infof("[%s] Query T_TCdpCloudclientInfoService success. flowId[%s] BizId[%d] mac[%s] cloudClientInfos: %v", sessionId, req.FlowId, req.BizId, req.Mac, cloudClientInfo)

	resp = &types.QueryCloudClientStrategyResp{
		CloudClientId: int64(cloudClientInfo.Id),
		Mac:           cloudClientInfo.CloudboxMac,
		ConfigInfo:    cloudClientInfo.ConfigInfo,
		AdminState:    cloudClientInfo.AdminState,
		HostIp:        cloudClientInfo.HostIp,
		HostName:      cloudClientInfo.Name,
	}

	// 调用无盘实例的接口 通过IP查询主机实例信息

	instListReq := &instance_types.ListInstancesRequestNew{
		Offset: 0,
		Length: 9999,
		Ips:    []string{cloudClientInfo.HostIp},
	}
	Instancelist, err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).SearcchInstanceList(int64(bizInfo.AreaId), sessionId, instListReq)
	if err != nil {
		l.Logger.Errorf("[%s] GetInstanceList failed. AreaId[%d]  instListReq:%s err:%+v", sessionId, bizInfo.AreaId, helper.ToJSON(instListReq), err)
	}
	if len(Instancelist) != 0 {
		instInfo := Instancelist[0]
		resp.HostMac = instInfo.BootMac
		resp.HostName = cloudClientInfo.Name
		resp.HostNetmask = instInfo.DefaultConfig.Netmask
		resp.HostGateway = instInfo.DefaultConfig.Gateway
	}
	return
}
