package biz

import (
	"context"
	"fmt"
	"time"

	"cdp-admin-service/internal/helper"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	diskless "cdp-admin-service/internal/helper/diskless"

	table "cdp-admin-service/internal/helper/dal"
	proto "cdp-admin-service/internal/proto/location_seat_service"

	gopublic "gitlab.vrviu.com/inviu/backend-go-public/gopublic"

	"github.com/zeromicro/go-zero/core/logx"
	"gitlab.vrviu.com/cloud_esport_backend/saas_user_server/pb/saas_user"
)

type CreateBizLogic struct {
	logx.Logger
	ctx                context.Context
	svcCtx             *svc.ServiceContext
	disklessWebGateway *diskless.DisklessWebGateway
	userRpc            saas_user.UserServiceClient
}

func NewCreateBizLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateBizLogic {
	return &CreateBizLogic{
		Logger:             logx.WithContext(ctx),
		ctx:                ctx,
		svcCtx:             svcCtx,
		disklessWebGateway: diskless.NewDisklessWebGateway(ctx, svcCtx),
		userRpc:            svcCtx.UserRpc,
	}
}

func (l *CreateBizLogic) CreateBiz(req *types.CreateBizReq) (resp *types.CreateBizResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)
	createBy := helper.GetUserName(l.ctx)

	areaInfo, _, err := table.T_TCdpAreaInfoService.Query(l.ctx, sessionId, fmt.Sprintf("area_id:%d$status:%d", req.AreaId, table.AreaStatusEnable), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpAreaInfoService Query err. AreaId[%d] err:%+v", sessionId, req.AreaId, err)
		return nil, errorx.NewDefaultError(errorx.AreaNotFound)
	}
	if req.ClientNumLimit > 100000 {
		return nil, errorx.NewDefaultCodeError("客户机授权数量过大")
	}

	newUser := saas_user.CreateUserInfoReq{
		SessionId:             sessionId,
		PrimaryId:             areaInfo.PrimaryId,
		AgentId:               areaInfo.AgentId,
		UserName:              req.BizName,
		TemplateId:            0,
		TotalInstances:        0,
		ExpectedEffectiveDate: time.Now().Format("2006-01-02T15:04:05-07:00"),
		Domain:                "",
		CreateBy:              createBy,
		SetId:                 3, //网吧
		TemplatePriority:      1,
		State:                 1, // 0-未生效 1-合作中 2-停止合作
		EffectiveDate:         time.Now().Format("2006-01-02T15:04:05-07:00"),
	}

	newUserResp, err := l.userRpc.CreateUserInfo(l.ctx, &newUser)
	if err != nil {
		l.Logger.Errorf("[%s] Rpc.CreateUserInfo err. PrimaryID[%d] AgentID[%d], name[%s] err:%+v", sessionId, areaInfo.PrimaryId, areaInfo.AgentId, req.BizName, err)
		return nil, errorx.NewCodeError(errorx.CreateUserInfoErrorCode, helper.ParseGRPCError(err).Error())
	}

	// 通知无盘建立门店
	locationReq := &proto.AddLocationRequest{
		FlowId:         sessionId,
		BizId:          int32(newUserResp.UserInfo.BizId),
		Name:           req.BizName,
		SeatNumLimit:   int32(req.ClientNumLimit),
		BoxVlanId:      0,
		PcVlanId:       0,
		InstanceVlanId: 0,
		ManagerState:   int32(proto.LocationManagerState_LocationManagerStateEnable),
	}
	_, err = diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).AddLocation(sessionId, int64(areaInfo.AreaId), locationReq)
	if err != nil {
		l.Logger.Errorf("[%s] diskless.NewDisklessWebGateway.AddLocation locationReq[%s] err. err:%+v", sessionId, gopublic.ToJSON(locationReq), err)
		return nil, errorx.NewDefaultError(errorx.CreateUserInfoErrorCode)
	}
	l.Logger.Infof("[%s] diskless.NewDisklessWebGateway.AddLocation success. locationReq[%s]", sessionId, gopublic.ToJSON(locationReq))

	_, _, err = table.T_TCdpBizInfoService.Insert(l.ctx, sessionId, table.TCdpBizInfo{
		BizId:          newUserResp.UserInfo.BizId,
		BizName:        req.BizName,
		AreaId:         int32(req.AreaId),
		CreateBy:       createBy,
		CreateTime:     time.Now(),
		UpdateTime:     time.Now(),
		ModifyTime:     time.Now(),
		Status:         table.BizStatusWaitWorking,
		Remark:         req.Remark,
		ContactPerson:  req.ContactPerson,
		Mobile:         req.Mobile,
		ClientNumLimit: int32(req.ClientNumLimit),
	})
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpBizInfoService.Insert err. BizId[%d] err:%+v", sessionId, newUserResp.UserInfo.BizId, err)
		return nil, errorx.NewDefaultError(errorx.CreateUserInfoErrorCode)
	}

	return &types.CreateBizResp{
		BizId: newUserResp.UserInfo.BizId,
	}, nil
}
