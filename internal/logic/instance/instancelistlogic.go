package instance

import (
	"context"

	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InstanceListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInstanceListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InstanceListLogic {
	return &InstanceListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InstanceListLogic) InstanceList(req *types.CommonPageRequest) (resp *types.InstanceListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
