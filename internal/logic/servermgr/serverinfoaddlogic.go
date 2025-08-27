package servermgr

import (
	"context"
	"time"

	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ServerInfoAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewServerInfoAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ServerInfoAddLogic {
	return &ServerInfoAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ServerInfoAddLogic) ServerInfoAdd(req *types.ServerInfoAddReq) (resp *types.ServerInfoAddResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)
	updateBy := helper.GetUserName(l.ctx)

	if req.Ip == "" || req.Type == "" || req.AreaId == 0 {
		return nil, errorx.NewDefaultError(errorx.ParamErrorCode)
	}

	if isValid, _, _ := helper.CheckIP(req.Ip); !isValid {
		l.Logger.Errorf("[%s] ServerInfoAdd err: %v", sessionId, errorx.NewDefaultError(errorx.ParamErrorCode))
		return nil, errorx.NewDefaultError(errorx.IPFormatErrorCode)
	}

	newServerInfo, _, err := table.T_TCdpServerManagementService.Insert(l.ctx, sessionId, &table.TCdpServerManagement{
		Ip:         req.Ip,
		AreaId:     req.AreaId,
		Type:       req.Type,
		Remark:     req.Remark,
		BootTime:   time.Now(),
		Status:     1,
		CreateBy:   updateBy,
		UpdateBy:   updateBy,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		ModifyTime: time.Now(),
	})
	if err != nil {
		l.Logger.Errorf("[%s] ServerInfoAdd  dal Insert err: %v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.DalInsertErrorCode)
	}

	l.Logger.Debugf("[%s] ServerInfoAdd success newServerInfo: %v", sessionId, newServerInfo)
	resp = &types.ServerInfoAddResp{}
	return
}
