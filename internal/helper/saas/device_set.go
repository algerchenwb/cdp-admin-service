package saas

import (
	"cdp-admin-service/internal/helper"
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

// 查询座位的列表
type GetBizAllDeviceSetReq struct {
	AgentID  int64    `json:"agent_id"`
	BizID    int64    `json:"biz_id"`
	Offset   int64    `json:"offset"`
	Limit    int64    `json:"limit"`
	CondList []string `json:"cond_list"`
	Sorts    string   `json:"sorts"`
	Orders   string   `json:"orders"`
	Header   string   `header:"X-System"`
}

type BizDeviceSetInfo struct {
	DeviceSetID     int64     `json:"device_set_id"`
	DeviceSetNumber string    `json:"device_set_number"`
	NetBarAddr      string    `json:"net_bar_addr"`
	NetBarGateway   string    `json:"net_bar_gateway"`
	NetBarMask      string    `json:"net_bar_mask"`
	DeviceNumber    string    `json:"device_number"`
	State           int64     `json:"state"`
	CreateBy        string    `json:"create_by"`
	CreateTime      time.Time `json:"create_time"`
	SpecNameList    []string  `json:"spec_name_list"`
	SpecIDList      string    `json:"spec_id_list"`
	PoolID          int       `json:"pool_id"`
	RoomName        string    `json:"room_name,omitempty"`
	DeviceMac       string    `json:"device_mac,omitempty"`
	DeviceName      string    `json:"device_name,omitempty"`
	ConfigInfo      string    `json:"config_info,omitempty"`
	AdminState      int       `json:"admin_state,omitempty"`
}

type GetBizAllDeviceSet struct {
	Total int64              `json:"total"`
	List  []BizDeviceSetInfo `json:"list"`
}

type GetBizAllDeviceSetResp struct {
	Code int                `json:"code"`
	Msg  string             `json:"msg"`
	Data GetBizAllDeviceSet `json:"data"`
}

// 更新座位信息
type UpdateDeviceSetReq struct {
	DeviceSetID     int64  `json:"device_set_id"`
	DeviceSetNumber string `json:"device_set_number"`
	NetBarAddr      string `json:"net_bar_addr"`
	NetBarGateway   string `json:"net_bar_gateway"`
	NetBarMask      string `json:"net_bar_mask"`
	DeviceNumber    string `json:"device_number"`
	Operator        string `json:"operator"`
	PoolID          int64  `json:"pool_id"`
	Header          string `header:"X-System"`
}

// 创建座位信息
type CreateDeviceSetReq struct {
	BizID           int64  `json:"biz_id"`
	DeviceSetNumber string `json:"device_set_number"`
	NetBarAddr      string `json:"net_bar_addr"`
	NetBarGateway   string `json:"net_bar_gateway"`
	NetBarMask      string `json:"net_bar_mask"`
	DeviceNumber    string `json:"device_number"`
	Operator        string `json:"operator"`
	PoolID          int64  `json:"pool_id"`
	//DeviceMac       string `json:"device_mac,omitempty"`
	Header string `header:"X-System"`
}

// 删除座位信息
type DeleteDeviceSetReq struct {
	DeviceSetID int64  `json:"device_set_id"`
	Header      string `header:"X-System"`
}

// 更新座位的超管信息
type UpdateDeviceSetDisklessReq struct {
	DeviceSetID int64  `json:"device_set_id"`
	ConfigInfo  string `json:"config_info"`
	AdminState  int    `json:"admin_state"`
	Header      string `header:"X-System"`
}

type UpdateDeviceSetDisklessResp struct {
}

type EsportRoomConfigInfo struct {
	VmId int `json:"vmid"` // 实例ID
}

// 查询座位的列表
func EsportGetBizAllDeviceSet(ctx context.Context, host, sessionId string, req *GetBizAllDeviceSetReq) (list *GetBizAllDeviceSet, err error) {

	url := fmt.Sprintf("%s/device_set/getBizAllDeviceSet", host)
	req.Header = "diskless-aggregator"

	httpResp, err := httpc.Do(context.Background(), http.MethodPost, url, req)
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] httpc.Do err bizId[%d] host[%s] req:%s err:%+v", sessionId, req.BizID, url, helper.ToJSON(req), err)
		return nil, errors.New("查询终端的设备列表失败")
	}
	if httpResp.StatusCode != 200 {
		logx.WithContext(ctx).Errorf("[%s] httpc.Do err url:%s, req:%s, resp:%s", url, helper.ToJSON(req), helper.ToJSON(httpResp))
		return nil, fmt.Errorf("请求失败, status:%d", httpResp.StatusCode)
	}

	saasResp := new(GetBizAllDeviceSetResp)
	body, _ := io.ReadAll(httpResp.Body)

	if err := json.Unmarshal(body, saasResp); err != nil {
		logx.WithContext(ctx).Errorf("[%s] json.Unmarshal err bizId[%d] host[%s],err:%+v", sessionId, req.BizID, url, err)
		return nil, errors.New("反序列化失败")
	}
	if saasResp.Code != 0 {
		logx.WithContext(ctx).Errorf("[%s] saasResp err bizId[%d] host[%s] saasResp:%+v", sessionId, req.BizID, url, helper.ToJSON(saasResp))
		return nil, errors.New(saasResp.Msg)
	}
	list = &saasResp.Data
	logx.WithContext(ctx).Debugf("[%s] EsportGetBizAllDeviceSet success url:%s, bizId[%d] req:%s, list:%s", sessionId, url, req.BizID, helper.ToJSON(req), helper.ToJSON(list))
	return
}

