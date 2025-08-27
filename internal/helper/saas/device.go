package saas

import (
	"cdp-admin-service/internal/helper"
	"cdp-admin-service/internal/types"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpc"
)

type CommonUserInfo struct {
	CommonBizName     string `json:"common_biz_name"`
	CommonAgentName   string `json:"common_agent_name"`
	CommonPrimaryName string `json:"common_primary_name"`
	CommonBizID       int64  `json:"common_biz_id"`
	CommonAgentID     int64  `json:"common_agent_id"`
	CommonPrimaryID   int64  `json:"common_primary_id"`
}

type TEsportsDeviceInfo struct {
	Id              int64     `orm:"column(device_id);auto" description:"规则Id" json:"device_id"`
	SetID           int64     `orm:"column(set_id)" description:"业务ID" json:"set_id"`
	AreaType        int64     `orm:"column(area_type)" description:"区域ID" json:"area_type"`
	AgentID         int64     `orm:"column(agent_id)" description:"二级代理商ID" json:"agent_id"`
	BizID           int64     `orm:"column(biz_id)" description:"租户ID" json:"biz_id"`
	RoomID          int64     `orm:"column(room_id)" description:"房间ID" json:"room_id"`
	RoomName        string    `orm:"column(room_name)" description:"房间名称" json:"room_name"`
	DeviceNumber    string    `orm:"column(device_number)" description:"设备号" json:"device_number"`
	MACAddress      string    `orm:"column(mac_address)" description:"MAC地址" json:"mac_address"`
	OrderNumber     string    `orm:"column(order_number)" description:"订单号" json:"order_number"`
	ResidualTime    int64     `orm:"column(residual_time)" description:"剩余时间" json:"residual_time"`
	LastActiveTime  time.Time `orm:"column(last_active_time);type(datetime);auto_now_add" description:"激活时间" json:"last_active_time"`
	RegisterTime    time.Time `orm:"column(register_time);type(datetime);auto_now_add" description:"注册时间" json:"register_time"`
	AllocatedTime   time.Time `orm:"column(allocated_time);type(datetime);auto_now_add" description:"分片时间" json:"allocated_time"`
	EndTime         time.Time `orm:"column(end_time);type(datetime);auto_now_add" description:"去激活时间" json:"end_time"`
	CurrentVersion  string    `orm:"column(current_version)" description:"当前版本" json:"current_version"`
	ExpectedVersion string    `orm:"column(expected_version)" description:"期待版本" json:"expected_version"`
	ConfigInfo      string    `orm:"column(config_info)" description:"配置信息" json:"config_info"`
	CreateBy        string    `orm:"column(create_by)" description:"创建者" json:"create_by"`
	UpdateBy        string    `orm:"column(update_by)" description:"更新者" json:"update_by"`
	Remark          string    `orm:"column(remark)" description:"备注" json:"remark"`
	State           int       `orm:"column(state)" description:"状态: 0-无效, 1-有效" json:"state"`
	UpdateState     int       `orm:"column(update_state)" description:"更新状态: 0-可更新, 1-不可更新, 2-强制更新" json:"update_state"`
	LogonToken      string    `orm:"column(logon_token)" description:"登录token,已登录时有效[cag]" json:"logon_token"`
	LogonUid        string    `orm:"column(logon_uid)" description:"登录用户id,已登录时有效[cag]" json:"logon_uid"`
	ActiveTs        int64     `orm:"column(active_ts)" description:"活跃时间，已登录时有效[cag]" json:"active_ts"`
	DeviceName      string    `orm:"column(device_name)" description:"设备名字" json:"device_name"`
	CreateTime      time.Time `orm:"column(create_time);type(datetime);auto_now_add" json:"create_time"`
	UpdateTime      time.Time `orm:"column(update_time);type(datetime);auto_now_add" json:"update_time"`
	ModifyTime      time.Time `orm:"column(modify_time);type(datetime);auto_now" json:"modify_time"`
}

type GetDeviceInfoListReq struct {
	Offset    int64    `json:"offset"`
	Limit     int64    `json:"limit"`
	PrimaryID int64    `json:"primary_id"`
	AgentID   int64    `json:"agent_id"`
	BizIDList []int64  `json:"biz_id_list"`
	CondList  []string `json:"cond_list"`
	Sorts     string   `json:"sorts"`
	Orders    string   `json:"orders"`
}

