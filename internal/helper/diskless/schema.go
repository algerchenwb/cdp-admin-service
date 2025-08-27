package diskless

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"cdp-admin-service/internal/proto/diskless_cloud_image"
	instance_types "cdp-admin-service/internal/proto/instance_service/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpc"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
	"google.golang.org/protobuf/encoding/protojson"
)

type CreateSchemaReq struct {
	FlowId       string `json:"flow_id"`
	Name         string `json:"name"`
	OsImageId    string `json:"os_image_id"`
	OsSnapshotId string `json:"os_snapshot_id"`
	StorageType  int32  `json:"storage_type"`
	WrbackType   int32  `json:"wrback_type"`
	BootpnpPath  string `json:"bootpnp_path"`
}

type CreateSchemaResp struct {
	SchemeId int64 `json:"scheme_id"`
}

func (l *DisklessWebGateway) CreateSchema(AreaId int64, flowId string, req map[string]any) (*CreateSchemaResp, error) {

	l.Api = "create_scheme"
	l.AreaId = AreaId

	request, err := json.Marshal(req)
	if err != nil {
		logx.WithContext(l.Ctx).Errorf("json.Marshal flowId[%s] areaId[%d] api[%s], req:%+v", flowId, l.AreaId, l.Api, req)
		return nil, err
	}
	info, err := l.PostAreaApi(flowId, string(request))
	if err != nil {
		return nil, err
	}
	resp := new(CreateSchemaResp)
	httpResp := &instance_types.HTTPResponse{
		Body: resp,
	}
	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		logx.WithContext(l.Ctx).Errorf("json.Unmarshal flowId[%s] areaId[%d] api[%s], err:%+v", flowId, l.AreaId, l.Api, err)
		return nil, err
	}
	if httpResp.Head.Code != 0 {
		logx.WithContext(l.Ctx).Errorf("l.PostAreaApi flowId[%s] areaId[%d] api[%s], httpResp:%+v", flowId, l.AreaId, l.Api, httpResp)
		return nil, errors.New(httpResp.Head.Msg)
	}

	logx.WithContext(l.Ctx).Debugf("CreateSchema flowId[%s] areaId[%d] api[%s], req:%+v, resp:%+v", flowId, l.AreaId, l.Api, req, info.String())
	return resp, nil

}

type UpdateSchemaReq struct {
	FlowId       string `json:"flow_id"`
	Name         string `json:"name"`
	OsImageId    string `json:"os_image_id"`
	OsSnapshotId string `json:"os_snapshot_id"`
	SchemeId     int64  `json:"scheme_id"`
	StorageType  int32  `json:"storage_type"`
	WrbackType   int32  `json:"wrback_type"`
	BootpnpPath  string `json:"bootpnp_path"`
}

type UpdateSchemaResp struct {
	SchemeId int64 `json:"scheme_id"`
}

func (l *DisklessWebGateway) UpdateSchema(AreaId int64, flowId string, req map[string]any) (*UpdateSchemaResp, error) {

	l.Api = "update_scheme"
	l.AreaId = AreaId
	request, err := json.Marshal(req)
	if err != nil {
		logx.WithContext(l.Ctx).Errorf("json.Marshal flowId[%s] areaId[%d] api[%s], req:%+v", flowId, l.AreaId, l.Api, req)
		return nil, err
	}
	info, err := l.PostAreaApi(flowId, string(request))
	if err != nil {
		return nil, err
	}
	resp := new(UpdateSchemaResp)
	httpResp := &instance_types.HTTPResponse{
		Body: resp,
	}
	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		logx.WithContext(l.Ctx).Errorf("json.Unmarshal flowId[%s] areaId[%d] api[%s], err:%+v", flowId, l.AreaId, l.Api, err)
		return nil, err
	}
	if httpResp.Head.Code != 0 {
		logx.WithContext(l.Ctx).Errorf("l.PostAreaApi flowId[%s] areaId[%d] api[%s], httpResp:%+v", flowId, l.AreaId, l.Api, httpResp)
		return nil, errors.New(httpResp.Head.Msg)
	}

	logx.WithContext(l.Ctx).Debugf("UpdateSchema flowId[%s] areaId[%d] api[%s], req:%+v, resp:%+v", flowId, l.AreaId, l.Api, req, info.String())
	return resp, nil

}

