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

type GetMenuListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMenuListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuListLogic {
	return &GetMenuListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMenuListLogic) GetMenuList(req *types.CommonPageRequest) (resp *types.MenuListResp, err error) {
	sessionId := helper.GetSessionId(l.ctx)
	resp = &types.MenuListResp{}
	req.CondList = append(req.CondList, fmt.Sprintf("status:%d", table.MenuStatusEnable))

	qry, err := helper.CheckCommQueryParam(l.ctx, sessionId, req.CondList, req.Orders, req.Sorts, table.TCdpSysPermMenu{})
	if err != nil {
		l.Logger.Errorf("[%s] CheckCommQueryParam failed, req:%+v err:%+v", sessionId, req, err)
		return nil, errorx.NewDefaultError(errorx.ParamErrorCode)
	}
	total, menus, _, err := table.T_TCdpSysPermMenuService.QueryPage(l.ctx, sessionId, qry, req.Offset, req.Limit, req.Orders, req.Sorts)
	if err != nil {
		l.Logger.Error("查询权限菜单失败", err)
		return nil, errorx.NewDefaultError(errorx.QuerySysPermMenuFailedErrorCode)
	}
	permConfigs, _, err := table.T_TCdpSysPermConfigService.QueryAll(l.ctx, sessionId, "", "", "")
	if err != nil {
		l.Logger.Errorf("[%s] 查询权限配置失败 err[%v]", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.QuerySysPermMenuFailedErrorCode)
	}
	permConfigsMap := make(map[string]table.TCdpSysPermConfig)
	for _, permConfig := range permConfigs {
		permConfigsMap[permConfig.Perm] = permConfig
	}

	resp.Total = int64(total)
	for _, menu := range menus {
		menuDto, err := common.MenuToDto(l.ctx, menu, permConfigsMap)
		if err != nil {
			l.Logger.Errorf("[%s] 转换菜单失败 err[%v]", sessionId, err)
			return nil, errorx.NewDefaultError(errorx.QuerySysPermMenuFailedErrorCode)
		}
		resp.List = append(resp.List, menuDto)
	}
	resp.List = common.FixMenu(resp.List)
	return resp, nil
}
