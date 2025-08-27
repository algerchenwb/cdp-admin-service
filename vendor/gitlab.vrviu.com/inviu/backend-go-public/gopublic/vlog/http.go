package vlog

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"gitlab.vrviu.com/inviu/backend-go-public/go.uber.org/zap"
	//"go.uber.org/zap"
	"gitlab.vrviu.com/inviu/backend-go-public/go.uber.org/zap/zapcore"
)

// HTTP adds the correct "HTTP" field.
func HTTP(req *HTTPPayload) zap.Field {
	return zap.Object("httpRequest", req)
}

// HTTP Result
type HTTPRet struct {
	// ret.code in response body
	RetCode int `json:"retcode"`

	// ret.msg in response body
	RetMsg string `json:"retmsg"`

	// ret.request_id in response body
	RetRequestID string `json:"retrequestid"`
}

// HTTPPayload is the complete payload that can be interpreted by
// vlog as a HTTP request.
type HTTPPayload struct {
	// The request method. Examples: "GET", "HEAD", "PUT", "POST".
	RequestMethod string `json:"requestMethod"`

	// The scheme (http, https), the host name, the path and the query portion of
	// the URL that was requested.
	//
	// Example: "http://example.com/some/info?color=red".
	RequestURL string `json:"requestUrl"`

	// The size of the HTTP request message in bytes, including the request
	// headers and the request body.
	RequestSize string `json:"requestSize"`

	// The response code indicating the status of response.
	//
	// Examples: 200, 404.
	Status int `json:"status"`

	// The size of the HTTP response message sent back to the client, in bytes,
	// including the response headers and the response body.
	ResponseSize string `json:"responseSize"`

	// The user agent sent by the client.
	//
	// Example: "Mozilla/4.0 (compatible; MSIE 6.0; Windows 98; Q312461; .NET CLR 1.0.3705)".
	UserAgent string `json:"userAgent"`

	// The IP address (IPv4 or IPv6) of the client that issued the HTTP request.
	//
	// Examples: "192.168.1.1", "FE80::0202:B3FF:FE1E:8329".
	RemoteIP string `json:"remoteIp"`

	// The IP address (IPv4 or IPv6) of the origin server that the request was
	// sent to.
	ServerIP string `json:"serverIp"`

	// The referrer URL of the request, as defined in HTTP/1.1 Header Field
	// Definitions.
	Referer string `json:"referer"`

	// The request processing latency on the server, from the time the request was
	// received until the response was sent.
	//
	// A duration in seconds with up to nine fractional digits, terminated by 's'.
	//
	// Example: "3.5s".
	Latency string `json:"latency"`

	// Whether or not a cache lookup was attempted.
	CacheLookup bool `json:"cacheLookup"`

	// Whether or not an entity was served from cache (with or without
	// validation).
	CacheHit bool `json:"cacheHit"`

	// Whether or not the response was validated with the origin server before
	// being served from cache. This field is only meaningful if cacheHit is True.
	CacheValidatedWithOriginServer bool `json:"cacheValidatedWithOriginServer"`

	// The number of HTTP response bytes inserted into cache. Set only when a
	// cache fill was attempted.
	CacheFillBytes string `json:"cacheFillBytes"`

	// Protocol used for the request.
	//
	// Examples: "HTTP/1.1", "HTTP/2", "websocket"
	Protocol string `json:"protocol"`

	// ret.code in response body
	RetCode int `json:"retcode"`

	// ret.msg in response body
	RetMsg string `json:"retmsg"`

	// ret.request_id in response body
	RetRequestID string `json:"retrequestid"`
}

