package diskless

import (
	"cdp-admin-service/internal/helper"
	"cdp-admin-service/internal/svc"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"cdp-admin-service/internal/proto/diskless_cloud_image"
	"cdp-admin-service/internal/proto/diskless_cloud_web"
	"cdp-admin-service/internal/proto/image_service"
	"cdp-admin-service/internal/proto/instance_scheduler"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpc"

	instance_types "cdp-admin-service/internal/proto/instance_service/types"
	"cdp-admin-service/internal/proto/network_service"

	"google.golang.org/protobuf/encoding/protojson"
)

type PostAreaApiReq struct {
	AreaId    int64  `json:"area_id"`
	Subsystem string `json:"subsystem"`
	Api       string `json:"api"`
	Request   string `json:"request"`
}

type DisklessWebGateway struct {
	Ctx       context.Context
	SvcCtx    *svc.ServiceContext
	Subsystem string
	AreaId    int64
	Api       string
}

func NewDisklessWebGateway(ctx context.Context, svcCtx *svc.ServiceContext) *DisklessWebGateway {

	return &DisklessWebGateway{
		Ctx:       ctx,
		Subsystem: "iaas",
		SvcCtx:    svcCtx,
	}
}

func (l *DisklessWebGateway) PostAreaApi(sessionId string, request string) (info *diskless_cloud_web.CallAreaApiBody, err error) {

	url := fmt.Sprintf("%s/v1/DisklessCloudWeb/PostAreaApi", l.SvcCtx.Config.OutSide.DisklessHost)
	req := PostAreaApiReq{
		AreaId:    l.AreaId,
		Subsystem: l.Subsystem,
		Api:       l.Api,
		Request:   request,
	}

	resp := new(diskless_cloud_web.CallAreaApiResponse)
	defer func() {
		logx.WithContext(l.Ctx).Debugf("[%s] PostAreaApi  host[%s]  areaId[%d] api[%s], req:%s,  resp.Body:%s", sessionId, url, l.AreaId, l.Api, helper.ToJSON(req), helper.ToJSON(resp.Body))
	}()

	code, err := helper.HttpPost(l.Ctx, sessionId, url, 120000, req, resp)
	if code != 0 || err != nil {
		logx.WithContext(l.Ctx).Errorf("[%s] helper.HttpPost HttpPost failure url[%s] err:%+v", sessionId, url, err)
		return nil, err
	}

	if resp.Ret.Code != 0 {
		logx.WithContext(l.Ctx).Errorf("[%s] HttpPost  resp err. host[%s] areaId[%d] api[%s], resp:%+v", sessionId, url, l.AreaId, l.Api, resp)
		return nil, errors.New(resp.Ret.Msg)
	}

	return resp.Body, nil
}

func (l *DisklessWebGateway) GetAreaApi(sessionId string, request string) (info *diskless_cloud_web.CallAreaApiBody, err error) {

	request1 := url.QueryEscape(request)
	req := fmt.Sprintf("area_id=%d&subsystem=%s&api=%s&request=%s", l.AreaId, l.Subsystem, l.Api, request1)
	url := fmt.Sprintf("%s/v1/DisklessCloudWeb/GetAreaApi?%s", l.SvcCtx.Config.OutSide.DisklessHost, req)
	resp := new(diskless_cloud_web.CallAreaApiResponse)

	defer func() {
		logx.WithContext(l.Ctx).Debugf("[%s] GetAreaApi host[%s] areaId[%d] api[%s], req:%s, resp.Body:%s", sessionId, url, l.AreaId, l.Api, helper.ToJSON(req), helper.ToJSON(resp.Body))
	}()

	data, err := httpc.Do(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		logx.WithContext(l.Ctx).Errorf("[%s] GetAreaApi  httpc.Do err host[%s] flowId[%s] areaId[%d] api[%s], err:%+v", sessionId, url, l.AreaId, l.Api, err)
		return nil, err
	}

	body, _ := io.ReadAll(data.Body)

	if err := json.Unmarshal(body, resp); err != nil {
		logx.WithContext(l.Ctx).Errorf("[%s] GetAreaApi json.Unmarshal err host[%s] flowId[%s] areaId[%d] api[%s], err:%+v", sessionId, url, l.AreaId, l.Api, err)
		return nil, errors.New("反序列化失败")
	}

	if resp.Ret.Code != 0 {
		logx.WithContext(l.Ctx).Errorf("[%s] GetAreaApihttpc.Do resp err. host[%s] flowId[%s] areaId[%d] api[%s], resp:%+v", sessionId, url, l.AreaId, l.Api, resp)
		return nil, errors.New(resp.Ret.Msg)
	}

	return resp.Body, nil
}

