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

	"github.com/zeromicro/go-zero/core/logx"
	"gitlab.vrviu.com/cloud_esport_backend/saas_user_server/pb/saas_user"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
)

type DeleteBizLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteBizLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteBizLogic {
	return &DeleteBizLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteBizLogic) DeleteBiz(req *types.DeleteBizReq) error {

	sessionId := helper.GetSessionId(l.ctx)
	userName := helper.GetUserName(l.ctx)

	biz, _, err := table.T_TCdpBizInfoService.Query(l.ctx, sessionId, fmt.Sprintf("biz_id:%d", req.BizId), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpBizInfoService query biz info failed, bizId[%d] err: %v", sessionId, req.BizId, err)
		return errorx.NewDefaultError(errorx.QueryBizFailedErrorCode)
	}

	if biz.Status == table.BizStatusDeleted {
		l.Logger.Errorf("[%s] T_TCdpBizInfoService biz info already deleted [%d] err: %v", sessionId, req.BizId, err)
		return errorx.NewDefaultError(errorx.BizAlreadyDeletedErrorCode)
	}

	saasuser, err := l.svcCtx.UserRpc.GetUserInfoList(l.ctx, &saas_user.GetUserInfoListReq{
		SessionId:  sessionId,
		Offset:     0,
		Limit:      1,
		Conditions: fmt.Sprintf("biz_id:%d", req.BizId),
		Sorts:      "",
		Orders:     "",
	})
	if err != nil {
		l.Logger.Errorf("[%s] UserRpc get user info failed,bizId[%d] err: %v", sessionId, req.BizId, err)
		return errorx.NewDefaultError(errorx.UserNotFoundErrorCode)
	}
	if len(saasuser.List) == 0 {
		l.Logger.Errorf("[%s] UserRpc user info not found, bizId[%d] err: %v", sessionId, req.BizId, err)
		return errorx.NewDefaultError(errorx.UserNotFoundErrorCode)
	}
	_, err = l.svcCtx.UserRpc.UpdateUserInfo(l.ctx, &saas_user.UpdateUserInfoReq{
		SessionId:             sessionId,
		UserId:                saasuser.List[0].UserId,
		UserName:              saasuser.List[0].UserName,
		Remark:                saasuser.List[0].Remark,
		AuthoritySet:          saasuser.List[0].AuthoritySet,
		Type:                  saasuser.List[0].Type,
		AreaType:              saasuser.List[0].AreaType,
		SetId:                 saasuser.List[0].SetId,
		PrimaryId:             saasuser.List[0].PrimaryId,
		AgentId:               saasuser.List[0].AgentId,
		BizId:                 saasuser.List[0].BizId,
		WallpaperId:           saasuser.List[0].WallpaperId,
		TemplateId:            saasuser.List[0].TemplateId,
		TemplatePriority:      saasuser.List[0].TemplatePriority,
		TotalInstances:        saasuser.List[0].TotalInstances,
		EffectiveDate:         saasuser.List[0].EffectiveDate,
		ExpectedEffectiveDate: saasuser.List[0].ExpectedEffectiveDate,
		Email:                 saasuser.List[0].Email,
		AccessKey:             saasuser.List[0].AccessKey,
		AccessSecret:          saasuser.List[0].AccessSecret,
		EnableQualityReport:   saasuser.List[0].EnableQualityReport,
		UpgradeGreyType:       saasuser.List[0].UpgradeGreyType,
		Domain:                saasuser.List[0].Domain,
		CreateBy:              saasuser.List[0].CreateBy,
		UpdateBy:              saasuser.List[0].UpdateBy,
		CreateTime:            saasuser.List[0].CreateTime,
		UpdateTime:            saasuser.List[0].UpdateTime,
		ModifyTime:            saasuser.List[0].ModifyTime,
		VlanId:                saasuser.List[0].VlanId,
		ShopMode:              saasuser.List[0].ShopMode,
		Tag:                   saasuser.List[0].Tag,
		State:                 3,
	})
	if err != nil {
		l.Logger.Errorf("[%s] UserRpc update user info failed, biz[%d] err: %v", sessionId, req.BizId, err)
		return errorx.NewCodeError(errorx.DeleteBizFailedErrorCode, helper.ParseGRPCError(err).Error())
	}
	_, err = l.svcCtx.UserRpc.DeleteUserInfo(l.ctx, &saas_user.DeleteUserInfoReq{
		SessionId: sessionId,
		UserId:    saasuser.List[0].UserId,
	})
	if err != nil {
		l.Logger.Errorf("[%s] UserRpc delete user info failed, biz[%d] err: %v", sessionId, req.BizId, err)
		return errorx.NewCodeError(errorx.DeleteBizFailedErrorCode, helper.ParseGRPCError(err).Error())
	}
	l.Logger.Infof("[%s] UserRpc delete user info success, biz[%d]", sessionId, req.BizId)
	disklessReq := &proto.DeleteLocationRequest{
		FlowId: sessionId,
		BizId:  int32(req.BizId),
	}
	_, err = diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).DeleteLocation(sessionId, int64(biz.AreaId), disklessReq)
	if err != nil {
		l.Logger.Errorf("[%s] diskless.NewDisklessWebGateway.DeleteLocation disklessReq[%s] err. err:%+v", sessionId, gopublic.ToJSON(disklessReq), err)
		return errorx.NewDefaultCodeError(err.Error())
	}
	l.Logger.Infof("[%s] diskless.NewDisklessWebGateway.DeleteLocation success. disklessReq[%s]", sessionId, gopublic.ToJSON(disklessReq))
	_, _, err = table.T_TCdpBizInfoService.Update(l.ctx, sessionId, biz.Id, map[string]interface{}{
		"status":      table.BizStatusDeleted,
		"update_by":   userName,
		"update_time": time.Now(),
	})
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpBizInfoService update biz info failed, bizId[%d] err: %v", sessionId, req.BizId, err)
		return errorx.NewDefaultError(errorx.DeleteBizFailedErrorCode)
	}

	l.Logger.Debugf("[%s] DeleteBiz success, bizId[%d]", sessionId, req.BizId)

	return nil
}
