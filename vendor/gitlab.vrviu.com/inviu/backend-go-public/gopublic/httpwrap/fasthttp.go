package httpwrap

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/buaazp/fasthttprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic/errwrap"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic/metricwrap"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic/metricwrap/fasthttpmetric"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic/otelwrap/fasthttpotel"
	vlog "gitlab.vrviu.com/inviu/backend-go-public/gopublic/vlog/tracing"
)

func NewFasthttpRouter() *fasthttprouter.Router {
	router := fasthttprouter.New()

	// 添加路由: metrics
	router.GET("/metrics", fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler()))

	return router
}

// Response 响应消息
func Response(ctx *fasthttp.RequestCtx, rsp interface{}) {
	bytes, _ := json.Marshal(rsp)

	// 返回结果
	ctx.Response.Header.AddBytesV("vrviu-mc-session-id", ctx.Request.Header.Peek("vrviu-mc-session-id"))
	ctx.Response.Header.Add("vrviu-mc-request-id", GenRequestID(ctx))
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetStatusCode(200)
	ctx.Response.SetBody(bytes)
	vlog.Infof(fasthttpotel.GetTraceCtx(ctx), "transmit response. rspbody(%s)", string(bytes))
}

func ResponseWithCode(ctx *fasthttp.RequestCtx, rsp interface{}, errcode int) {
	Response(ctx, rsp)

	builder := metricwrap.MetricLabelsBuilderFromContext(fasthttpotel.GetTraceCtx(ctx))
	builder.SetKV(metricwrap.MetricLabelCode, strconv.Itoa(errcode))
}

// Response2 响应消息
func Response2(ctx *fasthttp.RequestCtx, errcode int, errmsg string, body interface{}) {
	ResponseWithCode(ctx, HTTPResponse{
		Head: HTTPCommonHead{
			Code:      errcode,
			Msg:       errmsg,
			RequestID: GenRequestID(ctx),
		},
		Body: body,
	}, errcode)
}

func ResponseWrapMiddleWares(fn fasthttp.RequestHandler) fasthttp.RequestHandler {
	return MiddleWareRecover(
		fasthttpotel.MiddleWareTraceSpan(
			fasthttpmetric.RequestMetricWrap(
				fn,
			),
		),
	)
}

// ResponseWrap 返回响应语法糖：响应中有body字段
func ResponseWrap(fn func(*fasthttp.RequestCtx) (errcode int, errmsg string, body interface{})) fasthttp.RequestHandler {
	return ResponseWrapMiddleWares(func(ctx *fasthttp.RequestCtx) {
		vlog.Infof(fasthttpotel.GetTraceCtx(ctx), "recieve request. url(%s) param(%s) session_id(%s) request_id(%s)",
			string(ctx.URI().String()), string(ctx.Request.Body()), string(ctx.Request.Header.Peek("vrviu-mc-session-id")), GenRequestID(ctx))
		errcode, errmsg, body := fn(ctx)
		Response2(ctx, errcode, errmsg, body)
	})
}

// ResponseWrap2 返回响应语法糖：响应结构自定义
// @param rsp: 响应结构体
func ResponseWrap2(fn func(*fasthttp.RequestCtx) (rsp interface{})) fasthttp.RequestHandler {
	return ResponseWrapMiddleWares(func(ctx *fasthttp.RequestCtx) {
		vlog.Infof(fasthttpotel.GetTraceCtx(ctx), "recieve request. url(%s) param(%s) session_id(%s) request_id(%s)",
			string(ctx.URI().String()), string(ctx.Request.Body()), string(ctx.Request.Header.Peek("vrviu-mc-session-id")), GenRequestID(ctx))
		Response(ctx, fn(ctx))
	})
}

