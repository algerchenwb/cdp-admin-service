package common

// import (
// 	"cdp-admin-service/internal/helper"
// 	table "cdp-admin-service/internal/helper/dal"
// 	"cdp-admin-service/internal/svc"
// 	"context"
// 	"fmt"

// 	"github.com/zeromicro/go-zero/core/logx"
// )

// // 1.0客户机
// func UpdateLocation1(ctx context.Context, svcCtx *svc.ServiceContext,
// 	clientId int64, hostIp string, hostName string, mac string, bootSchemaId int64) (err error) {
// 	sessionId := helper.GetSessionId(ctx)
// 	client, _, err := table.T_TCdpCloudclientInfoService.Query(ctx, sessionId, fmt.Sprintf("id:%d$status:%d", clientId, table.CloudClientStatusValid), nil, nil)
// 	if err != nil {
// 		logx.WithContext(ctx).Errorf("[%s] T_TCdpCloudclientInfoService Query err. clientId[%d] err:%+v", sessionId, clientId, err)
// 		return err
// 	}

// 	return nil
// }

// // 2.0客户机
// func UpdateLocation2(ctx context.Context, svcCtx *svc.ServiceContext,
// 	clientId int64, hostIp string, hostName string, mac string, bootSchemaId int64) (err error) {
// 	sessionId := helper.GetSessionId(ctx)
// 	client, _, err := table.T_TCdpCloudclientInfoService.Query(ctx, sessionId, fmt.Sprintf("id:%d$status:%d", clientId, table.CloudClientStatusValid), nil, nil)
// 	if err != nil {
// 		logx.WithContext(ctx).Errorf("[%s] T_TCdpCloudclientInfoService Query err. clientId[%d] err:%+v", sessionId, clientId, err)
// 		return err
// 	}
// 	return nil
// }

// // 新增1.0客户机
// func AddLocation1(ctx context.Context, svcCtx *svc.ServiceContext,
// 	clientId int64, hostIp string, hostName string, mac string, bootSchemaId int64) (err error) {
// 	sessionId := helper.GetSessionId(ctx)
// 	client, _, err := table.T_TCdpCloudclientInfoService.Query(ctx, sessionId, fmt.Sprintf("id:%d$status:%d", clientId, table.CloudClientStatusValid), nil, nil)
// 	if err != nil {
// 		logx.WithContext(ctx).Errorf("[%s] T_TCdpCloudclientInfoService Query err. clientId[%d] err:%+v", sessionId, clientId, err)
// 		return err
// 	}
// 	return nil
// }

// // 新增2.0客户机
// func AddLocation2(ctx context.Context, svcCtx *svc.ServiceContext,
// 	clientId int64, hostIp string, hostName string, mac string, bootSchemaId int64) (err error) {
// 	sessionId := helper.GetSessionId(ctx)
// 	client, _, err := table.T_TCdpCloudclientInfoService.Query(ctx, sessionId, fmt.Sprintf("id:%d$status:%d", clientId, table.CloudClientStatusValid), nil, nil)
// 	if err != nil {
// 		logx.WithContext(ctx).Errorf("[%s] T_TCdpCloudclientInfoService Query err. clientId[%d] err:%+v", sessionId, clientId, err)
// 		return err
// 	}
// 	return nil
// }
