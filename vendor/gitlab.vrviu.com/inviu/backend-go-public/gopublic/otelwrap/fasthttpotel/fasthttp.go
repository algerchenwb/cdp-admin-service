package fasthttpotel

import (
	"fmt"
	"net/http"

	"gitlab.vrviu.com/inviu/backend-go-public/gopublic/otelwrap"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic/otelwrap/internal/semconvutil"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	ReadBytesKey  = attribute.Key("http.read_bytes")  // if anything was read from the request body, the total number of bytes read
	WroteBytesKey = attribute.Key("http.wrote_bytes") // if anything was written to the response writer, the total number of bytes written
)

func ConvertRequest(ctx *fasthttp.RequestCtx) (*http.Request, error) {
	req := &http.Request{}
	err := fasthttpadaptor.ConvertRequest(ctx, req, true)
	if err != nil {
		return nil, err
	}
	return req.Clone(GetTraceCtx(ctx)), nil
}

// MiddleWareTraceSpan
// *> set traceCtx to fasthttp.RequestCtx, get traceCtx by fasthttpotel.GetTraceCtx
// *> traceCtx := fasthttpotel.GetTraceCtx(ctx)
func MiddleWareTraceSpan(fn fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(rctx *fasthttp.RequestCtx) {
		stdRequest, err := ConvertRequest(rctx)
		if err != nil {
			fn(rctx)
			return
		}

		for _, f := range otelwrap.GetFilters() {
			if !f(stdRequest) {
				fn(rctx)
				return
			}
		}

		traceCtx := GetTraceCtx(rctx)
		opts := []trace.SpanStartOption{
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(semconvutil.HTTPServerRequest(otelwrap.GetServiceName(), stdRequest)...),
		}
		spanName := otelwrap.GetSpanNameFormatter()(stdRequest)
		if spanName == "" {
			spanName = fmt.Sprintf("HTTP %s route not found", stdRequest.Method)
		} else {
			rAttr := semconv.HTTPRoute(spanName)
			opts = append(opts, trace.WithAttributes(rAttr))
		}
		traceCtx, span := otelwrap.GetTracer().Start(traceCtx, spanName, opts...)
		defer span.End()
		SetTraceCtx(rctx, traceCtx)

		fn(rctx)

		setAfterServeAttributes(span, int64(len(rctx.Request.Body())), int64(len(rctx.Response.Body())), rctx.Response.StatusCode())
	}
}

func setAfterServeAttributes(span trace.Span, read, wrote int64, statusCode int) {
	attributes := []attribute.KeyValue{}

	if read > 0 {
		attributes = append(attributes, ReadBytesKey.Int64(read))
	}
	if wrote > 0 {
		attributes = append(attributes, WroteBytesKey.Int64(wrote))
	}
	if statusCode > 0 {
		attributes = append(attributes, semconv.HTTPStatusCode(statusCode))
	}
	span.SetStatus(semconvutil.HTTPServerStatus(statusCode))

	span.SetAttributes(attributes...)
}
