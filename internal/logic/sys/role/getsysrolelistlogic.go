package role

import (
	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/logic/common"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSysRoleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSysRoleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSysRoleListLogic {
	return &GetSysRoleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSysRoleListLogic) GetSysRoleList(req *types.CommonPageRequest) (resp *types.RoleListResp, err error) {
	sessionId := helper.GetSessionId(l.ctx)
	isAdmin := helper.GetIsAdmin(l.ctx)

	if !isAdmin {
		regions, err := common.GetRegions(l.ctx)
		if err != nil {
			l.Logger.Error("[%s] 获取区域失败 err[%v]", sessionId, err)
			return nil, errorx.NewDefaultError(errorx.QuerySysRoleFailedErrorCode)
		}
		req.CondList = append(req.CondList, fmt.Sprintf("region_id__in:%s", strings.Join(regions, ",")))
	}
	req.CondList = append(req.CondList, fmt.Sprintf("status:%d", table.RoleStatusEnable))
	qry, err := helper.CheckCommQueryParam(l.ctx, sessionId, req.CondList, req.Orders, req.Sorts, table.TCdpSysRole{})
	if err != nil {
		l.Logger.Errorf("[%s] CheckCommQueryParam failed, req:%+v err:%+v", sessionId, req, err)
		return nil, errorx.NewDefaultError(errorx.ParamErrorCode)
	}
	total, sysRoleList, _, err := table.T_TCdpSysRoleService.QueryPage(l.ctx, sessionId, qry, req.Offset, req.Limit, req.Sorts, req.Orders)
	if err != nil {
		l.Logger.Error("[%s] 查询角色失败 err[%v]", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.QuerySysRoleFailedErrorCode)
	}
	resp = &types.RoleListResp{
		Total: int64(total),
	}
	for _, sysRole := range sysRoleList {
		resp.List = append(resp.List, types.Role{
			Id:         int64(sysRole.Id),
			Name:       sysRole.Name,
			Remark:     sysRole.Remark,
			RegionId:   int64(sysRole.RegionId),
			Platform:   sysRole.Platform,
			IsAdmin:    sysRole.IsAdmin,
			CreateBy:   sysRole.CreateBy,
			UpdateBy:   sysRole.UpdateBy,
			CreateTime: sysRole.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateTime: sysRole.UpdateTime.Format("2006-01-02 15:04:05"),
		})
	}
	return
}
