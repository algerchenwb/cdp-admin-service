package httpwrap

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic/otelwrap/fasthttpotel"
	vlog "gitlab.vrviu.com/inviu/backend-go-public/gopublic/vlog/tracing"
)

// HandlerFunc1 -
type HandlerFunc1 func(*fasthttp.RequestCtx) (errcode int, errmsg string, body interface{})

// HandlerFunc2 -
type HandlerFunc2 func(*fasthttp.RequestCtx) (rsp interface{})

// HandlerFunc3 -
type HandlerFunc3 func(*fasthttp.RequestCtx) (errcode int, errmsg string, rspbody interface{})

type MiddlewareWrap func(fasthttp.RequestHandler) fasthttp.RequestHandler

func (m MiddlewareWrap) HandlerFunc(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return m(handler)
}

// MiddleWareChain -
type MiddleWareChain struct {
	l      sync.Mutex
	layers []Middleware
}

// Append -
func (c *MiddleWareChain) Append(m ...Middleware) {
	c.l.Lock()
	defer c.l.Unlock()
	c.layers = append(c.layers, m...)
}

// Apply -
func (c *MiddleWareChain) Apply(fn fasthttp.RequestHandler) fasthttp.RequestHandler {
	c.l.Lock()
	defer c.l.Unlock()

	for i := len(c.layers) - 1; i > -1; i-- {
		fn = c.layers[i].HandlerFunc(fn)
	}

	return fn
}

// Response1 -
func (c *MiddleWareChain) Response1(h HandlerFunc1) fasthttp.RequestHandler {
	return c.Apply(func(ctx *fasthttp.RequestCtx) {
		errcode, errmsg, body := h(ctx)
		Response2(ctx, errcode, errmsg, body)
	})
}

// Response2 -
func (c *MiddleWareChain) Response2(h HandlerFunc2) fasthttp.RequestHandler {
	return c.Apply(func(ctx *fasthttp.RequestCtx) {
		Response(ctx, h(ctx))
	})
}

// Response3 -
func (c *MiddleWareChain) Response3(h HandlerFunc3) fasthttp.RequestHandler {
	return c.Apply(func(ctx *fasthttp.RequestCtx) {
		errcode, errmsg, rspbody := h(ctx)

		response := make(map[string]interface{})

		if rspbody != nil {
			json.Unmarshal([]byte(gopublic.ToJSON(rspbody)), &response)
		}

		response["ret"] = HTTPCommonHead{
			Code:      errcode,
			Msg:       errmsg,
			RequestID: GenRequestID(ctx),
		}

		ResponseWithCode(ctx, response, errcode)
	})
}

// Middleware -
type Middleware interface {
	HandlerFunc(fasthttp.RequestHandler) fasthttp.RequestHandler
}

// NewTraceRequestMiddleWare -
func NewTraceRequestMiddleWare() *TraceRequestMiddleWare {
	return &TraceRequestMiddleWare{}
}

// TraceRequestMiddleWare 记录请求中间层
type TraceRequestMiddleWare struct{}

// HandlerFunc -
func (w *TraceRequestMiddleWare) HandlerFunc(fn fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		vlog.Infof(fasthttpotel.GetTraceCtx(ctx), "recieve request. url(%s) param(%s) session_id(%s) request_id(%s)",
			string(ctx.URI().String()), string(ctx.Request.Body()), string(ctx.Request.Header.Peek("vrviu-mc-session-id")), GenRequestID(ctx))

		fn(ctx)
	})
}

// NewResponseWrapMiddleWares -
func NewResponseWrapMiddleWares() Middleware {
	return MiddlewareWrap(ResponseWrapMiddleWares)
}

// NewAuthMiddleWare -
func NewAuthMiddleWare(helper SecretHelper) *AuthMiddleWare {
	return &AuthMiddleWare{secretHelper: helper}
}

// AuthMiddleWare AuthMiddleWare 鉴权中间层
type AuthMiddleWare struct {
	secretHelper SecretHelper
}

// extractAuthArgs 提取鉴权请求参数
func (w *AuthMiddleWare) extractAuthArgs(ctx *fasthttp.RequestCtx) map[string]string {
	m := make(map[string]string)

	// 遍历所有查询参数
	ctx.QueryArgs().VisitAll(func(key, value []byte) {
		m[string(key)] = string(value)
	})

	return m
}

// HandlerFunc -
func (w *AuthMiddleWare) HandlerFunc(fn fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		queryargs := w.extractAuthArgs(ctx)

		if v, ok := queryargs["AccessKey"]; !ok && v == "" {
			ResponseWithCode(ctx, HTTPResponse{Head: HTTPCommonHead{Code: 3013, Msg: "empty `AccessKey`"}}, 3013)
			return
		}

		var signature string
		if v, ok := queryargs["Signature"]; ok {
			signature = v
		}
		delete(queryargs, "Signature")

		// 秘钥计算版本
		var urlPath = "/"
		if v, ok := queryargs["SignatureVersion"]; ok && v == SignatureVersionV1_1 {
			urlPath = string(ctx.Path())
		}

		// 获取秘钥
		var apargs AuthPublicParam
		if err := ParseQueryArgs(ctx, &apargs); err != nil {
			ResponseWithCode(ctx, HTTPResponse{Head: HTTPCommonHead{Code: 3013, Msg: err.Error()}}, 3013)
			return
		}

		secret, err := w.secretHelper.GetSecret(&apargs)
		if err != nil {
			ResponseWithCode(ctx, HTTPResponse{Head: HTTPCommonHead{Code: 3012, Msg: err.Error()}}, 3012)
			return
		}

		if secret.CheckExpireTime == 1 && time.Now().After(secret.ExpireTime) {
			ResponseWithCode(ctx, HTTPResponse{Head: HTTPCommonHead{Code: 3012, Msg: fmt.Sprintf("auth info expire. time limit:%v", secret.ExpireTime)}}, 3012)
			return
		}

		// 修正BizID使用请求
		if apargs.BizId != "" {
			secret.BizId = apargs.BizId
		}

		// 设置权限信息
		ctx.SetUserValue(VrviuPermInfoCtxCacheKey, secret)

		// 校验权限
		if secret.Permission != 1 {
			ResponseWithCode(ctx, HTTPResponse{Head: HTTPCommonHead{Code: 3012, Msg: "no permission"}}, 3012)
			return
		}

		// 判断是否忽略签名校验
		if secret.IgnoreSignature == 1 {
			fn(ctx)
			return
		}

		if err := Authenticate(ctx, string(ctx.Request.Header.Method()), secret.AccessKeySecret, signature, urlPath, queryargs, ctx.Request.Body()); err != nil {
			ResponseWithCode(ctx, HTTPResponse{Head: HTTPCommonHead{Code: 3012, Msg: err.Error()}}, 3012)
			return
		}

		fn(ctx)
	}
}

// GetPermInfo 获取权限信息
func GetPermInfo(ctx *fasthttp.RequestCtx) *PermInfo {
	switch value := ctx.UserValue(VrviuPermInfoCtxCacheKey).(type) {
	case *PermInfo:
		return value
	default:
		return nil
	}
}

// GetPermBizType 获取权限信息：业务类型
func GetPermBizType(ctx *fasthttp.RequestCtx, defaultValue int) int {
	permInfo := GetPermInfo(ctx)
	if permInfo == nil {
		return defaultValue
	}

	return permInfo.BizType
}
