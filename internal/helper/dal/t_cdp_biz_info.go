package table

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type TCdpBizInfo struct {
	Id int64 `orm:"column(id);auto" description:"自增ID（字符串形式转成BIGINT）" json:"id"`

	BizId int64 `orm:"column(biz_id)" description:"租户ID" json:"biz_id"`

	BizName string `orm:"column(biz_name)" description:"租户名称" json:"biz_name"`

	AreaId int32 `orm:"column(area_id)" description:"区域ID" json:"area_id"`

	RegionId int32 `orm:"column(region_id)" description:"地域ID" json:"region_id"`

	ContactPerson string `orm:"column(contact_person)" description:"联系人" json:"contact_person"`

	Mobile string `orm:"column(mobile)" description:"手机号" json:"mobile"`

	Serverinfo string `orm:"column(serverinfo)" description:"服务器列表" json:"serverinfo"`

	StrategyIds string `orm:"column(strategy_ids)" description:"算力策略列表" json:"strategy_ids"`

	BootSchemaId int64 `orm:"column(boot_schema_id)" description:"启动方杂" json:"boot_schema_id"`

	VlanId int32 `orm:"column(vlan_id)" description:"客户机vlanID" json:"vlan_id"`

	BoxVlanId int32 `orm:"column(box_vlan_id)" description:"云盒vlanID" json:"box_vlan_id"`

	ClientNumLimit int32 `orm:"column(client_num_limit)" description:"终端数量限制" json:"client_num_limit"`

	Remark string `orm:"column(remark)" description:"备注" json:"remark"`

	Status uint32 `orm:"column(status)" description:"0=禁用 1=待审核 2=已上线" json:"status"`

	CreateBy string `orm:"column(create_by)" description:"创建租户账号" json:"create_by"`

	UpdateBy string `orm:"column(update_by)" description:"更新租户账号" json:"update_by"`

	CreateTime time.Time `orm:"column(create_time)" description:"创建时间" json:"create_time"`

	UpdateTime time.Time `orm:"column(update_time)" description:"更新时间" json:"update_time"`

	ModifyTime time.Time `orm:"column(modify_time)" description:"更新时间" json:"modify_time"`
}

const (
	BizStatusDeleted     = 0 // 已删除 未生效
	BizStatusWaitWorking = 1 // 待审核
	BizStatusWorking     = 2 // 已审核
	BizStatusOnline      = 3 // 已上线  已运营
)

type TCdpBizInfoService struct {
	tableInfo *TableInfo
}

var T_TCdpBizInfoService *TCdpBizInfoService = &TCdpBizInfoService{
	tableInfo: &TableInfo{
		TableName: "t_cdp_biz_info",
		Tpy:       reflect.TypeOf(TCdpBizInfo{}),
	},
}

func init() {
	_TableMap["t_cdp_biz_info"] = T_TCdpBizInfoService.tableInfo
}

func (s *TCdpBizInfoService) Query(ctx context.Context, sessionId, query string, sortby interface{}, ascending interface{}) (*TCdpBizInfo, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Query(query, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Query[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpBizInfo), errcode, err
}

func (s *TCdpBizInfoService) QueryPage(ctx context.Context, sessionId, query string, offset int, limit int, sortby interface{}, ascending interface{}) (int, []TCdpBizInfo, int, error) {
	total, info, errcode, err := s.tableInfo.DBWrap.QueryPage(query, offset, limit, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] QueryPage[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return 0, nil, errcode, err
	}
	return total, info.([]TCdpBizInfo), errcode, err
}

func (s *TCdpBizInfoService) QueryAll(ctx context.Context, sessionId, query string, sortby interface{}, ascending interface{}) ([]TCdpBizInfo, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.QueryAll(query, sortby, ascending)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] QueryAll[%v] errcode[%v] err[%v] resp[%v]", sessionId, query, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.([]TCdpBizInfo), errcode, err
}

func (s *TCdpBizInfoService) Update(ctx context.Context, sessionId string, key interface{}, info interface{}) (*TCdpBizInfo, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Update(key, info)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Update[%v] errcode[%v] err[%v] resp[%v]", sessionId, key, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpBizInfo), errcode, err
}

func (s *TCdpBizInfoService) Insert(ctx context.Context, sessionId string, info interface{}) (*TCdpBizInfo, int, error) {
	info, errcode, err := s.tableInfo.DBWrap.Insert(info)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Insert errcode[%v] err[%v] resp[%v]", sessionId, errcode, err, info))
	}()
	if err != nil {
		return nil, errcode, err
	}
	return info.(*TCdpBizInfo), errcode, err
}

func (s *TCdpBizInfoService) Delete(ctx context.Context, sessionId string, key interface{}) (int, error) {
	errcode, err := s.tableInfo.DBWrap.Delete(key)
	defer func() {
		logx.WithContext(ctx).Debug(fmt.Sprintf("[%v] Delete[%v] errcode[%v] err[%v]", sessionId, key, errcode, err))
	}()

	return errcode, err
}
