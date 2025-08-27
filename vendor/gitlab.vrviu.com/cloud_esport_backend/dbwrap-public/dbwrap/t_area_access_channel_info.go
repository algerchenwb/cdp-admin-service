package dbwrap

import (
	"reflect"
	"time"
)

// CREATE TABLE `t_area_access_channel_info` (
// 	`aacid` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '区域接入通道id',
// 	`area_type` smallint(5) unsigned NOT NULL DEFAULT '0' COMMENT '分区，用于接入管理',
// 	`mgr_state` smallint(5) unsigned NOT NULL DEFAULT '0' COMMENT '管理状态，0-未知；1-允许使用；2-暂停使用；',
// 	`channel_ipv4_address` char(16) NOT NULL DEFAULT '' COMMENT '通道地址',
// 	`mgr_ipv4_address` char(16) NOT NULL DEFAULT '' COMMENT '通道管理地址',
// 	`weight` smallint(5) NOT NULL DEFAULT '0' COMMENT '通道使用权重',
// 	`usage_upper_limit` smallint(5) NOT NULL DEFAULT '0' COMMENT '通道使用上限',
// 	`portseg_base`  int(11) NOT NULL DEFAULT '0' COMMENT '端口组基数',
// 	`portseg_limit`  int(11) NOT NULL DEFAULT '0' COMMENT '端口组数量',
// 	`create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
// 	`update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
// 	`modify_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
// 	PRIMARY KEY (`aapid`),
// 	UNIQUE KEY `channel_ipv4_address` (`channel_ipv4_address`),
// ) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8

// DBAreaAccessChannelInfo 区域接入通道信息
type DBAreaAccessChannelInfo struct {
	AACID              uint64    `json:"aacid"`
	AreaType           int       `json:"area_type"`
	MgrState           int       `json:"mgr_state"`
	ChannelIpv4Address string    `json:"channel_ipv4_address"`
	MgrIpv4Address     string    `json:"mgr_ipv4_address"`
	Weight             int       `json:"weight"`
	UsageUpperLimit    uint64    `json:"usage_upper_limit"`
	PortSegBase        int       `json:"portseg_base"`
	PortSegLimit       int       `json:"portseg_limit"`
	CreateTime         time.Time `json:"create_time"`
	UpdateTime         time.Time `json:"update_time"`
	ModifyTime         time.Time `json:"modify_time"`
}

type AACIWrap struct {
	AreaDBWrap
}

func CreateAACIWrap(host string, callerID int, flowID string) IAreaDBWrap {
	return &AACIWrap{
		AreaDBWrap{
			_table:    _tAreaAccessChannelInfo,
			_host:     host,
			_callerID: callerID,
			_flowID:   flowID,
			_typ:      reflect.TypeOf(DBAreaAccessChannelInfo{}),
		},
	}
}
