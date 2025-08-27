package rediswrap

import (
	"context"

	"github.com/go-redis/redis/v8"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic/otelwrap"
)

// RedisSRem 移除集合中一个或多个成员
// @param    key: 键
// @param member: 元素
// @return <arg_1>: 移除元素个数，类型为int64
func RedisSRem(key string, member []interface{}, clients ...redis.UniversalClient) (int64, error) {
	return RedisSRemWithCtx(otelwrap.Skip(), key, member, clients...)
}

func RedisSRemWithCtx(ctx context.Context, key string, member []interface{}, clients ...redis.UniversalClient) (int64, error) {
	return GetRedisCli(clients...).SRem(ctx, key, member...).Result()
}

// RedisSMembers 返回集合中的所有的成员。 不存在的集合 key 被视为空集合
// @param    key: 键
// @return <arg_1>: 集合成员
func RedisSMembers(key string, clients ...redis.UniversalClient) ([]string, error) {
	return RedisSMembersWithCtx(otelwrap.Skip(), key, clients...)
}

func RedisSMembersWithCtx(ctx context.Context, key string, clients ...redis.UniversalClient) ([]string, error) {
	return GetRedisCli(clients...).SMembers(ctx, key).Result()
}
