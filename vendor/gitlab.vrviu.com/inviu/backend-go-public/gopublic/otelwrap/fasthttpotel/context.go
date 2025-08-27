package fasthttpotel

import (
	"context"

	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"
)

const (
	traceCtxKey = "trace-ctx-key"
)

func SetTraceCtx(ctx *fasthttp.RequestCtx, traceCtx context.Context) {
	ctx.SetUserValue(traceCtxKey, traceCtx)
}

func GetTraceCtx(ctx *fasthttp.RequestCtx) context.Context {
	traceCtx, ok := ctx.UserValue(traceCtxKey).(context.Context)
	if ok {
		return traceCtx
	}
	traceCtx = otel.GetTextMapPropagator().Extract(ctx, NewHeaderCarrier(&ctx.Request.Header))
	SetTraceCtx(ctx, traceCtx)
	return traceCtx
}

// ContextWithSpan returns a copy of parent with span set as the current Span.
func ContextWithSpan(ctx *fasthttp.RequestCtx, span trace.Span) *fasthttp.RequestCtx {
	traceCtx := trace.ContextWithSpan(GetTraceCtx(ctx), span)
	SetTraceCtx(ctx, traceCtx)
	return ctx
}

// ContextWithSpanContext returns a copy of parent with sc as the current
// Span. The Span implementation that wraps sc is non-recording and performs
// no operations other than to return sc as the SpanContext from the
// SpanContext method.
func ContextWithSpanContext(ctx *fasthttp.RequestCtx, sc trace.SpanContext) *fasthttp.RequestCtx {
	traceCtx := trace.ContextWithSpanContext(GetTraceCtx(ctx), sc)
	SetTraceCtx(ctx, traceCtx)
	return ctx
}

// ContextWithRemoteSpanContext returns a copy of parent with rsc set explicly
// as a remote SpanContext and as the current Span. The Span implementation
// that wraps rsc is non-recording and performs no operations other than to
// return rsc as the SpanContext from the SpanContext method.
func ContextWithRemoteSpanContext(ctx *fasthttp.RequestCtx, rsc trace.SpanContext) *fasthttp.RequestCtx {
	traceCtx := trace.ContextWithRemoteSpanContext(GetTraceCtx(ctx), rsc)
	SetTraceCtx(ctx, traceCtx)
	return ctx
}

// SpanFromContext returns the current Span from ctx.
//
// If no Span is currently set in ctx an implementation of a Span that
// performs no operations is returned.
func SpanFromContext(ctx *fasthttp.RequestCtx) trace.Span {
	return trace.SpanFromContext(GetTraceCtx(ctx))
}

// SpanContextFromContext returns the current Span's SpanContext.
func SpanContextFromContext(ctx *fasthttp.RequestCtx) trace.SpanContext {
	return trace.SpanContextFromContext(GetTraceCtx(ctx))
}

// ContextWithBaggage returns a copy of parent with baggage.
func ContextWithBaggage(ctx *fasthttp.RequestCtx, b baggage.Baggage) *fasthttp.RequestCtx {
	// Delegate so any hooks for the OpenTracing bridge are handled.
	traceCtx := baggage.ContextWithBaggage(GetTraceCtx(ctx), b)
	SetTraceCtx(ctx, traceCtx)
	return ctx
}

// ContextWithoutBaggage returns a copy of parent with no baggage.
func ContextWithoutBaggage(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
	// Delegate so any hooks for the OpenTracing bridge are handled.
	traceCtx := baggage.ContextWithoutBaggage(GetTraceCtx(ctx))
	SetTraceCtx(ctx, traceCtx)
	return ctx
}

// BaggageFromContext returns the baggage contained in ctx.
func BaggageFromContext(ctx *fasthttp.RequestCtx) baggage.Baggage {
	// Delegate so any hooks for the OpenTracing bridge are handled.
	return baggage.FromContext(GetTraceCtx(ctx))
}
