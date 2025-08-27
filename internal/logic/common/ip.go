package common

import (
	"cdp-admin-service/internal/helper"
	"cdp-admin-service/internal/helper/diskless"
	"cdp-admin-service/internal/svc"
	"context"
	"fmt"

	proto "cdp-admin-service/internal/proto/network_service"

	"github.com/zeromicro/go-zero/core/logx"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
)

type Vlan struct {
	VlanId  int64
	FreeIps map[string]struct{}
}

func LoadVlanIP(ctx context.Context, svcCtx *svc.ServiceContext, vlanId int64, areaId int64, ipType int32) (vlan Vlan, err error) {
	sessionId := helper.GetSessionId(ctx)
	vlan = Vlan{
		VlanId:  vlanId,
		FreeIps: make(map[string]struct{}),
	}
	vlanReq := &proto.QueryOpSimNetInfoReq{
		FlowId:   sessionId,
		Offset:   0,
		Limit:    4096,
		CondList: []string{fmt.Sprintf("vlan_id__eq:%d", vlanId), "bind_mac__eq:", fmt.Sprintf("bind_type__eq:%d", ipType)},
		Operator: "cdp-admin-service",
	}
	vlanRsp, err := diskless.NewDisklessWebGateway(ctx, svcCtx).VlanIPList(sessionId, areaId, vlanReq)
	if err != nil {
		logx.WithContext(ctx).Errorf("load vlan ip failed, areaId: %d, vlanReq: %s, err: %v", areaId, gopublic.ToJSON(vlanReq), err)
		return vlan, err
	}
	logx.WithContext(ctx).Infof("load vlan ip success, areaId: %d, vlanReq: %s", areaId, gopublic.ToJSON(vlanReq))
	for _, ip := range vlanRsp.List {
		vlan.FreeIps[ip.Ip] = struct{}{}
	}
	return vlan, nil
}

// 新增云盒
// 编辑云盒
// 新增客户机
// 编辑客户机
func (v *Vlan) FreeIp(ip string) bool {
	if _, ok := v.FreeIps[ip]; ok {
		return true
	}
	return false
}
