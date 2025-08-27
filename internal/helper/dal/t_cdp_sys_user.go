package table

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type TCdpSysUser struct {
	Id uint32 `orm:"column(id);auto" description:"编号" json:"id"`

	Account string `orm:"column(account)" description:"账号" json:"account"`

	Password string `orm:"column(password)" description:"密码" json:"password"`

	Nickname string `orm:"column(nickname)" description:"昵称" json:"nickname"`

	Avatar string `orm:"column(avatar)" description:"头像" json:"avatar"`

	Mobile string `orm:"column(mobile)" description:"手机号" json:"mobile"`

	RoleId int32 `orm:"column(role_id)" description:"角色ID" json:"role_id"`

	AreaIds string `orm:"column(area_ids)" description:"用户有权限区域ID列表" json:"area_ids"`

	AreaRegions string `orm:"column(area_regions)" description:"可管域" json:"area_regions"`

	BizIds string `orm:"column(biz_ids)" description:"授权合约列表" json:"biz_ids"`

	BizRegions string `orm:"column(biz_regions)" description:"授权合约域" json:"biz_regions"`

	Platform int32 `orm:"column(platform)" description:"平台 1-算力平台 2-施工平台" json:"platform"`

	IsAdmin uint32 `orm:"column(is_admin)" description:"0-否 1-是" json:"is_admin"`

	Status uint32 `orm:"column(status)" description:"0=禁用 1=开启" json:"status"`

	Remark string `orm:"column(remark)" description:"备注" json:"remark"`

	CreateBy string `orm:"column(create_by)" description:"创建人" json:"create_by"`

	UpdateBy string `orm:"column(update_by)" description:"更新人" json:"update_by"`

	CreateTime time.Time `orm:"column(create_time)" description:"创建时间" json:"create_time"`

	UpdateTime time.Time `orm:"column(update_time)" description:"更新时间" json:"update_time"`

	ModifyTime time.Time `orm:"column(modify_time)" description:"修改时间" json:"modify_time"`
}

const (
	SysUserStatusEnable  = 1
	SysUserStatusDisable = 0
)

const (
	PlatformPower        = 1
	PlatformConstruction = 2
)

const (
	IsAdminYes = 1
	IsAdminNo  = 0
)

func (u *TCdpSysUser) UserIsAdmin() bool {
	return u.IsAdmin == IsAdminYes
}

type TCdpSysUserService struct {
	tableInfo *TableInfo
}

var T_TCdpSysUserService *TCdpSysUserService = &TCdpSysUserService{
	tableInfo: &TableInfo{
		TableName: "t_cdp_sys_user",
		Tpy:       reflect.TypeOf(TCdpSysUser{}),
	},
}

func init() {
	_TableMap["t_cdp_sys_user"] = T_TCdpSysUserService.tableInfo
}

func (s *TCdpSysUserService) Query(ctx context.Context, sessionId, query string, sortby interface{}, ascending interface{}) (*TCdpSysUser, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Query(query, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Query[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpSysUser), errcode, err
}

func (s *TCdpSysUserService) QueryPage(ctx context.Context, sessionId, query string, offset int, limit int, sortby interface{}, ascending interface{}) (int, []TCdpSysUser, int, error) {
	total, info, errcode, err := s.tableInfo.DBWrap.QueryPage(query, offset, limit, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] QueryPage[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return 0, nil, errcode, err
	}
	return total, info.([]TCdpSysUser), errcode, err
}

func (s *TCdpSysUserService) QueryAll(ctx context.Context, sessionId, query string, sortby interface{}, ascending interface{}) ([]TCdpSysUser, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.QueryAll(query, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] QueryAll[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.([]TCdpSysUser), errcode, err
}

func (s *TCdpSysUserService) Update(ctx context.Context, sessionId string, key interface{}, info interface{}) (*TCdpSysUser, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Update(key, info)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Update[%v] errcode[%v] err[%v] resp[%v]", sessionId, key, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpSysUser), errcode, err
}

func (s *TCdpSysUserService) Insert(ctx context.Context, sessionId string, info interface{}) (*TCdpSysUser, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Insert(info)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Insert errcode[%v] err[%v] resp[%v]", sessionId, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpSysUser), errcode, err
}

func (s *TCdpSysUserService) Delete(ctx context.Context, sessionId string, key interface{}) (int, error) {
	errcode, err := s.tableInfo.DBWrap.Delete(key)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Delete[%v] errcode[%v] err[%v]", sessionId, key, errcode, err))
	}()

	return errcode, err
}
