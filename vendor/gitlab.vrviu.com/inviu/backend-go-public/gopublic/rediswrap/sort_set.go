package rediswrap

import (
	"context"

	"github.com/go-redis/redis/v8"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic/otelwrap"
)

// RedisZAdd 用于将一个成员元素及其分数值加入到有序集当中
// @param   key: 键
// @param score: 分数
// @param member: 元素
// @return <arg_1>: 值列表，类型为int
func RedisZAdd(key string, score float64, member interface{}, clients ...redis.UniversalClient) (int64, error) {
	return RedisZAddWithCtx(otelwrap.Skip(), key, score, member, clients...)
}

func RedisZAddWithCtx(ctx context.Context, key string, score float64, member interface{}, clients ...redis.UniversalClient) (int64, error) {
	memberInfo := &redis.Z{
		Score:  score,
		Member: member,
	}

	return GetRedisCli(clients...).ZAdd(ctx, key, memberInfo).Result()
}

// RedisZAddMultiSameScoreWithCtx 用于将多个成员元素使用相同的分数值加入到有序集当中
// @param   key: 键
// @param score: 分数
// @param members: 元素
// @return <arg_1>: 值列表，类型为int
func RedisZAddMultiSameScoreWithCtx(ctx context.Context, key string, score float64, members []interface{}, clients ...redis.UniversalClient) (int64, error) {
	memberInfos := make([]*redis.Z, 0, len(members))
	for _, v := range members {
		memberInfos = append(memberInfos, &redis.Z{
			Score:  score,
			Member: v,
		})
	}

	return GetRedisCli(clients...).ZAdd(ctx, key, memberInfos...).Result()
}

// RedisZRem 用于移除有序集中的一个或多个成员，不存在的成员将被忽略。
// @param   key: 键
// @param members: 要删除的元素
// @return <arg_1>: 值列表，类型为int
func RedisZRem(key string, members []interface{}, clients ...redis.UniversalClient) (int64, error) {
	return RedisZRemWithCtx(otelwrap.Skip(), key, members, clients...)
}

func RedisZRemWithCtx(ctx context.Context, key string, members []interface{}, clients ...redis.UniversalClient) (int64, error) {
	return GetRedisCli(clients...).ZRem(ctx, key, members...).Result()
}

// RedisZCard 计算集合中元素的数量
// @param   key: 键
// @return <arg_1>: 元素数量，类型为int64
func RedisZCard(key string, clients ...redis.UniversalClient) (int64, error) {
	return RedisZCardWithCtx(otelwrap.Skip(), key, clients...)
}

func RedisZCardWithCtx(ctx context.Context, key string, clients ...redis.UniversalClient) (int64, error) {
	return GetRedisCli(clients...).ZCard(ctx, key).Result()
}

// RedisZRank 返回有序集中指定成员的排名
// @param    key: 键
// @param member: 元素
// @return <arg_1>: 元素的排名，类型为int64
func RedisZRank(key, member string, clients ...redis.UniversalClient) (int64, error) {
	return RedisZRankWithCtx(otelwrap.Skip(), key, member, clients...)
}

func RedisZRankWithCtx(ctx context.Context, key, member string, clients ...redis.UniversalClient) (int64, error) {
	return GetRedisCli(clients...).ZRank(ctx, key, member).Result()
}

// RedisZScore 返回有序集中，成员的分数值
// @param    key: 键
// @param member: 元素
// @return <arg_1>: 成员的分数值，类型为float64
func RedisZScore(key, member string, clients ...redis.UniversalClient) (float64, error) {
	return RedisZScoreWithCtx(otelwrap.Skip(), key, member, clients...)
}

func RedisZScoreWithCtx(ctx context.Context, key, member string, clients ...redis.UniversalClient) (float64, error) {
	return GetRedisCli(clients...).ZScore(ctx, key, member).Result()
}

// RedisZRangeIsRev 返回有序集中，指定区间内的成员,其中成员的位置按分数来排序。
// @param   key: 键
// @param start stop: start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推, 也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。
// @param   isRev: 是否倒叙 false (从小到大)来排序，true (从大到小)来排序
// @return <arg_1>: 值列表，类型为int
func RedisZRangeIsRev(key string, start, stop int64, isRev bool, clients ...redis.UniversalClient) ([]string, error) {
	return RedisZRangeIsRevWithCtx(otelwrap.Skip(), key, start, stop, isRev, clients...)
}

