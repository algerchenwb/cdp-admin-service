package user

import (
	"context"
	"fmt"
	"time"

	"cdp-admin-service/internal/helper"
	"cdp-admin-service/internal/helper/cache"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteUserLogic) DeleteUser(req *types.DeleteUserReq) error {
	sessionId := helper.GetSessionId(l.ctx)
	userName := helper.GetUserName(l.ctx)
	updateUser, _, err := table.T_TCdpSysUserService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d", req.Id), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] table.TCdpSysUserService.Query err. Id[%v] err:%+v", sessionId, req.Id, err)
		return errorx.NewDefaultError(errorx.ParamErrorCode)
	}
	modifier, _, err := table.T_TCdpSysUserService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d", helper.GetUserId(l.ctx)), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] table.TCdpSysUserService.Query err. Id[%v] err:%+v", sessionId, helper.GetUserId(l.ctx), err)
		return errorx.NewDefaultError(errorx.ParamErrorCode)
	}

	access, err := checkDeleteUser(l.ctx, l.svcCtx.Cache, modifier, updateUser)
	if err != nil {
		l.Logger.Errorf("[%s] checkDeleteUser err. Id[%v] err:%+v", sessionId, req.Id, err)
		return err
	}
	if !access {
		l.Logger.Errorf("[%s] checkDeleteUser err. Id[%v] err:%+v", sessionId, req.Id, errorx.DeleteUserRefusedErrorCode)
		return errorx.NewDefaultError(errorx.DeleteUserRefusedErrorCode)
	}

	_, err = l.svcCtx.Cache.Del(l.ctx, sessionId, cache.UserPermKey(req.Id))
	if err != nil {
		l.Logger.Errorf("[%s] svc.Cache.Del err. Id[%v] err:%+v", sessionId, req.Id, err)
	}
	_, _ = l.svcCtx.Cache.Del(l.ctx, sessionId, cache.UserOnlineKey(req.Id))
	if err != nil {
		l.Logger.Errorf("[%s] svc.Cache.Del err. Id[%v] err:%+v", sessionId, req.Id, err)
	}

	_, _, err = table.T_TCdpSysUserService.Update(l.ctx, sessionId, req.Id, map[string]interface{}{
		"status":    table.SysUserStatusDisable,
		"account":   fmt.Sprintf("del_%s_%s", updateUser.Account, time.Now().Format("20060102150405")),
		"update_by": userName,
	})
	if err != nil {
		l.Logger.Errorf("[%s] table.TCdpSysUserService.Update err. Id[%v] err:%+v", sessionId, req.Id, err)
		return errorx.NewDefaultError(errorx.ParamErrorCode)
	}

	return nil
}

func checkDeleteUser(ctx context.Context, c *cache.Cache, modifier *table.TCdpSysUser, waitUpdateUser *table.TCdpSysUser) (bool, error) {
	sessionId := helper.GetSessionId(ctx)
	//  不能删除自己
	if modifier.Id == waitUpdateUser.Id {
		logx.WithContext(ctx).Errorf("[%s] modifier[%v] cannot delete self", sessionId, modifier)
		return false, nil
	}
	// 超级管理员可以删除所有账号
	if modifier.IsAdmin == table.IsAdminYes {
		logx.WithContext(ctx).Infof("[%s] modifier[%v] is admin, can delete all user", sessionId, modifier)
		return true, nil
	}

	// 不能删除管理员
	roleIsAdmin, err := c.Sismember(ctx, sessionId, cache.RoleAdminKey(), fmt.Sprint(waitUpdateUser.RoleId))
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] svc.Cache.Sismember err. Id[%v] err:%+v", sessionId, waitUpdateUser.Id, err)
		return false, errorx.NewDefaultError(errorx.QueryRoleFailedErrorCode)
	}
	if roleIsAdmin {
		logx.WithContext(ctx).Errorf("[%s] waitUpdateUser is admin. Id[%v]", sessionId, waitUpdateUser.Id)
		return false, nil
	}

	waitUpdateAreaIds, err := getUserAreaIds(ctx, waitUpdateUser)
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] getUserAreaIds err. Id[%v] err:%+v", sessionId, waitUpdateUser.Id, err)
		return false, errorx.NewDefaultError(errorx.QueryAreaFailedErrorCode)
	}
	modifierAreaIds, err := getUserAreaIds(ctx, modifier)
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] getUserAreaIds err. Id[%v] err:%+v", sessionId, modifier.Id, err)
		return false, errorx.NewDefaultError(errorx.QueryAreaFailedErrorCode)
	}
	logx.WithContext(ctx).Infof("[%s] modifierAreaIds[%+v] waitUpdateAreaIds[%+v]", sessionId, modifierAreaIds, waitUpdateAreaIds)

	return checkAreaIdsAccessV2(modifierAreaIds, waitUpdateAreaIds), nil
}
