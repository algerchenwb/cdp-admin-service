package rediswrap

import (
	"context"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego"
	"github.com/go-redis/redis/v8"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic/otelwrap"
)

const Nil = redis.Nil

var (
	// RedisClient - 通用客户端
	_client         redis.UniversalClient
	_OnceInitClient sync.Once
)

func initClient() {
	_client = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:            strings.Split(beego.AppConfig.String("redis::hosts"), ","),
		Username:         beego.AppConfig.DefaultString("redis::username", ""),
		Password:         beego.AppConfig.DefaultString("redis::password", ""),
		ReadTimeout:      time.Duration(beego.AppConfig.DefaultInt("redis::read_timeout", 300)) * time.Millisecond,
		WriteTimeout:     time.Duration(beego.AppConfig.DefaultInt("redis::write_timeout", 300)) * time.Millisecond,
		PoolSize:         beego.AppConfig.DefaultInt("redis::conn_pool_capacity", 0),
		IdleTimeout:      time.Duration(beego.AppConfig.DefaultInt("redis::conn_pool_idle_timeout", 240)) * time.Second,
		MinIdleConns:     beego.AppConfig.DefaultInt("redis::conn_pool_min_idle", 0),
		MasterName:       beego.AppConfig.DefaultString("redis::mastername", ""),
		SentinelPassword: beego.AppConfig.DefaultString("redis::sentinel_password", ""),
	})
	if beego.AppConfig.DefaultBool("redis::enable_otel", false) {
		_client.AddHook(NewTracingHook())
	}
	// rdb.HMSet()
}

// GetRedisCli -
func GetRedisCli(clients ...redis.UniversalClient) redis.UniversalClient {
	if len(clients) != 0 && clients[0] != nil {
		return clients[0]
	}

	if _client == nil {
		_OnceInitClient.Do(func() {
			initClient()
		})
	}
	return _client
}

// RedisSet 设置键值
// @param     key: 键
// @param   value: 值
func RedisSet(key string, value interface{}, clients ...redis.UniversalClient) error {
	return RedisSetWithCtx(otelwrap.Skip(), key, value, clients...)
}

func RedisSetWithCtx(ctx context.Context, key string, value interface{}, clients ...redis.UniversalClient) error {
	set := GetRedisCli(clients...).Set(ctx, key, value, 0)
	if _, err := set.Result(); err != nil {
		return err
	}

	return nil
}

func RedisSetEX(key string, value interface{}, expirationSec int, clients ...redis.UniversalClient) error {
	return RedisSetEXWithCtx(otelwrap.Skip(), key, value, expirationSec, clients...)
}

func RedisSetEXWithCtx(ctx context.Context, key string, value interface{}, expirationSec int, clients ...redis.UniversalClient) error {
	//var sec time.Duration = time.Duration(expirationSec) * time.Second
	set := GetRedisCli(clients...).Set(ctx, key, value, time.Duration(expirationSec)*time.Second)
	if _, err := set.Result(); err != nil {
		return err
	}

	return nil
}

func RedisExists(keys []string, clients ...redis.UniversalClient) (int64, error) {
	return RedisExistsWithCtx(otelwrap.Skip(), keys, clients...)
}

func RedisExistsWithCtx(ctx context.Context, keys []string, clients ...redis.UniversalClient) (int64, error) {
	return GetRedisCli(clients...).Exists(ctx, keys...).Result()
}

func RedisDel(keys []string, clients ...redis.UniversalClient) (int64, error) {
	return RedisDelWithCtx(otelwrap.Skip(), keys, clients...)
}

func RedisDelWithCtx(ctx context.Context, keys []string, clients ...redis.UniversalClient) (int64, error) {
	return GetRedisCli(clients...).Del(ctx, keys...).Result()
}

// RedisGet 获取键值
// @param     key: 键
// @param   value: 值(string/[]byte/int/int8/int16/int32/int64//uint/uint8/uint16/uint32/float32/float64/bool/BinaryUnmarshaler类型的指针）

func RedisGet(key string, value interface{}, clients ...redis.UniversalClient) error {
	return RedisGetWithCtx(otelwrap.Skip(), key, value, clients...)
}

func RedisGetWithCtx(ctx context.Context, key string, value interface{}, clients ...redis.UniversalClient) error {
	get := GetRedisCli(clients...).Get(ctx, key)
	return get.Scan(value)
}

// RedisGetString 获取键值（字符串）
// @param key: 键
func RedisGetString(key string, clients ...redis.UniversalClient) (string, error) {
	return RedisGetStringWithCtx(otelwrap.Skip(), key, clients...)
}

