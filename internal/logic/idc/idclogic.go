package idc

import (
	"context"

	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IdcLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIdcLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IdcLogic {
	return &IdcLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IdcLogic) Idc(req *types.CommonPageRequest) (resp *types.IdcListResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)

	qry, err := helper.CheckCommQueryParam(l.ctx, sessionId, req.CondList, req.Orders, req.Sorts, table.TCdpAreaInfo{})
	if err != nil {
		l.Logger.Errorf("[%s] CheckCommQueryParam failed, req:%+v err:%+v", sessionId, req, err)
		return nil, errorx.NewDefaultError(errorx.ParamErrorCode)
	}
	total, list, _, err := table.T_TCdpAreaInfoService.QueryPage(l.ctx, sessionId, qry, req.Offset, req.Limit, req.Sorts, req.Orders)
	if err != nil {
		l.Logger.Errorf("[%s] QueryPage qry:%s. err:%v", sessionId, qry, err)
		return nil, errorx.NewDefaultError(errorx.QueryAreaInfoError)
	}

	resp = &types.IdcListResp{
		Total: total,
	}
	for _, v := range list {
		resp.List = append(resp.List, types.Idc{
			Id:             int(v.Id),
			PrimaryId:      int(v.PrimaryId),
			AgentId:        int(v.AgentId),
			AreaId:         v.AreaId,
			Name:           v.Name,
			RegionId:       int(v.RegionId),
			DeploymentType: int(v.DeploymentType),
			Remark:         v.Remark,
			ProxyAddr:      v.ProxyAddr,
			CreateBy:       v.CreateBy,
			UpdateBy:       v.UpdateBy,
			CreateTime:     v.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateTime:     v.UpdateTime.Format("2006-01-02 15:04:05"),
			ModifyTime:     v.ModifyTime.Format("2006-01-02 15:04:05"),
		})
	}
	return
}
