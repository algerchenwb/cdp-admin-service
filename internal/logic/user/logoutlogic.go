package user

import (
	"context"

	"cdp-admin-service/internal/helper"
	"cdp-admin-service/internal/helper/cache"
	"cdp-admin-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout() error {
	userId := helper.GetUserId(l.ctx)
	sessionId, _ := helper.GenSessionId(l.ctx)
	// _, _ = l.svcCtx.Redis.Del(globalkey.SysPermMenuCachePrefix + userId)
	_, _ = l.svcCtx.Cache.Del(l.ctx, sessionId, cache.UserOnlineKey(userId))
	// _, _ = l.svcCtx.Redis.Del(globalkey.SysUserIdCachePrefix + userId)

	return nil
}
