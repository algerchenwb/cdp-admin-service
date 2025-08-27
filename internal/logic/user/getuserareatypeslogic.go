package user

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserAreaTypesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserAreaTypesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserAreaTypesLogic {
	return &GetUserAreaTypesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserAreaTypesLogic) GetUserAreaTypes(req *types.AreaInfosReq) (resp *types.AreaInfosResp, err error) {

	resp = &types.AreaInfosResp{}
	sessionId := helper.GetSessionId(l.ctx)
	userId := helper.GetUserId(l.ctx)
	isAdmin := helper.GetIsAdmin(l.ctx)
	if isAdmin {
		areas, _, err := table.T_TCdpAreaInfoService.QueryAll(l.ctx, sessionId, "", nil, nil)
		if err != nil {
			l.Logger.Errorf("[%s] table.T_AreaInfoService.QueryAll err. err:%+v", sessionId, err)
			return nil, errorx.NewDefaultError(errorx.AreaNotFound)
		}
		for _, area := range areas {
			resp.AreaInfos = append(resp.AreaInfos, types.AreaInfo{
				PrimaryId:      area.PrimaryId,
				AgentId:        area.AgentId,
				AreaId:         area.AreaId,
				Name:           area.Name,
				RegionId:       int64(area.RegionId),
				DeploymentType: int64(area.DeploymentType),
				Remark:         area.Remark,
				ProxyAddr:      area.ProxyAddr,
			})
		}
		return resp, nil
	}
	user, _, err := table.T_TCdpSysUserService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d", userId), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] table.T_SysUserService.Query err. Id[%d] err:%+v", sessionId, userId, err)
		return nil, errorx.NewDefaultError(errorx.UserNotFoundErrorCode)
	}
	var uniqueArea map[int32]types.AreaInfo = make(map[int32]types.AreaInfo)
	userAreaIds := strings.Split(user.AreaIds, ",")
	// 避免末尾有逗号查询失败
	userAreaRegions := strings.Split(user.AreaRegions, ",")

	areas, _, err := table.T_TCdpAreaInfoService.QueryAll(l.ctx, sessionId, fmt.Sprintf("area_id__in:%s", strings.Join(userAreaIds, ",")), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] table.T_AreaInfoService.QueryAll err. err:%+v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.UserNotFoundErrorCode)
	}
	for _, area := range areas {
		uniqueArea[area.AreaId] = types.AreaInfo{
			PrimaryId:      area.PrimaryId,
			AgentId:        area.AgentId,
			AreaId:         area.AreaId,
			Name:           area.Name,
			RegionId:       int64(area.RegionId),
			DeploymentType: int64(area.DeploymentType),
			Remark:         area.Remark,
			ProxyAddr:      area.ProxyAddr,
		}
	}
	areas, _, err = table.T_TCdpAreaInfoService.QueryAll(l.ctx, sessionId, fmt.Sprintf("region_id__in:%s", strings.Join(userAreaRegions, ",")), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] table.T_AreaInfoService.QueryAll err. err:%+v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.UserNotFoundErrorCode)
	}
	for _, area := range areas {
		uniqueArea[area.AreaId] = types.AreaInfo{
			PrimaryId:      area.PrimaryId,
			AgentId:        area.AgentId,
			AreaId:         area.AreaId,
			Name:           area.Name,
			RegionId:       int64(area.RegionId),
			DeploymentType: int64(area.DeploymentType),
			Remark:         area.Remark,
			ProxyAddr:      area.ProxyAddr,
		}
	}

	for _, area := range uniqueArea {
		resp.AreaInfos = append(resp.AreaInfos, area)
	}
	sort.Slice(resp.AreaInfos, func(i, j int) bool {
		return resp.AreaInfos[i].AreaId < resp.AreaInfos[j].AreaId
	})

	return
}
