package middleware

import (
	"bytes"
	"cdp-admin-service/internal/helper"
	"cdp-admin-service/internal/helper/cache"
	table "cdp-admin-service/internal/helper/dal"
	"cdp-admin-service/internal/helper/queue"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
)

type LogwayHandleMiddleware struct {
	cache *cache.Cache
	queue *queue.Queue
}

func NewLogwayHandleMiddleware(cache *cache.Cache, queue *queue.Queue) *LogwayHandleMiddleware {
	return &LogwayHandleMiddleware{
		cache: cache,
		queue: queue,
	}
}

func (m *LogwayHandleMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		sessionId, ctx := helper.GenSessionId(r.Context())

		startTime := time.Now()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logx.WithContext(ctx).Errorf("[%s] Failed to read request body: %v", sessionId, err)
		}

		logx.WithContext(ctx).Debugf("[%s][Request]: %s %s %+v %s", sessionId, r.Method, r.RequestURI, r.Header, body)
		// 创建一个新的请求主体用于后续读取
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			body:           make([]byte, 0),
		}
		next(recorder, r.WithContext(ctx))

		logx.WithContext(ctx).Infof("[%s][Response]: cost:%d, body:%s", helper.GetSessionId(ctx), time.Since(startTime).Milliseconds(), string(recorder.body))

		perm := helper.PermStandardBase64(r.URL.Path)
		isLogging, err := m.cache.Get(ctx, sessionId, cache.PermIsLoggingKey(perm))
		if err != nil && !errors.Is(err, redis.Nil) {
			logx.WithContext(ctx).Debugf("[%s] 日志记录失败 err[%v]", sessionId, err)
			return
		}
		if isLogging == fmt.Sprintf("%d", table.LoggingEnable) {
			account := GetUserAccount(recorder)
			platform := GetUserPlatform(recorder)
			userId := GetUserId(recorder)
			m.queue.Producer(table.TCdpSysLog{
				UserId:     uint32(userId),
				Account:    account,
				Platform:   platform,
				Ip:         r.RemoteAddr,
				Uri:        r.RequestURI,
				Type:       2,
				Request:    string(body),
				Response:   string(recorder.body),
				CreateTime: time.Now(),
				UpdateTime: time.Now(),
				ModifyTime: time.Now(),
				Status:     1,
			})
		}

	}
}

// 自定义的 ResponseWriter
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

// WriteHeader 重写 WriteHeader 方法，捕获状态码
func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

// 重写 Write 方法，捕获响应数据
func (r *responseRecorder) Write(body []byte) (int, error) {
	r.body = body
	return r.ResponseWriter.Write(body)
}
