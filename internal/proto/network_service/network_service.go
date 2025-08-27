package network_service

import (
	common "gitlab.vrviu.com/diskless_stack/diskless_stack/protocol/common"
)

type QueryOpSimNetInfoReq struct {
	FlowId   string   `protobuf:"bytes,1,opt,name=flow_id,proto3" json:"flow_id,omitempty"`
	Offset   int32    `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
	Limit    int32    `protobuf:"varint,3,opt,name=limit,proto3" json:"limit,omitempty"`
	CondList []string `protobuf:"bytes,4,rep,name=cond_list,proto3" json:"cond_list,omitempty"`
	Sorts    string   `protobuf:"bytes,5,opt,name=sorts,proto3" json:"sorts,omitempty"`
	Orders   string   `protobuf:"bytes,6,opt,name=orders,proto3" json:"orders,omitempty"`
	Operator string   `protobuf:"bytes,7,opt,name=operator,proto3" json:"operator,omitempty"`
}
type QueryOpSimNetInfoRsp struct {
	Ret  *common.RspInfo        `protobuf:"bytes,1,opt,name=ret,proto3" json:"ret,omitempty"`
	Body *QueryOpSimNetInfoBody `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
}

type QueryOpSimNetInfoBody struct {
	FlowId string               `protobuf:"bytes,1,opt,name=flow_id,proto3" json:"flow_id,omitempty"`
	Total  int32                `protobuf:"varint,2,opt,name=total,proto3" json:"total,omitempty"`
	List   []*OperateSimNetInfo `protobuf:"bytes,3,rep,name=list,proto3" json:"list,omitempty"`
}
type OperateSimNetInfo struct {
	Id float64 `protobuf:"fixed64,1,opt,name=id,proto3" json:"id,omitempty"`
	Ip string  `protobuf:"bytes,2,opt,name=ip,proto3" json:"ip,omitempty"`
}


const (
	IPTypeInstance = 1
	IPTypeClient   = 2
	IpTypeBox = 3
)