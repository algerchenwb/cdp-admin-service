package cloudbox

import (
	"context"
	"fmt"
	"strings"
	"time"

	"cdp-admin-service/internal/helper"
	"cdp-admin-service/internal/helper/cdp_cache"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/helper/diskless"
	"cdp-admin-service/internal/helper/saas"
	"cdp-admin-service/internal/logic/common"
	"cdp-admin-service/internal/model/errorx"
	disklessType "cdp-admin-service/internal/proto/instance_service/types"
	instance_types "cdp-admin-service/internal/proto/instance_service/types"
	proto "cdp-admin-service/internal/proto/location_seat_service"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
)

type CloudBoxBatchAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCloudBoxBatchAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CloudBoxBatchAddLogic {
	return &CloudBoxBatchAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CloudBoxBatchAddLogic) CloudBoxBatchAdd(req *types.CloudBoxBatchAddReq) (resp *types.CloudBoxBatchAddResp, err error) {
	sessionId := helper.GetSessionId(l.ctx)
	userName := helper.GetUserName(l.ctx)
	resp = &types.CloudBoxBatchAddResp{}

	if !helper.CheckAreaId(l.ctx, fmt.Sprintf("%d", req.AreaId)) {
		return nil, errorx.NewDefaultCodeError("无权限操作该节点数据")
	}
	if req.FirstStrategyId == 0 || req.FirstStrategyId == req.SecondStrategyId {
		return nil, errorx.NewDefaultCodeError("主算力策略不为空且与从算力策略不能相同")
	}

	biz := cdp_cache.GetBizCache(l.ctx, sessionId, req.BizId)
	if biz == nil {
		return nil, errorx.NewDefaultCodeError("租户不存在")
	}
	if biz.BoxVlanId == 0 {
		return nil, errorx.NewDefaultCodeError("租户云盒vlan未配置")
	}
	vlan, err := common.LoadVlanIP(l.ctx, l.svcCtx, int64(biz.VlanId), req.AreaId, diskless.IpTypeBox)
	if err != nil {
		l.Logger.Errorf("[%s] LoadVlanIP failed. areaId[%d] vlanId[%d] err: %v", sessionId, req.AreaId, biz.VlanId, err)
		return nil, errorx.NewDefaultCodeError("获取vlan ip失败")
	}

	var ips, boxNames, macList, cloudClientIds []string
	ipNum, err := helper.IpToNum(req.StartIp)
	if err != nil {
		return nil, errorx.NewDefaultError(errorx.ParamErrorCode)
	}
	for offset := uint32(0); offset < uint32(len(req.MacList)); offset++ {

		ip := helper.NumToIP(ipNum + offset)
		mac := req.MacList[offset].Mac
		if !vlan.FreeIp(ip) {
			resp.FailedItem = append(resp.FailedItem, types.CloudBoxBatchAddItem{
				Mac:      mac,
				Ip:       ip,
				ErrorMsg: "IP不存在或已被使用",
			})
			continue
		}
		// TODO 批量添加云盒 需要优化
		err = common.UniqueMac("", mac, l.ctx)
		if err != nil {
			l.Logger.Errorf("[%s] UniqueMac failed. mac[%s] ip[%s] err: %v", sessionId, mac, ip, err)
			resp.FailedItem = append(resp.FailedItem, types.CloudBoxBatchAddItem{
				Mac:      mac,
				Ip:       ip,
				ErrorMsg: err.Error(),
			})
			continue
		}
		ips = append(ips, ip)
		boxNames = append(boxNames, fmt.Sprintf("N-%s-%s", strings.ReplaceAll(ip, ".", ""), helper.CreatCloudBoxName()))
		macList = append(macList, mac)
		if req.MacList[offset].CloudClientId != 0 && gopublic.StringInArray(fmt.Sprintf("%d", req.MacList[offset].CloudClientId), cloudClientIds) {
			l.Logger.Errorf("[%s] CloudClientId[%d]已存在", sessionId, req.MacList[offset].CloudClientId)
			return nil, errorx.NewDefaultCodeError(fmt.Sprintf("云客户机[%d]多次绑定云盒", req.MacList[offset].CloudClientId))
		}
		cloudClientIds = append(cloudClientIds, fmt.Sprintf("%d", req.MacList[offset].CloudClientId))
	}

	originCloudClients, _, err := table.T_TCdpCloudclientInfoService.QueryAll(l.ctx, sessionId,
		fmt.Sprintf("id__in:%s$status:%d", strings.Join(cloudClientIds, ","), table.CloudClientStatusValid), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s]获取云客户机信息失败: %s ", userName, err.Error())
		return nil, errorx.NewDefaultCodeError(fmt.Sprintf("获取云客户机信息失败: %s", err.Error()))
	}

	var cloudClientMap = make(map[string]table.TCdpCloudclientInfo)
	for _, cloudClient := range originCloudClients {
		cloudClientMap[fmt.Sprintf("%d", cloudClient.Id)] = cloudClient
	}

	var cloudClients []table.TCdpCloudclientInfo = make([]table.TCdpCloudclientInfo, len(req.MacList))
	for idx, id := range cloudClientIds {
		cloudClients[idx] = cloudClientMap[id]
	}

	// 主算力策略
	_, _, err = table.T_TCdpBizStrategyService.Query(l.ctx, sessionId,
		fmt.Sprintf("inst_strategy_id:%d$biz_id:%d$status:%d", req.FirstStrategyId, req.BizId, table.BizStrategyStatusValid), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s]获取主算力策略失败: %s ", userName, err.Error())
		return nil, errorx.NewDefaultCodeError(fmt.Sprintf("获取主算力策略失败: %s", err.Error()))
	}

	if req.SecondStrategyId != 0 {
		if req.FirstStrategyId == req.SecondStrategyId {
			return nil, errorx.NewDefaultCodeError("主算力策略与从算力策略不能相同")
		}
		_, _, err = table.T_TCdpBizStrategyService.QueryAll(l.ctx, sessionId,
			fmt.Sprintf("inst_strategy_id:%d$biz_id:%d$status:%d", req.SecondStrategyId, req.BizId, table.BizStrategyStatusValid), nil, nil)
		if err != nil {
			l.Logger.Errorf("[%s]获取从算力策略失败: %s", userName, err.Error())
			return nil, errorx.NewDefaultCodeError(fmt.Sprintf("获取从算力策略失败: %s", err.Error()))
		}
	}

	// 配置方案
	bootSchema := new(table.TCdpBootSchemaInfo)
	if req.StartMode == int(proto.BootType_BOOTTYPE_DISKLESS_UPGRADE) {
		bootSchema, _, err = table.T_TCdpBootSchemaInfoService.Query(l.ctx, sessionId,
			fmt.Sprintf("id:%d$biz_id:%d$status:%d", req.BootSchemaId, req.BizId, table.TCdpBootSchemaInfoStatusEnable), nil, nil)
		if err != nil {
			l.Logger.Errorf("[%s]获取配置方案失败: %v ", userName, err)
			return nil, errorx.NewDefaultCodeError(fmt.Sprintf("获取配置方案失败: %s", err.Error()))
		}
	}

	instListReq := &instance_types.ListInstancesRequestNew{
		Offset: 0,
		Length: 9999,
		Ips:    ips,
	}
	Instancelist2, err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).SearcchInstanceList(req.AreaId, sessionId, instListReq)
	if err != nil {
		l.Logger.Errorf("[%s] GetInstanceList failed. AreaId[%d]  err:%+v", sessionId, req.AreaId, err)
		return nil, errorx.NewDefaultCodeError("查询云主机信息失败")
	}
	var instanceMap = make(map[string]disklessType.InstanceDetail)
	for _, instance := range Instancelist2 {
		instanceMap[instance.NetInfo.Ip] = instance
	}

	for idx := 0; idx < len(ips); idx++ {
		cloudClient := cloudClients[idx] // 已初始化容量和box一致不panic
		if cloudClient.ClientType == table.ClientType2 && cloudClient.CloudboxMac != "" {
			resp.FailedItem = append(resp.FailedItem, types.CloudBoxBatchAddItem{
				Mac:      macList[idx],
				Ip:       ips[idx],
				ErrorMsg: "云客户机已绑定云盒",
			})
			continue
		}
		ip := ips[idx]
		cloudBoxName := boxNames[idx]
		mac := macList[idx]
		// 调用saas接口创建设备
		saasReq := saas.CreateDeviceInfoReq{
			MACAddress:  helper.ConvertMacAddress(mac),
			UpdateState: 0,
			BizID:       req.BizId,
			DeviceName:  cloudBoxName,
			UpdateBy:    userName,
		}
		_, err := saas.CreateEsportDeviceInfo(l.ctx, sessionId, l.svcCtx.Config.OutSide.SaasHost, saasReq)
		if err != nil {
			l.Logger.Errorf("[%s] CreateEsportDeviceInfo failed. bizId[%d] err: %v", sessionId, req.BizId, err)
			resp.FailedItem = append(resp.FailedItem, types.CloudBoxBatchAddItem{
				Mac:      mac,
				Ip:       ip,
				ErrorMsg: err.Error(),
			})
			continue
		}
		l.Logger.Infof("[%s] CreateEsportDeviceInfo success. saasReq:%s", sessionId, helper.ToJSON(saasReq))

		instReq := &instance_types.CreateInstanceRequest{
			FlowId:     sessionId,
			DeviceType: 2, // 0-云主机 1-本地主机 2-云盒 3-本地盒子
			Vlan:       int(biz.VlanId),
			InstanceInfo: instance_types.InstanceInfo{
				BootMac:  mac,
				BootType: req.StartMode,
				NetInfo: instance_types.NetInfo{
					Ip:       ip,
					Hostname: cloudBoxName,
				},
				DefaultConfig: instance_types.DefaultConfig{
					NetInfo: instance_types.NetInfo{
						Ip:       ip,
						Hostname: cloudBoxName,
					},
				},
				State:    100,
				Remark:   cloudBoxName,
				SchemeId: int64(bootSchema.DisklessSchemaId),
			},
		}

		if err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).CreateInstance(req.AreaId, instReq); err != nil {
			l.Logger.Errorf("[%s] diskless CreateInstance failed. areaId[%d] instReq:%s err: %v", sessionId, req.AreaId, helper.ToJSON(instReq), err)
			return nil, errorx.NewDefaultError(errorx.DisklessCreateInstanceErrorCode)
		}
		l.Logger.Infof("[%s] diskless CreateInstance success. areaId[%d] instReq:%s", sessionId, req.AreaId, helper.ToJSON(instReq))

		// 添加到cdp 数据库
		newCloudBoxInfo, _, err := table.T_TCdpCloudboxInfoService.Insert(l.ctx, sessionId, table.TCdpCloudboxInfo{
			Name:             cloudBoxName,
			BizId:            req.BizId,
			AreaId:           int64(biz.AreaId),
			Mac:              mac,
			Ip:               ip,
			Status:           1,
			CreateBy:         userName,
			UpdateBy:         userName,
			CreateTime:       time.Now(),
			UpdateTime:       time.Now(),
			ModifyTime:       time.Now(),
			FirstStrategyId:  req.FirstStrategyId,
			SecondStrategyId: req.SecondStrategyId,
			BootSchemaId:     int32(req.BootSchemaId),
			BootType:         int32(req.StartMode),
		})
		if err != nil {
			l.Logger.Errorf("[%s] T_TCdpCloudboxInfoService Insert failed. bizId[%d] err: %v", sessionId, req.BizId, err)
			resp.FailedItem = append(resp.FailedItem, types.CloudBoxBatchAddItem{
				Mac:      mac,
				Ip:       ip,
				ErrorMsg: err.Error(),
			})
			continue
		}
		l.Logger.Infof("[%s] T_TCdpCloudboxInfoService Insert success. bizId[%d] newCloudBoxInfo:%s", sessionId, req.BizId, helper.ToJSON(newCloudBoxInfo))

		if cloudClient.Id != 0 {
			err = l.DisklessClientBindBox(l.ctx, cloudClient, newCloudBoxInfo.Mac, req.FirstStrategyId, bootSchema)
			if err != nil {
				l.Logger.Errorf("[%s] ClientBindBox failed. bizId[%d] err: %v", sessionId, req.BizId, err)
				resp.FailedItem = append(resp.FailedItem, types.CloudBoxBatchAddItem{
					Mac:      mac,
					Ip:       ip,
					ErrorMsg: err.Error(),
				})
				continue
			}
			l.Logger.Infof("[%s] T_TCdpCloudboxInfoService Insert success. bizId[%d] newCloudBoxInfo:%s", sessionId, req.BizId, helper.ToJSON(newCloudBoxInfo))

			updateInfo := map[string]interface{}{
				"cloudbox_mac":       mac,
				"first_strategy_id":  req.FirstStrategyId,
				"second_strategy_id": req.SecondStrategyId,
				"client_type":        table.ClientType2,
				"update_by":          userName,
				"update_time":        time.Now(),
			}
			// 添加到云客户机
			_, _, err = table.T_TCdpCloudclientInfoService.Update(l.ctx, sessionId, cloudClient.Id, updateInfo)
			if err != nil {
				l.Logger.Errorf("[%s] T_TCdpCloudclientInfoService Update failed. bizId[%d] err: %v", sessionId, req.BizId, err)
				resp.FailedItem = append(resp.FailedItem, types.CloudBoxBatchAddItem{
					Mac:      mac,
					Ip:       ip,
					ErrorMsg: "绑定云客户机失败",
				})
				continue
			}
			resp.SuccessItem = append(resp.SuccessItem, types.CloudBoxBatchAddItem{
				Mac:      mac,
				Ip:       ip,
				ErrorMsg: "绑定云客户机成功",
			})
			l.Logger.Infof("[%s] T_TCdpCloudclientInfoService Update success. bizId[%d] updateInfo:%s", sessionId, req.BizId, helper.ToJSON(updateInfo))
		}
	}

	return
}

