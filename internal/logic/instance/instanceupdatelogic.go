package instance

import (
	"context"

	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InstanceUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInstanceUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InstanceUpdateLogic {
	return &InstanceUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InstanceUpdateLogic) InstanceUpdate(req *types.InstanceUpdateReq) (resp *types.InstanceUpdateResp, err error) {

	return
}
