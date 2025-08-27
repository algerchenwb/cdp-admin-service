package cdp_cache

import (
	table "cdp-admin-service/internal/helper/dal"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// bizCache 用于缓存 bizId 到 TCdpBizInfo 的映射
var bizCache sync.Map
var defaultExpireTime int64 = 60

// BizCacheItem 缓存项结构
type BizCacheItem struct {
	Info       *table.TCdpBizInfo
	ExpireTime time.Time
}

// SetBizCache 设置业务信息缓存
func SetBizCache(bizId int64, second int64, info *table.TCdpBizInfo) {
	if second == 0 {
		second = defaultExpireTime // 默认缓存1分钟
	}

	bizCache.Store(bizId, &BizCacheItem{
		Info:       info,
		ExpireTime: time.Now().Add(time.Duration(second) * time.Second), // 缓存1小时
	})
}

// GetBizCache 获取业务信息缓存
func GetBizCache(ctx context.Context, sessionId string, bizId int64) *table.TCdpBizInfo {
	bizInfo := new(table.TCdpBizInfo)

	if value, ok := bizCache.Load(bizId); ok {
		if item, ok := value.(*BizCacheItem); ok {
			// 检查是否过期
			if time.Now().Before(item.ExpireTime) {
				logx.WithContext(ctx).Infof("[%s] get bizInfo from cache success, bizId:%d, bizInfo: %v", sessionId, bizId, item.Info)
				return item.Info
			}
			// 已过期，删除缓存
			bizCache.Delete(bizId)
			// 缓存未命中，从数据库查询
		}
	}
	bizInfo, _, err := table.T_TCdpBizInfoService.Query(ctx, sessionId, fmt.Sprintf("biz_id:%d$status__ex:%d", bizId, 0), nil, nil)
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] GetBizCache T_TCdpBizInfoService Query biz info failed, bizId:%d, err: %v", sessionId, bizId, err)
		return bizInfo
	} else {
		// 更新缓存
		logx.WithContext(ctx).Infof("[%s] get bizInfo from db and set cache success,, bizId:%d, bizInfo: %v", sessionId, bizId, bizInfo)
		SetBizCache(bizId, defaultExpireTime, bizInfo)
		return bizInfo
	}
}

// DeleteBizCache 删除业务信息缓存
func DeleteBizCache(bizId int64) {
	bizCache.Delete(bizId)
}

// GetBizName 获取业务名称
func GetBizName(ctx context.Context, sessionId string, bizId int64) string {
	// 先尝试从缓存获取
	bizInfo := GetBizCache(ctx, sessionId, bizId)
	if bizInfo != nil {
		return bizInfo.BizName
	}
	return ""
}