// NewHTTP returns a new HTTPPayload struct, based on the passed
// in http.Request and http.Response objects.
func NewHTTP(req *http.Request, res *http.Response, ret *HTTPRet) *HTTPPayload {
	if req == nil {
		req = &http.Request{}
	}

	if res == nil {
		res = &http.Response{}
	}

	if ret == nil {
		ret = &HTTPRet{}
	}

	sdreq := &HTTPPayload{
		RequestMethod: req.Method,
		Status:        res.StatusCode,
		UserAgent:     req.UserAgent(),
		RemoteIP:      req.RemoteAddr,
		Referer:       req.Referer(),
		Protocol:      req.Proto,
		RetCode:       ret.RetCode,
		RetMsg:        ret.RetMsg,
		RetRequestID:  ret.RetRequestID,
	}

	if req.URL != nil {
		sdreq.RequestURL = req.URL.String()
	}

	buf := &bytes.Buffer{}
	if req.Body != nil {
		n, _ := io.Copy(buf, req.Body) // nolint: gas
		sdreq.RequestSize = strconv.FormatInt(n, 10)
	}

	if res.Body != nil {
		buf.Reset()
		n, _ := io.Copy(buf, res.Body) // nolint: gas
		sdreq.ResponseSize = strconv.FormatInt(n, 10)
	}

	return sdreq
}

// NewHTTP returns a new HTTPPayload struct, based on the passed
// in http.Request and http.Response objects.
func NewHTTPWithLatency(req *http.Request, res *http.Response, ret *HTTPRet, latency string, responseSize string) *HTTPPayload {
	sdreq := NewHTTP(req, res, ret)
	sdreq.Latency = latency
	sdreq.ResponseSize = responseSize

	return sdreq
}

// MarshalLogObject implements zapcore.ObjectMarshaller interface.
func (req HTTPPayload) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("requestMethod", req.RequestMethod)
	enc.AddString("requestUrl", req.RequestURL)
	enc.AddString("requestSize", req.RequestSize)
	enc.AddInt("status", req.Status)
	enc.AddString("responseSize", req.ResponseSize)
	enc.AddString("userAgent", req.UserAgent)
	enc.AddString("remoteIp", req.RemoteIP)
	enc.AddString("serverIp", req.ServerIP)
	enc.AddString("referer", req.Referer)
	enc.AddString("latency", req.Latency)
	enc.AddBool("cacheLookup", req.CacheLookup)
	enc.AddBool("cacheHit", req.CacheHit)
	enc.AddBool("cacheValidatedWithOriginServer", req.CacheValidatedWithOriginServer)
	enc.AddString("cacheFillBytes", req.CacheFillBytes)
	enc.AddString("protocol", req.Protocol)
	enc.AddInt("retCode", req.RetCode)
	enc.AddString("retMsg", req.RetMsg)
	enc.AddString("retRequestID", req.RetRequestID)

	return nil
}

func (req HTTPPayload) String() (str string) {
	str = str + "requestMethod" + ": " + req.RequestMethod + ", "
	str = str + "requestUrl" + ": " + req.RequestURL + ", "
	str = str + "requestSize" + ": " + req.RequestSize + ", "
	str = str + "status" + ": " + strconv.Itoa(req.Status) + ", "
	str = str + "responseSize" + ": " + req.ResponseSize + ", "
	str = str + "userAgent" + ": " + req.UserAgent + ", "
	str = str + "remoteIp" + ": " + req.RemoteIP + ", "
	str = str + "serverIp" + ": " + req.ServerIP + ", "
	str = str + "referer" + ": " + req.Referer + ", "
	str = str + "latency" + ": " + req.Latency + ", "
	str = str + "cacheLookup" + ": " + strconv.FormatBool(req.CacheLookup) + ", "
	str = str + "cacheHit" + ": " + strconv.FormatBool(req.CacheHit) + ", "
	str = str + "cacheValidatedWithOriginServer" + ": " + strconv.FormatBool(req.CacheValidatedWithOriginServer) + ", "
	str = str + "cacheFillBytes" + ": " + req.CacheFillBytes + ", "
	str = str + "protocol" + ": " + req.Protocol + ", "
	str = str + "retCode" + ": " + strconv.Itoa(req.RetCode) + ", "
	str = str + "retMsg" + ": " + req.RetMsg + ", "
	str = str + "retRequestID" + ": " + req.RetRequestID

	return
}
