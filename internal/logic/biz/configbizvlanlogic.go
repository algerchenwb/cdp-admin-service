package biz

import (
	"context"
	"fmt"
	"time"

	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	diskless "cdp-admin-service/internal/helper/diskless"
	"cdp-admin-service/internal/model/errorx"
	proto "cdp-admin-service/internal/proto/location_seat_service"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"gitlab.vrviu.com/cloud_esport_backend/saas_user_server/pb/saas_user"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
)

type ConfigBizVlanLogic struct {
	logx.Logger
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	userRpc saas_user.UserServiceClient
}

func NewConfigBizVlanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigBizVlanLogic {
	return &ConfigBizVlanLogic{
		Logger:  logx.WithContext(ctx),
		ctx:     ctx,
		svcCtx:  svcCtx,
		userRpc: svcCtx.UserRpc,
	}
}

func (l *ConfigBizVlanLogic) ConfigBizVlan(req *types.ConfigBizVlanReq) (resp *types.ConfigBizVlanResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)
	updateBy := helper.GetUserName(l.ctx)

	if req.VlanId == 0 || req.BizId == 0 {
		return nil, errorx.NewDefaultError(errorx.ParamErrorCode)
	}

	// 查询租户信息
	bizInfo, _, err := table.T_TCdpBizInfoService.Query(l.ctx, sessionId, fmt.Sprintf("biz_id:%d", req.BizId), nil, nil)
	if err != nil && err != gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] T_TCdpBizInfoService Query err. bizId[%d] err:%+v", sessionId, req.BizId, err)
		return nil, errorx.NewDefaultError(errorx.DalGetErrorCode)
	}
	if err == gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] Table Query err. ErrNotExist bizId[%d] err:%+v", sessionId, req.BizId, err)
		return nil, errorx.NewDefaultError(errorx.BizNotFoundErrorCode)
	}

	//  查询vlan是否被占用
	if req.VlanId != int64(bizInfo.VlanId) {
		clientTotal, _, _, err := table.T_TCdpCloudclientInfoService.QueryPage(l.ctx, sessionId, fmt.Sprintf("biz_id:%d$status:%d", req.BizId, table.CloudClientStatusValid), 0, 1, nil, nil)
		if err != nil {
			l.Logger.Errorf("[%s] T_TCdpCloudclientInfoService QueryPage err. bizId[%d] err:%+v", sessionId, req.BizId, err)
			return nil, errorx.NewDefaultError(errorx.DalGetErrorCode)
		}
		if clientTotal > 0 {
			return nil, errorx.NewDefaultCodeError("vlan已被使用，卸载客户机后重试")
		}
	}
	if req.BoxVlanId != int64(bizInfo.BoxVlanId) {
		boxTotal, _, _, err := table.T_TCdpCloudboxInfoService.QueryPage(l.ctx, sessionId, fmt.Sprintf("biz_id:%d$status:%d", req.BizId, table.CloudBoxStatusValid), 0, 1, nil, nil)
		if err != nil {
			l.Logger.Errorf("[%s] T_TCdpCloudboxInfoService QueryPage err. bizId[%d] err:%+v", sessionId, req.BizId, err)
			return nil, errorx.NewDefaultError(errorx.DalGetErrorCode)
		}
		if boxTotal > 0 {
			return nil, errorx.NewDefaultCodeError("vlan已被使用，卸载云盒后重试")
		}
	}

	newBizInfo, _, err := table.T_TCdpBizInfoService.Update(l.ctx, sessionId, bizInfo.Id, map[string]any{
		"vlan_id":   req.VlanId,
		"box_vlan_id": req.BoxVlanId,
		"update_by": updateBy,
		"update_at": time.Now(),
	})
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpBizInfoService Update err. bizId[%d] err:%+v", sessionId, req.BizId, err)
		return nil, errorx.NewDefaultError(errorx.DalUpdateErrorCode)
	}
	updateBizReq := &proto.UpdateLocationRequest{
		FlowId:         sessionId,
		BizId:          int32(req.BizId),
		BoxVlanId:      int32(req.BoxVlanId),
		PcVlanId:       int32(req.VlanId),
		InstanceVlanId: int32(req.VlanId),
	}
	_, err = diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).UpdateLocation(sessionId, int64(bizInfo.AreaId), updateBizReq)
	if err != nil {
		l.Logger.Errorf("[%s] diskless UpdateLocation failed. BizId[%d] err:%+v", sessionId, req.BizId, err)
		return nil, errorx.NewDefaultError(errorx.UpdateBizFailedErrorCode)
	}

	l.Logger.Debugf("[%s] T_TCdpBizInfoService Update success. bizId[%d] newBizInfo:%s", sessionId, req.BizId, helper.ToJSON(newBizInfo))

	saasuser, err := l.userRpc.GetUserInfoList(l.ctx, &saas_user.GetUserInfoListReq{
		SessionId:  sessionId,
		Offset:     0,
		Limit:      1,
		Conditions: fmt.Sprintf("biz_id:%d", req.BizId),
		Sorts:      "",
		Orders:     "",
	})
	if err != nil {
		l.Logger.Errorf("[%s] userRpc get user info failed, bizId[%d] err: %v", sessionId, req.BizId, err)
		return nil, errorx.NewDefaultError(errorx.QueryBizFailedErrorCode)
	}

	if len(saasuser.List) == 0 {
		l.Logger.Errorf("[%s] userRpc user info not found, bizId[%d] err: %v", sessionId, req.BizId, err)
		return nil, errorx.NewDefaultError(errorx.QueryBizFailedErrorCode)
	}

	saasUser := saasuser.List[0]

	var newUser = new(saas_user.UpdateUserInfoReq)

	if err := copier.Copy(newUser, saasUser); err != nil {
		l.Logger.Errorf("[%s] userRpc copier.Copy BizId[%d] name[%s] err:%+v", req.BizId, saasUser.UserName, err)
		return nil, errorx.NewDefaultError(errorx.UpdateBizFailedErrorCode)
	}

	newUser.VlanId = req.VlanId
	newUser.State = saasUser.State
	newUser.UserName = saasUser.UserName

	if _, err := l.userRpc.UpdateUserInfo(l.ctx, newUser); err != nil {
		l.Logger.Errorf("[%s] userRpc UpdateUserInfo failed. BizId[%d] name[%s] err:%+v", sessionId, req.BizId, saasUser.UserName, err)
		return nil, errorx.NewDefaultError(errorx.UpdateBizFailedErrorCode)
	}

	l.Logger.Infof("[%s] userRpc Update success. bizId[%d] newUser:%s", sessionId, req.BizId, helper.ToJSON(newUser))



	return
}