// 更新座位信息
func EsportUpdateDeviceSet(ctx context.Context, host, sessionId string, bizId int64, req *UpdateDeviceSetReq) (err error) {

	url := fmt.Sprintf("%s/device_set/update", host)
	req.Header = "diskless-aggregator"

	httpResp, err := httpc.Do(context.Background(), http.MethodPost, url, req)
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] httpc.Do err bizId[%d] host[%s] err:%+v", sessionId, bizId, url, err)
		return errors.New("查询终端的设备列表失败")
	}
	if httpResp.StatusCode != 200 {
		logx.WithContext(ctx).Errorf("[%s] httpc.Do err url:%s, req:%s, resp:%s", sessionId, url, helper.ToJSON(req), helper.ToJSON(httpResp))
		return fmt.Errorf("请求失败, status:%d", httpResp.StatusCode)
	}

	saasResp := new(SaasCommonResp)
	body, _ := io.ReadAll(httpResp.Body)

	if err := json.Unmarshal(body, saasResp); err != nil {
		logx.WithContext(ctx).Errorf("[%s] json.Unmarshal err host[%s],err:%+v", sessionId, bizId, url, err)
		return errors.New("反序列化失败")
	}
	if saasResp.Code != 0 {
		return errors.New(saasResp.Msg)
	}

	logx.WithContext(ctx).Debugf("[%s] EsportUpdateDeviceSet success url:%s, bizId:%d, saasResp:%s", sessionId, url, bizId, helper.ToJSON(saasResp))
	return
}

// 创建座位信息
func EsportCreateDeviceSet(ctx context.Context, host, sessionId string, bizId int64, req *CreateDeviceSetReq) (err error) {

	url := fmt.Sprintf("%s/device_set/create", host)
	req.Header = "diskless-aggregator"

	httpResp, err := httpc.Do(context.Background(), http.MethodPost, url, req)
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] httpc.Do err bizId[%d] host[%s] err:%+v", sessionId, bizId, url, err)
		return errors.New("查询终端的设备列表失败")
	}
	if httpResp.StatusCode != 200 {
		logx.WithContext(ctx).Errorf("[%s] httpc.Do err url:%s, req:%s, resp:%s", sessionId, url, helper.ToJSON(req), helper.ToJSON(httpResp))
		return fmt.Errorf("请求失败, status:%d", httpResp.StatusCode)
	}

	saasResp := new(SaasCommonResp)
	body, _ := io.ReadAll(httpResp.Body)

	if err := json.Unmarshal(body, saasResp); err != nil {
		logx.WithContext(ctx).Errorf("[%s] json.Unmarshal err bizId[%d] host[%s],err:%+v", sessionId, bizId, url, err)
		return errors.New("反序列化失败")
	}
	if saasResp.Code != 0 {
		logx.WithContext(ctx).Errorf("[%s] saasResp code err  bizId[%d] host[%s],saasResp:%s", sessionId, bizId, url, helper.ToJSON(saasResp))
		return errors.New(saasResp.Msg)
	}

	logx.WithContext(ctx).Debugf("[%s] EsportCreateDeviceSet success url:%s,,bizId[%d], saasResp:%+v", sessionId, url, bizId, saasResp)
	return
}

