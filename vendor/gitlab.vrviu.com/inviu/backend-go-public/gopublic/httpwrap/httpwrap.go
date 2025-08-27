package httpwrap

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	neturl "net/url"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/google/go-querystring/query"
	"github.com/levigross/grequests"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic/modcallwrap"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic/otelwrap"
	vlog "gitlab.vrviu.com/inviu/backend-go-public/gopublic/vlog/tracing"
)

func init() {
	if _sharaTp == nil {
		_sharaTp = &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   time.Duration(beego.AppConfig.DefaultInt("http::dail_timeout", 3000)) * time.Millisecond,
				KeepAlive: time.Duration(beego.AppConfig.DefaultInt("http::dail_keepalive", 30000)) * time.Millisecond,
			}).Dial,

			TLSClientConfig: &tls.Config{InsecureSkipVerify: beego.AppConfig.DefaultBool("http::insecure_skip_verify", false)},
		}
	}

	_sharaTp.DisableKeepAlives = false
	// 每个host最大连接数
	_sharaTp.MaxConnsPerHost = beego.AppConfig.DefaultInt("http::max_conns_per_host", 10)
	// 每个host最大空闲连接数
	_sharaTp.MaxIdleConnsPerHost = beego.AppConfig.DefaultInt("http::max_idle_conns_per_host", 10)
	// 最大空闲连接数
	_sharaTp.MaxIdleConns = beego.AppConfig.DefaultInt("http::max_idle_conns", 10)
	// 空闲连接超时时长
	_sharaTp.IdleConnTimeout = time.Duration(beego.AppConfig.DefaultInt("http::idle_conn_timeout", 180000)) * time.Millisecond
}

var _sharaTp *http.Transport

// HTTPCommonHead  响应消息头
type HTTPCommonHead struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	RequestID string `json:"request_id,omitempty"`
}

// HTTPResponse  响应消息
type HTTPResponse struct {
	Head HTTPCommonHead `json:"ret"`
	Body interface{}    `json:"body,omitempty"`
}

// MCParam 模调参数
type MCParam struct {
	Ctx         context.Context
	FlowID      string // 请求业务ID
	SessionID   string // 请求会话ID, 模调中使用，自动生成 gopublic.GenerateRandonString(12)。不用填，保留这个字段，以便原有代码能编译通过
	CallerID    int    // 主调服务ID
	CalleeID    int    // 被调服务ID
	InterfaceID int    // 被调接口ID
	VMID        int    // vmid
	AreaType    int    // 被调edge
}

type RequestOtelOptions struct {
	IgnoreResponseBody bool // 请求成功时不上报response body(head正常上报)
}

// RequestOptions 请求参数
type RequestOptions struct {
	grequests.RequestOptions                // http请求参数
	RequestOtelOptions                      // otel参数
	RequestName              string         // 请求名称
	MCParam                  MCParam        // 模调参数
	DisableReuseConn         bool           // 禁止复用长连接
	DisableModCall           bool           // 禁止上报模调
	DisableProcRspOnlyOK     bool           // 禁止http状态码非2xx时, 不解析返回体直接返回errcode=3008
	RC                       RetryCondition // 重试条件
	RT                       int            // 重试次数
	Method                   HTTPMethod     // 请求方法
	Verbose                  bool           // 打印详细日志
}

