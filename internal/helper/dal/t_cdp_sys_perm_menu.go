package table

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type TCdpSysPermMenu struct {
	Id uint32 `orm:"column(id);auto" description:"编号" json:"id"`

	ParentId uint32 `orm:"column(parent_id)" description:"父级id" json:"parent_id"`

	Name string `orm:"column(name)" description:"名称" json:"name"`

	Router string `orm:"column(router)" description:"路由" json:"router"`

	Perms string `orm:"column(perms)" description:"权限" json:"perms"`

	Type uint32 `orm:"column(type)" description:"0=目录 1=菜单 2=权限" json:"type"`

	Icon string `orm:"column(icon)" description:"图标" json:"icon"`

	OrderNum uint32 `orm:"column(order_num)" description:"排序值" json:"order_num"`

	ViewPath string `orm:"column(view_path)" description:"页面路径" json:"view_path"`

	IsShow uint32 `orm:"column(is_show)" description:"0=隐藏 1=显示" json:"is_show"`

	ActiveRouter string `orm:"column(active_router)" description:"当前激活的菜单" json:"active_router"`

	SystemHost string `orm:"column(system_host)" description:"接口对应的系统host, 在字典中配置" json:"system_host"`

	IsPrivate int32 `orm:"column(is_private)" description:"是否只在内网展示 0-否，1-是" json:"is_private"`

	IsAdmin int32 `orm:"column(is_admin)" description:"是否超管独有 0-否，1-是" json:"is_admin"`

	Status uint32 `orm:"column(status)" description:"0=禁用 1=开启" json:"status"`

	CreateBy string `orm:"column(create_by)" description:"创建人" json:"create_by"`

	UpdateBy string `orm:"column(update_by)" description:"更新人" json:"update_by"`

	CreateTime time.Time `orm:"column(create_time)" description:"创建时间" json:"create_time"`

	UpdateTime time.Time `orm:"column(update_time)" description:"更新时间" json:"update_time"`

	ModifyTime time.Time `orm:"column(modify_time)" description:"更新时间" json:"modify_time"`
}

const (
	MenuStatusEnable  = 1
	MenuStatusDisable = 0
)

const (
	MenuTypeDirectory  = 0
	MenuTypeMenu       = 1
	MenuTypePermission = 2
)

const (
	IsPrivate = 1
	IsPublic  = 0
)

type TCdpSysPermMenuService struct {
	tableInfo *TableInfo
}

var T_TCdpSysPermMenuService *TCdpSysPermMenuService = &TCdpSysPermMenuService{
	tableInfo: &TableInfo{
		TableName: "t_cdp_sys_perm_menu",
		Tpy:       reflect.TypeOf(TCdpSysPermMenu{}),
	},
}

func init() {
	_TableMap["t_cdp_sys_perm_menu"] = T_TCdpSysPermMenuService.tableInfo
}
func (s *TCdpSysPermMenuService) Query(ctx context.Context, sessionId, query string, sortby interface{}, ascending interface{}) (*TCdpSysPermMenu, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Query(query, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Query[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpSysPermMenu), errcode, err
}

func (s *TCdpSysPermMenuService) QueryPage(ctx context.Context, sessionId, query string, offset int, limit int, sortby interface{}, ascending interface{}) (int, []TCdpSysPermMenu, int, error) {
	total, info, errcode, err := s.tableInfo.DBWrap.QueryPage(query, offset, limit, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] QueryPage[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return 0, nil, errcode, err
	}
	return total, info.([]TCdpSysPermMenu), errcode, err
}

func (s *TCdpSysPermMenuService) QueryAll(ctx context.Context, sessionId, query string, sortby interface{}, ascending interface{}) ([]TCdpSysPermMenu, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.QueryAll(query, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] QueryAll[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.([]TCdpSysPermMenu), errcode, err
}

func (s *TCdpSysPermMenuService) Update(ctx context.Context, sessionId string, key interface{}, info interface{}) (*TCdpSysPermMenu, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Update(key, info)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Update[%v] errcode[%v] err[%v] resp[%v]", sessionId, key, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpSysPermMenu), errcode, err
}

func (s *TCdpSysPermMenuService) Insert(ctx context.Context, sessionId string, info interface{}) (*TCdpSysPermMenu, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Insert(info)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Insert errcode[%v] err[%v] resp[%v]", sessionId, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpSysPermMenu), errcode, err
}

func (s *TCdpSysPermMenuService) Delete(ctx context.Context, sessionId string, key interface{}) (int, error) {
	errcode, err := s.tableInfo.DBWrap.Delete(key)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Delete[%v] errcode[%v] err[%v]", sessionId, key, errcode, err))
	}()

	return errcode, err
}
