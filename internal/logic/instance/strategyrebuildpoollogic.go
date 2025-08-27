package instance

import (
	"context"
	"fmt"

	"cdp-admin-service/internal/helper"
	"cdp-admin-service/internal/helper/diskless"
	"cdp-admin-service/internal/model/errorx"
	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type StrategyReBuildPoolLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStrategyReBuildPoolLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StrategyReBuildPoolLogic {
	return &StrategyReBuildPoolLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StrategyReBuildPoolLogic) StrategyReBuildPool(req *types.StrategyReBuildPoolReq) (resp *types.StrategyReBuildPoolResp, err error) {
	sessionId := helper.GetSessionId(l.ctx)

	if req.AreaType == 0 || req.ResourceId == 0 {
		return nil, errorx.NewDefaultCodeError("参数错误")
	}

	reBuidReq := &diskless.RebuildPoolRequest{
		AreaType:   fmt.Sprintf("%d", req.AreaType),
		ResourceId: fmt.Sprintf("%d", req.ResourceId),
	}
	if err = diskless.NewDisklessWebGateway(l.ctx, l.svcCtx).RebuidPool(sessionId, int64(req.AreaType), reBuidReq); err != nil {
		l.Logger.Errorf("[%s] diskless ReBuildPool req:%s err:%+v", sessionId, helper.ToJSON(req), err)
		return nil, errorx.NewDefaultCodeError("重建策略资源池失败")
	}

	return
}
