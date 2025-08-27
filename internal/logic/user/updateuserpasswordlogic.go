package user

import (
	"context"
	"fmt"

	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserPasswordLogic {
	return &UpdateUserPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserPasswordLogic) UpdateUserPassword(req *types.UpdatePasswordReq) error {
	sessionId := helper.GetSessionId(l.ctx)
	userId := helper.GetUserId(l.ctx)
	modifier, _, err := table.T_TCdpSysUserService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d", userId), nil, nil)
	if err != nil {
		l.Logger.Errorf("table.T_TCdpSysUserService.Query err. Id[%d] err:%+v", userId, err)
		return errorx.NewDefaultError(errorx.UserNotFoundErrorCode)
	}
	waitUpdateUser, _, err := table.T_TCdpSysUserService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d", req.Id), nil, nil)
	if err != nil {
		l.Logger.Errorf("table.T_TCdpSysUserService.Query err. Id[%d] err:%+v", req.Id, err)
		return errorx.NewDefaultError(errorx.UserNotFoundErrorCode)
	}
	access, err := checkUpdateUser(l.ctx, modifier, waitUpdateUser, l.svcCtx.Cache)
	if err != nil {
		l.Logger.Errorf("checkUpdateUser err. modifier[%v] waitUpdateUser[%v] err:%+v", modifier, waitUpdateUser, err)
		return errorx.NewDefaultError(errorx.UpdateUserRefusedErrorCode)
	}
	if !access {
		l.Logger.Errorf("checkUpdateUser no access. modifier[%v] waitUpdateUser[%v]", modifier, waitUpdateUser)
		return errorx.NewDefaultError(errorx.UpdateUserRefusedErrorCode)
	}
	oldPassword, err := helper.Decode(req.OldPassword)
	if err != nil {
		l.Logger.Errorf("helper.Decode err. Password[%s] err:%+v", req.OldPassword, err)
		return errorx.NewDefaultError(errorx.DecodePasswordError)
	}
	if waitUpdateUser.Password != helper.HashPassword(oldPassword, l.svcCtx.Config.AES.Salt) {
		l.Logger.Errorf("oldPassword[%s] not match. modifier[%v] waitUpdateUser[%v]", oldPassword, modifier, waitUpdateUser)
		return errorx.NewDefaultError(errorx.OldPasswordError)
	}
	newPassword, err := helper.Decode(req.NewPassword)
	if err != nil {
		l.Logger.Errorf("helper.Decode err. Password[%s] err:%+v", req.NewPassword, err)
		return errorx.NewDefaultError(errorx.DecodePasswordError)
	}

	password := helper.HashPassword(newPassword, l.svcCtx.Config.AES.Salt)
	_, _, err = table.T_TCdpSysUserService.Update(l.ctx, sessionId, req.Id, map[string]interface{}{
		"password":  password,
		"update_by": modifier.Nickname,
	})
	if err != nil {
		l.Logger.Errorf("table.T_TCdpSysUserService.Update err. Id[%d] err:%+v", userId, err)
		return errorx.NewDefaultError(errorx.UpdateUserFailedErrorCode)
	}

	return nil
}
