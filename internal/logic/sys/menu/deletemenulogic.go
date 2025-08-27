package menu

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cdp-admin-service/internal/helper"
	"cdp-admin-service/internal/helper/cache"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"cdp-admin-service/internal/logic/common"

	"github.com/zeromicro/go-zero/core/logx"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
)

type DeleteMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMenuLogic {
	return &DeleteMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMenuLogic) DeleteMenu(req *types.DeleteMenuReq) (resp *types.DeleteMenuResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)
	username := helper.GetUserName(l.ctx)
	menu, _, err := table.T_TCdpSysPermMenuService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d$status:%d", req.Id, table.MenuStatusEnable), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] 查询菜单失败 err[%v]", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.QueryMenuFailedErrorCode)
	}

	childMenu, _, err := table.T_TCdpSysPermMenuService.Query(l.ctx, sessionId, fmt.Sprintf("parent_id:%d$status:%d", req.Id, table.MenuStatusEnable), nil, nil)
	if err != nil && !errors.Is(err, gopublic.ErrNotExist) {
		l.Logger.Errorf("[%s] 查询子菜单失败 err[%v]", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.QueryMenuFailedErrorCode)
	}
	if err == nil && childMenu != nil {
		l.Logger.Errorf("[%s] 子菜单存在不可删除 childMenu[%+v]", sessionId, childMenu)
		return nil, errorx.NewDefaultError(errorx.ExistChildMenuErrorCode)
	}

	// 移除角色对应菜单
	roles, _, err := table.T_TCdpSysRoleService.QueryAll(l.ctx, sessionId, fmt.Sprintf("perm_menu_ids__contains:%d$status:%d", req.Id, table.RoleStatusEnable), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] 查询角色失败 err[%v]", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.QueryRoleFailedErrorCode)
	}

	for _, role := range roles {
		newRole, _, err := table.T_TCdpSysRoleService.Update(l.ctx, sessionId, role.Id, map[string]interface{}{
			"perm_menu_ids": helper.RemoveFromStr(role.PermMenuIds, fmt.Sprintf("%d", req.Id)),
			"update_by":     username,
			"update_time":   time.Now(),
		})
		if err != nil {
			l.Logger.Errorf("[%s] 更新角色失败 err[%v]", sessionId, err)
			return nil, errorx.NewDefaultError(errorx.UpdateRoleFailedErrorCode)
		}
		err = common.UpdateRolePermCache(l.ctx, l.svcCtx, newRole)
		if err != nil {
			l.Logger.Errorf("[%s] 更新角色权限失败 err[%v]", sessionId, err)
			return nil, errorx.NewDefaultError(errorx.UpdateUserPermFailedErrorCode)
		}
	}

	if menu.IsPrivate == table.IsPrivate {
		perm, err := common.ParseMenuPerm(menu)
		if err != nil {
			l.Logger.Errorf("[%s] 解析菜单权限失败 err[%v]", sessionId, err)
			return nil, errorx.NewDefaultError(errorx.ParseMenuPermFailedErrorCode)
		}
		for _, p := range perm {
			err = l.svcCtx.Cache.SRem(l.ctx, sessionId, cache.PermIsPrivateKey(), helper.PermStandardBase64(p))
			if err != nil {
				l.Logger.Errorf("[%s] 删除内网菜单权限失败 err[%v]", sessionId, err)
			}
		}
	}

	// 删除菜单
	_, _, err = table.T_TCdpSysPermMenuService.Update(l.ctx, sessionId, menu.Id, map[string]interface{}{
		"status": 0,
	})
	if err != nil {
		l.Logger.Errorf("[%s] 删除菜单失败 err[%v]", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.DeleteMenuFailedErrorCode)
	}

	return &types.DeleteMenuResp{
		Id: int64(menu.Id),
	}, nil
}