type DeviceInfo struct {
	TEsportsDeviceInfo
	SetName         string `json:"set_name"`
	DeviceSetNumber string `json:"device_set_number"`
	CommonUserInfo
}

type GetDeviceInfoList struct {
	Total int64        `json:"total"`
	List  []DeviceInfo `json:"list"`
}

// 查询瘦终端设备信息列表
type EsportDeviceInfoListResp struct {
	Code int               `json:"code"`
	Msg  string            `json:"msg"`
	Data GetDeviceInfoList `json:"data"`
}

// 查询瘦终端设备信息
type GetEsportDeviceInfoResp struct {
	Code int        `json:"code"`
	Msg  string     `json:"msg"`
	Data DeviceInfo `json:"data"`
}

// 更新瘦终端设备信息
type UpdateDeviceInfoReq struct {
	BizID       int64  `json:"biz_id,omitempty"`        // 租户ID
	UpdateState int    `json:"update_state,,omitempty"` // 更新状态
	MACAddress  string `json:"mac_address,omitempty"`   // mac 地址
	DeviceName  string `json:"device_name,omitempty"`   // 更新名字
	UpdateBy    string `json:"update_by,omitempty"`     // 更新人
	Header      string `header:"X-System"`
}

// 删除瘦终端设备信息
type DeleteDeviceInfoResp struct {
	Code int                `json:"code"`
	Msg  string             `json:"msg"`
	Data TEsportsDeviceInfo `json:"data"`
}

// 创建瘦设备信息
type CreateDeviceInfoReq struct {
	MACAddress  string `json:"mac_address"`
	BizID       int64  `json:"biz_id"`
	UpdateState int    `json:"update_state"`
	DeviceName  string `json:"device_name,omitempty"`
	UpdateBy    string `json:"update_by,omitempty"`
	Header      string `header:"X-System"`
}

// 创建瘦终端设备信息
type CreateDeviceInfoResp struct {
	Code int                `json:"code"`
	Msg  string             `json:"msg"`
	Data TEsportsDeviceInfo `json:"data"`
}

// 查询瘦终端的设备列表
func GetEsportDeviceInfoList(ctx context.Context, sessionId, host string, req *types.GetCloudBoxListReq) (list *GetDeviceInfoList, err error) {

	url := fmt.Sprintf("%s/device/getDeviceInfoList", host)

	httpResp, err := httpc.Do(context.Background(), http.MethodPost, url, req)
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s]httpc.Do err host[%s] err:%+v", sessionId, url, err)
		return nil, errors.New("查询终端的设备列表失败")
	}
	if httpResp.StatusCode != 200 {
		logx.WithContext(ctx).Errorf("[%s] httpc.Do err url:%s, req:%s, resp:%s", sessionId, url, helper.ToJSON(req), helper.ToJSON(httpResp))
		return nil, fmt.Errorf("查询终端的设备列表失败, httpResp.StatusCode:%d", httpResp.StatusCode)
	}

	saasResp := new(EsportDeviceInfoListResp)
	body, _ := io.ReadAll(httpResp.Body)

	if err := json.Unmarshal(body, saasResp); err != nil {
		logx.WithContext(ctx).Errorf("[%s] json.Unmarshal err host[%s],err:%+v", sessionId, url, err)
		return nil, errors.New("反序列化失败")
	}
	if saasResp.Code != 0 {
		logx.WithContext(ctx).Errorf("[%s] saasResp err host[%s], saasResp:%+v", sessionId, url, helper.ToJSON(saasResp))
		return nil, errors.New(saasResp.Msg)
	}

	list = &saasResp.Data

	return
}

