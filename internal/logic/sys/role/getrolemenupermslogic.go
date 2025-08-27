package role

import (
	"context"

	"cdp-admin-service/internal/helper"
	"cdp-admin-service/internal/logic/common"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoleMenuPermsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRoleMenuPermsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleMenuPermsLogic {
	return &GetRoleMenuPermsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRoleMenuPermsLogic) GetRoleMenuPerms(req *types.RoleMenuPermsReq) (resp *types.RoleMenuPermsResp, err error) {
	resp = &types.RoleMenuPermsResp{}
	resp.Menus, err = common.GetRoleMenuPerms(l.ctx, l.svcCtx, req.Id)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("[%s] get role menu perms failed. err:%+v", helper.GetSessionId(l.ctx), err)
		return nil, err
	}
	resp.Menus = common.FixMenu(resp.Menus)
	return
}