func (r *RequestOptions) Attributes() []attribute.KeyValue {
	attrs := []attribute.KeyValue{}

	if r.Data != nil {
		attrs = append(attrs, RequestOptionsDataKey.String(gopublic.ToJSON(r.Data)))
	}

	if r.Params != nil {
		attrs = append(attrs, RequestOptionsParamsKey.String(gopublic.ToJSON(r.Params)))
	}

	if r.QueryStruct != nil {
		attrs = append(attrs, RequestOptionsQueryStructKey.String(gopublic.ToJSON(r.QueryStruct)))
	}

	if r.JSON != nil {
		attrs = append(attrs, RequestOptionsJSONKey.String(gopublic.ToJSON(r.JSON)))
	}

	if r.XML != nil {
		attrs = append(attrs, RequestOptionsXMLKey.String(gopublic.ToJSON(r.XML)))
	}

	if r.Headers != nil {
		attrs = append(attrs, RequestOptionsHeadersKey.String(gopublic.ToJSON(r.Headers)))
	}

	if r.TLSHandshakeTimeout != 0 {
		attrs = append(attrs, RequestOptionsTLSHandshakeTimeoutKey.String(r.TLSHandshakeTimeout.String()))
	}

	if r.DialTimeout != 0 {
		attrs = append(attrs, RequestOptionsDialTimeoutKey.String(r.DialTimeout.String()))
	}

	if r.DialKeepAlive != 0 {
		attrs = append(attrs, RequestOptionsDialKeepAliveKey.String(r.DialKeepAlive.String()))
	}

	if r.RequestTimeout != 0 {
		attrs = append(attrs, RequestOptionsRequestTimeoutKey.String(r.RequestTimeout.String()))
	}

	attrs = append(attrs,
		RequestOptionsDisableReuseConnKey.Bool(r.DisableReuseConn),
		RequestOptionsDisableModCallKey.Bool(r.DisableModCall),
		RequestOptionsRTKey.Int(r.RT),
		RequestOptionsVerboseKey.Bool(r.Verbose),
	)

	return attrs
}

func (r *RequestOptions) SessionID() string {
	if r.MCParam.SessionID == "" {
		r.MCParam.SessionID = gopublic.GenerateRandonString(12)
	}
	return r.MCParam.SessionID
}

func (r *RequestOptions) Context() context.Context {
	if r.MCParam.Ctx == nil {
		r.MCParam.Ctx = otelwrap.NewSkipTraceCtx("RequestOptions_ctx_nil")
		vlog.Errorf(r.MCParam.Ctx, "RequestOptions.Ctx(). ctx is null")
	}
	return r.MCParam.Ctx
}

func (r *RequestOptions) SetCtx(ctx context.Context) {
	r.MCParam.Ctx = ctx
}

// procrsp 处理请求结果
func procrsp(grsp *grequests.Response, head interface{}, body interface{}) (errcode int, err error) {
	var rsp HTTPResponse
	err = grsp.JSON(&rsp)
	if err != nil && err != io.EOF {
		return 3009, err
	} else if err == io.EOF {
		return 3010, err
	}

	if body != nil && rsp.Body != nil {
		if err := json.Unmarshal([]byte(gopublic.ToJSON(rsp.Body)), body); err != nil {
			return 1101, err
		}
	}

	if head != nil {
		head.(*HTTPCommonHead).Code = rsp.Head.Code
		head.(*HTTPCommonHead).Msg = rsp.Head.Msg
		head.(*HTTPCommonHead).RequestID = rsp.Head.RequestID
	}

	if rsp.Head.Code != 0 {
		return 3011, fmt.Errorf(rsp.Head.Msg)
	}

	return 0, nil
}

// getCore get请求
func getCore(url string, ro *RequestOptions, head interface{}, body interface{}) (errcode int, err error) {
	// 发送请求
	grsp, e := grequests.Get(url, &ro.RequestOptions)
	if e != nil {
		err = e
		errcode = 3000
		return
	}
	defer grsp.Close()

	if !grsp.Ok && !ro.DisableProcRspOnlyOK {
		errcode = 3008
		err = fmt.Errorf("http get fail, [%d]%s %v", grsp.StatusCode, url, grsp.Error)
		return
	}

	return procrsp(grsp, head, body)
}

// postCore post请求
func postCore(url string, ro *RequestOptions, head interface{}, body interface{}) (errcode int, err error) {
	// 发送请求
	grsp, e := grequests.Post(url, &ro.RequestOptions)
	if e != nil {
		err = e
		errcode = 3001
		return
	}
	defer grsp.Close()

	if !grsp.Ok && !ro.DisableProcRspOnlyOK {
		errcode = 3008
		err = fmt.Errorf("http post fail, [%d]%s %v", grsp.StatusCode, url, grsp.Error)
		return
	}

	return procrsp(grsp, head, body)
}

