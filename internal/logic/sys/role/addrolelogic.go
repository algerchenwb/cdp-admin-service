package role

import (
	"context"
	"fmt"
	"time"

	"cdp-admin-service/internal/helper"
	"cdp-admin-service/internal/helper/cache"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/logic/common"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddRoleLogic {
	return &AddRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddRoleLogic) AddRole(req *types.AddRoleReq) (resp *types.AddRoleResp, err error) {
	sessionId := helper.GetSessionId(l.ctx)
	username := helper.GetUserName(l.ctx)
	//  施工平台只能有非超管角色
	if req.Role.Platform == table.PlatformConstruction {
		req.Role.IsAdmin = table.RoleIsAdminNo
	}
	isAdmin := helper.GetIsAdmin(l.ctx)
	if !isAdmin {
		access, err := common.CheckRegionAccess(l.ctx, req.Role.RegionId)
		if err != nil {
			l.Logger.Errorf("[%s] 检查区域权限失败 err[%v]", sessionId, err)
			return nil, errorx.NewDefaultError(errorx.CheckRegionAccessErrorCode)
		}
		if !access {
			l.Logger.Debugf("[%s] 无权限操作该节点域数据", sessionId)
			return nil, errorx.NewDefaultError(errorx.CheckRegionAccessErrorCode)
		}
	}

	role := &table.TCdpSysRole{
		Name:       req.Role.Name,
		Remark:     req.Role.Remark,
		Platform:   req.Role.Platform,
		IsAdmin:    req.Role.IsAdmin,
		RegionId:   uint32(req.Role.RegionId),
		Status:     table.RoleStatusEnable,
		CreateBy:   username,
		UpdateBy:   username,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		ModifyTime: time.Now(),
	}
	role, _, err = table.T_TCdpSysRoleService.Insert(l.ctx, sessionId, role)
	if err != nil {
		l.Logger.Error("[%s] 添加角色失败 err[%v]", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.InsertRoleFailedErrorCode)
	}
	if role.RoleIsAdmin() {
		err = l.svcCtx.Cache.SAdd(l.ctx, sessionId, cache.RoleAdminKey(), fmt.Sprint(role.Id))
		if err != nil {
			l.Logger.Error("[%s] 添加角色失败 err[%v]", sessionId, err)
			return nil, errorx.NewDefaultError(errorx.InsertRoleFailedErrorCode)
		}
	}
	return &types.AddRoleResp{
		Id: int64(role.Id),
	}, nil
}
