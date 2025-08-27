package menu

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

type AddMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddMenuLogic {
	return &AddMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddMenuLogic) AddMenu(req *types.AddMenuReq) (resp *types.AddMenuResp, err error) {
	sessionId := helper.GetSessionId(l.ctx)
	userName := helper.GetUserName(l.ctx)

	if req.Menu.ParentId != 0 {
		parentPermMenu, _, err := table.T_TCdpSysPermMenuService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d$status:%d", req.Menu.ParentId, table.MenuStatusEnable), nil, nil)
		if err != nil {
			l.Logger.Errorf("[%s] 查询父级菜单失败 id[%d] err[%v]", sessionId, req.Menu.ParentId, err)
			return nil, errorx.NewDefaultError(errorx.ParentPermMenuIdErrorCode)
		}

		if parentPermMenu.Type == table.MenuTypePermission {
			l.Logger.Errorf("[%s] 父级菜单类型错误 id[%d] type[%d]", sessionId, req.Menu.ParentId, parentPermMenu.Type)
			return nil, errorx.NewDefaultError(errorx.SetParentTypeErrorCode)
		}
	}

	if req.Menu.Type != table.MenuTypePermission && len(req.Menu.Perms) > 0 {
		l.Logger.Errorf("[%s] 菜单类型错误  type[%d]", sessionId, req.Menu.Type)
		return nil, errorx.NewDefaultError(errorx.NotPermissionWithPermRefuseErrorCode)
	}

	menu := &table.TCdpSysPermMenu{
		ParentId:     uint32(req.Menu.ParentId),
		Name:         req.Menu.Name,
		Router:       req.Menu.Router,
		Type:         uint32(req.Menu.Type),
		Icon:         req.Menu.Icon,
		OrderNum:     uint32(req.Menu.OrderNum),
		ViewPath:     req.Menu.ViewPath,
		IsShow:       uint32(req.Menu.IsShow),
		ActiveRouter: req.Menu.ActiveRouter,
		SystemHost:   req.Menu.SystemHost,
		IsPrivate:    int32(req.Menu.IsPrivate),
		IsAdmin:      int32(req.Menu.IsAdmin),
		Perms:        PermsToJson(req.Menu.Perms),
		Status:       table.MenuStatusEnable,
		CreateBy:     userName,
		UpdateBy:     userName,
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
		ModifyTime:   time.Now(),
	}
	menu, _, err = table.T_TCdpSysPermMenuService.Insert(l.ctx, sessionId, menu)
	if err != nil {
		l.Logger.Errorf("[%s] 插入菜单失败 menu[%s] err[%v]", sessionId, helper.ToJSON(menu), err)
		return nil, errorx.NewDefaultError(errorx.InsertMenuErrorCode)
	}
	l.Logger.Infof("[%s] 插入菜单成功 menu[%+v]", sessionId, menu)

	if menu.IsPrivate == table.IsPrivate {
		for _, perm := range req.Menu.Perms {
			err = l.svcCtx.Cache.SAdd(l.ctx, sessionId, cache.PermIsPrivateKey(), helper.PermStandardBase64(perm.Perm))
			if err != nil {
				l.Logger.Errorf("[%s] 更新内网菜单失败 err[%v]", sessionId, err)
				return nil, errorx.NewDefaultError(errorx.UpdatePrivateIpMenuErrorCode)
			}
		}
	}

	// 插入菜单配置
	UpdateOrInsertMenuConfig(l.ctx, l.svcCtx, userName, req.Menu.Perms)

	return &types.AddMenuResp{
		Id: int64(menu.Id),
	}, nil
}
