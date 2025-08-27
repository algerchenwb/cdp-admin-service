package common

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"cdp-admin-service/internal/helper"
	"cdp-admin-service/internal/helper/cache"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// 获取平台菜单
// rootIds 平台根ID
// isPrivate 是否 0-公共 1-私密    超管账号（0，1） 普通账号（0）
// isAdmin 是否超管独有 0-共有 1-超管独有    超管角色（0，1） 非超管角色（0）
func GetPlatformMenu(ctx context.Context, svcCtx *svc.ServiceContext, rootIds []string, isPrivate []string, isAdmin []string) ([]table.TCdpSysPermMenu, error) {
	ansMenus := make([]table.TCdpSysPermMenu, 0)
	sessionId := helper.GetSessionId(ctx)
	//  避免死循环
	total, _, _, err := table.T_TCdpSysPermMenuService.QueryPage(ctx, sessionId,
		fmt.Sprintf("status:%d$is_private__in:%s$is_admin__in:%s", table.MenuStatusEnable, strings.Join(isPrivate, ","), strings.Join(isAdmin, ",")),
		0, -1, nil, nil)
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] query menus failed. err:%+v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.QuerySysPermMenuFailedErrorCode)
	}
	// 加载根菜单
	menus, _, err := table.T_TCdpSysPermMenuService.QueryAll(ctx, sessionId,
		fmt.Sprintf("id__in:%s$status:%d$is_private__in:%s$is_admin__in:%s", strings.Join(rootIds, ","),
			table.MenuStatusEnable, strings.Join(isPrivate, ","), strings.Join(isAdmin, ",")), "id", "desc")
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] query menus failed. err:%+v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.QuerySysPermMenuFailedErrorCode)
	}
	for _, menu := range menus {
		ansMenus = append(ansMenus, menu)
	}

	for len(rootIds) > 0 && len(ansMenus) < total {
		menus, _, err := table.T_TCdpSysPermMenuService.QueryAll(ctx, sessionId,
			fmt.Sprintf("parent_id__in:%s$status:%d$is_private__in:%s$is_admin__in:%s", strings.Join(rootIds, ","),
				table.MenuStatusEnable, strings.Join(isPrivate, ","), strings.Join(isAdmin, ",")), "id", "desc")
		if err != nil {
			logx.WithContext(ctx).Errorf("[%s] query menus failed. err:%+v", sessionId, err)
			return nil, errorx.NewDefaultError(errorx.QuerySysPermMenuFailedErrorCode)
		}
		rootIds = make([]string, 0)
		for _, menu := range menus {
			ansMenus = append(ansMenus, menu)
			//  权限是最低级菜单，不需要继续遍历
			if menu.Type == table.MenuTypePermission {
				continue
			}
			rootIds = append(rootIds, fmt.Sprintf("%d", menu.Id))
		}
	}
	return ansMenus, nil
}

// 角色菜单
func GetRoleMenuPerms(ctx context.Context, svcCtx *svc.ServiceContext, roleId int64) ([]types.Menu, error) {
	sessionId := helper.GetSessionId(ctx)
	res := make([]types.Menu, 0)
	permConfigs, _, err := table.T_TCdpSysPermConfigService.QueryAll(ctx, sessionId, "", "id", "desc")
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] query perm configs failed. err:%+v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.QueryMenuFailedErrorCode)
	}
	permConfigsMap := make(map[string]table.TCdpSysPermConfig)
	for _, permConfig := range permConfigs {
		permConfigsMap[permConfig.Perm] = permConfig
	}
	role, _, err := table.T_TCdpSysRoleService.Query(ctx, sessionId, fmt.Sprintf("id:%d", roleId), nil, nil)
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] query role failed. err:%+v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.QueryRoleFailedErrorCode)
	}

	// 超管只有 算力平台测权限
	var menus []table.TCdpSysPermMenu
	if role.RoleIsAdmin() {
		menus, err = GetPlatformMenu(ctx, svcCtx,
			[]string{fmt.Sprint(svcCtx.Config.Menu.SuanliRootId)},
			[]string{fmt.Sprint(table.IsPublic)},
			[]string{fmt.Sprint(table.IsAdminNo), fmt.Sprint(table.IsAdminYes)})
		if err != nil {
			logx.WithContext(ctx).Errorf("[%s] query menus failed. err:%+v", sessionId, err)
			return nil, errorx.NewDefaultError(errorx.QueryMenuFailedErrorCode)
		}
	} else {
		menus, _, err = table.T_TCdpSysPermMenuService.QueryAll(ctx, sessionId, fmt.Sprintf("id__in:%s", strings.Join(strings.Split(role.PermMenuIds, ","), ",")), "id", "desc")
		if err != nil {
			logx.WithContext(ctx).Errorf("[%s] query menus failed. err:%+v", sessionId, err)
			return nil, errorx.NewDefaultError(errorx.QueryRoleFailedErrorCode)
		}
	}

	for _, menu := range menus {
		_menu, err := MenuToDto(ctx, menu, permConfigsMap)
		if err != nil {
			logx.WithContext(ctx).Errorf("[%s] menu to dto failed. err:%+v", sessionId, err)
			return nil, errorx.NewDefaultError(errorx.QueryMenuFailedErrorCode)
		}
		res = append(res, _menu)
	}
	return res, nil

}