// 更新 vlan 信息
func (l *DisklessWebGateway) UpdateVlan(areaId int64, flowId string, vlanInfo *network_service.UpdateVlanRequest) error {

	l.Api = "update_vlan"
	l.AreaId = areaId

	request, err := protojson.Marshal(vlanInfo)
	if err != nil {
		logx.WithContext(l.Ctx).Errorf("protojson.Marshal flowId[%s] areaId[%d] api[%s], vlanInfo:%+v", flowId, l.AreaId, l.Api, vlanInfo.String())
		return err
	}

	info, err := l.PostAreaApi(flowId, string(request))
	if err != nil {
		return err
	}

	vlanResp := new(network_service.UpdateVlanResponse)
	if err := protojson.Unmarshal([]byte(info.Response), vlanResp); err != nil {
		logx.WithContext(l.Ctx).Errorf("protojson.Unmarshal flowId[%s] areaId[%d] api[%s], err:%+v", flowId, l.AreaId, l.Api, err)
		return errors.New("反序列化失败")
	}

	if vlanResp.Ret.GetCode() != 0 {
		logx.WithContext(l.Ctx).Errorf("l.PostAreaApi flowId[%s] areaId[%d] api[%s], vlanInfo:%+v", flowId, l.AreaId, l.Api, vlanResp.String())
		return errors.New(vlanResp.Ret.GetMsg())
	}

	logx.WithContext(l.Ctx).Debugf("UpdateVlan flowId[%s] areaId[%d] api[%s], vlanInfo:%+v, resp:%+v", flowId, l.AreaId, l.Api, info.String(), vlanResp)
	return nil
}

// 创建实例
func (l *DisklessWebGateway) CreateInstance(areaId int64, instReq *instance_types.CreateInstanceRequest) error {

	l.Api = "create_instance"
	l.AreaId = areaId

	request, err := json.Marshal(instReq)
	if err != nil {
		logx.WithContext(l.Ctx).Errorf(" json.Marshal flowId[%s] areaId[%d] api[%s], instReq:%+v", instReq.FlowId, l.AreaId, l.Api, instReq)
		return err
	}

	info, err := l.PostAreaApi(instReq.FlowId, string(request))
	if err != nil {
		return err
	}

	httpResp := new(instance_types.HTTPResponse)
	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		logx.WithContext(l.Ctx).Errorf("json.Unmarshal flowId[%s] areaId[%d] api[%s], err:%+v", instReq.FlowId, l.AreaId, l.Api, err)
		return errors.New("反序列化失败")
	}
	if httpResp.Head.Code != 0 {
		logx.WithContext(l.Ctx).Errorf("l.PostAreaApi flowId[%s] areaId[%d] api[%s], httpResp:%+v", instReq.FlowId, l.AreaId, l.Api, httpResp)
		return errors.New(httpResp.Head.Msg)
	}

	logx.WithContext(l.Ctx).Debugf("CreateInstance flowId[%s] areaId[%d] api[%s], instReq:%+v, httpResp:%+v", instReq.FlowId, l.AreaId, l.Api, instReq, httpResp)
	return nil
}

// 查询实例列表
func (l *DisklessWebGateway) GetInstanceList(areaId int64, flowId string, macList []string) (list []instance_types.InstanceDetail, err error) {
	if len(macList) == 0 {
		return list, nil
	}
	l.Api = "list_instances"
	l.AreaId = areaId
	request := fmt.Sprintf("macs=%s", strings.Join(macList, ","))

	logx.WithContext(l.Ctx).Infof("l.PostAreaApi flowId[%s] areaId[%d] api[%s], request:%s", flowId, l.AreaId, l.Api, request)

	info, err := l.GetAreaApi(flowId, request)
	if err != nil {
		return nil, err
	}

	logx.WithContext(l.Ctx).Infof("l.PostAreaApi flowId[%s] areaId[%d] api[%s], Response:%s", flowId, l.AreaId, l.Api, info.Response)
	httpResp := new(GetInstanceList)
	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		logx.WithContext(l.Ctx).Errorf("json.Unmarshal flowId[%s] areaId[%d] api[%s], err:%+v", flowId, l.AreaId, l.Api, err)
		return nil, errors.New("反序列化失败")
	}
	if httpResp.Head.Code != 0 {
		logx.WithContext(l.Ctx).Errorf("l.PostAreaApi flowId[%s] areaId[%d] api[%s], httpResp:%+v", flowId, l.AreaId, l.Api, httpResp)
		return nil, errors.New(httpResp.Head.Msg)
	}

	list = httpResp.Body.Instances

	logx.WithContext(l.Ctx).Debugf("GetInstanceList flowId[%s] areaId[%d] api[%s], httpResp:%+v", flowId, l.AreaId, l.Api, httpResp)
	return
}

