package table

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type TCdpCloudclientInfo struct {
	Id uint32 `orm:"column(id);auto" description:"自增长ID" json:"id"`

	Name string `orm:"column(name)" description:"云客户机名称" json:"name"`

	BizId int64 `orm:"column(biz_id)" description:"租户ID" json:"biz_id"`

	AreaId int64 `orm:"column(area_id)" description:"区域ID" json:"area_id"`

	CloudboxMac string `orm:"column(cloudbox_mac)" description:"mac地址" json:"cloudbox_mac"`

	Vmid int64 `orm:"column(vmid)" description:"实例Id v1.0.4新增" json:"vmid"`

	FlowId string `orm:"column(flow_id)" description:"分实例后flowId v1.0.4新增" json:"flow_id"`

	HostIp string `orm:"column(host_ip)" description:"主机IP" json:"host_ip"`

	FirstStrategyId int64 `orm:"column(first_strategy_id)" description:"算力池ID主" json:"first_strategy_id"`

	FirstBootSchemaId int64 `orm:"column(first_boot_schema_id)" description:"启动方案ID主" json:"first_boot_schema_id"`

	SecondStrategyId int64 `orm:"column(second_strategy_id)" description:"算力池ID备" json:"second_strategy_id"`

	SecondBootSchemaId int64 `orm:"column(second_boot_schema_id)" description:"启动方案ID备" json:"second_boot_schema_id"`

	ConfigInfo string `orm:"column(config_info)" description:"配置信息" json:"config_info"`

	AdminState uint32 `orm:"column(admin_state)" description:"超管状态" json:"admin_state"`

	ClientType int32 `orm:"column(client_type)" description:"1-1.0客户机 2-2.0客户机" json:"client_type"`

	DisklessSeatId int64 `orm:"column(diskless_seat_id)" description:"无盘座位ID" json:"diskless_seat_id"`

	Status uint32 `orm:"column(status)" description:"0=失效 1=生效" json:"status"`

	Remark string `orm:"column(remark)" description:"备注" json:"remark"`

	CreateBy string `orm:"column(create_by)" description:"创建人" json:"create_by"`

	UpdateBy string `orm:"column(update_by)" description:"更新人" json:"update_by"`

	CreateTime time.Time `orm:"column(create_time)" description:"创建时间" json:"create_time"`

	UpdateTime time.Time `orm:"column(update_time)" description:"更新时间" json:"update_time"`

	ModifyTime time.Time `orm:"column(modify_time)" description:"修改时间" json:"modify_time"`
}

const (
	CloudClientStatusValid   = 1
	CloudClientStatusInvalid = 0
)

const (
	ClientType1 = 1
	ClientType2 = 2
)

type TCdpCloudclientInfoService struct {
	tableInfo *TableInfo
}

var T_TCdpCloudclientInfoService *TCdpCloudclientInfoService = &TCdpCloudclientInfoService{
	tableInfo: &TableInfo{
		TableName: "t_cdp_cloudclient_info",
		Tpy:       reflect.TypeOf(TCdpCloudclientInfo{}),
	},
}

func init() {
	_TableMap["t_cdp_cloudclient_info"] = T_TCdpCloudclientInfoService.tableInfo
}

func (s *TCdpCloudclientInfoService) Query(ctx context.Context, sessionId, query string, sortby interface{}, ascending interface{}) (*TCdpCloudclientInfo, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Query(query, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Query[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpCloudclientInfo), errcode, err
}

func (s *TCdpCloudclientInfoService) QueryPage(ctx context.Context, sessionId, query string, offset int, limit int, sortby interface{}, ascending interface{}) (int, []TCdpCloudclientInfo, int, error) {
	total, info, errcode, err := s.tableInfo.DBWrap.QueryPage(query, offset, limit, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] QueryPage[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return 0, nil, errcode, err
	}
	return total, info.([]TCdpCloudclientInfo), errcode, err
}

func (s *TCdpCloudclientInfoService) QueryAll(ctx context.Context, sessionId, query string, sortby interface{}, ascending interface{}) ([]TCdpCloudclientInfo, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.QueryAll(query, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] QueryAll[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.([]TCdpCloudclientInfo), errcode, err
}

func (s *TCdpCloudclientInfoService) Update(ctx context.Context, sessionId string, key interface{}, info interface{}) (*TCdpCloudclientInfo, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Update(key, info)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Update[%v] errcode[%v] err[%v] resp[%v]", sessionId, key, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpCloudclientInfo), errcode, err
}

func (s *TCdpCloudclientInfoService) Insert(ctx context.Context, sessionId string, info interface{}) (*TCdpCloudclientInfo, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Insert(info)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Insert errcode[%v] err[%v] resp[%v]", sessionId, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpCloudclientInfo), errcode, err
}

func (s *TCdpCloudclientInfoService) Delete(ctx context.Context, sessionId string, key interface{}) (int, error) {
	errcode, err := s.tableInfo.DBWrap.Delete(key)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Delete[%v] errcode[%v] err[%v]", sessionId, key, errcode, err))
	}()

	return errcode, err
}