// 获取超管账号可见菜单
func GetAdminAccountMenus(ctx context.Context, svcCtx *svc.ServiceContext) ([]types.Menu, error) {
	sessionId := helper.GetSessionId(ctx)
	menus, err := GetPlatformMenu(ctx, svcCtx,
		[]string{fmt.Sprint(svcCtx.Config.Menu.SuanliRootId)},
		[]string{fmt.Sprintf("%d", table.IsPublic), fmt.Sprintf("%d", table.IsPrivate)},
		[]string{fmt.Sprintf("%d", table.IsAdminNo), fmt.Sprintf("%d", table.IsAdminYes)})
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] query menus failed. err:%+v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.QueryMenuFailedErrorCode)
	}
	permConfigs, _, err := table.T_TCdpSysPermConfigService.QueryAll(ctx, sessionId, "", "id", "desc")
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] query perm configs failed. err:%+v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.QueryMenuFailedErrorCode)
	}
	permConfigsMap := make(map[string]table.TCdpSysPermConfig)
	for _, permConfig := range permConfigs {
		permConfigsMap[permConfig.Perm] = permConfig
	}
	ansMenus := make([]types.Menu, 0)
	for _, menu := range menus {
		_menu, err := MenuToDto(ctx, menu, permConfigsMap)
		if err != nil {
			logx.WithContext(ctx).Errorf("[%s] menu to dto failed. err:%+v", sessionId, err)
			return nil, errorx.NewDefaultError(errorx.QueryMenuFailedErrorCode)
		}
		ansMenus = append(ansMenus, _menu)
	}
	return ansMenus, nil
}

// 获取菜单列表全部权限
func GetMenuAllPerms(ctx context.Context, menu table.TCdpSysPermMenu, permConfigs map[string]table.TCdpSysPermConfig) ([]types.Perm, error) {
	perms := make([]types.Perm, 0)
	sessionId := helper.GetSessionId(ctx)
	_perms := make([]string, 0)
	err := json.Unmarshal([]byte(menu.Perms), &_perms)
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] unmarshal perms failed. err:%+v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.QueryMenuFailedErrorCode)
	}
	for _, perm := range _perms {
		permConfig, ok := permConfigs[perm]
		if !ok {
			continue
		}
		perms = append(perms, types.Perm{
			Perm:          permConfig.Perm,
			LoggingEnable: int64(permConfig.LoggingEnable),
		})
	}
	return perms, nil
}

func MenuToDto(ctx context.Context, menu table.TCdpSysPermMenu, permConfigsMap map[string]table.TCdpSysPermConfig) (ans types.Menu, err error) {
	sessionId := helper.GetSessionId(ctx)
	ans = types.Menu{
		Id:           int64(menu.Id),
		ParentId:     int64(menu.ParentId),
		Name:         menu.Name,
		Router:       menu.Router,
		Type:         int64(menu.Type),
		Icon:         menu.Icon,
		OrderNum:     int64(menu.OrderNum),
		ViewPath:     menu.ViewPath,
		SystemHost:   menu.SystemHost,
		IsShow:       int64(menu.IsShow),
		IsPrivate:    int64(menu.IsPrivate),
		IsAdmin:      int64(menu.IsAdmin),
		ActiveRouter: menu.ActiveRouter,
		CreateBy:     menu.CreateBy,
		CreateTime:   menu.CreateTime.Format(types.DefaultTimeFormat),
		UpdateBy:     menu.UpdateBy,
		UpdateTime:   menu.UpdateTime.Format(types.DefaultTimeFormat),
		ModifyTime:   menu.ModifyTime.Format(types.DefaultTimeFormat),
	}
	perms, err := GetMenuAllPerms(ctx, menu, permConfigsMap)
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] get menu all perms failed. err:%+v", sessionId, err)
		return ans, errorx.NewDefaultError(errorx.QueryMenuFailedErrorCode)
	}
	ans.Perms = perms
	return ans, nil
}

