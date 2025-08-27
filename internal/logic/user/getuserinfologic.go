package user

import (
	"context"
	"fmt"

	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo() (resp *types.UserInfoResp, err error) {
	userId := helper.GetUserId(l.ctx)
	sessionId := helper.GetSessionId(l.ctx)
	user, _, err := table.T_TCdpSysUserService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d", userId), nil, nil)
	if err != nil {
		l.Logger.Errorf("[%s] GetUserInfoLogic GetUserInfo Query error %+v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.UserIdErrorCode)
	}

	return &types.UserInfoResp{
		Id:          int64(user.Id),
		Account:     user.Account,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Mobile:      user.Mobile,
		RoleId:      int64(user.RoleId),
		AreaIds:     helper.StringToInt64Slice(user.AreaIds),
		AreaRegions: helper.StringToInt64Slice(user.AreaRegions),
		BizIds:      helper.StringToInt64Slice(user.BizIds),
		BizRegions:  helper.StringToInt64Slice(user.BizRegions),
		Platform:    int64(user.Platform),
		Status:      int64(user.Status),
		Remark:      user.Remark,
	}, nil
}
