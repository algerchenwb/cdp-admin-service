package cloudbox

import (
	"context"

	"cdp-admin-service/internal/helper"
	"cdp-admin-service/internal/helper/diskless"
	"cdp-admin-service/internal/model/errorx"
	instance_types "cdp-admin-service/internal/proto/instance_service/types"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CloudBoxCancelAdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCloudBoxCancelAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CloudBoxCancelAdminLogic {
	return &CloudBoxCancelAdminLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CloudBoxCancelAdminLogic) CloudBoxCancelAdmin(req *types.CloudBoxCancelAdminReq) (resp *types.CloudBoxCancelAdminResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)

	// 该云盒是否已经是超管状态
	var macList = []string{req.MAC}

	Instancelist, err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).GetInstanceList(req.AreaId, sessionId, macList)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("[%s] GetInstanceList AreaId[%d] err:%+v", sessionId, req.AreaId, err)
		return nil, errorx.NewDefaultError(errorx.DisklessListInstanceErrorCode)
	}

	if len(Instancelist) > 0 && (Instancelist[0].UserMode == instance_types.RegularUser || Instancelist[0].UserMode == instance_types.RegularUser2) {
		return nil, errorx.NewDefaultError(errorx.CloudBoxUserModeErrorCode)
	}

	// 调用制作镜像接口
	if req.OsVersion != "" && req.Name != "" {

		if err := diskless.CreateImageFromAreaInstance(l.ctx,
			l.svcCtx.Config.OutSide.DisklessCloudImageHost,
			sessionId,
			req.ImageId,
			req.Name,
			req.OsVersion,
			req.Remark,
			req.ManagerState,
			req.BizId,
			req.AreaId,
			req.InstanceId,
			int32(req.FlattenFlag)); err != nil {
			return nil, errorx.NewDefaultError(errorx.DisklessCreateImageErrorCode)
		}
	}

	// 取消超管
	setAdminReq := &instance_types.SetAdminRequest{
		FlowID:     sessionId,
		AppID:      "diskless-aggregator",
		Mac:        req.MAC,
		InstanceID: req.InstanceId,
		UserMode:   instance_types.RegularUser2, // 设置取消超管
	}
	if err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).SetAdmin(req.AreaId, setAdminReq); err != nil {
		return nil, errorx.NewDefaultCodeError(err.Error())
	}

	return
}
