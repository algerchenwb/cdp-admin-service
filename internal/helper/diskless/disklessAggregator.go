package diskless

import (
	"cdp-admin-service/internal/config"
	"cdp-admin-service/internal/helper"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpc"
)

type CloudClientAssignNoticeBody struct {
}

type CloudClientAssignNoticeReq struct {
	FlowID        string `json:"flowId"           validate:"required,min=1,max=999"         label:"流水ID"`
	BizId         int64  `json:"bizId"            validate:"required,number,gte=1"          label:"租户ID"`
	AreaId        int64  `json:"areaId"           validate:"required,number,gte=1"          label:"节点区域ID"`
	VmId          int64  `json:"vmId"             validate:"required,number,gte=1"          label:"实例Id"`
	MAC           string `json:"mac"              validate:"required,mac"                   label:"Mac地址"`
	CloudHostName string `json:"cloudHostName"    validate:"required,min=0,max=191"          label:"云主机名"` //  串流的实例 room name
}

type CloudClientAssignNoticeResp struct {
	CommonRet
	Body CloudClientAssignNoticeBody `json:"body"`
}

type CommonRet struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

func CloudClientAssignNotice(ctx context.Context, config config.Config, sessionId string, req CloudClientAssignNoticeReq) (respData *CloudClientAssignNoticeResp, err error) {

	url := config.DisklessConfig.AggregatorHost + "/v1/diskless/aggregator/cloudclientassignnotice"
	logx.WithContext(ctx).Infof("[%s] CloudClientAssignNotice url:%s, req: %v", sessionId, url, helper.ToJSON(req))

	resp, err := httpc.Do(context.Background(), http.MethodPost, url, req)
	if err != nil || resp.StatusCode != http.StatusOK {
		logx.WithContext(ctx).Errorf("[%s] CloudClientAssignNotice request req:%s resp: %+v error:%+v", sessionId, helper.ToJSON(req), resp, err)
		return nil, err
	}
	respData = &CloudClientAssignNoticeResp{}
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, respData)
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] CloudClientAssignNotice.ParseJsonBody req:%s, error %s ", sessionId, helper.ToJSON(req), err)
		return nil, err
	}
	logx.WithContext(ctx).Infof("[%s]  CloudClientAssignNotice req:%s Response:%s", sessionId, helper.ToJSON(req), helper.ToJSON(respData))

	return
}
