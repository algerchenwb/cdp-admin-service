package table

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	ServerTypeImage int = 1 // 镜像服务器
	ServerTypeWrite int = 2 // 回写服务器
	ServerTypeData  int = 3 // 数据服务器
)

type TCdpServerManagement struct {
	Id int64 `orm:"column(id);auto" description:"自增ID" json:"id"`

	AreaId int64 `orm:"column(area_id)" description:"区域ID" json:"area_id"`

	Ip string `orm:"column(ip)" description:"服务器IP地址" json:"ip"`

	Type string `orm:"column(type)" description:"服务器类型" json:"type"`

	BootTime time.Time `orm:"column(boot_time)" description:"开机赶时间" json:"boot_time"`

	Status uint32 `orm:"column(status)" description:"0=禁用 1=启用 2=在线" json:"status"`

	Remark string `orm:"column(remark)" description:"备注" json:"remark"`

	CreateBy string `orm:"column(create_by)" description:"创建账号" json:"create_by"`

	UpdateBy string `orm:"column(update_by)" description:"更新账号" json:"update_by"`

	CreateTime time.Time `orm:"column(create_time)" description:"创建时间" json:"create_time"`

	UpdateTime time.Time `orm:"column(update_time)" description:"更新时间" json:"update_time"`

	ModifyTime time.Time `orm:"column(modify_time)" description:"更新时间" json:"modify_time"`
}

type TCdpServerManagementService struct {
	tableInfo *TableInfo
}

var T_TCdpServerManagementService *TCdpServerManagementService = &TCdpServerManagementService{
	tableInfo: &TableInfo{
		TableName: "t_cdp_server_management",
		Tpy:       reflect.TypeOf(TCdpServerManagement{}),
	},
}

func init() {
	_TableMap["t_cdp_server_management"] = T_TCdpServerManagementService.tableInfo
}

func (s *TCdpServerManagementService) Query(ctx context.Context, sessionId, query string, sortby interface{}, ascending interface{}) (*TCdpServerManagement, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Query(query, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Query[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpServerManagement), errcode, err
}

func (s *TCdpServerManagementService) QueryPage(ctx context.Context, sessionId, query string, offset int, limit int, sortby interface{}, ascending interface{}) (int, []TCdpServerManagement, int, error) {
	total, info, errcode, err := s.tableInfo.DBWrap.QueryPage(query, offset, limit, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] QueryPage[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return 0, nil, errcode, err
	}
	return total, info.([]TCdpServerManagement), errcode, err
}

func (s *TCdpServerManagementService) QueryAll(ctx context.Context, sessionId, query string, sortby interface{}, ascending interface{}) ([]TCdpServerManagement, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.QueryAll(query, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] QueryAll[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.([]TCdpServerManagement), errcode, err
}

func (s *TCdpServerManagementService) Update(ctx context.Context, sessionId string, key interface{}, info interface{}) (*TCdpServerManagement, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Update(key, info)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Update[%v] errcode[%v] err[%v] resp[%v]", sessionId, key, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpServerManagement), errcode, err
}

func (s *TCdpServerManagementService) Insert(ctx context.Context, sessionId string, info interface{}) (*TCdpServerManagement, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Insert(info)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Insert errcode[%v] err[%v] resp[%v]", sessionId, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpServerManagement), errcode, err
}

func (s *TCdpServerManagementService) Delete(ctx context.Context, sessionId string, key interface{}) (int, error) {
	errcode, err := s.tableInfo.DBWrap.Delete(key)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Delete[%v] errcode[%v] err[%v]", sessionId, key, errcode, err))
	}()

	return errcode, err
}
