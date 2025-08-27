package table

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	BizStrategyStatusInvalid = 0 // 0=禁用
	BizStrategyStatusValid   = 1 // 1=启用
)

type TCdpBizStrategy struct {
	Id int64 `orm:"column(id);auto" description:"自增ID" json:"id"`

	BizId int64 `orm:"column(biz_id)" description:"租户ID" json:"biz_id"`

	AreaId int64 `orm:"column(area_id)" description:"区域ID" json:"area_id"`

	InstStrategyId int64 `orm:"column(inst_strategy_id)" description:"算力策略ID" json:"inst_strategy_id"`

	Remark string `orm:"column(remark)" description:"备注" json:"remark"`

	Status uint32 `orm:"column(status)" description:"0=禁用 1=启用" json:"status"`

	CreateBy string `orm:"column(create_by)" description:"创建账号" json:"create_by"`

	UpdateBy string `orm:"column(update_by)" description:"更新账号" json:"update_by"`

	CreateTime time.Time `orm:"column(create_time)" description:"创建时间" json:"create_time"`

	UpdateTime time.Time `orm:"column(update_time)" description:"更新时间" json:"update_time"`

	ModifyTime time.Time `orm:"column(modify_time)" description:"更新时间" json:"modify_time"`
}

type TCdpBizStrategyService struct {
	tableInfo *TableInfo
}

var T_TCdpBizStrategyService *TCdpBizStrategyService = &TCdpBizStrategyService{
	tableInfo: &TableInfo{
		TableName: "t_cdp_biz_strategy",
		Tpy:       reflect.TypeOf(TCdpBizStrategy{}),
	},
}

func init() {
	_TableMap["t_cdp_biz_strategy"] = T_TCdpBizStrategyService.tableInfo
}

func (s *TCdpBizStrategyService) Query(ctx context.Context, sessionId, query string, sortby interface{}, ascending interface{}) (*TCdpBizStrategy, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Query(query, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Query[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpBizStrategy), errcode, err
}

func (s *TCdpBizStrategyService) QueryPage(ctx context.Context, sessionId, query string, offset int, limit int, sortby interface{}, ascending interface{}) (int, []TCdpBizStrategy, int, error) {
	total, info, errcode, err := s.tableInfo.DBWrap.QueryPage(query, offset, limit, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] QueryPage[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return 0, nil, errcode, err
	}
	return total, info.([]TCdpBizStrategy), errcode, err
}

func (s *TCdpBizStrategyService) QueryAll(ctx context.Context, sessionId, query string, sortby interface{}, ascending interface{}) ([]TCdpBizStrategy, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.QueryAll(query, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] QueryAll[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.([]TCdpBizStrategy), errcode, err
}

func (s *TCdpBizStrategyService) Update(ctx context.Context, sessionId string, key interface{}, info interface{}) (*TCdpBizStrategy, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Update(key, info)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Update[%v] errcode[%v] err[%v] resp[%v]", sessionId, key, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpBizStrategy), errcode, err
}

func (s *TCdpBizStrategyService) Insert(ctx context.Context, sessionId string, info interface{}) (*TCdpBizStrategy, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Insert(info)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Insert errcode[%v] err[%v] resp[%v]", sessionId, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpBizStrategy), errcode, err
}

func (s *TCdpBizStrategyService) Delete(ctx context.Context, sessionId string, key interface{}) (int, error) {
	errcode, err := s.tableInfo.DBWrap.Delete(key)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Delete[%v] errcode[%v] err[%v]", sessionId, key, errcode, err))
	}()

	return errcode, err
}
