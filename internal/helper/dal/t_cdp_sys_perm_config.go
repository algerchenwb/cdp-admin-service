package table

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type TCdpSysPermConfig struct {
	Id uint32 `orm:"column(id);auto" description:"编号" json:"id"`

	Perm string `orm:"column(perm)" description:"权限" json:"perm"`

	LoggingEnable uint32 `orm:"column(logging_enable)" description:"0=禁用 1=开启" json:"logging_enable"`

	CreateBy string `orm:"column(create_by)" description:"创建人" json:"create_by"`

	UpdateBy string `orm:"column(update_by)" description:"更新人" json:"update_by"`

	CreateTime time.Time `orm:"column(create_time)" description:"创建时间" json:"create_time"`

	UpdateTime time.Time `orm:"column(update_time)" description:"更新时间" json:"update_time"`

	ModifyTime time.Time `orm:"column(modify_time)" description:"修改时间" json:"modify_time"`
}

const (
	LoggingEnable  = 1
	LoggingDisable = 0
)

type TCdpSysPermConfigService struct {
	tableInfo *TableInfo
}

var T_TCdpSysPermConfigService *TCdpSysPermConfigService = &TCdpSysPermConfigService{
	tableInfo: &TableInfo{
		TableName: "t_cdp_sys_perm_config",
		Tpy:       reflect.TypeOf(TCdpSysPermConfig{}),
	},
}

func init() {
	_TableMap["t_cdp_sys_perm_config"] = T_TCdpSysPermConfigService.tableInfo
}
func (s *TCdpSysPermConfigService) Query(ctx context.Context, sessionId, query string, sortby interface{}, ascending interface{}) (*TCdpSysPermConfig, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Query(query, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Query[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpSysPermConfig), errcode, err
}

func (s *TCdpSysPermConfigService) QueryPage(ctx context.Context, sessionId, query string, offset int, limit int, sortby interface{}, ascending interface{}) (int, []TCdpSysPermConfig, int, error) {
	total, info, errcode, err := s.tableInfo.DBWrap.QueryPage(query, offset, limit, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] QueryPage[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return 0, nil, errcode, err
	}
	return total, info.([]TCdpSysPermConfig), errcode, err
}

func (s *TCdpSysPermConfigService) QueryAll(ctx context.Context, sessionId, query string, sortby interface{}, ascending interface{}) ([]TCdpSysPermConfig, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.QueryAll(query, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] QueryAll[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.([]TCdpSysPermConfig), errcode, err
}

func (s *TCdpSysPermConfigService) Update(ctx context.Context, sessionId string, key interface{}, info interface{}) (*TCdpSysPermConfig, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Update(key, info)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Update[%v] errcode[%v] err[%v] resp[%v]", sessionId, key, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpSysPermConfig), errcode, err
}

func (s *TCdpSysPermConfigService) Insert(ctx context.Context, sessionId string, info interface{}) (*TCdpSysPermConfig, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Insert(info)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Insert errcode[%v] err[%v] resp[%v]", sessionId, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpSysPermConfig), errcode, err
}

func (s *TCdpSysPermConfigService) Delete(ctx context.Context, sessionId string, key interface{}) (int, error) {
	errcode, err := s.tableInfo.DBWrap.Delete(key)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Delete[%v] errcode[%v] err[%v]", sessionId, key, errcode, err))
	}()

	return errcode, err
}
