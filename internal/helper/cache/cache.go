package cache

import (
	"cdp-admin-service/internal/config"
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic/rediswrap"
)

type Cache struct {
	redis  redis.UniversalClient
	prefix string
}

func NewCache(c config.Config) *Cache {
	return &Cache{
		redis:  rediswrap.GetRedisCli(),
		prefix: "cdp-admin-service-prefix:",
	}
}
func (c *Cache) Set(ctx context.Context, sessionId string, key string, value string) error {
	key = c.prefix + key
	logx.WithContext(ctx).Debugf("[%s] 设置缓存 key[%s] value[%s]", sessionId, key, value)
	return c.redis.Set(ctx, key, value, 0).Err()
}

func (c *Cache) SetEX(ctx context.Context, sessionId string, key string, value string, second int) error {
	key = c.prefix + key
	logx.WithContext(ctx).Debugf("[%s] 设置缓存 key[%s] value[%s] 过期时间[%d]秒", sessionId, key, value, second)
	return c.redis.SetEX(ctx, key, value, time.Duration(second)*time.Second).Err()
}

func (c *Cache) Exists(ctx context.Context, sessionId string, key string) bool {
	key = c.prefix + key
	logx.WithContext(ctx).Debugf("[%s] 检查缓存 key[%s]", sessionId, key)
	exist, err := c.redis.Exists(ctx, key).Result()
	if err != nil || exist == 0 {
		return false
	}
	return true
}

func (c *Cache) Get(ctx context.Context, sessionId string, key string) (string, error) {
	key = c.prefix + key
	logx.WithContext(ctx).Debugf("[%s] 获取缓存 key[%s]", sessionId, key)
	return c.redis.Get(ctx, key).Result()
}
func (c *Cache) Del(ctx context.Context, sessionId string, key string) (int64, error) {

	key = c.prefix + key
	logx.WithContext(ctx).Debugf("[%s] 删除缓存 key[%s]", sessionId, key)
	return c.redis.Del(ctx, key).Result()
}

func (c *Cache) SAdd(ctx context.Context, sessionId string, key string, value string) error {
	key = c.prefix + key
	logx.WithContext(ctx).Debugf("[%s] 添加缓存 key[%s] value[%s]", sessionId, key, value)
	_, err := c.redis.SAdd(ctx, key, value).Result()
	return err
}

func (c *Cache) SRem(ctx context.Context, sessionId string, key string, value string) error {
	key = c.prefix + key
	logx.WithContext(ctx).Debugf("[%s] 删除缓存 key[%s] value[%s]", sessionId, key, value)
	_, err := c.redis.SRem(ctx, key, value).Result()
	return err
}

func (c *Cache) Sismember(ctx context.Context, sessionId string, key string, value string) (bool, error) {
	key = c.prefix + key
	logx.WithContext(ctx).Debugf("[%s] 检查缓存 key[%s] value[%s]", sessionId, key, value)
	return c.redis.SIsMember(ctx, key, value).Result()
}
