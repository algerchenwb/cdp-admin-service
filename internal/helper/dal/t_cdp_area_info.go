package table

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type TCdpAreaInfo struct {
	Id int64 `orm:"column(id);auto" description:"自增ID（字符串形式转成BIGINT）" json:"id"`

	PrimaryId int64 `orm:"column(primary_id)" description:"一代ID" json:"primary_id"`

	AgentId int64 `orm:"column(agent_id)" description:"二代ID" json:"agent_id"`

	AreaId int32 `orm:"column(area_id)" description:"区域ID" json:"area_id"`

	Name string `orm:"column(name)" description:"区域名称" json:"name"`

	Status int32 `orm:"column(status)" description:"状态：=0 未生效 -1 生效" json:"status"`

	RegionId int32 `orm:"column(region_id)" description:"地域ID" json:"region_id"`

	DeploymentType int32 `orm:"column(deployment_type)" description:"部署类型" json:"deployment_type"`

	Remark string `orm:"column(remark)" description:"备注" json:"remark"`

	ProxyAddr string `orm:"column(proxy_addr)" description:"代理地址" json:"proxy_addr"`

	SchemaConfig string `orm:"column(schema_config)" description:"区域编排方案配置" json:"schema_config"`

	ResetSchemaConfig string `orm:"column(reset_schema_config)" description:"区域重置编排方案配置" json:"reset_schema_config"`

	CreateBy string `orm:"column(create_by)" description:"创建人" json:"create_by"`

	UpdateBy string `orm:"column(update_by)" description:"更新人" json:"update_by"`

	CreateTime time.Time `orm:"column(create_time)" description:"创建时间" json:"create_time"`

	UpdateTime time.Time `orm:"column(update_time)" description:"更新时间" json:"update_time"`

	ModifyTime time.Time `orm:"column(modify_time)" description:"修改时间" json:"modify_time"`
}

const (
	AreaStatusEnable  = 1
	AreaStatusDisable = 0
)

type TCdpAreaInfoService struct {
	tableInfo *TableInfo
}

var T_TCdpAreaInfoService *TCdpAreaInfoService = &TCdpAreaInfoService{
	tableInfo: &TableInfo{
		TableName: "t_cdp_area_info",
		Tpy:       reflect.TypeOf(TCdpAreaInfo{}),
	},
}

func init() {
	_TableMap["t_cdp_area_info"] = T_TCdpAreaInfoService.tableInfo
}

func (s *TCdpAreaInfoService) Query(ctx context.Context, sessionId, query string, sortby interface{}, ascending interface{}) (*TCdpAreaInfo, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Query(query, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Query[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpAreaInfo), errcode, err
}

func (s *TCdpAreaInfoService) QueryPage(ctx context.Context, sessionId, query string, offset int, limit int, sortby interface{}, ascending interface{}) (int, []TCdpAreaInfo, int, error) {
	total, info, errcode, err := s.tableInfo.DBWrap.QueryPage(query, offset, limit, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] QueryPage[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return 0, nil, errcode, err
	}
	return total, info.([]TCdpAreaInfo), errcode, err
}

func (s *TCdpAreaInfoService) QueryAll(ctx context.Context, sessionId, query string, sortby interface{}, ascending interface{}) ([]TCdpAreaInfo, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.QueryAll(query, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] QueryAll[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.([]TCdpAreaInfo), errcode, err
}

func (s *TCdpAreaInfoService) Update(ctx context.Context, sessionId string, key interface{}, info interface{}) (*TCdpAreaInfo, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Update(key, info)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Update[%v] errcode[%v] err[%v] resp[%v]", sessionId, key, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpAreaInfo), errcode, err
}

func (s *TCdpAreaInfoService) Insert(ctx context.Context, sessionId string, info interface{}) (*TCdpAreaInfo, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Insert(info)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Insert errcode[%v] err[%v] resp[%v]", sessionId, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpAreaInfo), errcode, err
}

func (s *TCdpAreaInfoService) Delete(ctx context.Context, sessionId string, key interface{}) (int, error) {
	errcode, err := s.tableInfo.DBWrap.Delete(key)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Delete[%v] errcode[%v] err[%v]", sessionId, key, errcode, err))
	}()

	return errcode, err
}