// 查询实例列表 用IP
func (l *DisklessWebGateway) GetInstanceListByIp(areaId int64, flowId string, ipList []string) (list []instance_types.InstanceDetail, err error) {
	if len(ipList) == 0 {
		return list, nil
	}
	l.Api = "list_instances"
	l.AreaId = areaId
	request := fmt.Sprintf("ips=%s", strings.Join(ipList, ","))

	logx.WithContext(l.Ctx).Infof("l.PostAreaApi flowId[%s] areaId[%d] api[%s], request:%s", flowId, l.AreaId, l.Api, request)

	info, err := l.GetAreaApi(flowId, request)
	if err != nil {
		return nil, err
	}

	logx.WithContext(l.Ctx).Infof("l.PostAreaApi flowId[%s] areaId[%d] api[%s], Response:%s", flowId, l.AreaId, l.Api, info.Response)
	httpResp := new(GetInstanceList)
	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		logx.WithContext(l.Ctx).Errorf("json.Unmarshal flowId[%s] areaId[%d] api[%s], err:%+v", flowId, l.AreaId, l.Api, err)
		return nil, errors.New("反序列化失败")
	}
	if httpResp.Head.Code != 0 {
		logx.WithContext(l.Ctx).Errorf("l.PostAreaApi flowId[%s] areaId[%d] api[%s], httpResp:%+v", flowId, l.AreaId, l.Api, httpResp)
		return nil, errors.New(httpResp.Head.Msg)
	}

	list = httpResp.Body.Instances

	logx.WithContext(l.Ctx).Debugf("GetInstanceList flowId[%s] areaId[%d] api[%s], httpResp:%+v", flowId, l.AreaId, l.Api, httpResp)
	return
}

// 查询实例列表
func (l *DisklessWebGateway) SearcchInstanceList(areaId int64, sessionId string, instListReq *instance_types.ListInstancesRequestNew) (list []instance_types.InstanceDetail, err error) {

	l.Api = "search_instances"
	l.AreaId = areaId
	request, err := json.Marshal(instListReq)
	if err != nil {
		logx.WithContext(l.Ctx).Errorf("[%s] SearcchInstanceList json.Marshal areaId[%d] subsystem[%s] api[%s], instReq:%+v", sessionId, l.AreaId, l.Subsystem, l.Api, instListReq)
		return nil, err
	}

	info, err := l.PostAreaApi(sessionId, string(request))
	if err != nil {
		return nil, err
	}

	httpResp := new(GetInstanceList)
	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		logx.WithContext(l.Ctx).Errorf("[%s]json.Unmarshal areaId[%d] api[%s], err:%+v", sessionId, l.AreaId, l.Api, err)
		return nil, errors.New("反序列化失败")
	}
	if httpResp.Head.Code != 0 {
		logx.WithContext(l.Ctx).Errorf("[%s] l.PostAreaApi areaId[%d] api[%s], httpResp:%+v", sessionId, l.AreaId, l.Api, httpResp)
		return nil, errors.New(httpResp.Head.Msg)
	}

	list = httpResp.Body.Instances

	logx.WithContext(l.Ctx).Debugf("[%s] SearcchInstanceList success areaId[%d] api[%s], httpResp:%+v", sessionId, l.AreaId, l.Api, httpResp)
	return
}

// 更新实例 update_instance
func (l *DisklessWebGateway) UpdateInstance(areaId int64, instReq *instance_types.UpdateInstanceRequest) error {

	l.Api = "update_instance"
	l.AreaId = areaId

	request, err := json.Marshal(instReq)
	if err != nil {
		logx.WithContext(l.Ctx).Errorf("[%s]json.Marshal areaId[%d] subsystem[%s] api[%s], instReq:%+v", instReq.FlowId, l.AreaId, l.Subsystem, l.Api, instReq)
		return err
	}

	info, err := l.PostAreaApi(instReq.FlowId, string(request))
	if err != nil {
		return err
	}
	httpResp := new(instance_types.HTTPResponse)
	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		logx.WithContext(l.Ctx).Errorf("[%s] json.Unmarshal areaId[%d] subsystem[%s] api[%s], err:%+v", instReq.FlowId, l.AreaId, l.Subsystem, l.Api, err)
		return errors.New("反序列化失败")
	}
	if httpResp.Head.Code != 0 {
		logx.WithContext(l.Ctx).Errorf("[%s]l.PostAreaApi areaId[%d] subsystem[%s] api[%s], httpResp:%+v", instReq.FlowId, l.AreaId, l.Subsystem, l.Api, httpResp)
		return errors.New(httpResp.Head.Msg)
	}

	logx.WithContext(l.Ctx).Debugf("[%s] UpdateInstance areaId[%d] subsystem[%s] api[%s], instReq:%+v, vlanInfo:%+v", instReq.FlowId, l.AreaId, l.Subsystem, l.Api, instReq, info.String())
	return nil
}