// putCore put请求
func putCore(url string, ro *RequestOptions, head interface{}, body interface{}) (errcode int, err error) {
	// 发送请求
	grsp, e := grequests.Put(url, &ro.RequestOptions)
	if e != nil {
		err = e
		errcode = 3002
		return
	}
	defer grsp.Close()

	if !grsp.Ok && !ro.DisableProcRspOnlyOK {
		errcode = 3008
		err = fmt.Errorf("http put fail, [%d]%s %v", grsp.StatusCode, url, grsp.Error)
		return
	}

	return procrsp(grsp, head, body)
}

// deleteCore delete请求
func deleteCore(url string, ro *RequestOptions, head interface{}, body interface{}) (errcode int, err error) {
	grsp, e := grequests.Delete(url, &ro.RequestOptions)
	if e != nil {
		err = e
		errcode = 3003
		return
	}
	defer grsp.Close()

	if !grsp.Ok && !ro.DisableProcRspOnlyOK {
		errcode = 3008
		err = fmt.Errorf("http delete fail, [%d]%s %v", grsp.StatusCode, url, grsp.Error)
		return
	}

	return procrsp(grsp, head, body)
}

// GET 请求（HTTPRetry封装）
// @param           url: 请求URL
// @param            ro: 请求配置(request options)
// @param [out] rspHead: 存放返回响应head的指针 (HTTPCommonHead)
// @param [out] rspBody: 存放返回响应body的指针
func GET(url string, ro *RequestOptions, param ...interface{}) (errcode int, err error) {
	if ro.RC == nil {
		ro.RC = DefaultRC
	}
	return HTTPRetry(ro.RC, ro.RT, ro, HTTPGet, url, param...)
}

// POST 请求（HTTPRetry封装）
// @param           url: 请求URL
// @param            ro: 请求配置(request options)
// @param [out] rspHead: 存放返回响应head的指针 (HTTPCommonHead)
// @param [out] rspBody: 存放返回响应body的指针
func POST(url string, ro *RequestOptions, param ...interface{}) (errcode int, err error) {
	if ro.RC == nil {
		ro.RC = DefaultRC
	}
	return HTTPRetry(ro.RC, ro.RT, ro, HTTPPost, url, param...)
}

// PUT 请求（HTTPRetry封装）
// @param           url: 请求URL
// @param            ro: 请求配置(request options)
// @param [out] rspHead: 存放返回响应head的指针 (HTTPCommonHead)
// @param [out] rspBody: 存放返回响应body的指针
func PUT(url string, ro *RequestOptions, param ...interface{}) (errcode int, err error) {
	if ro.RC == nil {
		ro.RC = DefaultRC
	}
	return HTTPRetry(ro.RC, ro.RT, ro, HTTPPut, url, param...)
}

// DELETE 请求（HTTPRetry封装）
// @param           url: 请求URL
// @param            ro: 请求配置(request options)
// @param [out] rspHead: 存放返回响应head的指针 (HTTPCommonHead)
// @param [out] rspBody: 存放返回响应body的指针
func DELETE(url string, ro *RequestOptions, param ...interface{}) (errcode int, err error) {
	if ro.RC == nil {
		ro.RC = DefaultRC
	}
	return HTTPRetry(ro.RC, ro.RT, ro, HttpDelete, url, param...)
}

// CreateLLiveClient 创建http长链接客户端
// @param timeout: 请求超时时长
func CreateLLiveClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Transport: _sharaTp,
		Timeout:   timeout,
	}
}

// RetryCondition 重试条件计算方法
// @param   count: 已请求次数；从1开始计数（包括首次请求）
// @param  header: 请求返回的通用头
// @param errcode: 请求错误码
// @return: 当方法返回false时，继续重试；否则不再重试
type RetryCondition func(count int, header HTTPCommonHead, errcode int) bool

// DefaultRC 默认判断方法
func DefaultRC(_ int, header HTTPCommonHead, errcode int) bool {
	return header.Code == 0 && errcode == 0
}

// DefaultNoRC 默认业务错误不重试判断方法
func DefaultNoRC(_ int, header HTTPCommonHead, errcode int) bool {
	return header.Code != 0 || header.Code == 0 && errcode == 0
}