func RedisZRangeIsRevWithCtx(ctx context.Context, key string, start, stop int64, isRev bool, clients ...redis.UniversalClient) ([]string, error) {
	if !isRev {
		return RedisZRangeWithCtx(ctx, key, start, stop, clients...)
	}

	return RedisZRevRangeWithCtx(ctx, key, start, stop, clients...)
}

// RedisZRange 返回有序集中，指定区间内的成员,其中成员的位置按分数值递增(从小到大)来排序。
// @param   key: 键
// @param start stop: start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推, 也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。
// @return <arg_1>: 值列表，类型为int
func RedisZRange(key string, start, stop int64, clients ...redis.UniversalClient) ([]string, error) {
	return RedisZRangeWithCtx(otelwrap.Skip(), key, start, stop, clients...)
}

func RedisZRangeWithCtx(ctx context.Context, key string, start, stop int64, clients ...redis.UniversalClient) ([]string, error) {
	return GetRedisCli(clients...).ZRange(ctx, key, start, stop).Result()
}

// RedisZRevRange 返回有序集中，指定区间内的成员,其中成员的位置按分数值递减(从大到小)来排列。
// @param   key: 键
// @param start stop: start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推, 也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。
// @return <arg_1>: 值列表，类型为int
func RedisZRevRange(key string, start, stop int64, clients ...redis.UniversalClient) ([]string, error) {
	return RedisZRevRangeWithCtx(otelwrap.Skip(), key, start, stop, clients...)
}

func RedisZRevRangeWithCtx(ctx context.Context, key string, start, stop int64, clients ...redis.UniversalClient) ([]string, error) {
	return GetRedisCli(clients...).ZRevRange(ctx, key, start, stop).Result()
}

// RedisZRangeWithScoresIsRev 返回有序集中，指定区间内的成员,其中成员的位置按分数值。
// @param   key: 键
// @param start stop: start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推, 也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。
// @param   isRev: 是否倒叙 false 按分数值递增(从小到大)来排序，true 按分数值递增(从大到小)来排序
// @return <arg_1>: 值列表，类型为 redis.Z, Score Member
func RedisZRangeWithScoresIsRev(key string, start, stop int64, isRev bool, clients ...redis.UniversalClient) ([]redis.Z, error) {
	return RedisZRangeWithScoresIsRevWithCtx(otelwrap.Skip(), key, start, stop, isRev, clients...)
}

func RedisZRangeWithScoresIsRevWithCtx(ctx context.Context, key string, start, stop int64, isRev bool, clients ...redis.UniversalClient) ([]redis.Z, error) {
	if !isRev {
		return RedisZRangeWithScoresWithCtx(ctx, key, start, stop, clients...)
	}

	return RedisZRevRangeWithScoresWithCtx(ctx, key, start, stop, clients...)
}

// RedisZRangeWithScores 返回有序集中，指定区间内的成员,其中成员的位置按分数值递增(从小到大)来排序。
// @param   key: 键
// @param start stop: start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推, 也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。
// @return <arg_1>: 值列表，类型为 redis.Z, Score Member
func RedisZRangeWithScores(key string, start, stop int64, clients ...redis.UniversalClient) ([]redis.Z, error) {
	return RedisZRangeWithScoresWithCtx(otelwrap.Skip(), key, start, stop, clients...)
}

func RedisZRangeWithScoresWithCtx(ctx context.Context, key string, start, stop int64, clients ...redis.UniversalClient) ([]redis.Z, error) {
	return GetRedisCli(clients...).ZRangeWithScores(ctx, key, start, stop).Result()
}

// RedisZRevRangeWithScores 返回有序集中，指定区间内的成员,其中成员的位置按分数值递减(从大到小)来排列。
// @param   key: 键
// @param start stop: start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推, 也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。
// @return <arg_1>: 值列表 类型为 redis.Z, Score Member
func RedisZRevRangeWithScores(key string, start, stop int64, clients ...redis.UniversalClient) ([]redis.Z, error) {
	return RedisZRevRangeWithScoresWithCtx(otelwrap.Skip(), key, start, stop, clients...)
}

func RedisZRevRangeWithScoresWithCtx(ctx context.Context, key string, start, stop int64, clients ...redis.UniversalClient) ([]redis.Z, error) {
	return GetRedisCli(clients...).ZRevRangeWithScores(ctx, key, start, stop).Result()
}