func UpdateRolePermCache(ctx context.Context, svcCtx *svc.ServiceContext, role *table.TCdpSysRole) error {
	sessionId := helper.GetSessionId(ctx)
	users, _, err := table.T_TCdpSysUserService.QueryAll(ctx, sessionId, fmt.Sprintf("role_id:%d$status:%d", role.Id, table.SysUserStatusEnable), nil, nil)
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] 查询用户失败 err[%v]", sessionId, err)
		return errorx.NewDefaultError(errorx.QueryUserFailedErrorCode)
	}
	var menus []table.TCdpSysPermMenu

	menus, _, err = table.T_TCdpSysPermMenuService.QueryAll(ctx, sessionId, fmt.Sprintf("id__in:%s$status:%d", role.PermMenuIds, table.MenuStatusEnable), nil, nil)
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] 查询菜单失败 err[%v]", sessionId, err)
		return errorx.NewDefaultError(errorx.QueryMenuFailedErrorCode)
	}

	rolePerms := make(map[string]struct{})
	for _, menu := range menus {
		perm, err := ParseMenuPerm(&menu)
		if err != nil {
			logx.WithContext(ctx).Errorf("[%s] 解析菜单权限失败 err[%v]", sessionId, err)
			return errorx.NewDefaultError(errorx.ParseMenuPermFailedErrorCode)
		}
		for _, p := range perm {
			rolePerms[helper.PermStandardBase64(p)] = struct{}{}
		}
	}

	for _, user := range users {
		svcCtx.Cache.Del(ctx, sessionId, cache.UserPermKey(int64(user.Id)))
		for p := range rolePerms {
			err = svcCtx.Cache.SAdd(ctx, sessionId, cache.UserPermKey(int64(user.Id)), p)
			if err != nil {
				logx.WithContext(ctx).Errorf("[%s] 添加用户权限失败 err[%v]", sessionId, err)
				return errorx.NewDefaultError(errorx.UpdateUserPermFailedErrorCode)
			}
		}
	}
	return nil
}

func UpdateUserPermCache(ctx context.Context, user *table.TCdpSysUser, cacheService *cache.Cache) error {
	sessionId := helper.GetSessionId(ctx)
	role, _, err := table.T_TCdpSysRoleService.Query(ctx, sessionId, fmt.Sprintf("id:%d", user.RoleId), nil, nil)
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] 查询角色失败 err[%v]", sessionId, err)
		return errorx.NewDefaultError(errorx.QueryRoleFailedErrorCode)
	}
	menus, _, err := table.T_TCdpSysPermMenuService.QueryAll(ctx, sessionId, fmt.Sprintf("id__in:%s", role.PermMenuIds), nil, nil)
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] 查询菜单失败 err[%v]", sessionId, err)
		return errorx.NewDefaultError(errorx.QueryMenuFailedErrorCode)
	}

	rolePerms := make(map[string]struct{})
	for _, menu := range menus {
		perm, err := ParseMenuPerm(&menu)
		if err != nil {
			logx.WithContext(ctx).Errorf("[%s] 解析菜单权限失败 err[%v]", sessionId, err)
			return errorx.NewDefaultError(errorx.ParseMenuPermFailedErrorCode)
		}
		for _, p := range perm {
			rolePerms[base64.StdEncoding.EncodeToString([]byte(p))] = struct{}{}
		}
	}

	cacheService.Del(ctx, sessionId, cache.UserPermKey(int64(user.Id)))
	for p := range rolePerms {
		err = cacheService.SAdd(ctx, sessionId, cache.UserPermKey(int64(user.Id)), p)
		if err != nil {
			logx.WithContext(ctx).Errorf("[%s] 添加用户权限失败 err[%v]", sessionId, err)
			return errorx.NewDefaultError(errorx.UpdateUserPermFailedErrorCode)
		}
		logx.WithContext(ctx).Infof("[%s] 添加用户权限成功. user[%v] perm[%s]", sessionId, user, p)
	}
	return nil
}

// 修复不合理菜单，上级菜单不存在，则移除该菜单
func FixMenu(menus []types.Menu) (ansMenus []types.Menu) {
	ansMenus = make([]types.Menu, 0)

	menuMap := make(map[int64]struct{})
	for _, menu := range menus {
		menuMap[menu.Id] = struct{}{}
	}
	for _, menu := range menus {
		if menu.ParentId == 0 {
			ansMenus = append(ansMenus, menu)
			continue
		}
		if _, ok := menuMap[menu.ParentId]; !ok {
			continue
		}
		ansMenus = append(ansMenus, menu)
	}
	return
}

func ParseMenuPerm(menu *table.TCdpSysPermMenu) (perms []string, err error) {
	perms = make([]string, 0)
	err = json.Unmarshal([]byte(menu.Perms), &perms)
	if err != nil {
		return nil, err
	}
	return
}
