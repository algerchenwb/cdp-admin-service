package httpwrap

import (
	"go.opentelemetry.io/otel/attribute"
)

const (
	URLKey                               = attribute.Key("httpwrap.url")
	MethodKey                            = attribute.Key("httpwrap.Method")
	RetryTimesKey                        = attribute.Key("httpwrap.RetryTimes")
	RequestOptionsDataKey                = attribute.Key("httpwrap.RequestOptions.Data")
	RequestOptionsParamsKey              = attribute.Key("httpwrap.RequestOptions.Params")
	RequestOptionsQueryStructKey         = attribute.Key("httpwrap.RequestOptions.QueryStruct")
	RequestOptionsJSONKey                = attribute.Key("httpwrap.RequestOptions.JSON")
	RequestOptionsXMLKey                 = attribute.Key("httpwrap.RequestOptions.XML")
	RequestOptionsHeadersKey             = attribute.Key("httpwrap.RequestOptions.Headers")
	RequestOptionsTLSHandshakeTimeoutKey = attribute.Key("httpwrap.RequestOptions.TLSHandshakeTimeout")
	RequestOptionsDialTimeoutKey         = attribute.Key("httpwrap.RequestOptions.DialTimeout")
	RequestOptionsDialKeepAliveKey       = attribute.Key("httpwrap.RequestOptions.DialKeepAlive")
	RequestOptionsRequestTimeoutKey      = attribute.Key("httpwrap.RequestOptions.RequestTimeout")
	RequestOptionsDisableReuseConnKey    = attribute.Key("httpwrap.RequestOptions.DisableReuseConn")
	RequestOptionsDisableModCallKey      = attribute.Key("httpwrap.RequestOptions.DisableModCall")
	RequestOptionsRTKey                  = attribute.Key("httpwrap.RequestOptions.RT")
	RequestOptionsVerboseKey             = attribute.Key("httpwrap.RequestOptions.Verbose")
	ErrCodeKey                           = attribute.Key("httpwrap.ErrCode")
	RetHeadKey                           = attribute.Key("httpwrap.RetHead")
	RetBodyKey                           = attribute.Key("httpwrap.RetBody")
)