// 删除实例 destroy_instance
func (l *DisklessWebGateway) DestroyInstance(areaId int64, instReq *instance_types.DestroyInstanceRequest) error {

	l.Api = "destroy_instance"
	l.AreaId = areaId

	request, err := json.Marshal(instReq)
	if err != nil {
		logx.WithContext(l.Ctx).Errorf(" json.Marshal flowId[%s] areaId[%d] subsystem[%s] api[%s], instReq----->:%+v", instReq.FlowId, l.AreaId, l.Subsystem, l.Api, instReq)
		return err
	}

	info, err := l.PostAreaApi(instReq.FlowId, string(request))
	if err != nil {
		return err
	}
	httpResp := new(instance_types.HTTPResponse)
	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		logx.WithContext(l.Ctx).Errorf("json.Unmarshal flowId[%s] areaId[%d] subsystem[%s] api[%s], err:%+v", instReq.FlowId, l.AreaId, l.Subsystem, l.Api, err)
		return errors.New("反序列化失败")
	}
	if httpResp.Head.Code != 0 {
		logx.WithContext(l.Ctx).Errorf("l.PostAreaApi flowId[%s] areaId[%d] subsystem[%s] api[%s], httpResp:%+v", instReq.FlowId, l.AreaId, l.Subsystem, l.Api, httpResp)
		return errors.New(httpResp.Head.Msg)
	}

	logx.WithContext(l.Ctx).Debugf("DestroyInstance flowId[%s] areaId[%d] subsystem[%s] api[%s], vlanInfo:%+v", instReq.FlowId, l.AreaId, l.Subsystem, l.Api, info.String())
	return nil
}

// 设置超管
func (l *DisklessWebGateway) UpdateInstanceStatus(sessionId string, areaId int64, statusReq instance_types.UpdateInstanceStatusRequest) error {

	l.Api = "update_instance_status"
	l.AreaId = areaId

	request, err := json.Marshal(statusReq)
	if err != nil {
		logx.WithContext(l.Ctx).Errorf("[%s] json.Marshal err, areaId[%d] subsystem[%s] api[%s], statusReq :%s", sessionId, l.AreaId, l.Subsystem, l.Api, helper.ToJSON(statusReq))
		return err
	}

	info, err := l.PostAreaApi(statusReq.FlowID, string(request))
	if err != nil {
		return err
	}

	httpResp := new(instance_types.HTTPResponse)
	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		logx.WithContext(l.Ctx).Errorf("[%s] json.Unmarshal err, areaId[%d] subsystem[%s] api[%s], err:%+v", sessionId, l.AreaId, l.Subsystem, l.Api, err)
		return errors.New("反序列化失败")
	}
	if httpResp.Head.Code != 0 {
		logx.WithContext(l.Ctx).Errorf("[%s] l.PostAreaApi areaId[%d] subsystem[%s] api[%s], httpResp:%s", sessionId, l.AreaId, l.Subsystem, l.Api, helper.ToJSON(httpResp))
		return errors.New(httpResp.Head.Msg)
	}

	logx.WithContext(l.Ctx).Debugf("[%s] SetAdmin areaId[%d] subsystem[%s] api[%s], statusReq:%+v, vlanInfo:%+v", sessionId, l.AreaId, l.Subsystem, l.Api, statusReq, info.String())
	return nil
}