// 通知无盘更新绑定云盒 2.0
func (l *CloudBoxBatchAddLogic) DisklessClientBindBox(ctx context.Context, cloudClient table.TCdpCloudclientInfo, boxMac string,
	strategyId int64, bootSchema *table.TCdpBootSchemaInfo) (err error) {
	sessionId := helper.GetSessionId(ctx)
	strategy, _, err := table.T_TCdpResourceStrategyService.Query(ctx, sessionId, fmt.Sprintf("id:%d$status:%d", strategyId, table.ResourceStrategyStatusValid), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] GetResourceStrategy failed. areaId[%d] err: %v", sessionId, cloudClient.AreaId, err)
		return err
	}
	updateSeatReq := &proto.UpdateSeatRequest{
		FlowId:              sessionId,
		Id:                  int32(cloudClient.DisklessSeatId),
		Type:                int32(proto.SeatType_SeatTypeBoxStreamCloud),
		LocalInstanceMac:    boxMac, // 本地方案ID
		StreamIp:            cloudClient.HostIp,
		StreamSchemeId:      int32(bootSchema.DisklessSchemaId),
		StreamSpecification: int32(strategy.InstPoolId),
	}
	_, err = diskless.NewDisklessWebGateway(ctx, l.svcCtx).UpdateSeat(sessionId, cloudClient.AreaId, updateSeatReq)
	if err != nil {
		l.Logger.Errorf("[%s] UpdateSeat failed. areaId[%d] updateSeatReq:%s err: %v", sessionId, cloudClient.AreaId, helper.ToJSON(updateSeatReq), err)
		return err
	}
	l.Logger.Infof("[%s] UpdateSeat success. areaId[%d] updateSeatReq:%s", sessionId, cloudClient.AreaId, helper.ToJSON(updateSeatReq))
	return nil
}
