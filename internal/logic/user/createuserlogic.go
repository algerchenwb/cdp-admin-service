package user

import (
	"context"
	"fmt"
	"strings"
	"time"

	"cdp-admin-service/internal/helper"
	"cdp-admin-service/internal/logic/common"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	table "cdp-admin-service/internal/helper/dal"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateUserLogic) CreateUser(req *types.CreateUserReq) (resp *types.CreateUserResp, err error) {
	sessionId := helper.GetSessionId(l.ctx)
	isAdmin := helper.GetIsAdmin(l.ctx)
	userId := helper.GetUserId(l.ctx)
	modifier, _, err := table.T_TCdpSysUserService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d", userId), nil, nil)
	if err != nil {
		l.Logger.Errorf("table.T_TCdpSysUserService.Query err. Id[%d] err:%+v", helper.GetUserId(l.ctx), err)
		return nil, errorx.NewDefaultError(errorx.QueryUserFailedErrorCode)
	}
	if !isAdmin {
		l.Logger.Debugf("[%s] modifier[%v] not admin", sessionId, modifier)
		modifierAreaIds, err := getUserAreaIds(l.ctx, modifier)
		if err != nil {
			l.Logger.Errorf("getUserAreaIds err. Id[%d] err:%+v", userId, err)
			return nil, errorx.NewDefaultError(errorx.QueryUserFailedErrorCode)
		}
		waitUpdateAreaIds, err := getUserAreaIds(l.ctx, &table.TCdpSysUser{
			AreaIds:     helper.SliceToString(req.AreaIds),
			AreaRegions: helper.SliceToString(req.AreaRegions),
		})
		if err != nil {
			logx.WithContext(l.ctx).Errorf("getUserAreaIds err. Id[%d] err:%+v", helper.GetUserId(l.ctx), err)
			return nil, errorx.NewDefaultError(errorx.QueryUserFailedErrorCode)
		}

		access := checkAreaIdsAccessV2(modifierAreaIds, waitUpdateAreaIds)
		if !access {
			l.Logger.Errorf("[%s] modifier[%v] not access, modifierAreaIds[%+v] waitUpdateAreaIds[%+v]", sessionId, modifier, modifierAreaIds, waitUpdateAreaIds)
			return nil, errorx.NewDefaultError(errorx.IdcRefuseErrorCode)
		}
	}

	password, err := helper.Decode(req.Password)
	if err != nil {
		l.Logger.Errorf("helper.Decode err. Password[%s] err:%+v", req.Password, err)
		return nil, errorx.NewDefaultError(errorx.DecodePasswordError)
	}
	l.Logger.Infof("[%s] origin password: %s", sessionId, password)
	password = helper.HashPassword(password, l.svcCtx.Config.AES.Salt)
	areaIds := helper.SliceToString(req.AreaIds)

	areaRegions := helper.SliceToString(req.AreaRegions)

	bizIds := helper.SliceToString(req.BizIds)

	bizRegions := helper.SliceToString(req.BizRegions)

	sysUser := &table.TCdpSysUser{
		Account:     req.Account,
		Password:    password,
		Nickname:    req.Nickname,
		Avatar:      req.Avatar,
		Mobile:      req.Mobile,
		RoleId:      int32(req.RoleId),
		AreaIds:     areaIds,
		AreaRegions: areaRegions,
		BizIds:      bizIds,
		BizRegions:  bizRegions,
		Status:      uint32(req.Status),
		Platform:    int32(req.Platform),
		Remark:      req.Remark,
		CreateBy:    helper.GetUserName(l.ctx),
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
		ModifyTime:  time.Now(),
	}
	sysUser, _, err = table.T_TCdpSysUserService.Insert(l.ctx, sessionId, sysUser)
	if err != nil {
		l.Logger.Errorf("[%s] table.T_SysUserService.Insert err. SysUser[%v] err:%+v", sessionId, sysUser, err)
		if strings.Contains(err.Error(), "Duplicate") {
			return nil, errorx.NewDefaultError(errorx.DuplicateAccountErrorCode)
		}
		return nil, errorx.NewDefaultError(errorx.CreateUserFailedErrorCode)
	}
	l.Logger.Infof("[%s] table.T_SysUserService.Insert success. SysUser[%v]", sessionId, sysUser)
	//  非管理员账号需更新权限
	if !sysUser.UserIsAdmin() {
		err = common.UpdateUserPermCache(l.ctx, sysUser, l.svcCtx.Cache)
		if err != nil {
			l.Logger.Errorf("[%s] ConfigUserRole err. SysUser[%v] err:%+v", sessionId, sysUser, err)
			return nil, errorx.NewDefaultError(errorx.UpdateUserPermFailedErrorCode)
		}
	}

	return &types.CreateUserResp{
		Id: int64(sysUser.Id),
	}, nil
}

func PermStandard(perm string) string {
	if strings.HasPrefix(perm, "/") {
		return perm
	}
	return "/" + perm
}
