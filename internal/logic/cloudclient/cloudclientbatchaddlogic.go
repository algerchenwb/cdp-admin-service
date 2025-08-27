package cloudclient

import (
	"context"
	"fmt"
	"strings"
	"time"

	"cdp-admin-service/internal/helper"
	"cdp-admin-service/internal/helper/cdp_cache"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/helper/diskless"
	"cdp-admin-service/internal/logic/common"
	"cdp-admin-service/internal/model/errorx"
	proto "cdp-admin-service/internal/proto/location_seat_service"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
)

type CloudClientBatchAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCloudClientBatchAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CloudClientBatchAddLogic {
	return &CloudClientBatchAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CloudClientBatchAddLogic) CloudClientBatchAdd(req *types.ClientBatchAddReq) (resp *types.ClientBatchAddResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)
	userName := helper.GetUserName(l.ctx)
	biz := cdp_cache.GetBizCache(l.ctx, sessionId, req.BizId)

	if biz == nil {
		return nil, errorx.NewDefaultError(errorx.BizNotFoundErrorCode)
	}
	if biz.VlanId == 0 {
		return nil, errorx.NewDefaultCodeError("当前租户未绑定客户机vlan")
	}

	resp = &types.ClientBatchAddResp{}
	var ips, clientNames []string
	ipNum, err := helper.IpToNum(req.StartIp)
	if err != nil {
		return nil, errorx.NewDefaultError(errorx.ParamErrorCode)
	}
	vlan, err := common.LoadVlanIP(l.ctx, l.svcCtx, int64(biz.VlanId), int64(biz.AreaId), diskless.IPTypeClient)
	if err != nil {
		l.Logger.Errorf("[%s] LoadVlanIP failed. biz_id:%d err: %v", sessionId, req.BizId, err)
		return nil, errorx.NewDefaultCodeError("查询vlan信息失败")
	}
	for offset := req.StartNumber; offset < req.StartNumber+req.ClientNum; offset++ {
		ip := helper.NumToIP(ipNum)
		ipNum++
		if !vlan.FreeIp(ip) {
			resp.FailedItem = append(resp.FailedItem, types.ClientBatchAddItem{
				HostIp:     ip,
				ClientName: fmt.Sprintf("%s%d", req.PrefixName, offset),
				ErrorMsg:   "IP不存在或已被使用",
			})
			continue
		}
		ips = append(ips, ip)
		clientName := fmt.Sprintf("%s%d", req.PrefixName, offset)
		clientNames = append(clientNames, clientName)
		if len(clientName) > 15 {
			return nil, errorx.NewDefaultCodeError("云主机名不能超过15个字符")
		}
	}

	err = l.checkClientNum(sessionId, biz, req.ClientNum)
	if err != nil {
		return nil, err
	}

	ips, clientNames, err = l.checkClientParams(sessionId, ips, clientNames, resp)
	if err != nil {
		return nil, err
	}

	bootSchema, _, err := table.T_TCdpBootSchemaInfoService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d$biz_id:%d$status:%d", req.BootSchemaId, req.BizId, table.TCdpBootSchemaInfoStatusEnable), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpBootSchemaInfoService Query err. bootSchemaId[%d] err:%+v", sessionId, req.BootSchemaId, err)
		return nil, errorx.NewDefaultError(errorx.QueryBizStrategyFailedErrorCode)
	}

	for index, ip := range ips {

		name := clientNames[index]
		disklessReq := &proto.AddSeatRequest{
			FlowId:           sessionId,
			Name:             name,
			LocationBizId:    int32(req.BizId),
			LocalInstanceMac: "",
			LocalIp:          ip,
			LocalSchemeId:    int32(bootSchema.DisklessSchemaId),
			LocalBootType:    int32(proto.BootType_BOOTTYPE_DISKLESS_UPGRADE),
			Type:             int32(proto.SeatType_SeatTypeDisklessComputer),
			ManagerState:     int32(proto.SeatManagerState_SeatManagerStateEnable),
		}

		disklessSeat, err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).AddSeat(sessionId, int64(req.AreaId), disklessReq)
		if err != nil {
			l.Logger.Errorf("[%s] diskless.NewDisklessWebGateway AddSeat [%s]err. err:%+v", sessionId, gopublic.ToJSON(disklessReq), err)
			resp.FailedItem = append(resp.FailedItem, types.ClientBatchAddItem{
				HostIp:     ip,
				ClientName: name,
				ErrorMsg:   fmt.Sprintf("添加客户机失败, err: %v", err),
			})
			continue
		}

		newCloudClientInfo, _, err := table.T_TCdpCloudclientInfoService.Insert(l.ctx, sessionId, &table.TCdpCloudclientInfo{
			Name:               name,
			BizId:              req.BizId,
			AreaId:             int64(biz.AreaId),
			FirstStrategyId:    0,
			SecondStrategyId:   0,
			FirstBootSchemaId:  req.BootSchemaId,
			SecondBootSchemaId: 0,
			CloudboxMac:        "",
			HostIp:             ip,
			ConfigInfo:         "",
			AdminState:         0,
			Status:             table.CloudClientStatusValid,
			Remark:             "",
			CreateBy:           userName,
			UpdateBy:           userName,
			UpdateTime:         time.Now(),
			CreateTime:         time.Now(),
			ModifyTime:         time.Now(),
			ClientType:         table.ClientType1,
			DisklessSeatId:     int64(disklessSeat.Id),
		})
		if err != nil {
			l.Logger.Errorf("[%s] Insert T_TCdpCloudclientInfo failed. CloudHostName[%s] CloudHostIP[%s]  err: %v", sessionId, name, ip, err)
			resp.FailedItem = append(resp.FailedItem, types.ClientBatchAddItem{
				HostIp:     ip,
				ClientName: name,
				ErrorMsg:   fmt.Sprintf("添加客户机失败, err: %v", err),
			})
			continue
		}
		resp.SuccessItem = append(resp.SuccessItem, types.ClientBatchAddItem{
			Id:         int64(newCloudClientInfo.Id),
			HostIp:     ip,
			ClientName: name,
			ErrorMsg:   "",
		})
		l.Logger.Infof("[%s] CloudClientAdd success. CloudHostName[%s] CloudHostIP[%s] newCloudClientInfo:%s", sessionId, name, ip, helper.ToJSON(newCloudClientInfo))
	}

	return
}

