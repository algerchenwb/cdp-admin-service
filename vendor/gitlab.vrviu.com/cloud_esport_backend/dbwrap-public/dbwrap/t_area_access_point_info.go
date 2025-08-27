package dbwrap

import (
	"reflect"
	"time"
)

// CREATE TABLE `t_area_access_point_info` (
// 	`aapid` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '区域接入点id',
// 	`area_type` smallint(5) unsigned NOT NULL DEFAULT '0' COMMENT '分区，用于接入管理',
// 	`online_state` smallint(5) unsigned NOT NULL DEFAULT '0' COMMENT '在线状态，0-未知；1-待初始化；2-初始化中；3-启动中；4-运行中；5-挂起中；6-已挂起；7-关机中；8-已关机；9-销毁中；10-已销毁',
// 	`mgr_state` smallint(5) unsigned NOT NULL DEFAULT '0' COMMENT '管理状态，0-未知；1-允许使用；2-暂停使用；',
// 	`inner_ipv4_address` char(16) NOT NULL DEFAULT '' COMMENT '供内网映射的ipv4地址',
// 	`outer_ipv4_address` char(16) NOT NULL DEFAULT '' COMMENT '供外网访问ipv4地址',
// 	`outer_ipv6_address` char(39) NOT NULL DEFAULT '' COMMENT '供外网访问ipv6地址',
// 	`outer_domain` varchar(128) NOT NULL DEFAULT '' COMMENT '供外网访问域名',
// 	`raw_ipv4_address` char(16) NOT NULL DEFAULT '' COMMENT '原始端口映射到ipv4地址',
// 	`biz_type` int(11) NOT NULL DEFAULT '0' COMMENT '业务ID',
// 	`create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
// 	`update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
// 	`modify_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
// 	PRIMARY KEY (`aapid`)
//   ) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8

// DBAreaAccessPointInfo 区域接入点信息
type DBAreaAccessPointInfo struct {
	AAPID            uint64    `json:"aapid"`
	AreaType         int       `json:"area_type"`
	OnlineState      int       `json:"online_state"`
	MgrState         int       `json:"mgr_state"`
	InnerIpv4Address string    `json:"inner_ipv4_address"`
	OuterIpv4Address string    `json:"outer_ipv4_address"`
	OuterIpv6Address string    `json:"outer_ipv6_address"`
	OuterDomain      string    `json:"outer_domain"`
	RawIpv4Address   string    `json:"raw_ipv4_address"`
	BizType          int       `json:"biz_type"`
	CreateTime       time.Time `json:"create_time"`
	UpdateTime       time.Time `json:"update_time"`
	ModifyTime       time.Time `json:"modify_time"`
}

type AAPIWrap struct {
	AreaDBWrap
}

func CreateAAPIWrap(host string, callerID int, flowID string) IAreaDBWrap {
	return &AAPIWrap{
		AreaDBWrap{
			_table:    _tAreaAccessPointInfo,
			_host:     host,
			_callerID: callerID,
			_flowID:   flowID,
			_typ:      reflect.TypeOf(DBAreaAccessPointInfo{}),
		},
	}
}
