package instance

import (
	"context"

	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InstanceAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInstanceAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InstanceAddLogic {
	return &InstanceAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InstanceAddLogic) InstanceAdd(req *types.InstanceAddReq) (resp *types.InstanceAddResp, err error) {

	return
}