// RedisZRangeByScoreOffsetCountIsRev 返回有序集合中指定分数区间的成员列表。有序集成员按分数值递增(从小到大)次序排列。
// @param   key: 键
// @param min, max: 区间的取值使用闭区间 (小于等于或大于等于)，你也可以通过给参数前增加 ( 符号来使用可选的开区间 (小于或大于)
// @param offset, count: 偏移 和 获取的 数
// @return <arg_1>: 值列表，类型为int
func RedisZRangeByScoreOffsetCountIsRev(key string, min, max string, offset, count int64, isRev bool, clients ...redis.UniversalClient) ([]string, error) {
	return RedisZRangeByScoreOffsetCountIsRevWithCtx(otelwrap.Skip(), key, min, max, offset, count, isRev, clients...)
}

func RedisZRangeByScoreOffsetCountIsRevWithCtx(ctx context.Context, key string, min, max string, offset, count int64, isRev bool, clients ...redis.UniversalClient) ([]string, error) {
	opt := &redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: offset,
		Count:  count,
	}

	if !isRev {
		return RedisZRangeByScoreWithCtx(ctx, key, opt, clients...)
	}

	return RedisZRevRangeByScoreWithCtx(ctx, key, opt, clients...)
}

// RedisZRangeByScoreIsRev 返回有序集合中指定分数区间的成员列表。有序集成员按分数值递增(从小到大)次序排列。
// @param   key: 键
// @param min, max: 区间的取值使用闭区间 (小于等于或大于等于)，你也可以通过给参数前增加 ( 符号来使用可选的开区间 (小于或大于)
// @return <arg_1>: 值列表，类型为int
func RedisZRangeByScoreIsRev(key string, min, max string, isRev bool, clients ...redis.UniversalClient) ([]string, error) {
	return RedisZRangeByScoreIsRevWithCtx(otelwrap.Skip(), key, min, max, isRev, clients...)
}

func RedisZRangeByScoreIsRevWithCtx(ctx context.Context, key string, min, max string, isRev bool, clients ...redis.UniversalClient) ([]string, error) {
	opt := &redis.ZRangeBy{
		Min: min,
		Max: max,
	}

	if !isRev {
		return RedisZRangeByScoreWithCtx(ctx, key, opt, clients...)
	}

	return RedisZRevRangeByScoreWithCtx(ctx, key, opt, clients...)
}

// RedisZRangeByScore 返回有序集合中指定分数区间的成员列表。有序集成员按分数值递增(从小到大)次序排列。
// @param   key: 键
// @param min, max: 区间的取值使用闭区间 (小于等于或大于等于)，你也可以通过给参数前增加 ( 符号来使用可选的开区间 (小于或大于)
// @return <arg_1>: 值列表，类型为int
func RedisZRangeByScore(key string, opt *redis.ZRangeBy, clients ...redis.UniversalClient) ([]string, error) {
	return RedisZRangeByScoreWithCtx(otelwrap.Skip(), key, opt, clients...)
}

func RedisZRangeByScoreWithCtx(ctx context.Context, key string, opt *redis.ZRangeBy, clients ...redis.UniversalClient) ([]string, error) {
	return GetRedisCli(clients...).ZRangeByScore(ctx, key, opt).Result()
}

// RedisZRevRangeByScore 返回有序集中指定分数区间内的所有的成员。有序集成员按分数值递减(从大到小)的次序排列。
// @param   key: 键
// @param min, max: 区间的取值使用闭区间 (小于等于或大于等于)，你也可以通过给参数前增加 ( 符号来使用可选的开区间 (小于或大于)
// @return <arg_1>: 值列表，类型为int
func RedisZRevRangeByScore(key string, opt *redis.ZRangeBy, clients ...redis.UniversalClient) ([]string, error) {
	return RedisZRevRangeByScoreWithCtx(otelwrap.Skip(), key, opt, clients...)
}

func RedisZRevRangeByScoreWithCtx(ctx context.Context, key string, opt *redis.ZRangeBy, clients ...redis.UniversalClient) ([]string, error) {
	return GetRedisCli(clients...).ZRevRangeByScore(ctx, key, opt).Result()
}

