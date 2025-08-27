package diskless

import (
	proto "cdp-admin-service/internal/proto/network_service"
	"encoding/json"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
)

func (d *DisklessWebGateway) VlanIPList(sessionId string, areaId int64, req *proto.QueryOpSimNetInfoReq) (*proto.QueryOpSimNetInfoBody, error) {
	d.Api = "query_op_sim_netinfo"
	d.AreaId = int64(areaId)
	d.Subsystem = "iaas"
	request, err := json.Marshal(req)
	if err != nil {
		logx.WithContext(d.Ctx).Errorf("json.Marshal flowId[%s] areaId[%d] subsystem[%s] api[%s], req :%+v", req.FlowId, d.AreaId, d.Subsystem, d.Api, req)
		return nil, err
	}
	info, err := d.PostAreaApi(req.FlowId, string(request))
	if err != nil {
		return nil, err
	}
	httpResp := new(proto.QueryOpSimNetInfoRsp)
	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		return nil, errors.New("反序列化失败")
	}
	if httpResp.Ret.Code != 0 {
		return nil, errors.New(httpResp.Ret.Msg)
	}

	return httpResp.Body, nil
}


const (
	IPTypeInstance = 1
	IPTypeClient   = 2
	IpTypeBox = 3
)