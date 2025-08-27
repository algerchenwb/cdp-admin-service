package table

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type TCdpSysDictionary struct {
	Id uint32 `orm:"column(id);auto" description:"编号" json:"id"`

	ParentId uint32 `orm:"column(parent_id)" description:"0=配置集 !0=父级id" json:"parent_id"`

	Name string `orm:"column(name)" description:"名称" json:"name"`

	Type uint32 `orm:"column(type)" description:"1文本 2数字 3数组 4单选 5多选 6下拉 7日期 8时间 9单图 10多图 11单文件 12多文件" json:"type"`

	UniqueKey string `orm:"column(unique_key)" description:"唯一值" json:"unique_key"`

	Value string `orm:"column(value)" description:"配置值" json:"value"`

	Status uint32 `orm:"column(status)" description:"0=禁用 1=开启" json:"status"`

	OrderNum uint32 `orm:"column(order_num)" description:"排序值" json:"order_num"`

	Remark string `orm:"column(remark)" description:"备注" json:"remark"`

	CreateBy string `orm:"column(create_by)" description:"创建人" json:"create_by"`

	UpdateBy string `orm:"column(update_by)" description:"更新人" json:"update_by"`

	CreateTime time.Time `orm:"column(create_time)" description:"创建时间" json:"create_time"`

	UpdateTime time.Time `orm:"column(update_time)" description:"更新时间" json:"update_time"`

	ModifyTime time.Time `orm:"column(modify_time)" description:"修改时间" json:"modify_time"`
}

type TCdpSysDictionaryService struct {
	tableInfo *TableInfo
}

var T_TCdpSysDictionaryService *TCdpSysDictionaryService = &TCdpSysDictionaryService{
	tableInfo: &TableInfo{
		TableName: "t_cdp_sys_dictionary",
		Tpy:       reflect.TypeOf(TCdpSysDictionary{}),
	},
}

func init() {
	_TableMap["t_cdp_sys_dictionary"] = T_TCdpSysDictionaryService.tableInfo
}

func (s *TCdpSysDictionaryService) Query(ctx context.Context, sessionId, query string, sortby interface{}, ascending interface{}) (*TCdpSysDictionary, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Query(query, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Query[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpSysDictionary), errcode, err
}

func (s *TCdpSysDictionaryService) QueryPage(ctx context.Context, sessionId, query string, offset int, limit int, sortby interface{}, ascending interface{}) (int, []TCdpSysDictionary, int, error) {
	total, info, errcode, err := s.tableInfo.DBWrap.QueryPage(query, offset, limit, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] QueryPage[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return 0, nil, errcode, err
	}
	return total, info.([]TCdpSysDictionary), errcode, err
}

func (s *TCdpSysDictionaryService) QueryAll(ctx context.Context, sessionId, query string, sortby interface{}, ascending interface{}) ([]TCdpSysDictionary, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.QueryAll(query, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] QueryAll[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.([]TCdpSysDictionary), errcode, err
}

func (s *TCdpSysDictionaryService) Update(ctx context.Context, sessionId string, key interface{}, info interface{}) (*TCdpSysDictionary, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Update(key, info)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Update[%v] errcode[%v] err[%v] resp[%v]", sessionId, key, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpSysDictionary), errcode, err
}

func (s *TCdpSysDictionaryService) Insert(ctx context.Context, sessionId string, info interface{}) (*TCdpSysDictionary, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Insert(info)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Insert errcode[%v] err[%v] resp[%v]", sessionId, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpSysDictionary), errcode, err
}

func (s *TCdpSysDictionaryService) Delete(ctx context.Context, sessionId string, key interface{}) (int, error) {
	errcode, err := s.tableInfo.DBWrap.Delete(key)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Delete[%v] errcode[%v] err[%v]", sessionId, key, errcode, err))
	}()

	return errcode, err
}