// RedisZRangeByScoreAndOffsetCountWithScoresIsRev 返回有序集合中指定分数区间的成员列表。有序集成员按分数值递增(从小到大)次序排列。
// @param   key: 键
// @param min, max: 区间的取值使用闭区间 (小于等于或大于等于)，你也可以通过给参数前增加 ( 符号来使用可选的开区间 (小于或大于)
// @param offset, count: 偏移 和 获取的 数
// @param   isRev: 是否倒叙 false 按分数值递增(从小到大)来排序，true 按分数值递增(从大到小)来排序
// @return <arg_1>: 值列表 类型为 redis.Z, Score Member
func RedisZRangeByScoreAndOffsetCountWithScoresIsRev(key string, min, max string, offset, count int64, isRev bool, clients ...redis.UniversalClient) ([]redis.Z, error) {
	return RedisZRangeByScoreAndOffsetCountWithScoresIsRevWithCtx(otelwrap.Skip(), key, min, max, offset, count, isRev, clients...)
}

func RedisZRangeByScoreAndOffsetCountWithScoresIsRevWithCtx(ctx context.Context, key string, min, max string, offset, count int64, isRev bool, clients ...redis.UniversalClient) ([]redis.Z, error) {
	opt := &redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: offset,
		Count:  count,
	}

	if !isRev {
		return RedisZRangeByScoreWithScoresWithCtx(ctx, key, opt, clients...)
	}

	return RedisZRevRangeByScoreWithScoresWithCtx(ctx, key, opt, clients...)
}

// RedisZRangeByScoreWithScoresIsRev 返回有序集合中指定分数区间的成员列表。有序集成员按分数值递增(从小到大)次序排列。
// @param   key: 键
// @param min, max: 区间的取值使用闭区间 (小于等于或大于等于)，你也可以通过给参数前增加 ( 符号来使用可选的开区间 (小于或大于)
// @param   isRev: 是否倒叙 false 按分数值递增(从小到大)来排序，true 按分数值递增(从大到小)来排序
// @return <arg_1>: 值列表 类型为 redis.Z, Score Member
func RedisZRangeByScoreWithScoresIsRev(key string, min, max string, isRev bool, clients ...redis.UniversalClient) ([]redis.Z, error) {
	return RedisZRangeByScoreWithScoresIsRevWithCtx(otelwrap.Skip(), key, min, max, isRev, clients...)
}

func RedisZRangeByScoreWithScoresIsRevWithCtx(ctx context.Context, key string, min, max string, isRev bool, clients ...redis.UniversalClient) ([]redis.Z, error) {
	opt := &redis.ZRangeBy{
		Min: min,
		Max: max,
	}

	if !isRev {
		return RedisZRangeByScoreWithScoresWithCtx(ctx, key, opt, clients...)
	}

	return RedisZRevRangeByScoreWithScoresWithCtx(ctx, key, opt, clients...)
}

// RedisZRangeByScoreWithScores 返回有序集合中指定分数区间的成员列表。有序集成员按分数值递增(从小到大)次序排列。
// @param   key: 键
// @param min, max: 区间的取值使用闭区间 (小于等于或大于等于)，你也可以通过给参数前增加 ( 符号来使用可选的开区间 (小于或大于)
// @return <arg_1>: 值列表 类型为 redis.Z, Score Member
func RedisZRangeByScoreWithScores(key string, opt *redis.ZRangeBy, clients ...redis.UniversalClient) ([]redis.Z, error) {
	return RedisZRangeByScoreWithScoresWithCtx(otelwrap.Skip(), key, opt, clients...)
}

func RedisZRangeByScoreWithScoresWithCtx(ctx context.Context, key string, opt *redis.ZRangeBy, clients ...redis.UniversalClient) ([]redis.Z, error) {
	return GetRedisCli(clients...).ZRangeByScoreWithScores(ctx, key, opt).Result()
}

// RedisZRevRangeByScoreWithScores 返回有序集中指定分数区间内的所有的成员。有序集成员按分数值递减(从大到小)的次序排列。
// @param   key: 键
// @param min, max: 区间的取值使用闭区间 (小于等于或大于等于)，你也可以通过给参数前增加 ( 符号来使用可选的开区间 (小于或大于)
// @return <arg_1>: 值列表 类型为 redis.Z, Score Member
func RedisZRevRangeByScoreWithScores(key string, opt *redis.ZRangeBy, clients ...redis.UniversalClient) ([]redis.Z, error) {
	return RedisZRevRangeByScoreWithScoresWithCtx(otelwrap.Skip(), key, opt, clients...)
}

func RedisZRevRangeByScoreWithScoresWithCtx(ctx context.Context, key string, opt *redis.ZRangeBy, clients ...redis.UniversalClient) ([]redis.Z, error) {
	return GetRedisCli(clients...).ZRevRangeByScoreWithScores(ctx, key, opt).Result()
}
