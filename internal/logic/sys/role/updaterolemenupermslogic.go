package role

import (
	"context"
	"fmt"
	"time"

	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/logic/common"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRoleMenuPermsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateRoleMenuPermsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleMenuPermsLogic {
	return &UpdateRoleMenuPermsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRoleMenuPermsLogic) UpdateRoleMenuPerms(req *types.UpdateRoleMenuPermsReq) (resp *types.UpdateRoleMenuPermsResp, err error) {
	sessionId := helper.GetSessionId(l.ctx)
	username := helper.GetUserName(l.ctx)
	isAdmin := helper.GetIsAdmin(l.ctx)

	role, _, err := table.T_TCdpSysRoleService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d", req.Id), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] 查询角色失败 err[%v]", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.QueryRoleFailedErrorCode)
	}
	if !isAdmin {
		access, err := common.CheckRegionAccess(l.ctx, int64(role.RegionId))
		if err != nil {
			l.Logger.Error("[%s] 检查区域权限失败 err[%v]", sessionId, err)
			return nil, errorx.NewDefaultError(errorx.CheckRegionAccessErrorCode)
		}
		if !access {
			l.Logger.Error("[%s] 无权限操作区域 regionId[%d]", sessionId, role.RegionId)
			return nil, errorx.NewDefaultError(errorx.CheckRegionAccessErrorCode)
		}
	}

	role, _, err = table.T_TCdpSysRoleService.Update(l.ctx, sessionId, role.Id, map[string]interface{}{
		"perm_menu_ids": helper.SliceToString(req.MenuIds),
		"update_by":     username,
		"update_time":   time.Now(),
	})
	if err != nil {
		l.Logger.Errorf("[%s] 更新角色失败 err[%v]", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.UpdateRoleFailedErrorCode)
	}

	// 非管理员角色配置菜单权限
	if !role.RoleIsAdmin() {
		err = common.UpdateRolePermCache(l.ctx, l.svcCtx, role)
		if err != nil {
			l.Logger.Errorf("[%s] 更新角色权限缓存失败 err[%v]", sessionId, err)
			return nil, errorx.NewDefaultError(errorx.UpdateRoleFailedErrorCode)
		}
	}

	return &types.UpdateRoleMenuPermsResp{
		Id: req.Id,
	}, nil
}