// 查询单体瘦终端的设备 // http://10.86.0.101:28246/device/getDeviceInfo/6537"
func GetEsportDeviceInfo(ctx context.Context, host string, sessionId string, deviceId int64) (info *DeviceInfo, err error) {

	url := fmt.Sprintf("%s/device/getDeviceInfo/%d", host, deviceId)

	httpResp, err := httpc.Do(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] httpc.Do err,  host[%s] err :%+v", sessionId, url, err)
		return nil, errors.New("查询终端的设备失败")
	}
	if httpResp.StatusCode != 200 {
		logx.WithContext(ctx).Errorf("[%s] httpc.Do err url:%s, deviceId:%d, resp:%s", sessionId, url, deviceId, helper.ToJSON(httpResp))
		return nil, fmt.Errorf("查询终端的设备失败, StatusCode:%d", httpResp.StatusCode)
	}

	saasResp := new(GetEsportDeviceInfoResp)
	body, _ := io.ReadAll(httpResp.Body)

	if err := json.Unmarshal(body, saasResp); err != nil {
		logx.WithContext(ctx).Errorf("[%s] json.Unmarshal err. host[%s],err:%+v", sessionId, url, err)
		return nil, errors.New("反序列化失败")
	}
	if saasResp.Code != 0 {
		logx.WithContext(ctx).Errorf("[%s] saasResp err. host[%s],saasResp:%+v", sessionId, url, helper.ToJSON(saasResp))
		return nil, errors.New(saasResp.Msg)
	}

	info = &saasResp.Data

	return
}

// 更新瘦终端的设备
func UpdateEsportDeviceInfo(ctx context.Context, host, sessionId string, deviceId int64, info *UpdateDeviceInfoReq) (err error) {

	url := fmt.Sprintf("%s/device/updateDeviceInfo/%d", host, deviceId)
	info.Header = "diskless-aggregator"
	httpResp, err := httpc.Do(context.Background(), http.MethodPut, url, info)
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] UpdateEsportDeviceInfo httpc.Do err host[%s],req:%s, err:%+v", sessionId, url, helper.ToJSON(info), err)
		return err
	}

	if httpResp.StatusCode != 200 {
		logx.WithContext(ctx).Errorf("[%s] UpdateEsportDeviceInfo httpc.Do err url:%s, deviceId:%d, resp:%s", sessionId, url, deviceId, helper.ToJSON(httpResp))
		return fmt.Errorf("更新瘦终端的设备失败, StatusCode:%d", httpResp.StatusCode)
	}

	body, _ := io.ReadAll(httpResp.Body)
	saasResp := new(SaasCommonResp)
	if err := json.Unmarshal(body, saasResp); err != nil {
		logx.WithContext(ctx).Errorf("[%s] UpdateEsportDeviceInfo json.Unmarshal err deviceId[%d] host[%s],err:%+v", sessionId, deviceId, url, err)
		return errors.New("反序列化失败")
	}
	if saasResp.Code != 0 {
		logx.WithContext(ctx).Errorf("[%s] UpdateEsportDeviceInfo saasResp err deviceId[%d] host[%s],saasResp:%s", sessionId, deviceId, url, helper.ToJSON(saasResp))
		return errors.New(saasResp.Msg)
	}

	logx.WithContext(ctx).Infof("[%s] httpc.Do url:%s, req:%s, saasResp	:%s", sessionId, url, helper.ToJSON(info), helper.ToJSON(httpResp))
	return
}

// 删除瘦终端的设备
func DeleteEsportDeviceInfo(ctx context.Context, sessionId, host string, deviceId int64) (err error) {

	url := fmt.Sprintf("%s/device/deleteDeviceInfo/%d", host, deviceId)

	httpResp, err := httpc.Do(context.Background(), http.MethodDelete, url, nil)
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] DeleteEsportDeviceInfo httpc.Do err host[%s],deviceId:%d, err:%+v", sessionId, url, deviceId, err)
		return err
	}
	if httpResp.StatusCode != 200 {
		logx.WithContext(ctx).Errorf("[%s] DeleteEsportDeviceInfo httpc.Do err url:%s, deviceId:%d, resp:%+v", sessionId, url, deviceId, helper.ToJSON(httpResp))
		return fmt.Errorf("删除瘦终端的设备失败,StatusCode:%d", httpResp.StatusCode)
	}

	body, _ := io.ReadAll(httpResp.Body)
	saasResp := new(DeleteDeviceInfoResp)
	if err := json.Unmarshal(body, saasResp); err != nil {
		logx.WithContext(ctx).Errorf("[%s] DeleteEsportDeviceInfo json.Unmarshal err deviceId:%d host[%s],err:%+v", sessionId, deviceId, url, err)
		return errors.New("反序列化失败")
	}
	if saasResp.Code != 0 {
		logx.WithContext(ctx).Errorf("[%s] DeleteEsportDeviceInfo saasResp err deviceId:%d host[%s],saasResp:%s", sessionId, deviceId, url, helper.ToJSON(saasResp))
		return errors.New(saasResp.Msg)
	}

	//logx.WithContext(ctx).Infof("httpc.Do url:%s, resp:%+v", url, resp)
	return
}

