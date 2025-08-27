package table

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type TCdpSysLog struct {
	Id uint32 `orm:"column(id);auto" description:"编号" json:"id"`

	UserId uint32 `orm:"column(user_id)" description:"操作账号" json:"user_id"`

	Account string `orm:"column(account)" description:"账号" json:"account"`

	Ip string `orm:"column(ip)" description:"ip" json:"ip"`

	Uri string `orm:"column(uri)" description:"请求路径" json:"uri"`

	Type uint32 `orm:"column(type)" description:"1=登录日志 2=操作日志" json:"type"`

	Request string `orm:"column(request)" description:"请求数据" json:"request"`

	Response string `orm:"column(response)" description:"响应数据" json:"response"`

	Platform int32 `orm:"column(platform)" description:"平台 1-算力平台 2-施工平台" json:"platform"`

	Status uint32 `orm:"column(status)" description:"0=失败 1=成功" json:"status"`

	CreateTime time.Time `orm:"column(create_time)" description:"创建时间" json:"create_time"`

	UpdateTime time.Time `orm:"column(update_time)" description:"更新时间" json:"update_time"`

	ModifyTime time.Time `orm:"column(modify_time)" description:"修改时间" json:"modify_time"`
}

type TCdpSysLogService struct {
	tableInfo *TableInfo
}

var T_TCdpSysLogService *TCdpSysLogService = &TCdpSysLogService{
	tableInfo: &TableInfo{
		TableName: "t_cdp_sys_log",
		Tpy:       reflect.TypeOf(TCdpSysLog{}),
	},
}

func init() {
	_TableMap["t_cdp_sys_log"] = T_TCdpSysLogService.tableInfo
}

func (s *TCdpSysLogService) Query(ctx context.Context, sessionId, query string, sortby interface{}, ascending interface{}) (*TCdpSysLog, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Query(query, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Query[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpSysLog), errcode, err
}

func (s *TCdpSysLogService) QueryPage(ctx context.Context, sessionId, query string, offset int, limit int, sortby interface{}, ascending interface{}) (int, []TCdpSysLog, int, error) {
	total, info, errcode, err := s.tableInfo.DBWrap.QueryPage(query, offset, limit, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] QueryPage[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return 0, nil, errcode, err
	}
	return total, info.([]TCdpSysLog), errcode, err
}

func (s *TCdpSysLogService) QueryAll(ctx context.Context, sessionId, query string, sortby interface{}, ascending interface{}) ([]TCdpSysLog, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.QueryAll(query, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] QueryAll[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.([]TCdpSysLog), errcode, err
}

func (s *TCdpSysLogService) Update(ctx context.Context, sessionId string, key interface{}, info interface{}) (*TCdpSysLog, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Update(key, info)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Update[%v] errcode[%v] err[%v] resp[%v]", sessionId, key, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpSysLog), errcode, err
}

func (s *TCdpSysLogService) Insert(ctx context.Context, sessionId string, info interface{}) (*TCdpSysLog, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Insert(info)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Insert errcode[%v] err[%v] resp[%v]", sessionId, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpSysLog), errcode, err
}

func (s *TCdpSysLogService) Delete(ctx context.Context, sessionId string, key interface{}) (int, error) {
	errcode, err := s.tableInfo.DBWrap.Delete(key)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Delete[%v] errcode[%v] err[%v]", sessionId, key, errcode, err))
	}()

	return errcode, err
}
