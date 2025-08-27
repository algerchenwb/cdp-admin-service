package menu

import (
	"context"
	"fmt"

	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/logic/common"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoleMenuTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRoleMenuTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleMenuTreeLogic {
	return &GetRoleMenuTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRoleMenuTreeLogic) GetRoleMenuTree(req *types.RoleMenuTreeReq) (resp *types.RoleMenuTreeResp, err error) {
	sessionId := helper.GetSessionId(l.ctx)
	role, _, err := table.T_TCdpSysRoleService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d", req.RoleId), "", "")

	if err != nil {
		l.Logger.Errorf("sessionId:%s,GetRoleMenuTree,GetRoleById,err:%v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.QueryRoleFailedErrorCode)
	}

	if role == nil {
		l.Logger.Errorf("[%s] role not found. roleId[%d]", sessionId, req.RoleId)
		return nil, errorx.NewDefaultError(errorx.RoleIdErrorCode)
	}
	permConfigs, _, err := table.T_TCdpSysPermConfigService.QueryAll(l.ctx, sessionId, "", "id", "desc")
	if err != nil {
		l.Logger.Errorf("[%s] query perm configs failed. err:%+v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.QueryMenuFailedErrorCode)
	}
	permConfigsMap := make(map[string]table.TCdpSysPermConfig)
	for _, permConfig := range permConfigs {
		permConfigsMap[permConfig.Perm] = permConfig
	}

	resp = &types.RoleMenuTreeResp{}
	var menus []table.TCdpSysPermMenu
	isPrivate := []string{fmt.Sprintf("%d", table.IsPublic)}
	isAdmin := []string{fmt.Sprintf("%d", table.IsAdminNo)}
	if role.RoleIsAdmin() {
		isAdmin = append(isAdmin, fmt.Sprintf("%d", table.IsAdminYes))
	}
	var rootIds []string
	if role.Platform == table.PlatformConstruction {
		rootIds = []string{fmt.Sprintf("%d", l.svcCtx.Config.Menu.ShigongRootId)}
	} else if role.Platform == table.PlatformPower {
		rootIds = []string{fmt.Sprintf("%d", l.svcCtx.Config.Menu.SuanliRootId)}
	}

	menus, err = common.GetPlatformMenu(l.ctx, l.svcCtx, rootIds, isPrivate, isAdmin)
	if err != nil {
		l.Logger.Errorf("sessionId:%s,GetRoleMenuTree,GetAdminPlatformMenu,err:%v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.QuerySysPermMenuFailedErrorCode)
	}

	for _, menu := range menus {
		menuDto, err := common.MenuToDto(l.ctx, menu, permConfigsMap)
		if err != nil {
			l.Logger.Errorf("sessionId:%s,GetRoleMenuTree,MenuToDto,err:%v", sessionId, err)
			return nil, errorx.NewDefaultError(errorx.QuerySysPermMenuFailedErrorCode)
		}
		resp.Menus = append(resp.Menus, menuDto)
	}
	resp.Menus = common.FixMenu(resp.Menus)

	return resp, nil
}