// 查询实例镜像进度
func (l *DisklessWebGateway) GetRestoreInstanceProcessing(flowId string, vmIds []int32) (mapTaskProcess map[int32]*image_service.TaskProcess, err error) {

	imageReq := &image_service.GetTaskProcessRequest{
		FlowId:   flowId,
		TaskIds:  vmIds,
		Resource: image_service.Resource_ManagerConsole,
	}

	request, err := protojson.Marshal(imageReq)
	if err != nil {
		logx.WithContext(l.Ctx).Errorf("flowId[%s] areaId[%d] subsystem[%s] api[%s], imageReq-->:%+v", flowId, l.AreaId, l.Subsystem, l.Api, imageReq.String())
		return nil, err
	}

	info, err := l.PostAreaApi(flowId, string(request))
	if err != nil {
		return nil, err
	}

	imageResp := new(image_service.GetTaskProcessHttpResponse)
	if err := protojson.Unmarshal([]byte(info.Response), imageResp); err != nil {
		logx.WithContext(l.Ctx).Errorf("protojson.Unmarshal flowId[%s] areaId[%d] subsystem[%s] api[%s], err:%+v", flowId, l.AreaId, l.Subsystem, l.Api, err)
		return nil, errors.New("反序列化失败")
	}
	if imageResp.Ret.GetCode() != 0 {
		logx.WithContext(l.Ctx).Errorf("l.PostAreaApi flowId[%s] areaId[%d] subsystem[%s] api[%s], vlanInfo:%+v", flowId, l.AreaId, l.Subsystem, l.Api, imageResp.String())
		return nil, errors.New(imageResp.Ret.GetMsg())
	}
	logx.WithContext(l.Ctx).Debugf("GetRestoreInstanceProcessing flowId[%s] areaId[%d] subsystem[%s] api[%s], vlanInfo:%+v", flowId, l.AreaId, l.Subsystem, l.Api, info.String())

	mapTaskProcess = imageResp.Body.Tasks

	return
}

// 设置超管
func (l *DisklessWebGateway) SetAdmin(areaId int64, setAdminReq *instance_types.SetAdminRequest) error {

	l.Api = "set_admin"
	l.AreaId = areaId

	request, err := json.Marshal(setAdminReq)
	if err != nil {
		logx.WithContext(l.Ctx).Errorf("json.Marshal flowId[%s] areaId[%d] subsystem[%s] api[%s], setAdminReq :%+v", setAdminReq.FlowID, l.AreaId, l.Subsystem, l.Api, setAdminReq)
		return err
	}
	info, err := l.PostAreaApi(setAdminReq.FlowID, string(request))
	if err != nil {
		return err
	}

	httpResp := new(instance_types.HTTPResponse)
	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		logx.WithContext(l.Ctx).Errorf("json.Unmarshal flowId[%s] areaId[%d] subsystem[%s] api[%s], err:%+v", setAdminReq.FlowID, l.AreaId, l.Subsystem, l.Api, err)
		return errors.New("反序列化失败")
	}
	if httpResp.Head.Code != 0 {
		logx.WithContext(l.Ctx).Errorf("l.PostAreaApi flowId[%s] areaId[%d] subsystem[%s] api[%s], httpResp:%+v", setAdminReq.FlowID, l.AreaId, l.Subsystem, l.Api, httpResp)
		return errors.New(httpResp.Head.Msg)
	}

	logx.WithContext(l.Ctx).Debugf("SetAdmin flowId[%s] areaId[%d] subsystem[%s] api[%s], instReq:%+v, info:%+v", setAdminReq.FlowID, l.AreaId, l.Subsystem, l.Api, setAdminReq, info.String())
	return nil
}

// 查询资源，或策略
func (l *DisklessWebGateway) SearchResource(sessionId string, areaId int64, listResReq *instance_scheduler.SearchResourceRequest) (list []*instance_scheduler.ResourceConfig, err error) {

	l.Api = "search_resource"
	l.AreaId = areaId

	request, err := json.Marshal(listResReq)
	if err != nil {
		logx.WithContext(l.Ctx).Errorf("[%s] SearchResource json.Marshal areaId[%d] api[%s], setAdminReq :%+v", sessionId, l.AreaId, l.Api, listResReq)
		return nil, err
	}
	info, err := l.PostAreaApi(sessionId, string(request))
	if err != nil {
		return nil, err
	}

	searchResResp := new(instance_scheduler.SearchResourceResponse)
	if err := protojson.Unmarshal([]byte(info.Response), searchResResp); err != nil {
		logx.WithContext(l.Ctx).Errorf("[%s] SearchResource json.Unmarshal err flowId[%s] areaId[%d] api[%s], searchResResp:%s, err:%+v", sessionId, l.AreaId, l.Api, searchResResp.String(), err)
		return nil, errors.New("反序列化失败")
	}
	if searchResResp.Ret.GetCode() != 0 {
		logx.WithContext(l.Ctx).Errorf("[%s] SearchResource http code err  areaId[%d] api[%s], searchResResp:%s", sessionId, l.AreaId, l.Api, searchResResp.String())
		return nil, errors.New(searchResResp.Ret.GetMsg())
	}

	list = searchResResp.Body.GetLists()

	logx.WithContext(l.Ctx).Debugf("[%s] SearchResource success areaId[%d] api[%s], instReq:%s, list:%+v", sessionId, l.AreaId, l.Api, helper.ToJSON(listResReq), list)
	return
}