func RedisGetStringWithCtx(ctx context.Context, key string, clients ...redis.UniversalClient) (string, error) {
	get := GetRedisCli(clients...).Get(ctx, key)
	return get.Result()
}

// RedisGetInt 获取键值（整形）
// @param key: 键
func RedisGetInt(key string, clients ...redis.UniversalClient) (int, error) {
	return RedisGetIntWithCtx(otelwrap.Skip(), key, clients...)
}

func RedisGetIntWithCtx(ctx context.Context, key string, clients ...redis.UniversalClient) (int, error) {
	get := GetRedisCli(clients...).Get(ctx, key)
	return get.Int()
}

// RedisHSet 设置哈希指定域的值
// @param   key: 键
// @param field: 域
// @param value: 值
func RedisHSet(key, field string, value interface{}, clients ...redis.UniversalClient) error {
	return RedisHSetWithCtx(otelwrap.Skip(), key, field, value, clients...)
}

func RedisHSetWithCtx(ctx context.Context, key, field string, value interface{}, clients ...redis.UniversalClient) error {
	if _, err := GetRedisCli(clients...).HSet(ctx, key, field, value).Result(); err != nil {
		return err
	}

	return nil
}

// RedisHSetRetCount 设置哈希指定域的值
// @param   key: 键
// @param field: 域
// @param value: 值
func RedisHSetRetCount(key, field string, value interface{}, clients ...redis.UniversalClient) (int64, error) {
	return RedisHSetRetCountWithCtx(otelwrap.Skip(), key, field, value, clients...)
}

func RedisHSetRetCountWithCtx(ctx context.Context, key, field string, value interface{}, clients ...redis.UniversalClient) (int64, error) {
	return GetRedisCli(clients...).HSet(ctx, key, field, value).Result()
}

// RedisHGet 获取哈希指定域的值
// @param     key: 键
// @param   field: 域
// @param   value: 值(string/[]byte/int/int8/int16/int32/int64//uint/uint8/uint16/uint32/float32/float64/bool/BinaryUnmarshaler类型的指针）
func RedisHGet(key, field string, value interface{}, clients ...redis.UniversalClient) error {
	return RedisHGetWithCtx(otelwrap.Skip(), key, field, value, clients...)
}

func RedisHGetWithCtx(ctx context.Context, key, field string, value interface{}, clients ...redis.UniversalClient) error {
	hget := GetRedisCli(clients...).HGet(ctx, key, field)
	return hget.Scan(value)
}

// RedisHGetString 获取哈希指定域的值（字符串）
// @param   key: 键
// @param field: 域
func RedisHGetString(key, field string, clients ...redis.UniversalClient) (string, error) {
	return RedisHGetStringWithCtx(otelwrap.Skip(), key, field, clients...)
}

func RedisHGetStringWithCtx(ctx context.Context, key, field string, clients ...redis.UniversalClient) (string, error) {
	hget := GetRedisCli(clients...).HGet(ctx, key, field)
	return hget.Result()
}

// RedisHGetInt 获取哈希指定域的值（字符串）
// @param   key: 键
// @param field: 域
func RedisHGetInt(key, field string, clients ...redis.UniversalClient) (int, error) {
	return RedisHGetIntWithCtx(otelwrap.Skip(), key, field, clients...)
}

func RedisHGetIntWithCtx(ctx context.Context, key, field string, clients ...redis.UniversalClient) (int, error) {
	hget := GetRedisCli(clients...).HGet(ctx, key, field)
	return hget.Int()
}

// RedisHMSet 批量设置哈希的域
// @param  key: 键
// @param pair: 域-值对
func RedisHMSet(key string, pair map[string]interface{}, clients ...redis.UniversalClient) error {
	return RedisHMSetWithCtx(otelwrap.Skip(), key, pair, clients...)
}

func RedisHMSetWithCtx(ctx context.Context, key string, pair map[string]interface{}, clients ...redis.UniversalClient) error {
	if _, err := GetRedisCli(clients...).HSet(ctx, key, pair).Result(); err != nil {
		return err
	}

	return nil
}

// RedisHMGet 批量获取哈希指定域的值
// @param   key: 键
// @param field: 域
// @return <arg_1>: 值列表，类型为string或nil
func RedisHMGet(key string, field []string, clients ...redis.UniversalClient) ([]interface{}, error) {
	return RedisHMGetWithCtx(otelwrap.Skip(), key, field, clients...)
}

func RedisHMGetWithCtx(ctx context.Context, key string, field []string, clients ...redis.UniversalClient) ([]interface{}, error) {
	return GetRedisCli(clients...).HMGet(ctx, key, field...).Result()
}

