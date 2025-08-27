package table

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type TCdpCloudboxInfo struct {
	Id uint32 `orm:"column(id);auto" description:"自增长ID" json:"id"`

	Name string `orm:"column(name)" description:"云盒名称" json:"name"`

	BizId int64 `orm:"column(biz_id)" description:"租户ID" json:"biz_id"`

	AreaId int64 `orm:"column(area_id)" description:"区域ID" json:"area_id"`

	Mac string `orm:"column(mac)" description:"mac地址" json:"mac"`

	Ip string `orm:"column(ip)" description:"云盒IP" json:"ip"`

	FirstStrategyId int64 `orm:"column(first_strategy_id)" description:"算力池ID主" json:"first_strategy_id"`

	SecondStrategyId int64 `orm:"column(second_strategy_id)" description:"算力池ID备" json:"second_strategy_id"`

	BootType int32 `orm:"column(boot_type)" description:"启动方式 2-本地启动 4-无盘启动" json:"boot_type"`

	BootSchemaId int32 `orm:"column(boot_schema_id)" description:"无盘启动启动方案" json:"boot_schema_id"`

	Status uint32 `orm:"column(status)" description:"0=失效 1=生效" json:"status"`

	Remark string `orm:"column(remark)" description:"备注" json:"remark"`

	CreateBy string `orm:"column(create_by)" description:"创建人" json:"create_by"`

	UpdateBy string `orm:"column(update_by)" description:"更新人" json:"update_by"`

	CreateTime time.Time `orm:"column(create_time)" description:"创建时间" json:"create_time"`

	UpdateTime time.Time `orm:"column(update_time)" description:"更新时间" json:"update_time"`

	ModifyTime time.Time `orm:"column(modify_time)" description:"修改时间" json:"modify_time"`
}

const (
	CloudBoxStatusValid   = 1
	CloudBoxStatusInvalid = 0
)

type TCdpCloudboxInfoService struct {
	tableInfo *TableInfo
}

var T_TCdpCloudboxInfoService *TCdpCloudboxInfoService = &TCdpCloudboxInfoService{
	tableInfo: &TableInfo{
		TableName: "t_cdp_cloudbox_info",
		Tpy:       reflect.TypeOf(TCdpCloudboxInfo{}),
	},
}

func init() {
	_TableMap["t_cdp_cloudbox_info"] = T_TCdpCloudboxInfoService.tableInfo
}

func (s *TCdpCloudboxInfoService) Query(ctx context.Context, sessionId, query string, sortby interface{}, ascending interface{}) (*TCdpCloudboxInfo, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Query(query, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Query[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpCloudboxInfo), errcode, err
}

func (s *TCdpCloudboxInfoService) QueryPage(ctx context.Context, sessionId, query string, offset int, limit int, sortby interface{}, ascending interface{}) (int, []TCdpCloudboxInfo, int, error) {
	total, info, errcode, err := s.tableInfo.DBWrap.QueryPage(query, offset, limit, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] QueryPage[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return 0, nil, errcode, err
	}
	return total, info.([]TCdpCloudboxInfo), errcode, err
}

func (s *TCdpCloudboxInfoService) QueryAll(ctx context.Context, sessionId, query string, sortby interface{}, ascending interface{}) ([]TCdpCloudboxInfo, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.QueryAll(query, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] QueryAll[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.([]TCdpCloudboxInfo), errcode, err
}

func (s *TCdpCloudboxInfoService) Update(ctx context.Context, sessionId string, key interface{}, info interface{}) (*TCdpCloudboxInfo, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Update(key, info)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Update[%v] errcode[%v] err[%v] resp[%v]", sessionId, key, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpCloudboxInfo), errcode, err
}

func (s *TCdpCloudboxInfoService) Insert(ctx context.Context, sessionId string, info interface{}) (*TCdpCloudboxInfo, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Insert(info)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Insert errcode[%v] err[%v] resp[%v]", sessionId, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpCloudboxInfo), errcode, err
}

func (s *TCdpCloudboxInfoService) Delete(ctx context.Context, sessionId string, key interface{}) (int, error) {
	errcode, err := s.tableInfo.DBWrap.Delete(key)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Delete[%v] errcode[%v] err[%v]", sessionId, key, errcode, err))
	}()

	return errcode, err
}