// 新创建资源，或策略
func (l *DisklessWebGateway) NewResource(sessionId string, areaId int64, newResReq *instance_scheduler.NewResourceRequest) error {

	l.Api = "new_resource"
	l.AreaId = areaId

	request, err := json.Marshal(newResReq)
	if err != nil {
		logx.WithContext(l.Ctx).Errorf("[%s] NewResource json.Marshal areaId[%d] api[%s], setAdminReq :%+v", sessionId, l.AreaId, l.Api, newResReq)
		return err
	}
	info, err := l.PostAreaApi(sessionId, string(request))
	if err != nil {
		return err
	}

	httpResp := new(instance_types.HTTPResponse)
	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		logx.WithContext(l.Ctx).Errorf("[%s] NewResource json.Unmarshal err flowId[%s] areaId[%d] api[%s], err:%+v", sessionId, l.AreaId, l.Api, err)
		return errors.New("反序列化失败")
	}
	if httpResp.Head.Code != 0 {
		logx.WithContext(l.Ctx).Errorf("[%s] NewResource http code err  areaId[%d] api[%s], httpResp:%+v", sessionId, l.AreaId, l.Api, httpResp)
		return errors.New(httpResp.Head.Msg)
	}

	logx.WithContext(l.Ctx).Debugf("[%s] NewResource success areaId[%d] api[%s], instReq:%s, info:%s", sessionId, l.AreaId, l.Api, helper.ToJSON(newResReq), info.String())
	return nil
}

// 新创建资源，或策略
func (l *DisklessWebGateway) UpdateResource(sessionId string, areaId int64, updateResReq *UpdateResourceRequest) error {

	l.Api = "update_resource"
	l.AreaId = areaId

	request, err := json.Marshal(updateResReq)
	if err != nil {
		logx.WithContext(l.Ctx).Errorf("[%s] UpdateResource json.Marshal areaId[%d] api[%s], updateResReq :%+v", sessionId, l.AreaId, l.Api, updateResReq)
		return err
	}
	info, err := l.PostAreaApi(sessionId, string(request))
	if err != nil {
		return err
	}

	httpResp := new(instance_types.HTTPResponse)
	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		logx.WithContext(l.Ctx).Errorf("[%s] UpdateResource json.Unmarshal err flowId[%s] areaId[%d] api[%s], err:%+v", sessionId, l.AreaId, l.Api, err)
		return errors.New("反序列化失败")
	}
	if httpResp.Head.Code != 0 {
		logx.WithContext(l.Ctx).Errorf("[%s] UpdateResource http code err  areaId[%d] api[%s], httpResp:%+v", sessionId, l.AreaId, l.Api, httpResp)
		return errors.New(httpResp.Head.Msg)
	}

	logx.WithContext(l.Ctx).Debugf("[%s] UpdateResource success areaId[%d] api[%s], instReq:%s, info:%s", sessionId, l.AreaId, l.Api, helper.ToJSON(updateResReq), info.String())
	return nil
}

// 查询无盘的资源池
func (l *DisklessWebGateway) SearchPool(sessionId string, areaId int64, req *SearchPoolRequest) (body SearchPoolBody, err error) {

	l.Api = "search_pool"
	l.AreaId = areaId
	body = SearchPoolBody{}
	request, err := json.Marshal(req)
	if err != nil {
		logx.WithContext(l.Ctx).Errorf("[%s] SearchPool json.Marshal areaId[%d] api[%s], setAdminReq :%+v", sessionId, l.AreaId, l.Api, helper.ToJSON(req))
		return body, err
	}
	info, err := l.PostAreaApi(sessionId, string(request))
	if err != nil {
		return body, err
	}

	httpResp := new(SearchPoolResponse)
	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		logx.WithContext(l.Ctx).Errorf("[%s] SearchPool json.Unmarshal err flowId[%s] areaId[%d] api[%s], err:%+v", sessionId, l.AreaId, l.Api, err)
		return body, errors.New("反序列化失败")
	}
	if httpResp.Head.Code != 0 {
		logx.WithContext(l.Ctx).Errorf("[%s] SearchPool http code err  areaId[%d] api[%s], httpResp:%+v", sessionId, l.AreaId, l.Api, httpResp)
		return body, errors.New(httpResp.Head.Msg)
	}
	body = httpResp.Body
	logx.WithContext(l.Ctx).Debugf("[%s] SearchPool success areaId[%d] api[%s], req:%s, httpResp:%s", sessionId, l.AreaId, l.Api, helper.ToJSON(req), helper.ToJSON(httpResp))
	return body, nil
}

