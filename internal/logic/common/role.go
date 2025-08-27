package common

import (
	"cdp-admin-service/internal/helper"
	"cdp-admin-service/internal/helper/cache"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"context"
	"fmt"
	"strings"

	table "cdp-admin-service/internal/helper/dal"

	"github.com/zeromicro/go-zero/core/logx"
)

func RoleIsAdmin(ctx context.Context, svcCtx *svc.ServiceContext, roleId int64) (bool, error) {
	sessionId := helper.GetSessionId(ctx)
	is, err := svcCtx.Cache.Sismember(ctx, sessionId, cache.RoleAdminKey(), fmt.Sprint(roleId))
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] 查询角色是否为管理员失败 err[%v]", sessionId, err)
		return false, errorx.NewDefaultError(errorx.QueryRoleFailedErrorCode)
	}
	return is, nil
}

func CheckRegionAccess(ctx context.Context, regionId int64) (bool, error) {
	sessionId := helper.GetSessionId(ctx)
	userId := helper.GetUserId(ctx)
	modifier, _, err := table.T_TCdpSysUserService.Query(ctx, sessionId, fmt.Sprintf("id:%d", userId), nil, nil)
	if err != nil {
		logx.WithContext(ctx).Errorf("table.TCdpSysUserService.Query err. err:%+v", err)
		return false, errorx.NewDefaultError(errorx.QueryUserFailedErrorCode)
	}
	areas, _, err := table.T_TCdpAreaInfoService.QueryAll(ctx, sessionId, fmt.Sprintf("area_id__in:%s$status__ex:%d", strings.Join(strings.Split(modifier.AreaIds, ","), ","), table.AreaStatusDisable), nil, nil)
	if err != nil {
		logx.WithContext(ctx).Errorf("table.TCdpAreaInfoService.QueryAll err. err:%+v", err)
		return false, errorx.NewDefaultError(errorx.AreaNotFound)
	}

	uniqueAreaRegion := make(map[string]struct{})

	for _, area := range areas {
		uniqueAreaRegion[fmt.Sprint(area.RegionId)] = struct{}{}
	}
	if modifier.AreaRegions != "" {
		for _, areaRegion := range strings.Split(modifier.AreaRegions, ",") {
			uniqueAreaRegion[areaRegion] = struct{}{}
		}
	}
	logx.WithContext(ctx).Infof("uniqueAreaRegion[%+v]", uniqueAreaRegion)

	_, exist := uniqueAreaRegion[fmt.Sprint(regionId)]
	return exist, nil
}

func GetRegions(ctx context.Context) (ans []string, err error) {
	sessionId := helper.GetSessionId(ctx)
	userId := helper.GetUserId(ctx)
	modifier, _, err := table.T_TCdpSysUserService.Query(ctx, sessionId, fmt.Sprintf("id:%d", userId), nil, nil)
	if err != nil {
		logx.WithContext(ctx).Errorf("table.TCdpSysUserService.Query err. err:%+v", err)
		return nil, errorx.NewDefaultError(errorx.QueryUserFailedErrorCode)
	}
	areas, _, err := table.T_TCdpAreaInfoService.QueryAll(ctx, sessionId, fmt.Sprintf("area_id__in:%s$status__ex:%d", strings.Join(strings.Split(modifier.AreaIds, ","), ","), table.AreaStatusDisable), nil, nil)
	if err != nil {
		logx.WithContext(ctx).Errorf("table.TCdpAreaInfoService.QueryAll err. err:%+v", err)
		return nil, errorx.NewDefaultError(errorx.AreaNotFound)
	}
	if modifier.AreaRegions != "" {
		ans = append(ans, strings.Split(modifier.AreaRegions, ",")...)
	}
	for _, area := range areas {
		ans = append(ans, fmt.Sprint(area.RegionId))
	}
	return ans, nil
}
