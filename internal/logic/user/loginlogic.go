package user

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cdp-admin-service/internal/helper/cache"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/model/errorx"

	"cdp-admin-service/internal/helper"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	Logger logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq, r *http.Request) (resp *types.LoginResp, err error) {
	sessionId := helper.GetSessionId(l.ctx)

	sysUser, _, err := table.T_TCdpSysUserService.Query(l.ctx, sessionId,
		fmt.Sprintf("account:%s$platform:%d$status__ex:%d", req.Account, req.Platform, table.SysUserStatusDisable), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] FindOneByAccount Account[%s] err :%+v", sessionId, req.Account, err)
		return nil, errorx.NewDefaultError(errorx.AccountErrorCode)
	}
	//  用户被禁用
	// if sysUser.Status != table.SysUserStatusEnable {
	// 	l.Logger.Errorf("[%s] Account[%s] is disabled", sessionId, req.Account)
	// 	return nil, errorx.NewDefaultError(errorx.AccountDisableErrorCode)
	// }

	password, err := helper.Decode(req.Password)
	if err != nil {
		l.Logger.Errorf("[%s] Decode Password[%s] err :%+v", sessionId, req.Password, err)
		return nil, errorx.NewDefaultError(errorx.PasswordErrorCode)
	}
	password = helper.MD5(password + l.svcCtx.Config.AES.Salt)
	if sysUser.Password != password {
		l.Logger.Errorf("[%s] Password[%s] is incorrect", sessionId, password)
		return nil, errorx.NewDefaultError(errorx.PasswordErrorCode)
	}

	token, err := helper.GenerateToken(int64(sysUser.Id), sysUser.Nickname,
		int64(sysUser.RoleId), int64(sysUser.Platform),
		sysUser.IsAdmin == table.IsAdminYes, l.svcCtx.Cache, int(l.svcCtx.Config.JwtAuth.AccessExpire))
	if err != nil {
		l.Logger.Errorf("[%s] Account[%s] GetJwtToken err :%+v", sessionId, req.Account, err)
		return nil, errorx.NewDefaultError(errorx.PasswordErrorCode)
	}
	err = l.svcCtx.Cache.SetEX(l.ctx, sessionId, cache.UserOnlineKey(int64(sysUser.Id)), token, int(l.svcCtx.Config.JwtAuth.AccessExpire))
	if err != nil {
		l.Logger.Errorf("[%s] Account[%s] Redis.Setex err :%+v", sessionId, req.Account, err)
		return nil, errorx.NewCodeError(errorx.ServerErrorCode, err.Error())
	}

	sysLog := &table.TCdpSysLog{
		UserId:     sysUser.Id,
		Account:    sysUser.Account,
		Platform:   sysUser.Platform,
		Ip:         r.RemoteAddr,
		Uri:        r.RequestURI,
		Type:       1,
		Request:    req.Account,
		Response:   token,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		ModifyTime: time.Now(),
		Status:     1,
	}

	_, _, err = table.T_TCdpSysLogService.Insert(l.ctx, sessionId, sysLog)
	if err != nil {
		l.Logger.Errorf("[%s] Account[%s] SysLogModel.Insert fail userId[%d] err:%+v", sessionId, req.Account, sysUser.Id, err)
	}

	return &types.LoginResp{
		Token: token,
		UserInfo: types.UserInfoResp{
			Id:          int64(sysUser.Id),
			Account:     sysUser.Account,
			Nickname:    sysUser.Nickname,
			Avatar:      sysUser.Avatar,
			Mobile:      sysUser.Mobile,
			RoleId:      int64(sysUser.RoleId),
			AreaIds:     helper.StringToInt64Slice(sysUser.AreaIds),
			AreaRegions: helper.StringToInt64Slice(sysUser.AreaRegions),
			BizIds:      helper.StringToInt64Slice(sysUser.BizIds),
			BizRegions:  helper.StringToInt64Slice(sysUser.BizRegions),
			Platform:    int64(sysUser.Platform),
			Status:      int64(sysUser.Status),
		},
	}, nil

}

// 确定系统默认角色
func confirmDefaultRoleIds() (defaultRoleIds string, defaultRoleId int) {
	return "[]", 0
}
