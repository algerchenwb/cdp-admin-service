package biz

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

type UpdateBizStrategyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateBizStrategyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBizStrategyLogic {
	return &UpdateBizStrategyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// todo 已废弃
func (l *UpdateBizStrategyLogic) UpdateBizStrategy(req *types.UpdateBizStrategyReq) (resp *types.UpdateBizStrategyResp, err error) {

	sessionId := helper.GetSessionId(l.ctx)
	updateBy := helper.GetUserName(l.ctx)

	if req.Id == 0 || req.BizId == 0 {
		return nil, errorx.NewDefaultError(errorx.ParamErrorCode)
	}
	// 查询租户信息
	_, _, err = table.T_TCdpBizInfoService.Query(l.ctx, sessionId, fmt.Sprintf("biz_id:%d", req.BizId), nil, nil)
	if err != nil && err != gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] T_TCdpBizInfoService Query err. bizId[%d] err:%+v", sessionId, req.BizId, err)
		return nil, errorx.NewDefaultError(errorx.DalGetErrorCode)
	}
	if err == gopublic.ErrNotExist {
		l.Logger.Errorf("[%s] T_TCdpBizInfoService Query err. ErrNotExist bizId[%d] err:%+v", sessionId, req.BizId, err)
		return nil, errorx.NewDefaultError(errorx.BizNotFoundErrorCode)
	}

	newBizStrategyInfo, _, err := table.T_TCdpBizStrategyService.Update(l.ctx, sessionId, req.Id, map[string]interface{}{
		// "boot_schema_id": req.BootSchemaId,
		"update_by": updateBy,
	})
	if err != nil {
		l.Logger.Errorf("[%s] T_TCdpBizStrategyService Update err. bizId[%d] err:%+v", sessionId, req.BizId, err)
		return nil, errorx.NewDefaultError(errorx.DalUpdateErrorCode)
	}

	l.Logger.Infof("[%s] T_TCdpBizStrategyService Update success. bizId[%d] newBizStrategyInfo:%+v", sessionId, req.BizId, newBizStrategyInfo)
	// 复制数据

	return
}
