package menu

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"cdp-admin-service/internal/helper"
	"cdp-admin-service/internal/helper/cache"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/logic/common"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"cdp-admin-service/internal/model/errorx"

	"github.com/zeromicro/go-zero/core/logx"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
)

type UpdateMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMenuLogic {
	return &UpdateMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMenuLogic) UpdateMenu(req *types.UpdateMenuReq) (resp *types.UpdateMenuResp, err error) {
	sessionId := helper.GetSessionId(l.ctx)
	userName := helper.GetUserName(l.ctx)

	menu, _, err := table.T_TCdpSysPermMenuService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d", req.Menu.Id), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] 查询菜单失败 err[%v]", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.QueryMenuFailedErrorCode)
	}

	if menu.IsPrivate != int32(req.Menu.IsPrivate) {
		for _, perm := range req.Menu.Perms {
			if menu.IsPrivate == table.IsPrivate {
				err = l.svcCtx.Cache.SRem(l.ctx, sessionId, cache.PermIsPrivateKey(), helper.PermStandardBase64(perm.Perm))
			} else {
				err = l.svcCtx.Cache.SAdd(l.ctx, sessionId, cache.PermIsPrivateKey(), helper.PermStandardBase64(perm.Perm))
			}
			if err != nil {
				l.Logger.Errorf("[%s] 更新内网菜单失败 err[%v]", sessionId, err)
			}
		}
	}

	roles, _, err := table.T_TCdpSysRoleService.QueryAll(l.ctx, sessionId, fmt.Sprintf("perm_menu_ids__contains:%d$status:%d", req.Menu.Id, table.RoleStatusEnable), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] 查询角色失败 err[%v]", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.QueryRoleFailedErrorCode)
	}
	for _, role := range roles {
		err = common.UpdateRolePermCache(l.ctx, l.svcCtx, &role)
		if err != nil {
			l.Logger.Errorf("[%s] 更新角色权限缓存失败 err[%v]", sessionId, err)
			return nil, errorx.NewDefaultError(errorx.UpdateRoleFailedErrorCode)
		}
	}

	menu.Name = req.Menu.Name
	menu.Router = req.Menu.Router
	menu.Type = uint32(req.Menu.Type)
	menu.Icon = req.Menu.Icon
	menu.OrderNum = uint32(req.Menu.OrderNum)
	menu.ViewPath = req.Menu.ViewPath
	menu.IsShow = uint32(req.Menu.IsShow)
	menu.ActiveRouter = req.Menu.ActiveRouter
	menu.SystemHost = req.Menu.SystemHost
	menu.IsPrivate = int32(req.Menu.IsPrivate)
	menu.IsAdmin = int32(req.Menu.IsAdmin)
	menu.Perms = PermsToJson(req.Menu.Perms)
	menu.UpdateTime = time.Now()
	menu.UpdateBy = userName

	menu, _, err = table.T_TCdpSysPermMenuService.Update(l.ctx, sessionId, menu.Id, menu)
	if err != nil {
		l.Logger.Errorf("[%s] 更新菜单失败 err[%v]", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.UpdateMenuErrorCode)
	}
	l.Logger.Infof("[%s] 更新菜单成功 menu[%+v]", sessionId, menu)

	// 更新菜单配置
	UpdateOrInsertMenuConfig(l.ctx, l.svcCtx, userName, req.Menu.Perms)

	return &types.UpdateMenuResp{
		Id: int64(menu.Id),
	}, nil
}

func PermsToJson(perms []types.Perm) string {
	permStrs := make([]string, 0)
	for _, p := range perms {
		permStrs = append(permStrs, helper.PermStandard(p.Perm))
	}
	permsJson, _ := json.Marshal(permStrs)
	return string(permsJson)
}

// 更新或插入菜单配置
func UpdateOrInsertMenuConfig(ctx context.Context, svcCtx *svc.ServiceContext, updateBy string, perms []types.Perm) {
	sessionId := helper.GetSessionId(ctx)
	for _, p := range perms {
		permConfig, _, err := table.T_TCdpSysPermConfigService.Query(ctx, sessionId, fmt.Sprintf("perm:%s", p.Perm), nil, nil)
		//  菜单配置不存在
		if errors.Is(err, gopublic.ErrNotExist) {
			permConfig = &table.TCdpSysPermConfig{
				Perm:          helper.PermStandard(p.Perm),
				LoggingEnable: uint32(p.LoggingEnable),
				CreateBy:      updateBy,
				CreateTime:    time.Now(),
				UpdateBy:      updateBy,
				UpdateTime:    time.Now(),
				ModifyTime:    time.Now(),
			}
			_, _, err = table.T_TCdpSysPermConfigService.Insert(ctx, sessionId, permConfig)
			if err != nil {
				logx.WithContext(ctx).Errorf("[%s] 插入菜单配置失败 perm[%s] err[%v]", sessionId, helper.ToJSON(permConfig), err)
				continue
			}
			logx.WithContext(ctx).Infof("[%s] 插入菜单配置成功 permConfig[%+v]", sessionId, permConfig)
		} else if err != nil {
			logx.WithContext(ctx).Errorf("[%s] 查询菜单配置失败 perm[%s] err[%v]", sessionId, helper.ToJSON(permConfig), err)
			continue
		} else {
			// 存在更新菜单配置
			permConfig, _, err = table.T_TCdpSysPermConfigService.Update(ctx, sessionId, permConfig.Id, map[string]interface{}{
				"logging_enable": uint32(p.LoggingEnable),
				"update_time":    time.Now(),
				"update_by":      updateBy,
			})
			if err != nil {
				logx.WithContext(ctx).Errorf("[%s] 更新菜单配置失败 err[%v]", sessionId, err)
				continue
			}
			logx.WithContext(ctx).Infof("[%s] 更新菜单配置成功 permConfig[%+v]", sessionId, permConfig)
		}

		err = svcCtx.Cache.Set(ctx, sessionId, cache.PermIsLoggingKey(helper.PermStandardBase64(p.Perm)), fmt.Sprintf("%d", p.LoggingEnable))
		if err != nil {
			logx.WithContext(ctx).Errorf("[%s] 配置菜单配置缓存失败 perm[%s] err[%v]", sessionId, helper.ToJSON(p), err)
			continue
		}
	}
}