// ResponseWrap3 返回响应语法糖：响应中无body字段
// @param errcode:
// @param  errmsg:
// @param rspbody:
func ResponseWrap3(fn func(*fasthttp.RequestCtx) (errcode int, errmsg string, rspbody interface{})) fasthttp.RequestHandler {
	return ResponseWrapMiddleWares(func(ctx *fasthttp.RequestCtx) {
		vlog.Infof(fasthttpotel.GetTraceCtx(ctx), "recieve request. url(%s) param(%s) session_id(%s) request_id(%s)",
			string(ctx.URI().String()), string(ctx.Request.Body()), string(ctx.Request.Header.Peek("vrviu-mc-session-id")), GenRequestID(ctx))
		errcode, errmsg, rspbody := fn(ctx)

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

// ResponseWrap4 返回响应语法糖：通过err获取code和msg
// @param rsp: 响应结构体
func ResponseWrap4(fn func(*fasthttp.RequestCtx) (rspbody interface{}, err error)) fasthttp.RequestHandler {
	return ResponseWrapMiddleWares(func(ctx *fasthttp.RequestCtx) {
		vlog.Infof(fasthttpotel.GetTraceCtx(ctx), "recieve request. url(%s) param(%s) session_id(%s) request_id(%s)",
			string(ctx.URI().String()), string(ctx.Request.Body()), string(ctx.Request.Header.Peek("vrviu-mc-session-id")), GenRequestID(ctx))
		rspbody, err := fn(ctx)
		Response2(ctx, errwrap.Code(err), errwrap.Msg(err), rspbody)
	})
}

// ForwardResponseWrap 返回响应语法糖：响应body为[]byte
func ForwardResponseWrap(fn func(*fasthttp.RequestCtx) (code int, msg string, body []byte)) fasthttp.RequestHandler {
	return ResponseWrapMiddleWares(func(ctx *fasthttp.RequestCtx) {
		vlog.Infof(fasthttpotel.GetTraceCtx(ctx), "recieve request. url(%s) param(%s) session_id(%s) request_id(%s)",
			string(ctx.URI().String()), string(ctx.Request.Body()),
			string(ctx.Request.Header.Peek("vrviu-mc-session-id")), GenRequestID(ctx))

		statusCode, responseMsg, rspbody := fn(ctx)

		// 返回结果
		ctx.Response.Header.AddBytesV("vrviu-mc-session-id", ctx.Request.Header.Peek("vrviu-mc-session-id"))
		ctx.Response.Header.Add("vrviu-mc-request-id", GenRequestID(ctx))
		ctx.Response.Header.Set("vrviu-response-msg", responseMsg)
		ctx.Response.SetStatusCode(statusCode)
		if len(rspbody) > 0 {
			ctx.Response.SetBody(rspbody)
		}

		vlog.Infof(fasthttpotel.GetTraceCtx(ctx), "transmit response. rspbody(%s) session_id(%s) request_id(%s)", string(rspbody),
			string(ctx.Request.Header.Peek("vrviu-mc-session-id")), GenRequestID(ctx))
	})
}

// GetTagFields 获取结构体中有指定tag的所有field(不支持嵌套结构)
// @param object: 结构体指针
// @param    tag:
func GetTagFields(object interface{}, tag string) []reflect.StructField {
	fields := make([]reflect.StructField, 0)
	rt := reflect.TypeOf(object).Elem()

	for i := 0; i < rt.NumField(); i++ {
		rf := rt.Field(i)
		switch rf.Type.Kind() {
		case reflect.Struct:
			// 忽略除time.Time外的所有结构体
			if rf.Type.String() != "time.Time" {
				continue
			}

			tvalue := rf.Tag.Get(tag)
			if tvalue != "" && tvalue != "-" {
				fields = append(fields, rf)
			}
		default:
			tvalue := rf.Tag.Get(tag)
			if tvalue != "" && tvalue != "-" {
				fields = append(fields, rf)
			}
		}
	}

	return fields
}

// String2Value 将字符串转换为指定类型的值
func String2Value(s string, rtype reflect.Type) (interface{}, error) {
	switch rtype.Kind() {
	case reflect.String:
		return s, nil
	case reflect.Int:
		v, e := strconv.ParseInt(s, 10, 64)
		if e != nil {
			return nil, e
		}
		return int(v), nil
	case reflect.Int8:
		v, e := strconv.ParseInt(s, 10, 8)
		if e != nil {
			return nil, e
		}
		return int8(v), nil
	case reflect.Int16:
		v, e := strconv.ParseInt(s, 10, 16)
		if e != nil {
			return nil, e
		}
		return int16(v), nil
	case reflect.Int32:
		v, e := strconv.ParseInt(s, 10, 32)
		if e != nil {
			return nil, e
		}
		return int32(v), nil
	case reflect.Int64:
		v, e := strconv.ParseInt(s, 10, 64)
		if e != nil {
			return nil, e
		}
		return int64(v), nil
	case reflect.Uint:
		v, e := strconv.ParseUint(s, 10, 64)
		if e != nil {
			return nil, e
		}
		return uint(v), nil
	case reflect.Uint8:
		v, e := strconv.ParseUint(s, 10, 8)
		if e != nil {
			return nil, e
		}
		return uint8(v), nil
	case reflect.Uint16:
		v, e := strconv.ParseUint(s, 10, 16)
		if e != nil {
			return nil, e
		}
		return uint16(v), nil
	case reflect.Uint32:
		v, e := strconv.ParseUint(s, 10, 32)
		if e != nil {
			return nil, e
		}
		return uint32(v), nil
	case reflect.Uint64:
		v, e := strconv.ParseUint(s, 10, 64)
		if e != nil {
			return nil, e
		}
		return uint64(v), nil
	case reflect.Float32:
		v, e := strconv.ParseFloat(s, 32)
		if e != nil {
			return nil, e
		}
		return float32(v), nil
	case reflect.Float64:
		v, e := strconv.ParseFloat(s, 64)
		if e != nil {
			return nil, e
		}
		return float64(v), nil
	case reflect.Bool:
		v, e := strconv.ParseBool(s)
		if e != nil {
			return nil, e
		}
		return v, nil
	default:
		return nil, errors.New("unsupport type:" + rtype.String())
	}
}

// ParseQueryArgs 解析请求的查询参数到指定结构
// 字段类型仅支持基础类型及基础类型的数组；不支持struct嵌套；
// @param ctx:
// @param obj: struct对象指针，存放解析结果
func ParseQueryArgs(ctx *fasthttp.RequestCtx, obj interface{}) (err error) {
	if reflect.Ptr != reflect.TypeOf(obj).Kind() {
		return errors.New("obj is not a pointer")
	}

	rv := reflect.ValueOf(obj).Elem()

	// 解析查询参数
	rsfs := GetTagFields(obj, "queryarg")
	for _, rsf := range rsfs {
		tagvalues := strings.Split(rsf.Tag.Get("queryarg"), ",")

		if len(tagvalues) < 1 {
			return fmt.Errorf("invalid query arg `%s`", rsf.Name)
		}

		tag := tagvalues[0]
		required := gopublic.StringInArray("required", tagvalues[1:])

		if rsf.Type.Kind() == reflect.Array {
			args := ctx.QueryArgs().PeekMulti(tag)
			rarr := reflect.New(reflect.ArrayOf(0, rsf.Type.Elem()))
			for _, arg := range args {
				for _, s := range strings.Split(string(arg), ",") {
					if value, err := String2Value(s, rsf.Type.Elem()); err == nil {
						reflect.Append(rarr, reflect.ValueOf(value))
					} else {
						return err
					}
				}
			}
			rv.FieldByName(rsf.Name).Set(rarr)
		} else if rsf.Type.String() == "time.Time" {
			value := ctx.QueryArgs().Peek(tag)
			if t, err := time.Parse(time.RFC3339Nano, string(value)); err != nil {
				return err
			} else {
				rv.FieldByName(rsf.Name).Set(reflect.ValueOf(t))
			}
		} else {
			value := ctx.QueryArgs().Peek(tag)

			if value == nil && required {
				return fmt.Errorf("required query arg `%s` no value", rsf.Name)
			} else if value == nil && !required {
				continue
			}

			if value, err := String2Value(string(value), rsf.Type); err == nil {
				if reflect.ValueOf(value).Type().AssignableTo(rsf.Type) {
					rv.FieldByName(rsf.Name).Set(reflect.ValueOf(value))
				} else {
					rv.FieldByName(rsf.Name).Set(reflect.ValueOf(value).Convert(rsf.Type))
				}
			} else {
				return err
			}
		}
	}

	// 解析命名参数
	rsfs = GetTagFields(obj, "namedarg")
	for _, rsf := range rsfs {
		tag := strings.Split(rsf.Tag.Get("namedarg"), ",")[0]
		if value, err := String2Value(ctx.UserValue(tag).(string), rsf.Type); err == nil {
			rv.FieldByName(rsf.Name).Set(reflect.ValueOf(value))
		} else {
			return err
		}
	}

	return
}

// ParseNamedArgs 解析请求中的命名参数到指定结构
// @param ctx:
// @param obj: struct对象指针，存放解析结果
func ParseNamedArgs(ctx *fasthttp.RequestCtx, obj interface{}) (err error) {
	if reflect.Ptr != reflect.TypeOf(obj).Kind() {
		return errors.New("obj is not a pointer")
	}

	rv := reflect.ValueOf(obj).Elem()
	rsfs := GetTagFields(obj, "namedarg")
	for _, rsf := range rsfs {
		tagvalues := strings.Split(rsf.Tag.Get("namedarg"), ",")

		if len(tagvalues) != 1 && len(tagvalues) != 2 {
			return fmt.Errorf("invalid named arg `%s`", rsf.Name)
		}

		tag := tagvalues[0]
		required := len(tagvalues) == 2 && tagvalues[1] == "required"

		value := ctx.UserValue(tag)
		if value == nil && required {
			return fmt.Errorf("required named arg `%s` no value", rsf.Name)
		} else if value == nil && !required {
			continue
		}

		if value, err := String2Value(value.(string), rsf.Type); err == nil {
			rv.FieldByName(rsf.Name).Set(reflect.ValueOf(value))
		} else {
			return err
		}
	}

	return
}

var RecoverHandleFunc = func(ctx *fasthttp.RequestCtx) {}

func RegisterRecoverHandle(fn func(ctx *fasthttp.RequestCtx)) {
	RecoverHandleFunc = fn
}

func MiddleWareRecover(fn fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if err := recover(); err != nil {
				ctx.Response.SetStatusCode(http.StatusInternalServerError)
				vlog.Errorf(fasthttpotel.GetTraceCtx(ctx), "panic: [%v]\n\n%s", err, debug.Stack())
				RecoverHandleFunc(ctx)
			}
		}()

		fn(ctx)
	}
}
