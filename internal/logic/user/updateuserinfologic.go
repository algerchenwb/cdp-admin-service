package user

import (
	"context"
	"fmt"

	"cdp-admin-service/internal/helper"
	"cdp-admin-service/internal/helper/cache"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/logic/common"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(req *types.UpdateUserInfoReq) error {
	sessionId := helper.GetSessionId(l.ctx)

	// 校验修改用户权限
	modifier, _, err := table.T_TCdpSysUserService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d", helper.GetUserId(l.ctx)), nil, nil)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("table.T_TCdpSysUserService.Query err. Id[%v] err:%+v", helper.GetUserId(l.ctx), err)
		return errorx.NewDefaultError(errorx.UpdateUserFailedErrorCode)
	}
	waitUpdateUser, _, err := table.T_TCdpSysUserService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d", req.Id), nil, nil)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("table.T_TCdpSysUserService.Query err. Id[%v] err:%+v", req.Id, err)
		return errorx.NewDefaultError(errorx.UpdateUserFailedErrorCode)
	}

	access, err := checkUpdateUser(l.ctx, modifier, waitUpdateUser, l.svcCtx.Cache)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("checkUpdateUser err. Id[%v] err:%+v", req.Id, err)
		return errorx.NewDefaultError(errorx.UpdateUserFailedErrorCode)
	}
	if !access {
		return errorx.NewDefaultError(errorx.UpdateUserRefusedErrorCode)
	}

	waitUpdateUser, _, err = table.T_TCdpSysUserService.Update(l.ctx, sessionId, req.Id, map[string]interface{}{
		"nickname":     req.Nickname,
		"avatar":       req.Avatar,
		"remark":       req.Remark,
		"role_id":      req.RoleId,
		"update_by":    helper.GetUserName(l.ctx),
		"area_ids":     helper.SliceToString(req.AreaIds),
		"area_regions": helper.SliceToString(req.AreaRegions),
		"biz_ids":      helper.SliceToString(req.BizIds),
		"biz_regions":  helper.SliceToString(req.BizRegions),
		"mobile":       req.Mobile,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("table.T_SysUserService.Update err. Id[%v] err:%+v", req.Id, err)
		return errorx.NewDefaultError(errorx.UpdateUserFailedErrorCode)
	}

	err = common.UpdateUserPermCache(l.ctx, waitUpdateUser, l.svcCtx.Cache)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("common.UpdateUserRole err. Id[%v] err:%+v", req.Id, err)
		return errorx.NewDefaultError(errorx.UpdateUserFailedErrorCode)
	}

	return nil
}

func checkUpdateUserPerm(ctx context.Context, modifier *table.TCdpSysUser, waitUpdateUser *table.TCdpSysUser, c *cache.Cache) (bool, error) {
	sessionId := helper.GetSessionId(ctx)

	// 非管理员账号不能编辑用户信息
	roleIsAdmin, err := c.Sismember(ctx, sessionId, cache.RoleAdminKey(), fmt.Sprint(modifier.RoleId))
	if err != nil {
		logx.WithContext(ctx).Errorf("c.Sismember err. RoleId[%d] err:%+v", modifier.RoleId, err)
		return false, errorx.NewDefaultError(errorx.UpdateUserRefusedErrorCode)
	}
	if !roleIsAdmin {
		logx.WithContext(ctx).Errorf("modifier[%+v] is not admin, no access", modifier)
		return false, nil
	}
	// 管理员角色只能编辑普通角色的账号
	waitUpdateRoleIsAdmin, err := c.Sismember(ctx, sessionId, cache.RoleAdminKey(), fmt.Sprint(waitUpdateUser.RoleId))
	if err != nil {
		logx.WithContext(ctx).Errorf("c.Sismember err. RoleId[%d] err:%+v", waitUpdateUser.RoleId, err)
		return false, errorx.NewDefaultError(errorx.UpdateUserRefusedErrorCode)
	}
	if waitUpdateRoleIsAdmin {
		logx.WithContext(ctx).Errorf("modifier[%+v] not admin, waitUpdateUser[%+v] is admin, no access", modifier, waitUpdateUser)
		return false, nil
	}
	return true, nil
}

func checkUpdateUserAreaIds(ctx context.Context, modifier *table.TCdpSysUser, waitUpdateUser *table.TCdpSysUser) (bool, error) {
	modifierAreaIds, err := getUserAreaIds(ctx, modifier)
	if err != nil {
		logx.WithContext(ctx).Errorf("getUserAreaIds err. Id[%v] err:%+v", helper.GetUserId(ctx), err)
		return false, errorx.NewDefaultError(errorx.UpdateUserFailedErrorCode)
	}
	waitUpdateAreaIds, err := getUserAreaIds(ctx, waitUpdateUser)
	if err != nil {
		logx.WithContext(ctx).Errorf("getUserAreaIds err. Id[%v] err:%+v", waitUpdateUser.Id, err)
		return false, errorx.NewDefaultError(errorx.UpdateUserFailedErrorCode)
	}
	return checkAreaIdsAccessV2(modifierAreaIds, waitUpdateAreaIds), nil
}

func checkUpdateUser(ctx context.Context, modifier *table.TCdpSysUser, waitUpdateUser *table.TCdpSysUser, c *cache.Cache) (bool, error) {
	// 能编辑修改自己信息
	if modifier.Id == waitUpdateUser.Id {
		return true, nil
	}
	// 超管账号能编辑修改所有账号
	if modifier.IsAdmin == table.IsAdminYes {
		return true, nil
	}

	access, err := checkUpdateUserPerm(ctx, modifier, waitUpdateUser, c)
	if err != nil {
		logx.WithContext(ctx).Errorf("checkUpdateUserPerm err. modifier[%v] waitUpdateUser[%v] err:%+v", modifier, waitUpdateUser, err)
		return false, err
	}
	if !access {
		logx.WithContext(ctx).Errorf("checkUpdateUserPerm no access. modifier[%v] waitUpdateUser[%v]", modifier, waitUpdateUser)
		return false, nil
	}
	return checkUpdateUserAreaIds(ctx, modifier, waitUpdateUser)
}