// 删除座位信息
func EsportDeleteDeviceSet(ctx context.Context, host, sessionId string, bizId int64, deviceSetId int64) (err error) {

	url := fmt.Sprintf("%s/device_set/delete", host)
	req := new(DeleteDeviceSetReq)
	req.DeviceSetID = deviceSetId
	req.Header = "diskless-aggregator"

	httpResp, err := httpc.Do(context.Background(), http.MethodPost, url, req)
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] httpc.Do err. bizId[%d] host[%s] err:%+v", sessionId, bizId, url, err)
		return errors.New("查询终端的设备列表失败")
	}
	if httpResp.StatusCode != 200 {
		logx.WithContext(ctx).Errorf("[%s] httpc.Do err url:%s, req:%s, resp:s", sessionId, url, helper.ToJSON(req), helper.ToJSON(httpResp))
		return fmt.Errorf("请求失败, status:%d", httpResp.StatusCode)

	}

	saasResp := new(SaasCommonResp)
	body, _ := io.ReadAll(httpResp.Body)

	if err := json.Unmarshal(body, saasResp); err != nil {
		logx.WithContext(ctx).Errorf("[%s] json.Unmarshal err. bizId[%d] host[%s],err:%+v", sessionId, bizId, url, err)
		return errors.New("反序列化失败")
	}
	if saasResp.Code != 0 {
		logx.WithContext(ctx).Errorf("[%s] saasResp code err agentId[%d] host[%s],saasResp:%s", sessionId, bizId, url, helper.ToJSON(saasResp))
		return errors.New(saasResp.Msg)
	}

	logx.WithContext(ctx).Debugf("[%s] EsportDeleteDeviceSet success url:%s,  bizId:%+v, saasResp:%s", sessionId, url, bizId, helper.ToJSON(saasResp))
	return
}

// 更新座位的超管信息
func EsportUpdateDeviceSetDiskless(ctx context.Context, host, sessionId string, bizId int64, deviceSetId int64, configInfo string, adminState int) (err error) {

	url := fmt.Sprintf("%s/device_set/updateDeviceSetDiskless", host)
	req := &UpdateDeviceSetDisklessReq{
		DeviceSetID: deviceSetId,
		ConfigInfo:  configInfo,
		AdminState:  adminState,
		Header:      "diskless-aggregator",
	}

	logx.WithContext(ctx).Debugf("[%s] httpc.Do url:%s, bizId[%d], req:%s", sessionId, url, bizId, helper.ToJSON(req))

	httpResp, err := httpc.Do(context.Background(), http.MethodPost, url, req)
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] httpc.Do err  bizId[%d] host[%s] deviceSetId[%d] err:%+v", sessionId, bizId, url, deviceSetId, err)
		return errors.New("查询终端的设备列表失败")
	}
	if httpResp.StatusCode != 200 {
		logx.WithContext(ctx).Errorf("[%s] httpc.Do err url:%s, req:%+v, resp:%+v", sessionId, url, req, helper.ToJSON(httpResp))
		return fmt.Errorf("请求失败, status:%d", httpResp.StatusCode)
	}

	saasResp := new(SaasCommonResp)
	body, _ := io.ReadAll(httpResp.Body)

	if err := json.Unmarshal(body, saasResp); err != nil {
		logx.WithContext(ctx).Errorf("[%s] json.Unmarshal err bizId[%d] host[%s], deviceSetId[%d] err:%+v", sessionId, bizId, url, deviceSetId, err)
		return errors.New("反序列化失败")
	}
	if saasResp.Code != 0 {
		logx.WithContext(ctx).Errorf("[%s] saasResp code err bizId[%d] host[%s], deviceSetId[%d] saasResp:%s", sessionId, bizId, url, deviceSetId, helper.ToJSON(saasResp))
		return errors.New(saasResp.Msg)
	}

	logx.WithContext(ctx).Debugf("[%s] EsportUpdateDeviceSetDiskless url:%s, bizId[%d], saasResp:%s", sessionId, url, bizId, helper.ToJSON(saasResp))
	return
}
