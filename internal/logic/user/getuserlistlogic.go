package user

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserListLogic {
	return &GetUserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserListLogic) GetUserList(req *types.CommonPageRequest) (resp *types.UserInfoListResp, err error) {
	sessionId := helper.GetSessionId(l.ctx)
	isAdmin := helper.GetIsAdmin(l.ctx)

	req.CondList = append(req.CondList, "status__ex:0")
	qry, err := helper.CheckCommQueryParam(l.ctx, sessionId, req.CondList, req.Orders, req.Sorts, table.TCdpSysUser{})
	if err != nil {
		l.Logger.Errorf("[%s] CheckCommQueryParam failed, req:%+v err:%+v", sessionId, req, err)
		return nil, errorx.NewDefaultError(errorx.ParamErrorCode)
	}
	users, _, err := table.T_TCdpSysUserService.QueryAll(l.ctx, sessionId, qry, req.Sorts, req.Orders)
	if err != nil {
		l.Logger.Errorf("[%s] query users failed, err: %v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.QueryUserFailedErrorCode)
	}

	modifier, _, err := table.T_TCdpSysUserService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d", helper.GetUserId(l.ctx)), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] query modifier failed, err: %v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.QueryUserFailedErrorCode)
	}
	modifierAreaIds, err := getUserAreaIds(l.ctx, modifier)
	if err != nil {
		l.Logger.Errorf("[%s] getModifyUserAreaIds failed, modifier[%v] err: %v", sessionId, modifier, err)
		return nil, errorx.NewDefaultError(errorx.QueryUserFailedErrorCode)
	}

	// 有权限用户列表
	accessUsers := make([]table.TCdpSysUser, 0)
	for _, user := range users {
		// 能加载全部账号（除超级管理员账号）
		if !isAdmin && user.IsAdmin == table.IsAdminYes {
			l.Logger.Infof("[%s] modifier[%v] user[%v] is admin, skip", sessionId, modifier, user)
			continue
		}
		// 自身或者超管账号不校验节点
		if user.Id == modifier.Id || isAdmin {
			l.Logger.Infof("[%s] modifier[%v] user[%v] is self or admin, access", sessionId, modifier, user)
			accessUsers = append(accessUsers, user)
			continue
		}

		waitUserAreaIds, err := getUserAreaIds(l.ctx, &user)
		if err != nil {
			l.Logger.Errorf("[%s] getUserAreaIds failed, err: %v", sessionId, err)
			return nil, errorx.NewDefaultError(errorx.QueryUserFailedErrorCode)
		}
		access := checkAreaIdsAccessV2(modifierAreaIds, waitUserAreaIds)
		if access {
			accessUsers = append(accessUsers, user)
		}
		l.Logger.Debugf("[%s] no access, modifierAreaIds[%+v] waitUserAreaIds[%+v]", sessionId, modifierAreaIds, waitUserAreaIds)
	}
	resp = &types.UserInfoListResp{
		Total: int64(len(accessUsers)),
	}
	for offset := req.Offset; offset < req.Offset+req.Limit; offset++ {
		if offset >= len(accessUsers) {
			break
		}
		resp.List = append(resp.List, types.UserInfoResp{
			Id:          int64(accessUsers[offset].Id),
			Account:     accessUsers[offset].Account,
			Nickname:    accessUsers[offset].Nickname,
			Avatar:      accessUsers[offset].Avatar,
			Mobile:      accessUsers[offset].Mobile,
			RoleId:      int64(accessUsers[offset].RoleId),
			AreaIds:     helper.StringToInt64Slice(accessUsers[offset].AreaIds),
			AreaRegions: helper.StringToInt64Slice(accessUsers[offset].AreaRegions),
			BizIds:      helper.StringToInt64Slice(accessUsers[offset].BizIds),
			BizRegions:  helper.StringToInt64Slice(accessUsers[offset].BizRegions),
			Platform:    int64(accessUsers[offset].Platform),
			Status:      int64(accessUsers[offset].Status),
			Remark:      accessUsers[offset].Remark,
		})
	}

	return
}

func checkAreaIdsAccessV2(modifierAreaIds map[string]struct{}, waitUpdateAreaIds map[string]struct{}) bool {
	for areaId := range waitUpdateAreaIds {
		if _, ok := modifierAreaIds[areaId]; !ok {
			return false
		}
	}
	return true
}

func getUserAreaIds(ctx context.Context, user *table.TCdpSysUser) (map[string]struct{}, error) {
	sessionId := helper.GetSessionId(ctx)
	var areaIdMap = make(map[string]struct{})
	for _, areaId := range strings.Split(user.AreaIds, ",") {
		areaIdMap[areaId] = struct{}{}
	}
	regionAreaIds, _, err := table.T_TCdpAreaInfoService.QueryAll(ctx, sessionId, fmt.Sprintf("region_id__in:%s", strings.Join(strings.Split(user.AreaRegions, ","), ",")), nil, nil)
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] query region area ids failed, err: %v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.QueryAreaFailedErrorCode)
	}
	for _, areaId := range regionAreaIds {
		areaIdMap[strconv.FormatInt(int64(areaId.AreaId), 10)] = struct{}{}
	}
	return areaIdMap, nil
}