// DeepRedisFields 获取结构体的所有field
func DeepRedisFields(object interface{}) []reflect.StructField {
	fields := make([]reflect.StructField, 0)
	rt := reflect.TypeOf(object).Elem()

	for i := 0; i < rt.NumField(); i++ {
		rf := rt.Field(i)
		switch rf.Type.Kind() {
		case reflect.Struct:
			if rf.Anonymous {
				rv := reflect.ValueOf(object).Elem()
				fields = append(fields, DeepRedisFields(rv.Field(i).Addr().Interface())...)
			} else if reflect.TypeOf((*time.Time)(nil)).Elem() == rf.Type {
				fields = append(fields, rf)
			}
		default:
			tag := rf.Tag.Get("redis")
			if tag != "" && tag != "-" {
				fields = append(fields, rf)
			}
		}
	}

	return fields
}

// RedisHMGetObject 从redis读取并存储到object中
// object字段支持boolean/string/integer/float/time.Time五基础类型，其中time.Time的格式必须是`2006-01-02 15:04:05.999999999 -0700 MST`
// 使用`redis`作为tag标记要查询的字段，支持嵌套结构体
// 若字段在redis中不存在，则默认字段值为zero value;
// 若字段的tag中指定`required`或`read_required`，则当字段在redis中不存在时，方法返回gopublic.ErrNotExist
// @param     key:
// @param  object: 存储返回结果的结构体对象指针
// @return errCode: 错误码
// @return     err: 错误
func RedisHMGetObject(key string, object interface{}, clients ...redis.UniversalClient) error {
	return RedisHMGetObjectWithCtx(otelwrap.Skip(), key, object, clients...)
}

func RedisHMGetObjectWithCtx(ctx context.Context, key string, object interface{}, clients ...redis.UniversalClient) error {
	rvo := reflect.ValueOf(object)
	if rvo.Kind() != reflect.Ptr || rvo.IsNil() {
		return errors.New("invalid pointer to store result")
	}

	var fieldTypes []reflect.Type
	var fieldNames []string
	var fieldTags []string
	var filedReqs []bool

	rsfs := DeepRedisFields(object)
	for _, rsf := range rsfs {
		tags := strings.Split(rsf.Tag.Get("redis"), ",")
		fieldTags = append(fieldTags, tags[0])
		filedReqs = append(filedReqs, (gopublic.StringInArray("required", tags) || gopublic.StringInArray("read_required", tags)))
		fieldTypes = append(fieldTypes, rsf.Type)
		fieldNames = append(fieldNames, rsf.Name)
	}

	values, err := RedisHMGetWithCtx(ctx, key, fieldTags, clients...)
	if err != nil {
		return err
	}

	rv := reflect.ValueOf(object).Elem()
	for i, value := range values {
		// 判断是否required字符不存在
		if value == nil && filedReqs[i] {
			return gopublic.ErrNotExist
		}

		switch fieldTypes[i].Kind() {
		case reflect.String:
			if value == nil {
				rv.FieldByName(fieldNames[i]).SetString("")
			} else {
				rv.FieldByName(fieldNames[i]).SetString(value.(string))
			}
		case reflect.Int,
			reflect.Int8,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64:
			if value == nil {
				rv.FieldByName(fieldNames[i]).SetInt(0)
			} else {
				if v, err := strconv.ParseInt(value.(string), 10, 64); err == nil {
					rv.FieldByName(fieldNames[i]).SetInt(v)
				} else {
					return err
				}
			}
		case reflect.Uint,
			reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64:
			if value == nil {
				rv.FieldByName(fieldNames[i]).SetUint(0)
			} else {
				if v, err := strconv.ParseUint(value.(string), 10, 64); err == nil {
					rv.FieldByName(fieldNames[i]).SetUint(v)
				} else {
					return err
				}
			}
		case reflect.Float32,
			reflect.Float64:
			if value == nil {
				rv.FieldByName(fieldNames[i]).SetFloat(0.0)
			} else {
				if v, err := strconv.ParseFloat(value.(string), 64); err == nil {
					rv.FieldByName(fieldNames[i]).SetFloat(v)
				} else {
					return err
				}
			}
		case reflect.Bool:
			if value == nil {
				rv.FieldByName(fieldNames[i]).SetBool(false)
			} else {
				if v, err := strconv.ParseBool(value.(string)); err == nil {
					rv.FieldByName(fieldNames[i]).SetBool(v)
				} else {
					return err
				}
			}
		case reflect.Struct:
			if reflect.TypeOf((*time.Time)(nil)).Elem() == fieldTypes[i] {
				if value == nil {
					rv.FieldByName(fieldNames[i]).Set(reflect.ValueOf(time.Unix(0, 0)))
				} else {
					t, _ := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", value.(string))
					rv.FieldByName(fieldNames[i]).Set(reflect.ValueOf(t))
				}
			} else {
				return errors.New("unsupported struct")
			}
		default:
			continue
		}
	}

	return nil
}