// 释放瘦终端的设备
func ReleaseEsportDeviceInfo(ctx context.Context, sessionId, host string, deviceId int64) (err error) {

	url := fmt.Sprintf("%s/device/releaseDevice/%d", host, deviceId)
	req := Request{
		Header: "diskless-aggregator",
	}

	httpResp, err := httpc.Do(context.Background(), http.MethodPut, url, req)
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] ReleaseEsportDeviceInfo httpc.Do err host[%s],deviceId:%d, err:%+v", sessionId, url, deviceId, err)
		return err
	}
	if httpResp.StatusCode != 200 {
		logx.WithContext(ctx).Errorf("[%s] ReleaseEsportDeviceInfo httpc.Do err url:%s, deviceId:%d, resp:%s", sessionId, url, deviceId, helper.ToJSON(httpResp))
		return fmt.Errorf("释放瘦终端的设备失败,StatusCode:%d", httpResp.StatusCode)
	}

	body, _ := io.ReadAll(httpResp.Body)
	saasResp := new(DeleteDeviceInfoResp)
	if err := json.Unmarshal(body, saasResp); err != nil {
		logx.WithContext(ctx).Errorf("[%s] ReleaseEsportDeviceInfo json.Unmarshal err deviceId[%d] host[%s],err:%+v", sessionId, deviceId, url, err)
		return errors.New("反序列化失败")
	}
	if saasResp.Code != 0 {
		logx.WithContext(ctx).Errorf("[%s] ReleaseEsportDeviceInfo saasResp err deviceId[%d] host[%s],saasResp:%s", sessionId, deviceId, url, helper.ToJSON(saasResp))
		return errors.New(saasResp.Msg)
	}

	//logx.WithContext(ctx).Infof("httpc.Do url:%s, resp:%+v", url, resp)
	return
}

// 创建瘦终端的设备
func CreateEsportDeviceInfo(ctx context.Context, sessionId string, host string, req CreateDeviceInfoReq) (info *TEsportsDeviceInfo, err error) {

	url := fmt.Sprintf("%s/device/createDeviceInfo", host)
	req.Header = "diskless-aggregator"
	httpResp, err := httpc.Do(context.Background(), http.MethodPost, url, req)
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] CreateEsportDeviceInfo httpc.Do err host[%s],req:%s, err:%+v", sessionId, url, helper.ToJSON(req), err)
		return nil, err
	}
	if httpResp.StatusCode != 200 {
		logx.WithContext(ctx).Errorf("[%s] CreateEsportDeviceInfo httpc.Do err url:%s, req:%s, resp:%s", sessionId, url, helper.ToJSON(req), helper.ToJSON(httpResp))
		return nil, fmt.Errorf("创建瘦终端的设备失败, StatusCode:%d", httpResp.StatusCode)
	}

	body, _ := io.ReadAll(httpResp.Body)
	saasResp := new(CreateDeviceInfoResp)
	if err := json.Unmarshal(body, saasResp); err != nil {
		logx.WithContext(ctx).Errorf("[%s]CreateEsportDeviceInfo  json.Unmarshal err  host[%s], req:%s, resp:%s", sessionId, url, helper.ToJSON(req), helper.ToJSON(httpResp))
		return nil, errors.New("反序列化失败")
	}

	if saasResp.Code != 0 {
		logx.WithContext(ctx).Errorf("[%s]CreateEsportDeviceInfo   saasResp error  host[%s], req:%s, resp:%s", sessionId, url, helper.ToJSON(req), helper.ToJSON(httpResp))
		return nil, errors.New(saasResp.Msg)
	}
	info = &saasResp.Data
	//logx.WithContext(ctx).Infof("httpc.Do url:%s, resp:%+v", url, resp)
	return
}
