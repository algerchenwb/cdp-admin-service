package diskless

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpc"
	instance_types "gitlab.vrviu.com/diskless_stack/protocol/instance_service/types"
)

type QueryAreaListReq struct {
	Offset   int64    `json:"offset"`
	Limit    int64    `json:"limit"`
	Order    string   `json:"order"`
	Sortby   string   `json:"sortby"`
	CondList []string `json:"cond_list"`
}
type QueryAreaListResp struct {
	Head HTTPCommonHead    `json:"ret"`
	Body QueryAreaListBody `json:"body,omitempty"`
}
type QueryAreaListBody struct {
	FlowId string     `json:"flow_id"`
	Total  int64      `json:"total"`
	List   []AreaInfo `json:"list"`
}
type AreaInfo struct {
	Id             string `json:"id"`
	AreaId         string `json:"area_id"`
	Name           string `json:"name"`
	RegionId       int64  `json:"region_id"`
	DeploymentType int64  `json:"deployment_type"`
	Remark         string `json:"remark"`
	ManagerState   int64  `json:"manager_state"`
	OnlineVm       int64  `json:"online_vm"`
	TotalVm        int64  `json:"total_vm"`
	Version        int64  `json:"version"`
	ProxyAddr      string `json:"proxy_addr"`
	ProxyOnline    int64  `json:"proxy_online"`
	CreateTime     string `json:"create_time"`
	UpdateTime     string `json:"update_time"`
	ModifyTime     string `json:"modify_time"`
}

func (d *DisklessWebGateway) QueryAreaList(req *QueryAreaListReq) (*QueryAreaListResp, error) {

	url := fmt.Sprintf("%s/v1/DisklessCloudWeb/QueryAreaList", d.SvcCtx.Config.OutSide.DisklessHost)

	data, err := httpc.Do(context.Background(), http.MethodPost, url, req)
	if err != nil {
		logx.WithContext(d.Ctx).Errorf("httpc.Do err req[%+v] host[%s] err:%+v", req, url, err)
		return nil, errors.New("查询失败")
	}
	logx.WithContext(d.Ctx).Debugf("QueryAreaList data[%+v]", data)

	body, _ := io.ReadAll(data.Body)
	resp := new(QueryAreaListResp)
	if err := json.Unmarshal([]byte(body), resp); err != nil {
		logx.WithContext(d.Ctx).Errorf("json.Unmarshal req[%+v] host[%s] err:%+v", req, url, err)
		return nil, errors.New("反序列化失败")
	}

	if resp.Head.Code != 0 {
		logx.WithContext(d.Ctx).Errorf("QueryAreaList err. req[%+v] host[%s] vlanInfo:%+v", req, url, resp)
		return nil, errors.New(resp.Head.Msg)
	}

	logx.WithContext(d.Ctx).Debugf("QueryAreaList succ req[%+v] host[%s] , resp:%+v", req, url, resp)

	return nil, nil
}

type SearchInstancesReq struct {
	FlowId      string  `json:"flow_id"`
	Offset      int64   `json:"offset"`
	Length      int64   `json:"length"`
	DeviceTypes []int64 `json:"device_types"`
	Ids         []int64 `json:"ids"`
}

type SearchInstancesResp struct {
	Instances []Instance `json:"instances"`
	Total     int64      `json:"total"`
}
type Instance struct {
	Id                 int64         `json:"id"`
	HostId             int64         `json:"host_id"`
	SchemeId           int64         `json:"scheme_id"`
	NetInfo            NetInfo       `json:"net_info"`
	BootMac            string        `json:"boot_mac"`
	BootType           int64         `json:"boot_type"`
	KeepType           int64         `json:"keep_type"`
	Options            Options       `json:"options"`
	State              int64         `json:"state"`
	Tags               []string      `json:"tags"`
	ActivityIp         string        `json:"activity_ip"`
	BootTime           string        `json:"boot_time"`
	ManageStatus       int64         `json:"manage_status"`
	PowerStatus        int64         `json:"power_status"`
	RunningStatus      int64         `json:"running_status"`
	BootStatus         int64         `json:"boot_status"`
	BusinessStatus     int64         `json:"business_status"`
	AssignStatus       int64         `json:"assign_status"`
	ManageStatusDesc   string        `json:"manage_status_desc"`
	PowerStatusDesc    string        `json:"power_status_desc"`
	RunningStatusDesc  string        `json:"running_status_desc"`
	BootStatusDesc     string        `json:"boot_status_desc"`
	BusinessStatusDesc string        `json:"business_status_desc"`
	AssignStatusDesc   string        `json:"assign_status_desc"`
	AssignSource       string        `json:"assign_source"`
	AssignOrder        string        `json:"assign_order"`
	StatusRemark       string        `json:"status_remark"`
	InstanceRemark     string        `json:"instance_remark"`
	UserMode           int64         `json:"user_mode"`
	HostInfo           HostInfo      `json:"host_info"`
	DefaultConfig      DefaultConfig `json:"default_config"`
	OsImage            string        `json:"os_image"`
	OsVolumeId         int64         `json:"os_volume_id"`
	DataImage          string        `json:"data_image"`
	Specification      int64         `json:"specification"`
}
type HostInfo struct {
	Name  string `json:"name"`
	Arch  string `json:"arch"`
	Net   string `json:"net"`
	Mem   string `json:"mem"`
	Disk0 string `json:"disk0"`
	Ip    string `json:"ip"`
}
type DefaultConfig struct {
	Ip      string `json:"ip"`
	Netmask string `json:"netmask"`
	Gateway string `json:"gateway"`
	Dns     string `json:"dns"`
}
type Options struct {
	OsVolumeOnly int64 `json:"os_volume_only"`
}
type NetInfo struct {
}

func (d *DisklessWebGateway) SearchInstances(areaId uint64, sessionId string, req *SearchInstancesReq) (*SearchInstancesResp, error) {
	d.Api = "search_instances"
	d.Subsystem = "iaas"
	d.AreaId = int64(areaId)
	request, err := json.Marshal(req)
	if err != nil {
		logx.WithContext(d.Ctx).Errorf("json.Marshal flowId[%s] areaId[%d] subsystem[%s] api[%s], req :%+v", req.FlowId, d.AreaId, d.Subsystem, d.Api, req)
		return nil, err
	}
	info, err := d.PostAreaApi(req.FlowId, string(request))
	if err != nil {
		return nil, err
	}
	resp := new(SearchInstancesResp)
	httpResp := &instance_types.HTTPResponse{
		Body: resp,
	}
	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		return nil, errors.New("反序列化失败")
	}
	if httpResp.Head.Code != 0 {
		return nil, errors.New(httpResp.Head.Msg)
	}

	return resp, nil
}