// RedisHMSetObject 将object存储到redis中
// 使用`redis`作为tag标记要存储的字段，支持嵌套结构体
// 若值为zero value，且未标记`required`或`write_required`，那么字段将被忽略不写入redis；否则以zero value将字段写入redis
// 仅支持boolean/string/integer/float/time.Time五种基础类型，其中time.Time的时间格式为`2006-01-02 15:04:05.999999999 -0700 MST`
// @param     key:
// @param  object: 存入redis的结构体对象
// @return errCode: 错误码
// @return     err: 错误
func RedisHMSetObject(key string, object interface{}, clients ...redis.UniversalClient) error {
	return RedisHMSetObjectWithCtx(otelwrap.Skip(), key, object, clients...)
}

func RedisHMSetObjectWithCtx(ctx context.Context, key string, object interface{}, clients ...redis.UniversalClient) error {
	rvo := reflect.ValueOf(object)
	var rvp reflect.Value

	if object == nil {
		return errors.New("nil object")
	} else if rvo.Kind() == reflect.Ptr {
		if rvo.IsNil() {
			return errors.New("nil pointer")
		}

		rvo = rvo.Elem()
		rvp = reflect.New(reflect.TypeOf(object).Elem())
	} else if rvo.Kind() != reflect.Struct {
		return errors.New("object is not a struct")
	} else {
		rvp = reflect.New(reflect.TypeOf(object))
	}

	pair := map[string]interface{}{}
	rvp.Elem().Set(rvo)
	rsfs := DeepRedisFields(rvp.Interface())
	for _, rsf := range rsfs {
		tags := strings.Split(rsf.Tag.Get("redis"), ",")
		rvf := rvo.FieldByName(rsf.Name)

		// ignore zero value filed if not tagged by `required`
		if (!gopublic.StringInArray("required", tags) && !gopublic.StringInArray("write_required", tags)) && rvf.Interface() == reflect.Zero(rvf.Type()).Interface() {
			continue
		}

		if reflect.TypeOf((*time.Time)(nil)).Elem() == rsf.Type {
			pair[tags[0]] = rvf.Interface().(time.Time).Format("2006-01-02 15:04:05.999999999 -0700 MST")
		} else {
			pair[tags[0]] = rvf.Interface()
		}
	}

	return RedisHMSetWithCtx(ctx, key, pair, clients...)
}

// RedisExpire 设置key过期时长
func RedisExpire(key string, expiration time.Duration, clients ...redis.UniversalClient) error {
	return RedisExpireWithCtx(otelwrap.Skip(), key, expiration, clients...)
}

func RedisExpireWithCtx(ctx context.Context, key string, expiration time.Duration, clients ...redis.UniversalClient) error {
	_, err := GetRedisCli(clients...).Expire(ctx, key, expiration).Result()
	return err
}

// RedisExpireAt 设置过期时间
// @param     key: 键
// @param      tm: 过期时间
func RedisExpireAt(key string, tm time.Time, clients ...redis.UniversalClient) (bool, error) {
	return RedisExpireAtWithCtx(otelwrap.Skip(), key, tm, clients...)
}

func RedisExpireAtWithCtx(ctx context.Context, key string, tm time.Time, clients ...redis.UniversalClient) (bool, error) {
	return GetRedisCli(clients...).ExpireAt(ctx, key, tm).Result()
}

// RedisPersist 移除过期时间
// @param     key: 键
func RedisPersist(key string, clients ...redis.UniversalClient) (bool, error) {
	return RedisPersistWithCtx(otelwrap.Skip(), key, clients...)
}

func RedisPersistWithCtx(ctx context.Context, key string, clients ...redis.UniversalClient) (bool, error) {
	return GetRedisCli(clients...).Persist(ctx, key).Result()
}

func RedisPipeline(clients ...redis.UniversalClient) redis.Pipeliner {
	return GetRedisCli(clients...).Pipeline()
}

type RunPipelineFunc func(context.Context, redis.Pipeliner) (interface{}, error)

func RedisRunPipeline(fn RunPipelineFunc, clients ...redis.UniversalClient) (interface{}, error) {
	return RedisRunPipelineWithCtx(otelwrap.Skip(), fn, clients...)
}

func RedisRunPipelineWithCtx(ctx context.Context, fn RunPipelineFunc, clients ...redis.UniversalClient) (interface{}, error) {
	return fn(ctx, RedisPipeline(clients...))
}
