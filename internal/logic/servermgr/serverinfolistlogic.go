package servermgr

import (
	"context"
	"fmt"

	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type ServerInfoListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewServerInfoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ServerInfoListLogic {
	return &ServerInfoListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ServerInfoListLogic) ServerInfoList(req *types.CommonPageRequest) (resp *types.ServerInfoListResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)
	areaIds := helper.GetAreaIds(l.ctx)
	req.CondList = append(req.CondList, fmt.Sprintf("area_id__in:%s", areaIds))
	qry, err := helper.CheckCommQueryParam(l.ctx, sessionId, req.CondList, req.Orders, req.Sorts, table.TCdpServerManagement{})
	if err != nil {
		l.Logger.Errorf("[%s] CheckCommQueryParam failed, req:%+v err:%+v", sessionId, req, err)
		return nil, errorx.NewDefaultError(errorx.ParamErrorCode)
	}

	total, list, _, err := table.T_TCdpServerManagementService.QueryPage(l.ctx, sessionId, qry, req.Offset, req.Limit, req.Sorts, req.Orders)
	if err != nil {
		l.Logger.Errorf("[%s] ServerInfoList err: %v", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.QueryServerInfoListErrorCode)
	}

	resp = &types.ServerInfoListResp{
		Total: total,
	}

	for _, item := range list {
		var serverInfo types.ServerInfo

		err = copier.Copy(&serverInfo, item)

		serverInfo.CreateTime = item.CreateTime.Format("2006-01-02 15:04:05")
		serverInfo.UpdateTime = item.UpdateTime.Format("2006-01-02 15:04:05")
		serverInfo.ModifyTime = item.ModifyTime.Format("2006-01-02 15:04:05")
		serverInfo.BootTime = item.BootTime.Format("2006-01-02 15:04:05")

		resp.List = append(resp.List, serverInfo)
	}

	l.Logger.Debugf("[%s] ServerInfoList list: %v", sessionId, resp.List)

	return
}
