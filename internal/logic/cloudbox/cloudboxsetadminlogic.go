package cloudbox

import (
	"context"
	"fmt"

	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/helper/diskless"
	"cdp-admin-service/internal/model/errorx"
	instance_types "cdp-admin-service/internal/proto/instance_service/types"
	proto "cdp-admin-service/internal/proto/location_seat_service"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
)

type CloudBoxSetAdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCloudBoxSetAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CloudBoxSetAdminLogic {
	return &CloudBoxSetAdminLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CloudBoxSetAdminLogic) CloudBoxSetAdmin(req *types.CloudBoxSetAdminReq) (resp *types.CloudBoxSetAdminResp, err error) {
	sessionId := helper.GetSessionId(l.ctx)

	resp = new(types.CloudBoxSetAdminResp)

	// 查询 cdp 云盒信息，如果不存在再查询saas 的设备信息
	qry := fmt.Sprintf("mac:%s$status:1", req.MAC)
	cloudBoxInfo, _, err := table.T_TCdpCloudboxInfoService.Query(l.ctx, sessionId, qry, nil, nil)
	if err != nil && err != gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] Query T_TCdpCloudboxInfo failed. err: %v", sessionId, err)
		return nil, errorx.NewDefaultCodeError("查询云盒信息失败")
	}

	if err == gopublic.ErrNotExist {
		l.Logger.Infof("[%s] Query T_TCdpCloudboxInfo failed. mac:%s err: %v", sessionId, req.MAC, err)
		return nil, errorx.NewDefaultCodeError("云盒信息不存在")
	}
	if cloudBoxInfo.BootType == int32(proto.BootType_BOOTTYPE_LOCAL) {
		l.Logger.Infof("[%s] CloudBoxSetAdmin failed. mac:%s err: %v", sessionId, req.MAC, err)
		return nil, errorx.NewDefaultCodeError("本地启动的云盒不能设置超管")
	}
	// 该云盒是否已经是超管状态
	var macList = []string{req.MAC}
	Instancelist, err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).GetInstanceList(req.AreaId, sessionId, macList)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("[%s] GetInstanceList failed,macList:%s  err:%+v", sessionId, macList, err)
		return nil, errorx.NewHandlerError(errorx.ServerErrorCode, "查询无盘实例接口失败")
	}

	if len(Instancelist) > 0 && (Instancelist[0].UserMode == 0 || Instancelist[0].UserMode == 2) {
		return nil, errorx.NewHandlerError(errorx.ServerErrorCode, "该云盒已经是超管状态")
	}


	// 设置超管
	setAdminReq := &instance_types.SetAdminRequest{
		FlowID:     sessionId,
		AppID:      "diskless-aggregator",
		Mac:        req.MAC,
		InstanceID: req.InstanceId,
		UserMode:   instance_types.AdminUser, // 设置超管
	}
	if err := diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).SetAdmin(req.AreaId, setAdminReq); err != nil {
		return nil, errorx.NewDefaultError(errorx.DisklessUserModeAdminUserErrorCode)
	}

	return
}
