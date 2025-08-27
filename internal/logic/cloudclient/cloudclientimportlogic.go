package cloudclient

import (
	"context"

	"cdp-admin-service/internal/svc"
	"cdp-admin-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CloudClientImportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCloudClientImportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CloudClientImportLogic {
	return &CloudClientImportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CloudClientImportLogic) CloudClientImport() (resp *types.CloudClientImportResp, err error) {

	return
}
