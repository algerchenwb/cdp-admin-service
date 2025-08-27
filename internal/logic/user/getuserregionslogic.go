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

type GetUserRegionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserRegionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserRegionsLogic {
	return &GetUserRegionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserRegionsLogic) GetUserRegions() (resp *types.RegionsResp, err error) {
	// todo: add your logic here and delete this line
	resp = &types.RegionsResp{}
	sessionId := helper.GetSessionId(l.ctx)
	userId := helper.GetUserId(l.ctx)
	isAdmin := helper.GetIsAdmin(l.ctx)
	if isAdmin {
		areaRegions, bizRegions, err := getAllRegions(l.ctx, sessionId)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("[%s] getAllRegions err. err:%+v", sessionId, err)
			return nil, err
		}
		resp.AreaRegions = areaRegions
		resp.BizRegions = bizRegions
		return resp, nil
	}
	user, _, err := table.T_TCdpSysUserService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d", userId), nil, nil)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("table.TCdpSysUserService.Query err. Id[%d] err:%+v", userId, err)
		return nil, errorx.NewDefaultError(errorx.UserNotFoundErrorCode)
	}
	//  user.AreaRegions  user.BizRegions  user.AreaIds user.BizIds   area_id-> area_region_id  biz_id-> biz_region_id
	areaRegions := helper.StringToInt64Slice(user.AreaRegions)
	bizRegions := helper.StringToInt64Slice(user.BizRegions)
	areaIds := strings.Split(user.AreaIds, ",")
	bizIds := strings.Split(user.BizIds, ",")
	areas, _, err := table.T_TCdpAreaInfoService.QueryAll(l.ctx, sessionId, fmt.Sprintf("area_id__in:%s", strings.Join(areaIds, ",")), nil, nil)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("table.TCdpAreaInfoService.QueryAll err. err:%+v", err)
		return nil, errorx.NewDefaultError(errorx.UserNotFoundErrorCode)
	}
	bizs, _, err := table.T_TCdpBizInfoService.QueryAll(l.ctx, sessionId, fmt.Sprintf("id__in:%s", strings.Join(bizIds, ",")), nil, nil)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("table.T_BizInfoService.QueryAll err. err:%+v", err)
		return nil, errorx.NewDefaultError(errorx.UserNotFoundErrorCode)
	}
	uniqueAreaRegion := make(map[int64]struct{})
	uniqueBizRegion := make(map[int64]struct{})
	for _, area := range areas {
		uniqueAreaRegion[int64(area.RegionId)] = struct{}{}
	}
	for _, biz := range bizs {
		uniqueBizRegion[int64(biz.RegionId)] = struct{}{}
	}
	for _, areaRegion := range areaRegions {
		uniqueAreaRegion[areaRegion] = struct{}{}
	}
	for _, bizRegion := range bizRegions {
		uniqueBizRegion[bizRegion] = struct{}{}
	}
	resp.AreaRegions = make([]int64, 0, len(uniqueAreaRegion))
	for regionId := range uniqueAreaRegion {
		resp.AreaRegions = append(resp.AreaRegions, regionId)
	}
	resp.BizRegions = make([]int64, 0, len(uniqueBizRegion))
	for regionId := range uniqueBizRegion {
		resp.BizRegions = append(resp.BizRegions, regionId)
	}

	return
}

func getAllRegions(ctx context.Context, sessionId string) ([]int64, []int64, error) {
	areas, _, err := table.T_TCdpAreaInfoService.QueryAll(ctx, sessionId, "", nil, nil)
	if err != nil {
		logx.WithContext(ctx).Errorf("table.TCdpAreaInfoService.QueryAll err. err:%+v", err)
		return nil, nil, errorx.NewDefaultError(errorx.AreaNotFound)
	}
	bizs, _, err := table.T_TCdpBizInfoService.QueryAll(ctx, sessionId, fmt.Sprintf("status__ex:%d", table.BizStatusDeleted), nil, nil)
	if err != nil {
		logx.WithContext(ctx).Errorf("table.T_BizInfoService.QueryAll err. err:%+v", err)
		return nil, nil, errorx.NewDefaultError(errorx.BizNotFoundErrorCode)
	}
	uniqueAreaRegion := make(map[int64]struct{})
	uniqueBizRegion := make(map[int64]struct{})
	for _, area := range areas {
		uniqueAreaRegion[int64(area.RegionId)] = struct{}{}
	}
	for _, biz := range bizs {
		uniqueBizRegion[int64(biz.RegionId)] = struct{}{}
	}
	areaRegions := make([]int64, 0, len(uniqueAreaRegion))
	for regionId := range uniqueAreaRegion {
		areaRegions = append(areaRegions, regionId)
	}
	bizRegions := make([]int64, 0, len(uniqueBizRegion))
	for regionId := range uniqueBizRegion {
		bizRegions = append(bizRegions, regionId)
	}
	sort.Slice(areaRegions, func(i, j int) bool {
		return areaRegions[i] < areaRegions[j]
	})
	sort.Slice(bizRegions, func(i, j int) bool {
		return bizRegions[i] < bizRegions[j]
	})
	return areaRegions, bizRegions, nil
}
