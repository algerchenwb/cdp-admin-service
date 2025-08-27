package cloudbox

import (
	"context"
	"fmt"
	"time"

	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/helper/diskless"
	"cdp-admin-service/internal/helper/saas"
	"cdp-admin-service/internal/model/errorx"
	instance_types "cdp-admin-service/internal/proto/instance_service/types"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CloudBoxDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCloudBoxDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CloudBoxDeleteLogic {
	return &CloudBoxDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CloudBoxDeleteLogic) CloudBoxDelete(req *types.CloudBoxDeleteReq) (resp *types.CloudBoxDeleteResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)
	updateBy := helper.GetUserName(l.ctx)
	resp = new(types.CloudBoxDeleteResp)

	for _, item := range req.List {
		cloudBoxInfo, _, err := table.T_TCdpCloudboxInfoService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d$status:1", item.CloudBoxId), nil, nil)
		if err != nil {
			l.Logger.Errorf("[%s] Query T_TCdpCloudboxInfo failed. err: %v", sessionId, err)
			return nil, errorx.NewDefaultCodeError("查询云盒信息失败")
		}

		bindClientCount, _, _, err := table.T_TCdpCloudclientInfoService.QueryPage(l.ctx, sessionId, fmt.Sprintf("cloudbox_mac:%s$status:%d", cloudBoxInfo.Mac, table.CloudClientStatusValid), 0, 1, nil, nil)
		if err != nil {
			l.Logger.Errorf("[%s] Query T_TCdpCloudclientInfo failed. err: %v", sessionId, err)
			return nil, errorx.NewDefaultCodeError("查询客户机信息失败")
		}

		if bindClientCount > 0 {
			l.Logger.Errorf("[%s] CloudBoxDelete failed. bindClientCount[%d] mac[%s]", sessionId, bindClientCount, item.MAC)
			return nil, errorx.NewDefaultCodeError("云盒已绑定云客户机")
		}

		saasMac := helper.ConvertMacAddress(cloudBoxInfo.Mac)
		saasReq := &types.GetCloudBoxListReq{
			Offset:   0,
			Limit:    100,
			Sorts:    "",
			Orders:   "",
			AreaId:   req.AreaId,
			CondList: []string{fmt.Sprintf("mac_address:%s$biz_id:%d", saasMac, req.BizId)},
		}
		saasDeviceInfos, err := saas.GetEsportDeviceInfoList(l.ctx, sessionId, l.svcCtx.Config.OutSide.SaasHost, saasReq)
		if err != nil {
			l.Logger.Errorf("[%s] CreateEsportDeviceInfo failed. bizId[%d]  saasReq:%s err: %v", sessionId, req.BizId, helper.ToJSON(saasReq), err)
			return nil, errorx.NewDefaultCodeError(err.Error())
		}

		l.Logger.Debugf("[%s] GetEsportDeviceInfoList success. saasReq:%s list:%s", sessionId, helper.ToJSON(saasReq), helper.ToJSON(saasDeviceInfos))

		instReq := &instance_types.DestroyInstanceRequest{
			FlowId:     sessionId,
			InstanceId: item.InstanceId,
		}
		if err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).DestroyInstance(req.AreaId, instReq); err != nil {
			l.Logger.Errorf("[%s] GetRestoreInstanceProcessing AreaId[%d] instReq:%s err:%+v", sessionId, req.AreaId, helper.ToJSON(instReq), err)
			return nil, errorx.NewDefaultError(errorx.DisklessDestroyInstanceErrorCode)
		}
		l.Logger.Infof("[%s] DestroyInstance success. instanceId[%d] instReq:%s", sessionId, item.InstanceId, helper.ToJSON(instReq))

		var deviceId int64 = 0
		for _, info := range saasDeviceInfos.List {
			if info.MACAddress == saasMac {
				deviceId = info.Id
				break
			}
		}
		if deviceId != 0 {
			if err = saas.ReleaseEsportDeviceInfo(l.ctx, sessionId, l.svcCtx.Config.OutSide.SaasHost, deviceId); err != nil {
				l.Logger.Errorf("[%s] ReleaseEsportDeviceInfo failed. deviceId[%d] mac[%s] err: %+v", sessionId, deviceId, item.MAC, err)
				return nil, errorx.NewCodeError(errorx.SaasReleaseDeviceErrorCode, err.Error())
			}
			if err = saas.DeleteEsportDeviceInfo(l.ctx, sessionId, l.svcCtx.Config.OutSide.SaasHost, deviceId); err != nil {
				l.Logger.Errorf("[%s] DeleteEsportDeviceInfo failed. deviceId[%d] mac[%s] err: %+vv", sessionId, deviceId, item.MAC, err)
				return nil, errorx.NewCodeError(errorx.SaasDeleteDeviceErrorCode, err.Error())
			}
		} else {
			l.Logger.Errorf("[%s] GetEsportDeviceInfo failed. deviceId=0 mac[%s] err: %v", sessionId, item.MAC, err)
		}
		// 删除云盒
		delTime := time.Now().Format("20060102150405")
		if _, _, err = table.T_TCdpCloudboxInfoService.Update(l.ctx, sessionId, cloudBoxInfo.Id, map[string]interface{}{
			"status":      0,
			"name":        fmt.Sprintf("%s-del-%s", cloudBoxInfo.Name, delTime),
			"mac":         fmt.Sprintf("%s-del-%s", cloudBoxInfo.Mac, delTime),
			"update_time": time.Now(),
			"update_by":   updateBy,
		}); err != nil {
			l.Logger.Errorf("[%s] Delete T_TCdpCloudboxInfo failed.  err: %v", sessionId, err)
			return nil, errorx.NewDefaultCodeError("删除云盒信息失败")
		}

		l.Logger.Infof("[%s] Delete T_TCdpCloudboxInfo success. mac[%s]", sessionId, cloudBoxInfo.Mac)
	}

	return
}
