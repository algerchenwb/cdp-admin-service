package saas

import (
	"cdp-admin-service/internal/helper"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type TEsportsOuterSpecificationInfo struct {
	Id          int64     `orm:"column(outer_spec_id);auto" description:"规格ID" json:"outer_spec_id"`
	Name        string    `orm:"column(name)" description:"名称" json:"name"`
	CPUClassID  int64     `orm:"column(cpu_class_id)" description:"CPU分类ID" json:"cpu_class_id"`
	GPUClassID  int64     `orm:"column(gpu_class_id)" description:"GPU分类ID" json:"gpu_class_id"`
	InnerSpecID int64     `orm:"column(inner_spec_id)" description:"内部规格ID" json:"inner_spec_id"`
	CreateBy    string    `orm:"column(create_by)" description:"创建者" json:"create_by"`
	UpdateBy    string    `orm:"column(update_by)" description:"修改者" json:"update_by"`
	Remark      string    `orm:"column(remark)" description:"备注" json:"remark"`
	State       int       `orm:"column(state)" description:"状态: 0-无效, 1-有效, 2-删除" json:"state"`
	CreateTime  time.Time `orm:"column(create_time);type(datetime);auto_now_add" json:"create_time"`
	UpdateTime  time.Time `orm:"column(update_time);type(datetime);auto_now_add" json:"update_time"`
	ModifyTime  time.Time `orm:"column(modify_time);type(datetime);auto_now" json:"modify_time"`
}

type GetOuterSpecificationInfoListReq struct {
	Offset   int64    `json:"offset"`
	Limit    int64    `json:"limit"`
	CondList []string `json:"cond_list"`
	Sorts    string   `json:"sorts"`
	Orders   string   `json:"orders"`
}

type GetOuterSpecificationInfoListResp struct {
	Code int                               `json:"code"`
	Msg  string                            `json:"msg"`
	Data GetOuterSpecificationInfoListBody `json:"data"`
}

type GetOuterSpecificationInfoListBody struct {
	Total int64                    `json:"total"`
	List  []OuterSpecificationInfo `json:"list"`
}
type OuterSpecificationInfo struct {
	TEsportsOuterSpecificationInfo
	UsedAreaList []int64 `json:"used_area_list"`
}

func (s *SaasServerService) GetSpecificationInfoList(ctx context.Context, sessionId string,
	condList []string) (total int64, list []OuterSpecificationInfo, err error) {

	url := fmt.Sprintf("%s/specification/getSpecificationInfoList", s.Host)
	req := &GetOuterSpecificationInfoListReq{
		Limit:    500,
		CondList: condList,
	}

	resp := new(GetOuterSpecificationInfoListResp)

	_, err = helper.HttpPost(ctx, sessionId, url, s.Timeout, req, resp)
	if err != nil {
		logx.WithContext(ctx).Errorf("helper.HttpPost err  host[%s] err:%+v", url, err)
		return 0, nil, errors.New("获取规格列表失败")
	}

	if resp.Code != 0 {
		return 0, nil, errors.New("获取规格列表失败")
	}

	logx.WithContext(ctx).Debugf("[%s] url[%s] req:%s resp:%s ", sessionId, url, helper.ToJSON(req), helper.ToJSON(resp))

	total = resp.Data.Total
	list = resp.Data.List
	return
}