// 释放资源池中的实例
func (l *DisklessWebGateway) ReleasePoolItem(sessionId string, areaId int64, req *ReleasePoolItemRequest) error {

	l.Api = "release_pool_item"
	l.AreaId = areaId

	request, err := json.Marshal(req)
	if err != nil {
		logx.WithContext(l.Ctx).Errorf("[%s] ReleasePoolItem json.Marshal areaId[%d] api[%s], req :%s", sessionId, l.AreaId, l.Api, helper.ToJSON(req))
		return err
	}
	info, err := l.PostAreaApi(sessionId, string(request))
	if err != nil {
		return err
	}

	httpResp := new(instance_types.HTTPResponse)
	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		logx.WithContext(l.Ctx).Errorf("[%s] ReleasePoolItem json.Unmarshal err areaId[%d] api[%s], err:%+v", sessionId, l.AreaId, l.Api, err)
		return errors.New("反序列化失败")
	}
	if httpResp.Head.Code != 0 {
		logx.WithContext(l.Ctx).Errorf("[%s] ReleasePoolItem http code err  areaId[%d] api[%s], httpResp:%+v", sessionId, l.AreaId, l.Api, httpResp)
		return errors.New(httpResp.Head.Msg)
	}

	logx.WithContext(l.Ctx).Debugf("[%s] ReleasePoolItem success areaId[%d] api[%s], req:%s, info:%s", sessionId, l.AreaId, l.Api, helper.ToJSON(req), info.String())
	return nil
}

// 释放资源池中的实例
func (l *DisklessWebGateway) RebuidPool(sessionId string, areaId int64, req *RebuildPoolRequest) error {

	l.Api = "rebuild_pool"
	l.AreaId = areaId

	request, err := json.Marshal(req)
	if err != nil {
		logx.WithContext(l.Ctx).Errorf("[%s] RebuidPool json.Marshal areaId[%d] api[%s], req :%s", sessionId, l.AreaId, l.Api, helper.ToJSON(req))
		return err
	}
	info, err := l.PostAreaApi(sessionId, string(request))
	if err != nil {
		return err
	}

	httpResp := new(instance_types.HTTPResponse)
	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		logx.WithContext(l.Ctx).Errorf("[%s] RebuidPool json.Unmarshal err flowId[%s] areaId[%d] api[%s], err:%+v", sessionId, l.AreaId, l.Api, err)
		return errors.New("反序列化失败")
	}
	if httpResp.Head.Code != 0 {
		logx.WithContext(l.Ctx).Errorf("[%s] RebuidPool http code err  areaId[%d] api[%s], httpResp:%+v", sessionId, l.AreaId, l.Api, httpResp)
		return errors.New(httpResp.Head.Msg)
	}

	logx.WithContext(l.Ctx).Debugf("[%s] RebuidPool success areaId[%d] api[%s], req:%s, info:%s", sessionId, l.AreaId, l.Api, helper.ToJSON(req), info.String())
	return nil
}

// 更新资源池信息
func (l *DisklessWebGateway) UpdatePoolItem(sessionId string, areaId int64, req *UpdatePoolItemRequest) error {

	l.Api = "update_pool_item"
	l.AreaId = areaId

	request, err := json.Marshal(req)
	if err != nil {
		logx.WithContext(l.Ctx).Errorf("[%s] UpdatePoolItem json.Marshal areaId[%d] api[%s], req :%s", sessionId, l.AreaId, l.Api, helper.ToJSON(req))
		return err
	}
	info, err := l.PostAreaApi(sessionId, string(request))
	if err != nil {
		return err
	}

	httpResp := new(instance_types.HTTPResponse)
	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		logx.WithContext(l.Ctx).Errorf("[%s] UpdatePoolItem json.Unmarshal err areaId[%d] api[%s], err:%+v", sessionId, l.AreaId, l.Api, err)
		return errors.New("反序列化失败")
	}
	if httpResp.Head.Code != 0 {
		logx.WithContext(l.Ctx).Errorf("[%s] UpdatePoolItem http code err  areaId[%d] api[%s], httpResp:%+v", sessionId, l.AreaId, l.Api, httpResp)
		return errors.New(httpResp.Head.Msg)
	}

	logx.WithContext(l.Ctx).Debugf("[%s] UpdatePoolItem success areaId[%d] api[%s], req:%s, info:%s", sessionId, l.AreaId, l.Api, helper.ToJSON(req), info.String())
	return nil
}

