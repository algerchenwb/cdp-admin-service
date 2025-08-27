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

type ResetPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetPasswordLogic) ResetPassword(req *types.ResetPasswordReq) (resp *types.ResetPasswordResp, err error) {
	sessionId := helper.GetSessionId(l.ctx)
	userName := helper.GetUserName(l.ctx)
	updateUser, _, err := table.T_TCdpSysUserService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d", req.Id), nil, nil)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("table.T_TCdpSysUserService.Query err. Id[%v] err:%+v", req.Id, err)
		return nil, errorx.NewDefaultError(errorx.QueryUserFailedErrorCode)
	}
	modifier, _, err := table.T_TCdpSysUserService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d", helper.GetUserId(l.ctx)), nil, nil)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("table.T_TCdpSysUserService.Query err. Id[%v] err:%+v", helper.GetUserId(l.ctx), err)
		return nil, errorx.NewDefaultError(errorx.QueryUserFailedErrorCode)
	}

	access, err := checkUpdateUser(l.ctx, modifier, updateUser, l.svcCtx.Cache)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("checkUpdateUser err. Id[%v] err:%+v", req.Id, err)
		return nil, errorx.NewDefaultError(errorx.UpdateUserFailedErrorCode)
	}
	if !access {
		l.Logger.Errorf("[%s] modifier[%v] updateUser[%v] no access", sessionId, modifier, updateUser)
		return nil, errorx.NewDefaultError(errorx.IdcRefuseErrorCode)
	}
	newPassword := helper.HashPassword(l.svcCtx.Config.Account.DefaultPassword, l.svcCtx.Config.AES.Salt)

	_, _, err = table.T_TCdpSysUserService.Update(l.ctx, sessionId, req.Id, map[string]interface{}{
		"password":  newPassword,
		"update_by": userName,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("table.T_TCdpSysUserService.Update err. Id[%v] err:%+v", req.Id, err)
		return nil, errorx.NewDefaultError(errorx.UpdateUserFailedErrorCode)
	}

	return &types.ResetPasswordResp{
		NewPassword: newPassword,
	}, nil
}
