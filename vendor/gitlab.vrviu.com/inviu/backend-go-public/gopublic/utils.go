package gopublic

import (
	"context"
	"encoding/json"
	"math/rand"
	"reflect"
	"sync"
	"time"

	vlog "gitlab.vrviu.com/inviu/backend-go-public/gopublic/vlog/tracing"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// UserTypeTest 测试用户类型
var UserTypeTest = 999 // 测试用户

// LetterRunes 随机字符串字符池
var LetterRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890")

// ToJSON 转json字符串包装
func ToJSON(object interface{}) string {
	bytes, _ := json.Marshal(object)
	return string(bytes)
}

// GenerateRandonString 生成随机字符串
func GenerateRandonString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = LetterRunes[rand.Intn(len(LetterRunes))]
	}
	return string(b)
}

// RegularTrigger 周期触发器
type RegularTrigger struct {
	ctx      context.Context
	name     string
	interval time.Duration
	ctrlCH   chan bool
	handler  func()
}

// CreateRegularTrigger 创建周期触发器
// @param     name: 触发器名称
// @param interval: 触发间隔
// @param  handler: 处理方法
func CreateRegularTrigger(name string, interval time.Duration, handler func()) *RegularTrigger {
	return &RegularTrigger{
		ctx:      vlog.NewTraceCtx(RandonStringWithPrefix("RegularTrigger", 10)),
		name:     name,
		interval: interval,
		handler:  handler,
	}
}

func (t *RegularTrigger) Ctx() context.Context {
	if t.ctx == nil {
		t.ctx = vlog.NewTraceCtx(RandonStringWithPrefix("RegularTrigger_ctx_nil", 10))
		vlog.Errorf(t.ctx, "RegularTrigger.Ctx(). ctx is null")
	}
	return t.ctx
}

// Start 启动触发器
func (t *RegularTrigger) Start() {
	// 判断任务是否已启动
	if t.ctrlCH != nil {
		vlog.Debugf(t.Ctx(), "RegularTrigger.Start(). task already running. [%s]", t.name)
		return
	}

	t.ctrlCH = make(chan bool, 1)
	vlog.Debugf(t.Ctx(), "RegularTrigger.Start() start rask. [%s]", t.name)

	// 定时检查
	go func() {
		t.handler()

		for {
			select {
			case <-time.After(t.interval):
				t.handler()
			case <-t.ctrlCH:
				vlog.Debugf(t.Ctx(), "RegularTrigger(). [%s] recv stop signal.", t.name)
				t.ctrlCH = nil
				return
			}
		}
	}()
}

// StartWithName 启动触发器
func (t *RegularTrigger) StartWithName(name string) {
	// 判断任务是否已启动
	if t.ctrlCH != nil {
		vlog.Debugf(t.Ctx(), "RegularTrigger.StartWithName(). task already running. [%s] ", t.name)
		return
	}

	t.name = name
	t.ctrlCH = make(chan bool, 1)
	vlog.Debugf(t.Ctx(), "RegularTrigger.StartWithName(). start task. [%s]", t.name)

	// 定时检查
	go func() {
		t.handler()

		for {
			select {
			case <-time.After(t.interval):
				t.handler()
			case <-t.ctrlCH:
				vlog.Debugf(t.Ctx(), "RegularTrigger(). recv stop signal. [%s] ", t.name)
				t.ctrlCH = nil
				return
			}
		}
	}()
}

// ChangeInterval 修改触发间隔
func (t *RegularTrigger) ChangeInterval(interval time.Duration) {
	t.interval = interval
	vlog.Debugf(t.Ctx(), "RegularTrigger.ChangeInterval(). [%s] interval [%d]", t.name, t.interval)
}

// Stop 停止触发器
func (t *RegularTrigger) Stop() {
	vlog.Debugf(t.Ctx(), "RegularTrigger.Stop(). stop task. [%s]", t.name)

	if t.ctrlCH == nil {
		return
	}

	select {
	case t.ctrlCH <- true:
	default:
		break
	}
}