func (l *CloudClientBatchAddLogic) checkClientNum(sessionId string, biz *table.TCdpBizInfo, clientNum int64) error {

	existClentTotal, _, _, err := table.T_TCdpCloudclientInfoService.QueryPage(l.ctx, sessionId, fmt.Sprintf("biz_id:%d$status:%d",
		biz.BizId, table.CloudClientStatusValid), 0, 0, nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpCloudclientInfoService QueryPage err. bizId[%d] err:%+v", sessionId, biz.BizId, err)
		return errorx.NewDefaultError(errorx.QueryBizStrategyFailedErrorCode)
	}

	if int64(existClentTotal)+clientNum > int64(biz.ClientNumLimit) {
		return errorx.NewDefaultError(errorx.ClientNumExceedErrorCode)
	}

	return nil
}

func (l *CloudClientBatchAddLogic) checkClientIp(sessionId string, ips []string, names []string, resp *types.ClientBatchAddResp) ([]string, []string, error) {

	existIpClient, _, err := table.T_TCdpCloudclientInfoService.QueryAll(l.ctx, sessionId, fmt.Sprintf("host_ip__in:%s$status:%d", strings.Join(ips, ","), table.CloudClientStatusValid), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpCloudclientInfoService QueryAll err. ip[%s] err:%+v", sessionId, strings.Join(ips, ","), err)
		return nil, nil, errorx.NewDefaultError(errorx.QueryBizStrategyFailedErrorCode)
	}

	var existIpClientMap = make(map[string]struct{})
	for _, existIpClient := range existIpClient {
		existIpClientMap[existIpClient.HostIp] = struct{}{}
	}

	validIps := make([]string, 0, len(ips))
	validNames := make([]string, 0, len(names))
	for index, ip := range ips {
		if _, ok := existIpClientMap[ip]; ok {
			resp.FailedItem = append(resp.FailedItem, types.ClientBatchAddItem{
				HostIp:     ip,
				ClientName: names[index],
				ErrorMsg:   "client ip already exists",
			})
		} else {
			validIps = append(validIps, ip)
			validNames = append(validNames, names[index])
		}
	}

	return validIps, validNames, nil
}

func (l *CloudClientBatchAddLogic) checkClientName(sessionId string, ips []string, names []string, resp *types.ClientBatchAddResp) ([]string, []string, error) {

	existClient, _, err := table.T_TCdpCloudclientInfoService.QueryAll(l.ctx, sessionId, fmt.Sprintf("name__in:%s$status:%d", strings.Join(names, ","), table.CloudClientStatusValid), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpCloudclientInfoService QueryAll err. err:%+v", sessionId, err)
		return nil, nil, errorx.NewDefaultError(errorx.QueryClientFailedErrorCode)
	}

	var existClientMap = make(map[string]struct{})
	for _, existClient := range existClient {
		existClientMap[existClient.Name] = struct{}{}
	}

	validIps := make([]string, 0, len(ips))
	validNames := make([]string, 0, len(names))
	for index, name := range names {
		if _, ok := existClientMap[name]; ok {
			resp.FailedItem = append(resp.FailedItem, types.ClientBatchAddItem{
				HostIp:     ips[index],
				ClientName: name,
				ErrorMsg:   "client name already exists",
			})
		} else {
			validIps = append(validIps, ips[index])
			validNames = append(validNames, name)
		}
	}

	return validIps, validNames, nil
}

func (l *CloudClientBatchAddLogic) checkClientParams(sessionId string, ips []string, names []string, resp *types.ClientBatchAddResp) ([]string, []string, error) {

	ips, names, err := l.checkClientIp(sessionId, ips, names, resp)
	if err != nil {
		return ips, names, err
	}

	ips, names, err = l.checkClientName(sessionId, ips, names, resp)
	if err != nil {
		return ips, names, err
	}
	return ips, names, nil
}
