package table

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	InstancePoolStatusInvalid = 0 // 0=禁用
	InstancePoolStatusValid   = 1 // 1=启用
)

const (
	InstancePoolDefaultId = 1000 // 0=禁用
)

type TCdpInstancePool struct {
	Id uint32 `orm:"column(id);auto" description:"自增长ID" json:"id"`

	PoolId int64 `orm:"column(pool_id)" description:"算力池ID" json:"pool_id"`

	InstPoolName string `orm:"column(inst_pool_name)" description:"算力池名" json:"inst_pool_name"`

	BizId int64 `orm:"column(biz_id)" description:"租户ID" json:"biz_id"`

	AreaId int64 `orm:"column(area_id)" description:"区域ID" json:"area_id"`

	Status uint32 `orm:"column(status)" description:"0=失效 1=生效" json:"status"`

	Remark string `orm:"column(remark)" description:"备注" json:"remark"`

	CreateBy string `orm:"column(create_by)" description:"创建人" json:"create_by"`

	UpdateBy string `orm:"column(update_by)" description:"更新人" json:"update_by"`

	CreateTime time.Time `orm:"column(create_time)" description:"创建时间" json:"create_time"`

	UpdateTime time.Time `orm:"column(update_time)" description:"更新时间" json:"update_time"`

	ModifyTime time.Time `orm:"column(modify_time)" description:"修改时间" json:"modify_time"`
}

type TCdpInstancePoolService struct {
	tableInfo *TableInfo
}

var T_TCdpInstancePoolService *TCdpInstancePoolService = &TCdpInstancePoolService{
	tableInfo: &TableInfo{
		TableName: "t_cdp_instance_pool",
		Tpy:       reflect.TypeOf(TCdpInstancePool{}),
	},
}

func init() {
	_TableMap["t_cdp_instance_pool"] = T_TCdpInstancePoolService.tableInfo
}

func (s *TCdpInstancePoolService) Query(ctx context.Context, sessionId, query string, sortby interface{}, ascending interface{}) (*TCdpInstancePool, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Query(query, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Query[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpInstancePool), errcode, err
}

func (s *TCdpInstancePoolService) QueryPage(ctx context.Context, sessionId, query string, offset int, limit int, sortby interface{}, ascending interface{}) (int, []TCdpInstancePool, int, error) {
	total, info, errcode, err := s.tableInfo.DBWrap.QueryPage(query, offset, limit, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] QueryPage[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return 0, nil, errcode, err
	}
	return total, info.([]TCdpInstancePool), errcode, err
}

func (s *TCdpInstancePoolService) QueryAll(ctx context.Context, sessionId, query string, sortby interface{}, ascending interface{}) ([]TCdpInstancePool, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.QueryAll(query, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] QueryAll[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.([]TCdpInstancePool), errcode, err
}

func (s *TCdpInstancePoolService) Update(ctx context.Context, sessionId string, key interface{}, info interface{}) (*TCdpInstancePool, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Update(key, info)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Update[%v] errcode[%v] err[%v] resp[%v]", sessionId, key, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpInstancePool), errcode, err
}

func (s *TCdpInstancePoolService) Insert(ctx context.Context, sessionId string, info interface{}) (*TCdpInstancePool, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Insert(info)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Insert errcode[%v] err[%v] resp[%v]", sessionId, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpInstancePool), errcode, err
}

func (s *TCdpInstancePoolService) Delete(ctx context.Context, sessionId string, key interface{}) (int, error) {
	errcode, err := s.tableInfo.DBWrap.Delete(key)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Delete[%v] errcode[%v] err[%v]", sessionId, key, errcode, err))
	}()

	return errcode, err
}