type ListSchemesResp struct {
	Schemes []Scheme `json:"schemes"`
	Total   int64    `json:"total"`
}

type Scheme struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	OsImageId  string `json:"os_image_id"`
	State      int32  `json:"state"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
	ModifyTime string `json:"modify_time"`
}



func (l *DisklessWebGateway) ListSchemes(AreaId int64, flowId string, offset int64, length int64, order string, sortBy string, ids string) (*ListSchemesResp, error) {
	l.Api = "list_schemes"
	l.AreaId = AreaId
	request := fmt.Sprintf("offset=%d&length=%d&order=%s&sort_by=%s", offset, length, order, sortBy)
	if ids != "" {
		request += fmt.Sprintf("&ids=%s", ids)
	}

	logx.WithContext(l.Ctx).Infof("l.ListSchemes flowId[%s] areaId[%d] api[%s], request:%s", flowId, l.AreaId, l.Api, request)

	info, err := l.GetAreaApi(flowId, request)
	if err != nil {
		return nil, err
	}

	logx.WithContext(l.Ctx).Infof("l.ListSchemes flowId[%s] areaId[%d] api[%s], Response:%s", flowId, l.AreaId, l.Api, info.Response)
	resp := new(ListSchemesResp)
	httpResp := &instance_types.HTTPResponse{
		Body: resp,
	}
	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		logx.WithContext(l.Ctx).Errorf("json.Unmarshal flowId[%s] areaId[%d] api[%s], err:%+v", flowId, l.AreaId, l.Api, err)
		return nil, errors.New("反序列化失败")
	}
	if httpResp.Head.Code != 0 {
		logx.WithContext(l.Ctx).Errorf("l.ListSchemes flowId[%s] areaId[%d] api[%s], httpResp:%+v", flowId, l.AreaId, l.Api, httpResp)
		return nil, errors.New(httpResp.Head.Msg)
	}

	logx.WithContext(l.Ctx).Debugf("ListSchemes flowId[%s] areaId[%d] api[%s], httpResp:%+v", flowId, l.AreaId, l.Api, httpResp)
	return resp, nil
}

type CreateScriptReq struct {
	FlowId     string `json:"flow_id"`
	Name       string `json:"name"`
	Script     string `json:"script"`
	ScriptType int32  `json:"script_type"`
}

type CreateScriptResp struct {
	FlowId string `json:"flow_id"`
	Script Script `json:"script"`
}

type Script struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Script        string `json:"script"`
	ScriptType    int32  `json:"script_type"`
	ExecutePath   string `json:"execute_path"`
	Desc          string `json:"desc"`
	ManagerStatus int32  `json:"manager_status"`
	Status        int32  `json:"status"`
	CreateTime    string `json:"create_time"`
	UpdateTime    string `json:"update_time"`
	ModifyTime    string `json:"modify_time"`
}

func (l *DisklessWebGateway) CreateScript(AreaId int64, flowId string, req CreateScriptReq) (*CreateScriptResp, error) {
	l.Api = "create_script"
	l.AreaId = AreaId
	request, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	info, err := l.PostAreaApi(flowId, string(request))
	if err != nil {
		return nil, err
	}
	resp := new(CreateScriptResp)
	httpResp := &instance_types.HTTPResponse{
		Body: resp,
	}
	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		logx.WithContext(l.Ctx).Errorf("json.Unmarshal flowId[%s] areaId[%d] api[%s], err:%+v", flowId, l.AreaId, l.Api, err)
		return nil, err
	}
	if httpResp.Head.Code != 0 {
		logx.WithContext(l.Ctx).Errorf("l.PostAreaApi flowId[%s] areaId[%d] api[%s], httpResp:%+v", flowId, l.AreaId, l.Api, httpResp)
		return nil, errors.New(httpResp.Head.Msg)
	}

	logx.WithContext(l.Ctx).Debugf("CreateScript flowId[%s] areaId[%d] api[%s], req:%+v, resp:%+v", flowId, l.AreaId, l.Api, req, info.String())
	return resp, nil
}

type ResetScriptBindReq struct {
	FlowId    string  `json:"flow_id"`
	SchemeIds []int64 `json:"scheme_ids"`
	ScriptId  string  `json:"script_id"`
}

type ResetScriptBindResp struct {
	FlowId string  `json:"flow_id"`
	Rets   []int64 `json:"rets"`
}

func (l *DisklessWebGateway) ResetScriptBind(AreaId int64, flowId string, req ResetScriptBindReq) (*ResetScriptBindResp, error) {
	l.Api = "reset_script_bind"
	l.AreaId = AreaId
	request, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	info, err := l.PostAreaApi(flowId, string(request))
	if err != nil {
		return nil, err
	}
	resp := new(ResetScriptBindResp)
	httpResp := &instance_types.HTTPResponse{
		Body: resp,
	}
	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		logx.WithContext(l.Ctx).Errorf("json.Unmarshal flowId[%s] areaId[%d] api[%s], err:%+v", flowId, l.AreaId, l.Api, err)
		return nil, err
	}
	if httpResp.Head.Code != 0 {
		logx.WithContext(l.Ctx).Errorf("l.PostAreaApi flowId[%s] areaId[%d] api[%s], httpResp:%+v", flowId, l.AreaId, l.Api, httpResp)
		return nil, errors.New(httpResp.Head.Msg)
	}

	logx.WithContext(l.Ctx).Debugf("ResetScriptBind flowId[%s] areaId[%d] api[%s], req:%+v, resp:%+v", flowId, l.AreaId, l.Api, req, info.String())
	return resp, nil
}

type UpdateScriptReq struct {
	FlowId        string `json:"flow_id"`
	Name          string `json:"name"`
	Script        string `json:"script"`
	ScriptType    int32  `json:"script_type"`
	Id            string `json:"id"`
	ManagerStatus int32  `json:"manager_status"`
}

type UpdateScriptResp struct {
	FlowId string `json:"flow_id"`
}

func (l *DisklessWebGateway) UpdateScript(AreaId int64, flowId string, req UpdateScriptReq) (*UpdateScriptResp, error) {
	l.Api = "update_script"
	l.AreaId = AreaId
	request, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	info, err := l.PostAreaApi(flowId, string(request))
	if err != nil {
		return nil, err
	}
	resp := new(UpdateScriptResp)
	httpResp := &instance_types.HTTPResponse{
		Body: resp,
	}
	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		logx.WithContext(l.Ctx).Errorf("json.Unmarshal flowId[%s] areaId[%d] api[%s], err:%+v", flowId, l.AreaId, l.Api, err)
		return nil, err
	}
	if httpResp.Head.Code != 0 {
		logx.WithContext(l.Ctx).Errorf("l.PostAreaApi flowId[%s] areaId[%d] api[%s], httpResp:%+v", flowId, l.AreaId, l.Api, httpResp)
		return nil, errors.New(httpResp.Head.Msg)
	}

	logx.WithContext(l.Ctx).Debugf("UpdateScript flowId[%s] areaId[%d] api[%s], req:%+v, resp:%+v", flowId, l.AreaId, l.Api, req, info.String())
	return resp, nil
}

type QuerySchemeScriptReq struct {
	FlowId   string   `json:"flow_id"`
	Offset   int64    `json:"offset"`
	Limit    int64    `json:"limit"`
	Orders   string   `json:"orders"`
	Sorts    string   `json:"sorts"`
	CondList []string `json:"cond_list"`
}

type QuerySchemeScriptResp struct {
	FlowId string         `json:"flow_id"`
	Total  int64          `json:"total"`
	List   []SchemeScript `json:"list"`
}

type SchemeScript struct {
	SchemeId string `json:"scheme_id"`
	ScriptId string `json:"script_id"`
}

type ListScriptReq struct {
	FlowId   string   `json:"flow_id"`
	Offset   int64    `json:"offset"`
	Limit    int64    `json:"limit"`
	Orders   string   `json:"orders"`
	Sorts    string   `json:"sorts"`
	CondList []string `json:"cond_list"`
}

type ListScriptResp struct {
	FlowId string       `json:"flow_id"`
	Total  int64        `json:"total"`
	List   []ScriptInfo `json:"list"`
}

type ScriptInfo struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Script        string `json:"script"`
	ScriptType    int32  `json:"script_type"`
	ExecutePath   string `json:"execute_path"`
	Desc          string `json:"desc"`
	ManagerStatus int32  `json:"manager_status"`
	Status        int32  `json:"status"`
	CreateTime    string `json:"create_time"`
	UpdateTime    string `json:"update_time"`
	ModifyTime    string `json:"modify_time"`
}

func (l *DisklessWebGateway) ListScript(AreaId int64, flowId string, req ListScriptReq) (*ListScriptResp, error) {
	l.Api = "list_script"
	l.AreaId = AreaId
	request, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	info, err := l.PostAreaApi(flowId, string(request))
	if err != nil {
		return nil, err
	}
	resp := new(ListScriptResp)
	httpResp := &instance_types.HTTPResponse{
		Body: resp,
	}
	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		logx.WithContext(l.Ctx).Errorf("json.Unmarshal flowId[%s] areaId[%d] api[%s], err:%+v", flowId, l.AreaId, l.Api, err)
		return nil, err
	}
	if httpResp.Head.Code != 0 {
		logx.WithContext(l.Ctx).Errorf("l.PostAreaApi flowId[%s] areaId[%d] api[%s], httpResp:%+v", flowId, l.AreaId, l.Api, httpResp)
		return nil, errors.New(httpResp.Head.Msg)
	}

	logx.WithContext(l.Ctx).Debugf("ListScript flowId[%s] areaId[%d] api[%s], req:%+v, resp:%+v", flowId, l.AreaId, l.Api, req, info.String())
	return resp, nil
}

type DeleteScriptReq struct {
	FlowId string `json:"flow_id"`
	Ids    string `json:"ids"`
}

type DeleteScriptResp struct {
	FlowId string `json:"flow_id"`
	Rets   []struct {
		Code int64  `json:"code"`
		Msg  string `json:"msg"`
	} `json:"rets"`
}

func (l *DisklessWebGateway) DeleteScript(AreaId int64, flowId string, req DeleteScriptReq) (*DeleteScriptResp, error) {
	l.Api = "delete_script"
	l.AreaId = AreaId
	request, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	info, err := l.PostAreaApi(flowId, string(request))
	if err != nil {
		return nil, err
	}
	resp := new(DeleteScriptResp)
	httpResp := &instance_types.HTTPResponse{
		Body: resp,
	}
	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		logx.WithContext(l.Ctx).Errorf("json.Unmarshal flowId[%s] areaId[%d] api[%s], err:%+v", flowId, l.AreaId, l.Api, err)
		return nil, err
	}
	if httpResp.Head.Code != 0 {
		logx.WithContext(l.Ctx).Errorf("l.PostAreaApi flowId[%s] areaId[%d] api[%s], httpResp:%+v", flowId, l.AreaId, l.Api, httpResp)
		return nil, errors.New(httpResp.Head.Msg)
	}
	for _, ret := range resp.Rets {
		if ret.Code != 0 {
			logx.WithContext(l.Ctx).Errorf("DeleteScript flowId[%s] areaId[%d] api[%s], req:%+v, ret:%+v", flowId, l.AreaId, l.Api, req, ret)
		}
	}

	logx.WithContext(l.Ctx).Debugf("DeleteScript flowId[%s] areaId[%d] api[%s], req:%+v, resp:%+v", flowId, l.AreaId, l.Api, req, info.String())
	return resp, nil
}

type DeleteSchemeReq struct {
	FlowId   string `json:"flow_id"`
	SchemeId int64  `json:"scheme_id"`
}

type DeleteSchemeResp struct {
}

func (l *DisklessWebGateway) DeleteScheme(AreaId int64, flowId string, req DeleteSchemeReq) (*DeleteSchemeResp, error) {
	l.Api = "delete_scheme"
	l.AreaId = AreaId
	request, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	info, err := l.PostAreaApi(flowId, string(request))
	if err != nil {
		return nil, err
	}
	resp := new(DeleteSchemeResp)
	httpResp := &instance_types.HTTPResponse{
		Body: resp,
	}
	if err := json.Unmarshal([]byte(info.Response), httpResp); err != nil {
		logx.WithContext(l.Ctx).Errorf("json.Unmarshal flowId[%s] areaId[%d] api[%s], err:%+v", flowId, l.AreaId, l.Api, err)
		return nil, err
	}
	if httpResp.Head.Code != 0 {
		logx.WithContext(l.Ctx).Errorf("l.PostAreaApi flowId[%s] areaId[%d] api[%s], httpResp:%+v", flowId, l.AreaId, l.Api, httpResp)
		return nil, errors.New(httpResp.Head.Msg)
	}

	logx.WithContext(l.Ctx).Debugf("DeleteScheme flowId[%s] areaId[%d] api[%s], req:%+v, resp:%+v", flowId, l.AreaId, l.Api, req, info.String())
	return resp, nil
}

type ResetSchemeImageBindReq struct {
	FlowId   string       `json:"flow_id"`
	SchemeId int64        `json:"scheme_id"`
	Images   []*ImageInfo `json:"images"`
}
type ImageInfo struct {
	MountPoint string `json:"mount_point"`
	ImageId    string `json:"image_id"`
	MntType    int32  `json:"mnt_type"`
	SchemeId   int64  `json:"scheme_id"`
}

func (l *DisklessWebGateway) ResetSchemeImageBind(AreaId int64, flowId string, req ResetSchemeImageBindReq) error {

	url := fmt.Sprintf("%s/v1/DisklessCloudWeb/DirectCallAreaApi/iaas/reset_scheme_image_bind", l.SvcCtx.Config.OutSide.DisklessEdge1)

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(gopublic.ToJSON(req))))
	if err != nil {
		logx.WithContext(l.Ctx).Errorf("http.NewRequest err bizId[%d] host[%s] err:%+v", l.AreaId, url, err)
		return errors.New("创建请求失败")
	}
	request.Header.Set("area", fmt.Sprint(AreaId))

	data, err := httpc.DoRequest(request)
	logx.WithContext(l.Ctx).Infof("ResetSchemeImageBind flowId[%s] url[%s], req:%s, resp:%+v", flowId, url, gopublic.ToJSON(req), data)
	if err != nil {
		logx.WithContext(l.Ctx).Errorf("httpc.Do err bizId[%d] host[%s] err:%+v", l.AreaId, url, err)
		return errors.New("重置方案镜像绑定失败")
	}

	body, _ := io.ReadAll(data.Body)
	resp := new(diskless_cloud_image.CreateImageFromAreaInstanceResponse)
	if err := protojson.Unmarshal([]byte(body), resp); err != nil {
		logx.WithContext(l.Ctx).Errorf("json.Unmarshal flowId[%s] areaId[%d], err:%+v", flowId, l.AreaId, err)
		return errors.New("反序列化失败")
	}

	if resp.Ret.GetCode() != 0 {
		logx.WithContext(l.Ctx).Errorf("CreateImageFromAreaInstance err. flowId[%s] areaId[%d] vlanInfo:%+v", flowId, l.AreaId, resp.String())
		return errors.New(resp.Ret.GetMsg())
	}

	logx.WithContext(l.Ctx).Debugf("ResetSchemeImageBind succ flowId[%s] areaId[%d] , resp:%+v", flowId, l.AreaId, resp.String())

	return nil

}
