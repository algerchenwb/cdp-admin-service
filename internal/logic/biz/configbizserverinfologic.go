package biz

import (
	"context"
	"fmt"
	"time"

	"cdp-admin-service/internal/helper"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
)

type ConfigBizServerInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfigBizServerInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigBizServerInfoLogic {
	return &ConfigBizServerInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfigBizServerInfoLogic) ConfigBizServerInfo(req *types.ConfigBizServerInfoReq) (resp *types.ConfigBizServerInfoResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)
	updateBy := helper.GetUserName(l.ctx)

	if len(req.Server_infos) == 0 || req.BizId == 0 {
		return nil, errorx.NewDefaultError(errorx.ParamErrorCode)
	}

	// 查询租户信息
	bizInfo, _, err := table.T_TCdpBizInfoService.Query(l.ctx, sessionId, fmt.Sprintf("biz_id:%d", req.BizId), nil, nil)
	if err != nil && err != gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] Table Query err. bizId[%d] err:%+v", sessionId, req.BizId, err)
		return nil, errorx.NewDefaultError(errorx.DalGetErrorCode)
	}
	if err == gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] Table Query err. ErrNotExist bizId[%d] err:%+v", sessionId, req.BizId, err)
		return nil, errorx.NewDefaultError(errorx.BizNotFoundErrorCode)
	}

	newBizInfo, _, err := table.T_TCdpBizInfoService.Update(l.ctx, sessionId, bizInfo.Id, map[string]interface{}{
		"serverinfo": req.Server_infos,
		"update_by":  updateBy,
		"update_at":  time.Now(),
	})
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpBizInfoService Update err. bizId[%d] err:%+v", sessionId, req.BizId, err)
		return nil, errorx.NewDefaultError(errorx.DalUpdateErrorCode)
	}

	l.Logger.Infof("[%s] ConfigBizServerInfo Table Update success. bizId[%d] newBizInfo:%s", sessionId, req.BizId, helper.ToJSON(newBizInfo))

	return
}