// 制作镜像
func CreateImageFromAreaInstance(ctx context.Context, host string, flowId, imageId, name, osVersion, remark string, managerState int, bizId, areaId, vmId int64, flattenFlag int32) error {

	req := &CreateImageFromAreaInstanceRequest{
		FlowId:       flowId,
		ImageId:      imageId,
		Name:         name,
		OsVersion:    osVersion,
		ManagerState: int32(managerState),
		AreaId:       areaId,
		VmId:         vmId,
		Remark:       remark,
		FlattenFlag:  flattenFlag,
	}

	url := fmt.Sprintf("%s/v1/DisklessCloudImage/CreateImageFromAreaInstance", host)

	data, err := httpc.Do(context.Background(), http.MethodPost, url, req)
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] CreateImageFromAreaInstance httpc.Do err bizId[%d] host[%s] err:%+v", flowId, bizId, url, err)
		return errors.New("制作镜像失败")
	}

	body, _ := io.ReadAll(data.Body)
	resp := new(diskless_cloud_image.CreateImageFromAreaInstanceResponse)
	if err := protojson.Unmarshal([]byte(body), resp); err != nil {
		logx.WithContext(ctx).Errorf("[%s] CreateImageFromAreaInstance protojson.Unmarshal err. areaId[%d], err:%+v", flowId, areaId, err)
		return errors.New("反序列化失败")
	}

	if resp.Ret.GetCode() != 0 {
		logx.WithContext(ctx).Errorf("[%s] CreateImageFromAreaInstance err areaId[%d] vlanInfo:%+v", flowId, areaId, resp.String())
		return errors.New(resp.Ret.GetMsg())
	}

	logx.WithContext(ctx).Debugf("[%s] CreateImageFromAreaInstance succ areaId[%d] , resp:%+v", flowId, areaId, resp.String())

	return nil

}

type QueryBootSessionRequest struct {
	FlowId     string  `json:"flow_id"`
	InstanceId []int64 `json:"instance_id"`
}
type QueryBootSessionResp struct {
	Head HTTPCommonHead       `json:"ret"`
	Body QueryBootSessionBody `json:"body,omitempty"`
}
type QueryBootSessionBody struct {
	FlowId string              `json:"flow_id"`
	Total  int                 `json:"total"`
	List   []BootSessionDetail `json:"list"`
}
type BootSessionDetail struct {
	BootSessionId    string `json:"boot_session_id"`
	Mac              string `json:"mac"`
	InstanceId       int64  `json:"instance_id"`
	LastEventType    string `json:"last_event_type"`
	LastEventTime    string `json:"last_event_time"`
	SessionBeginTime string `json:"session_begin_time"`
	BootTime         int    `json:"boot_time"`
	Process          int    `json:"process"`
	Status           int    `json:"status"`
	Ip               string `json:"ip"`
	HostName         string `json:"host_name"`
	CreateTime       string `json:"create_time"`
	UpdateTime       string `json:"update_time"`
	ModifyTime       string `json:"modify_time"`
	LastEventDesc    string `json:"last_event_desc"`
}

// 释放资源池中的实例
func (l *DisklessWebGateway) QueryBootSession(sessionId string, areaId int64, req QueryBootSessionRequest) (resp QueryBootSessionBody, err error) {

	l.Api = "query_boot_session"
	l.AreaId = areaId

	request, err := json.Marshal(req)
	if err != nil {
		logx.WithContext(l.Ctx).Errorf("[%s] QueryBootSession areaId[%d] api[%s], req :%s", sessionId, l.AreaId, l.Api, helper.ToJSON(req))
		return resp, err
	}
	info, err := l.PostAreaApi(sessionId, string(request))
	if err != nil {
		return resp, err
	}

	httpResp := new(QueryBootSessionResp)
	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		logx.WithContext(l.Ctx).Errorf("[%s] QueryBootSession err flowId[%s] areaId[%d] api[%s], err:%+v", sessionId, l.AreaId, l.Api, err)
		return resp, errors.New("反序列化失败")
	}
	if httpResp.Head.Code != 0 {
		logx.WithContext(l.Ctx).Errorf("[%s] QueryBootSession http code err  areaId[%d] api[%s], httpResp:%+v", sessionId, l.AreaId, l.Api, httpResp)
		return resp, errors.New(httpResp.Head.Msg)
	}

	logx.WithContext(l.Ctx).Debugf("[%s] QueryBootSession success areaId[%d] api[%s], req:%s, info:%s", sessionId, l.AreaId, l.Api, helper.ToJSON(req), info.String())
	return httpResp.Body, nil
}
