package table

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type TCdpSysRole struct {
	Id uint32 `orm:"column(id);auto" description:"编号" json:"id"`

	Name string `orm:"column(name)" description:"名称" json:"name"`

	Remark string `orm:"column(remark)" description:"备注" json:"remark"`

	PermMenuIds string `orm:"column(perm_menu_ids)" description:"菜单集合" json:"perm_menu_ids"`

	Platform uint32 `orm:"column(platform)" description:"平台 1-算力平台2-施工平台" json:"platform"`

	IsAdmin uint32 `orm:"column(is_admin)" description:"0-否 1-是" json:"is_admin"`

	RegionId uint32 `orm:"column(region_id)" description:"" json:"region_id"`

	Status uint32 `orm:"column(status)" description:"0=禁用 1=开启" json:"status"`

	CreateBy string `orm:"column(create_by)" description:"创建人" json:"create_by"`

	UpdateBy string `orm:"column(update_by)" description:"更新人" json:"update_by"`

	CreateTime time.Time `orm:"column(create_time)" description:"创建时间" json:"create_time"`

	UpdateTime time.Time `orm:"column(update_time)" description:"更新时间" json:"update_time"`

	ModifyTime time.Time `orm:"column(modify_time)" description:"更新时间" json:"modify_time"`
}

const (
	RoleStatusEnable  = 1
	RoleStatusDisable = 0
)

const (
	RoleIdAdmin             = 1
	RoleIdPowerAdmin        = 2
	RoleIdConstructionAdmin = 3
)

const (
	RoleIsAdminYes = 1
	RoleIsAdminNo  = 0
)

func (r *TCdpSysRole) RoleIsAdmin() bool {
	return r.IsAdmin == 1
}

type TCdpSysRoleService struct {
	tableInfo *TableInfo
}

var T_TCdpSysRoleService *TCdpSysRoleService = &TCdpSysRoleService{
	tableInfo: &TableInfo{
		TableName: "t_cdp_sys_role",
		Tpy:       reflect.TypeOf(TCdpSysRole{}),
	},
}

func init() {
	_TableMap["t_cdp_sys_role"] = T_TCdpSysRoleService.tableInfo
}

func (s *TCdpSysRoleService) Query(ctx context.Context, sessionId, query string, sortby interface{}, ascending interface{}) (*TCdpSysRole, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Query(query, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Query[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpSysRole), errcode, err
}

func (s *TCdpSysRoleService) QueryPage(ctx context.Context, sessionId, query string, offset int, limit int, sortby interface{}, ascending interface{}) (int, []TCdpSysRole, int, error) {
	total, info, errcode, err := s.tableInfo.DBWrap.QueryPage(query, offset, limit, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] QueryPage[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return 0, nil, errcode, err
	}
	return total, info.([]TCdpSysRole), errcode, err
}

func (s *TCdpSysRoleService) QueryAll(ctx context.Context, sessionId, query string, sortby interface{}, ascending interface{}) ([]TCdpSysRole, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.QueryAll(query, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] QueryAll[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.([]TCdpSysRole), errcode, err
}

func (s *TCdpSysRoleService) Update(ctx context.Context, sessionId string, key interface{}, info interface{}) (*TCdpSysRole, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Update(key, info)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Update[%v] errcode[%v] err[%v] resp[%v]", sessionId, key, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpSysRole), errcode, err
}

func (s *TCdpSysRoleService) Insert(ctx context.Context, sessionId string, info interface{}) (*TCdpSysRole, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Insert(info)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Insert errcode[%v] err[%v] resp[%v]", sessionId, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpSysRole), errcode, err
}

func (s *TCdpSysRoleService) Delete(ctx context.Context, sessionId string, key interface{}) (int, error) {
	errcode, err := s.tableInfo.DBWrap.Delete(key)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Delete[%v] errcode[%v] err[%v]", sessionId, key, errcode, err))
	}()

	return errcode, err
}
