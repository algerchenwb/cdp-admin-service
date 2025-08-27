package cloudbox

import (
	"context"

	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CloudBoxAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCloudBoxAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CloudBoxAddLogic {
	return &CloudBoxAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CloudBoxAddLogic) CloudBoxAdd(req *types.CloudBoxAddReq) (resp *types.CloudBoxAddResp, err error) {

	// sessionId := helper.GetSessionId(l.ctx)
	// updateBy := helper.GetUserName(l.ctx)
	// resp = new(types.CloudBoxAddResp)

	// for _, item := range req.List {

	// 	// 调用saas接口创建设备
	// 	cloudBoxName := fmt.Sprintf("N-%s-%s", strings.ReplaceAll(item.Ip, ".", ""), helper.CreatCloudBoxName())
	// 	macString := helper.ConvertMacAddress(item.MAC)
	// 	saasReq := saas.CreateDeviceInfoReq{
	// 		MACAddress:  macString,
	// 		UpdateState: 0,
	// 		BizID:       req.BizId,
	// 		DeviceName:  cloudBoxName,
	// 		UpdateBy:    updateBy,
	// 	}
	// 	_, err := saas.CreateEsportDeviceInfo(l.ctx, sessionId, l.svcCtx.Config.OutSide.SaasHost, saasReq)
	// 	if err != nil {
	// 		l.Logger.Errorf("[%s] CreateEsportDeviceInfo failed. bizId[%d] err: %v", sessionId, req.BizId, err)
	// 		return nil, errorx.NewDefaultCodeError(err.Error())
	// 	}

	// 	instReq := &instance_types.CreateInstanceRequest{
	// 		FlowId:     sessionId,
	// 		DeviceType: 2, // 0-云主机 1-本地主机 2-云盒 3-本地盒子
	// 		Vlan:       int(req.VlanId),
	// 		InstanceInfo: instance_types.InstanceInfo{
	// 			BootMac:  item.MAC,
	// 			BootType: item.StartMode,
	// 			NetInfo: instance_types.NetInfo{
	// 				Ip:       item.Ip,
	// 				Hostname: item.Name,
	// 			},
	// 			DefaultConfig: instance_types.DefaultConfig{
	// 				NetInfo: instance_types.NetInfo{
	// 					Ip:       item.Ip,
	// 					Hostname: item.Name,
	// 				},
	// 			},
	// 			State:    100,
	// 			Remark:   item.Name,
	// 			SchemeId: schemeId,
	// 		},
	// 	}

	// 	if err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).CreateInstance(req.AreaId, instReq); err != nil {
	// 		l.Logger.Errorf("[%s] diskless CreateInstance failed. areaId[%d] instReq:%s err: %v", sessionId, req.AreaId, helper.ToJSON(instReq), err)
	// 		return nil, errorx.NewDefaultError(errorx.DisklessCreateInstanceErrorCode)
	// 	}

	// 	// 添加到cdp 数据库
	// 	newCloudBoxInfo, _, err := table.T_TCdpCloudboxInfoService.Insert(l.ctx, sessionId, table.TCdpCloudboxInfo{
	// 		Name:       cloudBoxName,
	// 		BizId:      req.BizId,
	// 		AreaId:     req.AreaId,
	// 		Mac:        item.MAC,
	// 		Ip:         item.Ip,
	// 		Status:     1,
	// 		CreateBy:   updateBy,
	// 		UpdateBy:   updateBy,
	// 		CreateTime: time.Now(),
	// 		UpdateTime: time.Now(),
	// 		ModifyTime: time.Now(),
	// 	})
	// 	if err != nil {
	// 		l.Logger.Errorf("[%s] T_TCdpCloudboxInfoService Insert failed. bizId[%d] err: %v", sessionId, req.BizId, err)
	// 		return nil, errorx.NewDefaultCodeError("添加云盒数据失败")
	// 	}

	// 	l.Logger.Infof("[%s] T_TCdpCloudboxInfoService Insert success. bizId[%d] newCloudBoxInfo:%s", sessionId, req.BizId, helper.ToJSON(newCloudBoxInfo))
	// }

	return
}
