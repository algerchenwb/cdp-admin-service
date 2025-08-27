package cache

import (
	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/model/globalkey"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"cdp-admin-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CacheSyncLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCacheSyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CacheSyncLogic {
	return &CacheSyncLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CacheSyncLogic) CacheSync() (err error) {
	sessionId := helper.GetSessionId(l.ctx)
	dictionaries, _, err := table.T_TCdpSysDictionaryService.QueryAll(l.ctx, sessionId, "", nil, nil)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("T_TCdpSysDictionaryService.QueryAll fail err:%+v", err)
	}
	for _, dictionary := range dictionaries {
		key := fmt.Sprintf("%s:%s", globalkey.SysConfigUniqueCachePrefix, dictionary.UniqueKey)
		if err := l.svcCtx.Cache.Set(l.ctx, sessionId, key, dictionary.Value); err != nil {
			logx.WithContext(l.ctx).Errorf("Redis.Set fail key[%s] value[%s] err:%+v", key, dictionary.Value, err)
		}
	}

	sysPermMenus, _, err := table.T_TCdpSysPermMenuService.QueryAll(l.ctx, sessionId, "", nil, nil)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("T_TCdpSysPermMenuService.QueryAll fail err:%+v", err)
	}
	permSystemHost := make(map[string]string)
	sysPermMenuMap := make(map[string]table.TCdpSysPermMenu)
	for _, sysPermMenu := range sysPermMenus {
		sysPermMenuMap[fmt.Sprint(sysPermMenu.Id)] = sysPermMenu
		var _perms []string
		json.Unmarshal([]byte(sysPermMenu.Perms), &_perms)
		for _, _p := range _perms {
			permSystemHost[_p] = sysPermMenu.SystemHost
		}
	}
	for p, systemHost := range permSystemHost {
		p = helper.PermStandardBase64(p)
		key := fmt.Sprintf("%s:%s", globalkey.SysSystemHostCachePrefix, helper.MD5(p))
		if err := l.svcCtx.Cache.Set(l.ctx, sessionId, key, systemHost); err != nil {
			logx.WithContext(l.ctx).Errorf("Redis.Setex fail key[%s] p[%s] err:%+v", key, p, err)
		}
	}

	// // 同步菜单
	roleMap := make(map[int32][]string)
	roles, _, err := table.T_TCdpSysRoleService.QueryAll(l.ctx, sessionId, "", nil, nil)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("T_TCdpSysRoleService.QueryAll fail err:%+v", err)
	}

	for _, role := range roles {
		var menuIds []string
		if role.RoleIsAdmin() {
			menus, _, err := table.T_TCdpSysPermMenuService.QueryAll(l.ctx, sessionId, "", nil, nil)
			if err != nil {
				logx.WithContext(l.ctx).Errorf("T_TCdpSysPermMenuService.QueryAll fail err:%+v", err)
			}
			for _, menu := range menus {
				menuIds = append(menuIds, fmt.Sprint(menu.Id))
			}
		} else {
			menuIds = strings.Split(role.PermMenuIds, ",")
		}
		for _, menuId := range menuIds {
			permMenu, ok := sysPermMenuMap[menuId]
			if !ok {
				continue
			}
			var perms []string
			json.Unmarshal([]byte(permMenu.Perms), &perms)
			if len(perms) == 0 {
				continue
			}
			roleMap[int32(role.Id)] = append(roleMap[int32(role.Id)], perms...)
		}
	}

	sysUsers, _, err := table.T_TCdpSysUserService.QueryAll(l.ctx, sessionId, "", nil, nil)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("T_TCdpSysUserService.QueryAll fail err:%+v", err)
	}
	for _, sysUser := range sysUsers {
		_, err := l.svcCtx.Cache.Del(l.ctx, sessionId, globalkey.UserPermPrefix+fmt.Sprint(sysUser.Id))
		if err != nil {
			logx.WithContext(l.ctx).Errorf("Redis.DelCtx fail key[%s] err:%+v", globalkey.UserPermPrefix+fmt.Sprint(sysUser.Id), err)
			continue
		}
		if sysUser.RoleId == 0 {
			continue
		}
		for _, perm := range roleMap[sysUser.RoleId] {
			perm = base64.StdEncoding.EncodeToString([]byte(perm))
			if err := l.svcCtx.Cache.SAdd(l.ctx, sessionId, fmt.Sprintf("%s%d", globalkey.UserPermPrefix, sysUser.Id), perm); err != nil {
				logx.WithContext(l.ctx).Errorf("Redis.SaddCtx fail key[%s] roleIds[%+v] err:%+v", globalkey.UserPermPrefix+fmt.Sprint(sysUser.Id), roleMap[sysUser.RoleId], err)
			}
		}

	}

	return
}
