package table

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type TCdpBootSchemaInfo struct {
	Id int64 `orm:"column(id);auto" description:"自增ID（字符串形式转成BIGINT）" json:"id"`

	BizId int64 `orm:"column(biz_id)" description:"租户ID" json:"biz_id"`

	AreaId int64 `orm:"column(area_id)" description:"区域ID" json:"area_id"`

	Name string `orm:"column(name)" description:"启动方案信息名" json:"name"`

	DisklessSchemaId int64 `orm:"column(diskless_schema_id)" description:"无盘编排方案ID" json:"diskless_schema_id"`

	BootCommandIds string `orm:"column(boot_command_ids)" description:"启动命令表ID" json:"boot_command_ids"`

	OsImageId string `orm:"column(os_image_id)" description:"镜像ID" json:"os_image_id"`

	Remark string `orm:"column(remark)" description:"备注" json:"remark"`

	Status uint32 `orm:"column(status)" description:"0=禁用 1=启用" json:"status"`

	CreateBy string `orm:"column(create_by)" description:"创建账号" json:"create_by"`

	UpdateBy string `orm:"column(update_by)" description:"更新账号" json:"update_by"`

	CreateTime time.Time `orm:"column(create_time)" description:"创建时间" json:"create_time"`

	UpdateTime time.Time `orm:"column(update_time)" description:"更新时间" json:"update_time"`

	ModifyTime time.Time `orm:"column(modify_time)" description:"更新时间" json:"modify_time"`
}

const (
	TCdpBootSchemaInfoStatusEnable  = 1
	TCdpBootSchemaInfoStatusDisable = 0
)

type TCdpBootSchemaInfoService struct {
	tableInfo *TableInfo
}

var T_TCdpBootSchemaInfoService *TCdpBootSchemaInfoService = &TCdpBootSchemaInfoService{
	tableInfo: &TableInfo{
		TableName: "t_cdp_boot_schema_info",
		Tpy:       reflect.TypeOf(TCdpBootSchemaInfo{}),
	},
}

func init() {
	_TableMap["t_cdp_boot_schema_info"] = T_TCdpBootSchemaInfoService.tableInfo
}

func (s *TCdpBootSchemaInfoService) Query(ctx context.Context, sessionId, query string, sortby interface{}, ascending interface{}) (*TCdpBootSchemaInfo, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Query(query, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Query[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpBootSchemaInfo), errcode, err
}

func (s *TCdpBootSchemaInfoService) QueryPage(ctx context.Context, sessionId, query string, offset int, limit int, sortby interface{}, ascending interface{}) (int, []TCdpBootSchemaInfo, int, error) {
	total, info, errcode, err := s.tableInfo.DBWrap.QueryPage(query, offset, limit, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] QueryPage[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return 0, nil, errcode, err
	}
	return total, info.([]TCdpBootSchemaInfo), errcode, err
}

func (s *TCdpBootSchemaInfoService) QueryAll(ctx context.Context, sessionId, query string, sortby interface{}, ascending interface{}) ([]TCdpBootSchemaInfo, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.QueryAll(query, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] QueryAll[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.([]TCdpBootSchemaInfo), errcode, err
}

func (s *TCdpBootSchemaInfoService) Update(ctx context.Context, sessionId string, key interface{}, info interface{}) (*TCdpBootSchemaInfo, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Update(key, info)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Update[%v] errcode[%v] err[%v] resp[%v]", sessionId, key, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpBootSchemaInfo), errcode, err
}

func (s *TCdpBootSchemaInfoService) Insert(ctx context.Context, sessionId string, info interface{}) (*TCdpBootSchemaInfo, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Insert(info)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Insert errcode[%v] err[%v] resp[%v]", sessionId, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpBootSchemaInfo), errcode, err
}

func (s *TCdpBootSchemaInfoService) Delete(ctx context.Context, sessionId string, key interface{}) (int, error) {
	errcode, err := s.tableInfo.DBWrap.Delete(key)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Delete[%v] errcode[%v] err[%v]", sessionId, key, errcode, err))
	}()

	return errcode, err
}
