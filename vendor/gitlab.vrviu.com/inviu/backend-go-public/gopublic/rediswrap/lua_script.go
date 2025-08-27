package rediswrap

import (
	"context"

	"github.com/go-redis/redis/v8"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic/otelwrap"
)

// RedisLuaScriptRun Redis 脚本使用 Lua 解释器来执行脚本
// @param   luaScript: 脚本内容
// @param keys: 脚本中使用的 key
// @param args: 脚本中使用的参数
// @return <arg_1>: 值列表，类型为interface{}
func RedisLuaScriptRun(luaScript string, keys []string, args []interface{}, clients ...redis.UniversalClient) (interface{}, error) {
	return RedisLuaScriptRunWithCtx(otelwrap.Skip(), luaScript, keys, args, clients...)
}

func RedisLuaScriptRunWithCtx(ctx context.Context, luaScript string, keys []string, args []interface{}, clients ...redis.UniversalClient) (interface{}, error) {
	script := redis.NewScript(luaScript)
	return script.Run(ctx, GetRedisCli(clients...), keys, args...).Result()
}

func RedisEval(script string, keys []string, args []interface{}, clients ...redis.UniversalClient) *redis.Cmd {
	return RedisEvalWithCtx(otelwrap.Skip(), script, keys, args, clients...)
}

func RedisEvalWithCtx(ctx context.Context, script string, keys []string, args []interface{}, clients ...redis.UniversalClient) *redis.Cmd {
	return GetRedisCli(clients...).Eval(ctx, script, keys, args...)
}

func RedisEvalInt64(script string, keys []string, args []interface{}, clients ...redis.UniversalClient) (int64, error) {
	return RedisEvalInt64WithCtx(otelwrap.Skip(), script, keys, args, clients...)
}

func RedisEvalInt64WithCtx(ctx context.Context, script string, keys []string, args []interface{}, clients ...redis.UniversalClient) (int64, error) {
	return RedisEvalWithCtx(ctx, script, keys, args, clients...).Int64()
}

func RedisEvalUint64(script string, keys []string, args []interface{}, clients ...redis.UniversalClient) (uint64, error) {
	return RedisEvalUint64WithCtx(otelwrap.Skip(), script, keys, args, clients...)
}

func RedisEvalUint64WithCtx(ctx context.Context, script string, keys []string, args []interface{}, clients ...redis.UniversalClient) (uint64, error) {
	return RedisEvalWithCtx(ctx, script, keys, args, clients...).Uint64()
}
