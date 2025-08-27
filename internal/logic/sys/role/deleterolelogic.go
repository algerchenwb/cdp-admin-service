package role

import (
	"context"
	"fmt"
	"strings"
	"time"

	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/logic/common"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"cdp-admin-service/internal/model/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRoleLogic {
	return &DeleteRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteRoleLogic) DeleteRole(req *types.DeleteRoleReq) (resp *types.DeleteRoleResp, err error) {
	sessionId := helper.GetSessionId(l.ctx)
	username := helper.GetUserName(l.ctx)

	isAdmin := helper.GetIsAdmin(l.ctx)

	users, _, err := table.T_TCdpSysUserService.QueryAll(l.ctx, sessionId, fmt.Sprintf("role_id:%d$status__ex:%d", req.Id, table.SysUserStatusDisable), nil, nil)
	if err != nil {
		l.Logger.Error("[%s] 查询角色失败 err[%v]", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.UserQueryFailedErrorCode)
	}

	if len(users) > 0 {
		var userIds []string
		for _, user := range users {
			userIds = append(userIds, fmt.Sprintf("%d", user.Id))
		}
		l.Logger.Error("[%s] 删除角色失败 角色下挂靠的用户[%s]", sessionId, strings.Join(userIds, ","))
		return nil, errorx.NewDefaultError(errorx.DeleteRoleRefusedErrorCode)
	}
	role, _, err := table.T_TCdpSysRoleService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d", req.Id), nil, nil)
	if err != nil {
		l.Logger.Error("[%s] 查询角色失败 err[%v]", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.DeleteRoleFailedErrorCode)
	}
	if !isAdmin {
		access, err := common.CheckRegionAccess(l.ctx, int64(role.RegionId))
		if err != nil {
			l.Logger.Errorf("[%s] 检查区域权限失败 err[%v]", sessionId, err)
			return nil, errorx.NewDefaultError(errorx.CheckRegionAccessErrorCode)
		}
		if !access {
			l.Logger.Debugf("[%s] 无权限操作该节点域数据", sessionId)
			return nil, errorx.NewDefaultError(errorx.CheckRegionAccessErrorCode)
		}
	}
	_, _, err = table.T_TCdpSysRoleService.Update(l.ctx, sessionId, req.Id, map[string]any{
		"status":      table.RoleStatusDisable,
		"name":        fmt.Sprintf("%s_deleted_%d", role.Name, time.Now().Unix()),
		"update_by":   username,
		"update_time": time.Now(),
	})
	if err != nil {
		l.Logger.Error("[%s] 删除角色失败 err[%v]", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.DeleteRoleFailedErrorCode)
	}

	return &types.DeleteRoleResp{
		Id: req.Id,
	}, nil
}
