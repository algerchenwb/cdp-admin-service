package table

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	ResourceStrategyStatusInvalid = 0 // 0=禁用
	ResourceStrategyStatusValid   = 1 // 1=启用
)

const (
	ResourceStrategyBootTypePre  = 0 // 启动模式  0-预开机
	ResourceStrategyBootTypeReal = 1 // 启动模式  1-实时开机
)

type TCdpResourceStrategy struct {
	Id int64 `orm:"column(id);auto" description:"自增ID" json:"id"`

	Name string `orm:"column(name)" description:"算力策略名" json:"name"`

	ApplicableLever int32 `orm:"column(applicable_lever)" description:"适用等级 1-全部 2-节点域 3-指定节点 4-全部门店 5-指定门店 规则1包含2包含3..." json:"applicable_lever"`

	AreaId int64 `orm:"column(area_id)" description:"区域ID" json:"area_id"`

	SpecId int64 `orm:"column(spec_id)" description:"规格Id" json:"spec_id"`

	SpecName string `orm:"column(spec_name)" description:"规格名称" json:"spec_name"`

	OuterSpecId int64 `orm:"column(outer_spec_id)" description:"外部规格ID" json:"outer_spec_id"`

	InstPoolId int64 `orm:"column(inst_pool_id)" description:"算力池ID" json:"inst_pool_id"`

	TotalInstances int32 `orm:"column(total_instances)" description:"实例数量" json:"total_instances"`

	BootType int32 `orm:"column(boot_type)" description:"启动模式  0-预开机 1-实时开机" json:"boot_type"`

	Remark string `orm:"column(remark)" description:"备注" json:"remark"`

	Status uint32 `orm:"column(status)" description:"0=禁用 1=启用" json:"status"`

	CreateBy string `orm:"column(create_by)" description:"创建账号" json:"create_by"`

	UpdateBy string `orm:"column(update_by)" description:"更新账号" json:"update_by"`

	CreateTime time.Time `orm:"column(create_time)" description:"创建时间" json:"create_time"`

	UpdateTime time.Time `orm:"column(update_time)" description:"更新时间" json:"update_time"`

	ModifyTime time.Time `orm:"column(modify_time)" description:"更新时间" json:"modify_time"`
}
type TCdpResourceStrategyService struct {
	tableInfo *TableInfo
}

var T_TCdpResourceStrategyService *TCdpResourceStrategyService = &TCdpResourceStrategyService{
	tableInfo: &TableInfo{
		TableName: "t_cdp_resource_strategy",
		Tpy:       reflect.TypeOf(TCdpResourceStrategy{}),
	},
}

func init() {
	_TableMap["t_cdp_resource_strategy"] = T_TCdpResourceStrategyService.tableInfo
}

func (s *TCdpResourceStrategyService) Query(ctx context.Context, sessionId, query string, sortby interface{}, ascending interface{}) (*TCdpResourceStrategy, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Query(query, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Query[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpResourceStrategy), errcode, err
}

func (s *TCdpResourceStrategyService) QueryPage(ctx context.Context, sessionId, query string, offset int, limit int, sortby interface{}, ascending interface{}) (int, []TCdpResourceStrategy, int, error) {
	total, info, errcode, err := s.tableInfo.DBWrap.QueryPage(query, offset, limit, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] QueryPage[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return 0, nil, errcode, err
	}
	return total, info.([]TCdpResourceStrategy), errcode, err
}

func (s *TCdpResourceStrategyService) QueryAll(ctx context.Context, sessionId, query string, sortby interface{}, ascending interface{}) ([]TCdpResourceStrategy, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.QueryAll(query, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] QueryAll[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.([]TCdpResourceStrategy), errcode, err
}

func (s *TCdpResourceStrategyService) Update(ctx context.Context, sessionId string, key interface{}, info interface{}) (*TCdpResourceStrategy, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Update(key, info)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Update[%v] errcode[%v] err[%v] resp[%v]", sessionId, key, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpResourceStrategy), errcode, err
}

func (s *TCdpResourceStrategyService) Insert(ctx context.Context, sessionId string, info interface{}) (*TCdpResourceStrategy, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Insert(info)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Insert errcode[%v] err[%v] resp[%v]", sessionId, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpResourceStrategy), errcode, err
}

func (s *TCdpResourceStrategyService) Delete(ctx context.Context, sessionId string, key interface{}) (int, error) {
	errcode, err := s.tableInfo.DBWrap.Delete(key)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Delete[%v] errcode[%v] err[%v]", sessionId, key, errcode, err))
	}()

	return errcode, err
}