// HTTPMethod 已支持HTTP方法定义
type HTTPMethod int

// 已支持HTTP方法枚举
const (
	_          HTTPMethod = iota
	HTTPGet    HTTPMethod = 1 // GET 请求
	HTTPPost   HTTPMethod = 2 // POST 请求
	HTTPPut    HTTPMethod = 3 // PUT 请求
	HttpDelete HTTPMethod = 4 // DELETE 请求
)

func (m HTTPMethod) String() string {
	switch m {
	case HTTPGet:
		return "GET"
	case HTTPPost:
		return "POST"
	case HTTPPut:
		return "PUT"
	case HttpDelete:
		return "DELETE"
	default:
		return "UNKNOWN"
	}
}

// HTTPRetry HTTP重试请求
// @param            rc: 重试计算方法
// @param            rt: 重试请求次数（retry times）（包括首次请求）
// @param            ro: 请求配置（request options）
// @param        method: HTTP方法, 包括：httpwrap.HTTPGet, httpwrap.HTTPPost
// @param           url: 请求URL
// @param [out] rspHead: 存放返回响应head的指针（HTTPCommonHead）
// @param [out] rspBody: 存放返回响应body的指针
func HTTPRetry(rc RetryCondition, rt int, ro *RequestOptions, method HTTPMethod, url string, param ...interface{}) (errcode int, err error) {
	var head HTTPCommonHead

	if len(param) != 0 && len(param) != 2 {
		return 1000, errors.New("param not enough")
	} else if len(param) == 0 {
		param = append(param, nil, nil)
	}

	if rc == nil {
		rc = DefaultRC
	}

	urlobj, err := neturl.Parse(url)
	if err != nil {
		return 3004, err
	}

	if !otelwrap.IsSkip(ro.Context()) {
		spanName := ro.RequestName
		if len(spanName) == 0 {
			spanName = "NoName(HTTPRetry)"
			if len(urlobj.Path) < 20 {
				spanName = "(http) " + urlobj.Path
			}
		}

		attrs := []attribute.KeyValue{
			URLKey.String(url),
			MethodKey.String(method.String()),
			RetryTimesKey.Int(rt),
		}
		attrs = append(attrs, ro.Attributes()...)

		opts := []trace.SpanStartOption{
			trace.WithSpanKind(trace.SpanKindClient),
			trace.WithAttributes(attrs...),
		}

		ctx, span := otelwrap.GetTracer().Start(ro.Context(), spanName, opts...)
		defer func() {
			if err != nil {
				span.SetStatus(codes.Error, "HTTPRetry failed")
				span.RecordError(err, trace.WithAttributes(ErrCodeKey.Int(errcode), RetBodyKey.String(gopublic.ToJSON(param))))
			} else {
				attrs := []attribute.KeyValue{}
				if ro.IgnoreResponseBody {
					// 只上报head, 不上报body
					attrs = append(attrs, RetHeadKey.String(gopublic.ToJSON(head)))
				} else {
					attrs = append(attrs, RetBodyKey.String(gopublic.ToJSON(param)))
				}

				opts := []trace.EventOption{}
				opts = append(opts, trace.WithAttributes(attrs...))

				span.AddEvent("HTTPRetry succ", opts...)
			}
			span.End()
		}()
		ro.SetCtx(ctx)
	}

	// 打印请求详细日志
	if ro.Verbose {
		if param[1] == nil {
			param[1] = new(map[string]interface{})
		}

		defer func() {
			var queryags string
			if ro.QueryStruct != nil {
				if v, err := query.Values(ro.QueryStruct); err == nil {
					queryags = v.Encode()
				}
			} else {
				pairs := make([]string, 0, len(ro.Params))
				for k, v := range ro.Params {
					pairs = append(pairs, k+"="+v)
				}

				if len(pairs) > 0 {
					queryags = strings.Join(pairs, "&")
				}
			}
			vlog.Debugf(ro.Context(), "[httpwrap][%s][%s][%d] url(%s) request.queryargs(%s) request.payload(%s) response.header(%s) response.body(%s) errmsg(%v)",
				method, ro.MCParam.FlowID, ro.MCParam.InterfaceID, url, queryags, gopublic.ToJSON(ro.JSON), gopublic.ToJSON(head), gopublic.ToJSON(param[1]), err)
		}()
	}

	gro := &ro.RequestOptions
	startTime := time.Now()
	if gro.Headers == nil {
		gro.Headers = map[string]string{
			"vrviu-mc-flow-id":      ro.MCParam.FlowID,
			"vrviu-mc-session-id":   ro.SessionID(),
			"vrviu-mc-caller-id":    strconv.Itoa(int(ro.MCParam.CallerID)),
			"vrviu-mc-callee-id":    strconv.Itoa(int(ro.MCParam.CalleeID)),
			"vrviu-mc-interface-id": strconv.Itoa(int(ro.MCParam.InterfaceID)),
			"vrviu-mc-start-time":   startTime.Format("20060102_150405.000000000"),
		}
	} else {
		gro.Headers["vrviu-mc-flow-id"] = ro.MCParam.FlowID
		gro.Headers["vrviu-mc-session-id"] = ro.SessionID()
		gro.Headers["vrviu-mc-caller-id"] = strconv.Itoa(int(ro.MCParam.CallerID))
		gro.Headers["vrviu-mc-callee-id"] = strconv.Itoa(int(ro.MCParam.CalleeID))
		gro.Headers["vrviu-mc-interface-id"] = strconv.Itoa(int(ro.MCParam.InterfaceID))
		gro.Headers["vrviu-mc-start-time"] = startTime.Format("20060102_150405.000000000")
	}

	otel.GetTextMapPropagator().Inject(ro.Context(), propagation.MapCarrier(gro.Headers))

	// 默认使用长连接
	if beego.AppConfig.DefaultBool("http::conn_keepalive", true) && !ro.DisableReuseConn {
		gro.HTTPClient = CreateLLiveClient(ro.RequestTimeout)
	}

	// 上报模调信息
	if !ro.DisableModCall &&
		ro.MCParam.CalleeID != 0 &&
		ro.MCParam.CallerID != 0 &&
		ro.MCParam.InterfaceID != 0 {
		defer func(t time.Time) {
			modcallwrap.ReportModCallWithAreaTypeWithCtx(
				ro.Context(),
				urlobj.Host,
				ro.MCParam.CallerID,
				ro.MCParam.CalleeID,
				ro.MCParam.InterfaceID,
				ro.MCParam.FlowID,
				ro.SessionID(),
				head.RequestID,
				func() int {
					if errcode == 3011 {
						return head.Code
					}
					return errcode
				}(),
				func() string {
					if err != nil {
						return err.Error()
					}
					return "ok"
				}(),
				t,
				ro.MCParam.VMID,
				ro.MCParam.AreaType)
		}(startTime)
	}

	// 重试请求
	for _tk, ri := true, 0; _tk || (ri < rt && !rc(ri, head, errcode)); _tk, ri = false, ri+1 {
		switch method {
		case HTTPGet:
			errcode, err = getCore(url, ro, &head, param[1])
		case HTTPPost:
			errcode, err = postCore(url, ro, &head, param[1])
		case HTTPPut:
			errcode, err = putCore(url, ro, &head, param[1])
		case HttpDelete:
			errcode, err = deleteCore(url, ro, &head, param[1])
		default:
			return 3005, errors.New("unsupport method")
		}
	}

	// 最后一次回调
	rc(rt, head, errcode)

	// 填充响应头
	if param[0] != nil {
		param[0].(*HTTPCommonHead).Code = head.Code
		param[0].(*HTTPCommonHead).Msg = head.Msg
		param[0].(*HTTPCommonHead).RequestID = head.RequestID
	}

	return
}

// HTTPRetry2 HTTP重试请求
func HTTPRetry2(url string, ro *RequestOptions, param ...interface{}) (errcode int, err error) {
	if ro.RC == nil {
		ro.RC = DefaultRC
	}

	return HTTPRetry(ro.RC, ro.RT, ro, ro.Method, url, param...)
}
