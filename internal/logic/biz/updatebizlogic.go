package biz

import (
	"context"
	"fmt"
	"time"

	"cdp-admin-service/internal/helper"
	"cdp-admin-service/internal/helper/cdp_cache"
	table "cdp-admin-service/internal/helper/dal"
	diskless "cdp-admin-service/internal/helper/diskless"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	proto "cdp-admin-service/internal/proto/location_seat_service"

	"github.com/zeromicro/go-zero/core/logx"
	"gitlab.vrviu.com/cloud_esport_backend/saas_user_server/pb/saas_user"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
)

type UpdateBizLogic struct {
	logx.Logger
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	userRpc saas_user.UserServiceClient
}

func NewUpdateBizLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBizLogic {
	return &UpdateBizLogic{
		Logger:  logx.WithContext(ctx),
		ctx:     ctx,
		svcCtx:  svcCtx,
		userRpc: svcCtx.UserRpc,
	}
}

func (l *UpdateBizLogic) UpdateBiz(req *types.UpdateBizReq) error {

	sessionId := helper.GetSessionId(l.ctx)
	userName := helper.GetUserName(l.ctx)
	if req.ClientNumLimit > 100000 {
		l.Logger.Errorf("[%s] client num limit[%d] is greater than 10000", sessionId, req.ClientNumLimit)
		return errorx.NewDefaultCodeError("客户机授权数量过大")
	}
	biz, _, err := table.T_TCdpBizInfoService.Query(l.ctx, sessionId, fmt.Sprintf("biz_id:%d$status__ex:%d", req.BizId, table.BizStatusDeleted), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpBizInfoService query biz info failed, bizId[%d] err: %v", sessionId, req.BizId, err)
		return errorx.NewDefaultError(errorx.BizNotFoundErrorCode)
	}
	clientTotal, _, _, err := table.T_TCdpCloudclientInfoService.QueryPage(l.ctx, sessionId, fmt.Sprintf("biz_id:%d$status:%d", req.BizId, table.CloudClientStatusValid), 0, 1, nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpCloudclientInfoService query client total failed, bizId[%d] err: %v", sessionId, req.BizId, err)
		return errorx.NewDefaultError(errorx.QueryBizFailedErrorCode)
	}
	if clientTotal > int(req.ClientNumLimit) {
		l.Logger.Errorf("[%s] client total[%d] is greater than client num limit[%d]", sessionId, clientTotal, req.ClientNumLimit)
		return errorx.NewDefaultCodeError("客户机授权数量不能低于当前客户机数量")
	}

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
		return errorx.NewDefaultError(errorx.QueryBizFailedErrorCode)
	}

	if len(saasuser.List) == 0 {
		l.Logger.Errorf("[%s] userRpc user info not found, bizId[%d] err: %v", req.BizId, err)
		return errorx.NewDefaultError(errorx.QueryBizFailedErrorCode)
	}
	saasUser := saasuser.List[0]
	_, err = l.userRpc.UpdateUserInfo(l.ctx, &saas_user.UpdateUserInfoReq{
		SessionId:             sessionId,
		UserId:                saasUser.UserId,
		UserName:              req.BizName,
		Remark:                req.Remark,
		AuthoritySet:          saasUser.AuthoritySet,
		Type:                  saasUser.Type,
		AreaType:              saasUser.AreaType,
		SetId:                 saasUser.SetId,
		PrimaryId:             saasUser.PrimaryId,
		AgentId:               saasUser.AgentId,
		BizId:                 saasUser.BizId,
		WallpaperId:           saasUser.WallpaperId,
		TemplateId:            saasUser.TemplateId,
		TemplatePriority:      saasUser.TemplatePriority,
		TotalInstances:        saasUser.TotalInstances,
		EffectiveDate:         saasUser.EffectiveDate,
		ExpectedEffectiveDate: saasUser.ExpectedEffectiveDate,
		Email:                 saasUser.Email,
		AccessKey:             saasUser.AccessKey,
		AccessSecret:          saasUser.AccessSecret,
		EnableQualityReport:   saasUser.EnableQualityReport,
		UpgradeGreyType:       saasUser.UpgradeGreyType,
		Domain:                saasUser.Domain,
		CreateBy:              saasUser.CreateBy,
		UpdateBy:              saasUser.UpdateBy,
		State:                 saasUser.State,
		CreateTime:            saasUser.CreateTime,
		UpdateTime:            saasUser.UpdateTime,
		ModifyTime:            saasUser.ModifyTime,
		VlanId:                saasUser.VlanId,
		ShopMode:              saasUser.ShopMode,
		AgentSharePercent:     uint64(saasUser.AgentSharePercent),
		Tag:                   saasUser.Tag,
	})
	if err != nil {
		l.Logger.Errorf("[%s] update user info failed, bizId[%d] err: %v", sessionId, req.BizId, err)
		return errorx.NewCodeError(errorx.UpdateBizFailedErrorCode, helper.ParseGRPCError(err).Error())
	}

	disklessReq := &proto.UpdateLocationRequest{
		FlowId:         sessionId,
		BizId:          int32(req.BizId),
		Name:           req.BizName,
		SeatNumLimit:   int32(req.ClientNumLimit),
		BoxVlanId:      int32(biz.VlanId),
		PcVlanId:       int32(biz.VlanId),
		InstanceVlanId: int32(biz.VlanId),
		ManagerState:   int32(proto.LocationManagerState_LocationManagerStateEnable),
	}
	_, err = diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).UpdateLocation(sessionId, int64(biz.AreaId), disklessReq)
	if err != nil {
		l.Logger.Errorf("[%s] diskless.NewDisklessWebGateway.UpdateLocation disklessReq[%s] err. err:%+v", sessionId, gopublic.ToJSON(disklessReq), err)
		return errorx.NewDefaultError(errorx.UpdateBizFailedErrorCode)
	}

	l.Logger.Infof("[%s] diskless.NewDisklessWebGateway.UpdateLocation success. disklessReq[%s]", sessionId, gopublic.ToJSON(disklessReq))
	_, _, err = table.T_TCdpBizInfoService.Update(l.ctx, sessionId, biz.Id, map[string]interface{}{
		"biz_name":         req.BizName,
		"contact_person":   req.ContactPerson,
		"mobile":           req.Mobile,
		"remark":           req.Remark,
		"update_by":        userName,
		"update_time":      time.Now(),
		"client_num_limit": req.ClientNumLimit,
	})
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpBizInfoService update biz info failed, err: %v", sessionId, err)
		return errorx.NewDefaultError(errorx.UpdateBizFailedErrorCode)
	}

	cdp_cache.DeleteBizCache(biz.Id)

	return nil
}
