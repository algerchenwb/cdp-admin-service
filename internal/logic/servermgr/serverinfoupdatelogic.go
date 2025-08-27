package servermgr

import (
	"context"
	"fmt"

	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
)

type ServerInfoUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewServerInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ServerInfoUpdateLogic {
	return &ServerInfoUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ServerInfoUpdateLogic) ServerInfoUpdate(req *types.ServerInfoUpdateReq) (resp *types.ServerInfoUpdateResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)
	updateBy := helper.GetUserName(l.ctx)

	if req.Ip == "" || req.Type == "" {
		return nil, errorx.NewDefaultError(errorx.ParamErrorCode)
	}

	if isValid, _, _ := helper.CheckIP(req.Ip); !isValid {
		l.Logger.Errorf("[%s] Check ip format  err: %v", sessionId, errorx.NewDefaultError(errorx.ParamErrorCode))
		return nil, errorx.NewDefaultError(errorx.IPFormatErrorCode)
	}

	_, _, err = table.T_TCdpServerManagementService.Query(l.ctx, sessionId, fmt.Sprintf("id:%d", req.Id), nil, nil)
	if err != nil && err != gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] Table Query err. serverId[%d] err:%+v", sessionId, req.Id, err)
		return nil, errorx.NewDefaultError(errorx.DalGetErrorCode)
	}
	if err == gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] Table Query err. ErrNotExist serverId[%d] err:%+v", sessionId, req.Id, err)
		return nil, errorx.NewDefaultError(errorx.QueryServerInfoEmptyErrorCode)
	}

	newServerInfo, _, err := table.T_TCdpServerManagementService.Update(l.ctx, sessionId, req.Id, map[string]interface{}{
		"type":      req.Type,
		"ip":        req.Ip,
		"remark":    req.Remark,
		"update_by": updateBy})
	if err != nil {
		l.Logger.Errorf("[%s] Table Update err. serverId[%d]  err:%+v", sessionId, req.Id, err)
		return nil, errorx.NewDefaultError(errorx.ListStrategyErrorCode)
	}

	l.Logger.Debugf("[%s] ServerInfoUpdate success newServerInfo: %v", sessionId, newServerInfo)
	resp = &types.ServerInfoUpdateResp{}
	return
}
