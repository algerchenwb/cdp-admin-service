package log

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

type GetLogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLogListLogic {
	return &GetLogListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLogListLogic) GetLogList(req *types.CommonPageRequest) (resp *types.LogListResp, err error) {
	sessionId := helper.GetSessionId(l.ctx)
	platform := helper.GetPlatform(l.ctx)
	req.CondList = append(req.CondList, fmt.Sprintf("platform:%d", platform))
	qry, err := helper.CheckCommQueryParam(l.ctx, sessionId, req.CondList, req.Orders, req.Sorts, table.TCdpSysLog{})
	if err != nil {
		l.Logger.Errorf("[%s] CheckCommQueryParam failed, req:%+v err:%+v", sessionId, req, err)
		return nil, errorx.NewDefaultError(errorx.ParamErrorCode)
	}
	total, logs, _, err := table.T_TCdpSysLogService.QueryPage(l.ctx, sessionId, qry, req.Offset, req.Limit, req.Sorts, req.Orders)
	if err != nil {
		l.Logger.Error("[%s] 查询日志失败", sessionId, err)
		return nil, errorx.NewDefaultError(errorx.QuerySysLogFailedErrorCode)
	}
	resp = &types.LogListResp{
		Total: int64(total),
		List:  make([]types.Log, 0),
	}
	for _, log := range logs {
		resp.List = append(resp.List, types.Log{
			Id:         int64(log.Id),
			UserId:     int64(log.UserId),
			Account:    log.Account,
			Ip:         log.Ip,
			Uri:        log.Uri,
			Type:       int64(log.Type),
			Request:    log.Request,
			Response:   log.Response,
			Status:     int64(log.Status),
			CreateTime: log.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateTime: log.UpdateTime.Format("2006-01-02 15:04:05"),
			ModifyTime: log.ModifyTime.Format("2006-01-02 15:04:05"),
		})
	}

	return
}
