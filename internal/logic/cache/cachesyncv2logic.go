package cache

import (
	"context"
	"fmt"

	"cdp-admin-service/internal/helper"
	"cdp-admin-service/internal/helper/cache"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/logic/common"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CacheSyncV2Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCacheSyncV2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *CacheSyncV2Logic {
	return &CacheSyncV2Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CacheSyncV2Logic) CacheSyncV2() (resp *types.CommonNilJson, err error) {

	sessionId := helper.GetSessionId(l.ctx)

	sysPermMenus, _, err := table.T_TCdpSysPermMenuService.QueryAll(l.ctx, sessionId, fmt.Sprintf("status:%d", table.MenuStatusEnable), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpSysPermMenuService.QueryAll fail err:%+v", sessionId, err)
		return nil, err
	}

	for _, sysPermMenu := range sysPermMenus {
		if sysPermMenu.IsPrivate != table.IsPrivate {
			continue
		}
		perms, err := common.ParseMenuPerm(&sysPermMenu)
		if err != nil {
			l.Logger.Errorf("[%s] common.ParseMenuPerm fail err:%+v", sessionId, err)
			return nil, err
		}
		for _, perm := range perms {
			perm = helper.PermStandardBase64(perm)
			if err := l.svcCtx.Cache.SAdd(l.ctx, sessionId, cache.PermIsPrivateKey(), perm); err != nil {
				l.Logger.Errorf("[%s] Redis.SaddCtx fail key[%s] p[%s] err:%+v", cache.PermIsPrivateKey(), sessionId, perm, err)
				return nil, err
			}
		}
	}
	roles, _, err := table.T_TCdpSysRoleService.QueryAll(l.ctx, sessionId, fmt.Sprintf("status:%d", table.RoleStatusEnable), nil, nil)
	if err != nil {
		l.Logger.Errorf("T_TCdpSysRoleService.QueryAll fail err:%+v", err)
		return nil, err
	}
	for _, role := range roles {
		if role.RoleIsAdmin() {
			err = l.svcCtx.Cache.SAdd(l.ctx, sessionId, cache.RoleAdminKey(), fmt.Sprint(role.Id))
			if err != nil {
				l.Logger.Errorf("Redis.SaddCtx fail key[%s] p[%s] err:%+v", cache.RoleAdminKey(), role.Id, err)
				return nil, err
			}
		}
		err = common.UpdateRolePermCache(l.ctx, l.svcCtx, &role)
		if err != nil {
			l.Logger.Errorf("common.UpdateRolePermCache fail err:%+v", err)
			return nil, err
		}
	}

	permConfigs, _, err := table.T_TCdpSysPermConfigService.QueryAll(l.ctx, sessionId, "", "", "")
	if err != nil {
		l.Logger.Errorf("T_TCdpSysPermConfigService.QueryAll fail err:%+v", err)
		return nil, err
	}
	for _, permConfig := range permConfigs {
		perm := helper.PermStandardBase64(permConfig.Perm)

		if err := l.svcCtx.Cache.Set(l.ctx, sessionId, cache.PermIsLoggingKey(perm), fmt.Sprint(permConfig.LoggingEnable)); err != nil {
			l.Logger.Errorf("Redis.SaddCtx fail key[%s] p[%s] err:%+v", cache.PermIsLoggingKey(perm), perm, err)
			return nil, err
		}

	}
	return
}
