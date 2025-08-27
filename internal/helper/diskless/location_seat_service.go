package diskless

import (
	proto "cdp-admin-service/internal/proto/location_seat_service"
	"encoding/json"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
)

func (d *DisklessWebGateway) AddLocation(sessionId string, areaId int64, req *proto.AddLocationRequest) (*proto.AddLocationBody, error) {
	d.Api = "add_location"
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
	httpResp := new(proto.AddLocationResponse)
	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		return nil, errors.New("反序列化失败")
	}
	if httpResp.Ret.Code != 0 {
		return nil, errors.New(httpResp.Ret.Msg)
	}

	return httpResp.Body, nil

}

func (d *DisklessWebGateway) DeleteLocation(sessionId string, areaId int64, req *proto.DeleteLocationRequest) (*proto.DeleteLocationBody, error) {
	d.Api = "delete_location"
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
	httpResp := new(proto.DeleteLocationResponse)
	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		return nil, errors.New("反序列化失败")
	}
	if httpResp.Ret.Code != 0 {
		return nil, errors.New(httpResp.Ret.Msg)
	}

	return httpResp.Body, nil
}

func (d *DisklessWebGateway) UpdateLocation(sessionId string, areaId int64, req *proto.UpdateLocationRequest) (*proto.UpdateLocationBody, error) {
	d.Api = "update_location"
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
	httpResp := new(proto.UpdateLocationResponse)
	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		return nil, errors.New("反序列化失败")
	}
	if httpResp.Ret.Code != 0 {
		return nil, errors.New(httpResp.Ret.Msg)
	}

	return httpResp.Body, nil
}

func (d *DisklessWebGateway) GetLocation(sessionId string, areaId int64, req *proto.GetLocationRequest) (*proto.GetLocationBody, error) {
	d.Api = "get_location"
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
	httpResp := new(proto.GetLocationResponse)
	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		return nil, errors.New("反序列化失败")
	}
	if httpResp.Ret.Code != 0 {
		return nil, errors.New(httpResp.Ret.Msg)
	}

	return httpResp.Body, nil
}

func (d *DisklessWebGateway) GetLocationList(sessionId string, areaId int64, req *proto.GetLocationListRequest) (*proto.GetLocationListBody, error) {
	d.Api = "get_location_list"
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
	httpResp := new(proto.GetLocationListResponse)
	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		return nil, errors.New("反序列化失败")
	}
	if httpResp.Ret.Code != 0 {
		return nil, errors.New(httpResp.Ret.Msg)
	}

	return httpResp.Body, nil
}

func (d *DisklessWebGateway) AddSeat(sessionId string, areaId int64, req *proto.AddSeatRequest) (*proto.AddSeatBody, error) {
	d.Api = "add_seat"
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
	httpResp := new(proto.AddSeatResponse)

	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		return nil, errors.New("反序列化失败")
	}
	if httpResp.Ret.Code != 0 {
		return nil, errors.New(httpResp.Ret.Msg)
	}

	return httpResp.Body, nil
}

func (d *DisklessWebGateway) DeleteSeat(sessionId string, areaId int64, req *proto.DeleteSeatRequest) (*proto.DeleteSeatBody, error) {
	d.Api = "delete_seat"
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
	httpResp := new(proto.DeleteSeatResponse)

	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		return nil, errors.New("反序列化失败")
	}
	if httpResp.Ret.Code != 0 {
		return nil, errors.New(httpResp.Ret.Msg)
	}

	return httpResp.Body, nil
}

func (d *DisklessWebGateway) UpdateSeat(sessionId string, areaId int64, req *proto.UpdateSeatRequest) (*proto.UpdateSeatBody, error) {
	d.Api = "update_seat"
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
	httpResp := new(proto.UpdateSeatResponse)
	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		return nil, errors.New("反序列化失败")
	}
	if httpResp.Ret.Code != 0 {
		return nil, errors.New(httpResp.Ret.Msg)
	}

	return httpResp.Body, nil
}

func (d *DisklessWebGateway) GetSeat(sessionId string, areaId int64, req *proto.GetSeatRequest) (*proto.GetSeatBody, error) {
	d.Api = "get_seat"
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
	httpResp := new(proto.GetSeatResponse)
	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		return nil, errors.New("反序列化失败")
	}
	if httpResp.Ret.Code != 0 {
		return nil, errors.New(httpResp.Ret.Msg)
	}

	return httpResp.Body, nil
}

func (d *DisklessWebGateway) GetSeatList(sessionId string, areaId int64, req *proto.GetSeatListRequest) (*proto.GetSeatListBody, error) {
	d.Api = "get_seat_list"
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
	httpResp := new(proto.GetSeatListResponse)

	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		return nil, errors.New("反序列化失败")
	}
	if httpResp.Ret.Code != 0 {
		return nil, errors.New(httpResp.Ret.Msg)
	}

	return httpResp.Body, nil
}