// ParallelRoutine 并发执行任务
// @param parallel: 并发任务数量
// @param interval: 任务时间间隔
// @param     task: 任务列表
// @param       fn:
func ParallelRoutine(parallel int, interval time.Duration, task []interface{}, fn func(interface{}) bool) {
	ch := make(chan interface{}, len(task))
	for _, item := range task {
		ch <- item
	}

	wg := &sync.WaitGroup{}
	wg.Add(parallel)

	for i := 0; i < parallel; i++ {
		go func() {
			defer wg.Done()

			for {
				select {
				case item, ok := <-ch:
					if !ok || !fn(item) {
						return
					}

					// 间隔时间
					time.Sleep(interval)
				default:
					return
				}
			}
		}()
	}

	wg.Wait()
}

// ParallelRoutineWithChannel 并发执行任务
// @param      ctx: context
// @param parallel: 并发任务数量
// @param interval: 任务时间间隔
// @param     task: 任务列表
// @param       fn:
// ParallelRoutineWithChannel 并发执行任务
func ParallelRoutineWithChannel(ctx context.Context, parallel int, interval time.Duration, task <-chan interface{}, fn func(interface{}, int)) {
	wg := &sync.WaitGroup{}
	wg.Add(parallel)

	for i := 0; i < parallel; i++ {
		c, cancel := context.WithCancel(ctx)
		defer cancel()

		go func(ctx context.Context, id int) {
			defer wg.Done()

			for {
				select {
				case v, ok := <-task:
					if !ok {
						return
					}

					fn(v, id)
					time.Sleep(interval)
				case <-ctx.Done():
					return
				}
			}
		}(c, i)
	}

	wg.Wait()
}

// IsZero 判断val是否为Zero Value
func IsZero(val interface{}) bool {
	return reflect.Zero(reflect.TypeOf(val)).Interface() == val
}

// InArray 判断是否在数组中
func InArray(value int, arr []int) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}

	return false
}

// Uint64InArray 判断是否在数组中
func Uint64InArray(value uint64, arr []uint64) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}

	return false
}

// StringInArray 判断是否在数组中
func StringInArray(value string, arr []string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}

	return false
}

// RemoveFromStringArray 从数组中移除指定元素
func RemoveFromStringArray(arr []string, value ...string) []string {
	result := make([]string, 0)
	for _, v := range arr {
		if !StringInArray(v, value) {
			result = append(result, v)
		}
	}

	return result
}

// RemoveFromUint64Array 从数组中移除指定元素
func RemoveFromUint64Array(arr []uint64, value ...uint64) []uint64 {
	result := make([]uint64, 0)
	for _, v := range arr {
		if !Uint64InArray(v, value) {
			result = append(result, v)
		}
	}

	return result
}

func EqualIntArray(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// Deprecated
func Debugf(ctx context.Context, format string, args ...interface{}) {
	vlog.Debugf(ctx, format, args...)
}

// Deprecated
func PrintTraceLog(format string, args ...interface{}) {
	// if config.GetConfigManager().F.Section("log").Key("trace").MustBool(false) {
	vlog.Debugf(vlog.NewTraceCtx("trace_log"), format, args...)
	// return
	// }
}

// SmoothWeightedRobin 平滑加权轮询；该方案线程不安全，由调用者保证线程安全
// @param weight: 权重映射表 map<key: weight>
// @param  cache: 计算缓存表 map<key: value>，初始为value值均为0
// -
// @return <arg-1>: 本次命中的key
func SmoothWeightedRobin(weight, cache map[int]int) int {
	if len(weight) == 0 {
		return 0
	}

	if len(weight) == 1 {
		for k := range weight {
			return k
		}
	}

	totalweight := 0
	for k, v := range weight {
		if _, ok := cache[k]; !ok {
			cache[k] = 0
		}

		totalweight += v
	}

	var maxkey = 0
	var maxval = -totalweight
	for k, v := range cache {
		if w, ok := weight[k]; !ok {
			delete(cache, k)
		} else {
			cache[k] = v + w
			if cache[k] > maxval {
				maxkey = k
				maxval = cache[k]
			}
		}
	}

	cache[maxkey] = cache[maxkey] - totalweight
	return maxkey
}
