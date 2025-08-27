package role

import (
	"context"
	"time"

	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/logic/common"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"cdp-admin-service/internal/model/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleLogic {
	return &UpdateRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRoleLogic) UpdateRole(req *types.UpdateRoleReq) (resp *types.UpdateRoleResp, err error) {
	sessionId := helper.GetSessionId(l.ctx)
	username := helper.GetUserName(l.ctx)
	// isAdmin := helper.GetIsAdmin(l.ctx)
	// roleId := helper.GetRoleId(l.ctx)
	// roleIsAdmin, err := common.RoleIsAdmin(l.ctx, l.svcCtx, int64(roleId))
	// if err != nil {
	// 	l.Logger.Error("[%s] 查询角色是否为管理员失败 err[%v]", sessionId, err)
	// 	return nil, errorx.NewDefaultError(errorx.QueryRoleFailedErrorCode)
	// }
	// if !isAdmin && !roleIsAdmin {
	// 	l.Logger.Error("[%s] 无权限操作角色 accoutIsAdmin[%v] roleIsAdmin[%v]", sessionId, isAdmin, roleIsAdmin)
	// 	return nil, errorx.NewDefaultError(errorx.HandlerRoleRefusedErrorCode)
	// }
	isAdmin := helper.GetIsAdmin(l.ctx)
	if !isAdmin {
		access, err := common.CheckRegionAccess(l.ctx, int64(req.Role.RegionId))
		if err != nil {
			l.Logger.Error("[%s] 检查区域权限失败 err[%v]", sessionId, err)
			return nil, errorx.NewDefaultError(errorx.CheckRegionAccessErrorCode)
		}
		if !access {
			l.Logger.Error("[%s] 无权限操作区域 regionId[%d]", sessionId, req.Role.RegionId)
			return nil, errorx.NewDefaultError(errorx.CheckRegionAccessErrorCode)
		}
	}
	role, _, err := table.T_TCdpSysRoleService.Update(l.ctx, sessionId, req.Role.Id, map[string]any{
		"name":        req.Role.Name,
		"remark":      req.Role.Remark,
		"update_by":   username,
		"update_time": time.Now(),
	})
	if err != nil {
		l.Logger.Error("[%s] 更新角色失败 err[%v]", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.UpdateRoleFailedErrorCode)
	}
	return &types.UpdateRoleResp{
		Id: int64(role.Id),
	}, nil
}
