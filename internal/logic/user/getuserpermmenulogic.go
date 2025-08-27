package user

import (
	"cdp-admin-service/internal/helper"
	"cdp-admin-service/internal/logic/common"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserPermMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserPermMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserPermMenuLogic {
	return &GetUserPermMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserPermMenuLogic) GetUserPermMenu() (resp *types.UserPermMenuResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)
	roleId := helper.GetRoleId(l.ctx)
	isAdmin := helper.GetIsAdmin(l.ctx)
	// 用户所属角色
	resp = &types.UserPermMenuResp{}
	if isAdmin {
		l.Logger.Infof("[%s] is admin, get admin account menus", sessionId)
		resp.Menus, err = common.GetAdminAccountMenus(l.ctx, l.svcCtx)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("[%s] get admin account menus failed. err:%+v", sessionId, err)
			return nil, errorx.NewDefaultError(errorx.QueryMenuFailedErrorCode)
		}
		resp.Menus = common.FixMenu(resp.Menus)
		return
	}

	resp.Menus, err = common.GetRoleMenuPerms(l.ctx, l.svcCtx, int64(roleId))
	if err != nil {
		logx.WithContext(l.ctx).Errorf("[%s] get admin account menus failed. err:%+v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.QueryMenuFailedErrorCode)
	}
	resp.Menus = common.FixMenu(resp.Menus)
	return
}
