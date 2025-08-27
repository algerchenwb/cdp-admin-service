package httpwrap

import (
	"context"

	beegoctx "github.com/astaxie/beego/context"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic/otelwrap/fasthttpotel"
	vlog "gitlab.vrviu.com/inviu/backend-go-public/gopublic/vlog/tracing"
	"go.opentelemetry.io/otel/trace"
)

// GenRequestID 申请唯一的request-id
func GenRequestID(ctx *fasthttp.RequestCtx) string {
	span := fasthttpotel.SpanContextFromContext(ctx)
	if span.HasTraceID() {
		return span.TraceID().String() + "-" + span.SpanID().String()
	}

	id := vlog.GetTraceID(ctx)
	if id != "" {
		return id
	}

	id = uuid.NewString()
	ctx.SetUserValue(vlog.KeyXRequestID, id)
	return id
}

// GenRequestID 申请唯一的request-id
func GenRequestID2(ctx context.Context) (context.Context, string) {
	span := trace.SpanContextFromContext(ctx)
	if span.HasTraceID() {
		return ctx, span.TraceID().String() + "-" + span.SpanID().String()
	}

	id := vlog.GetTraceID(ctx)
	if id != "" {
		return ctx, id
	}

	id = uuid.NewString()
	return vlog.WithTraceID(ctx, id), id
}

// GenRequestID 申请唯一的request-id
func GenRequestID3(ctx *beegoctx.Context) string {
	span := trace.SpanContextFromContext(ctx.Request.Context())
	if span.HasTraceID() {
		return span.TraceID().String() + "-" + span.SpanID().String()
	}

	id := vlog.GetTraceID(ctx.Request.Context())
	if id != "" {
		return id
	}

	id = uuid.NewString()
	ctx.Request = ctx.Request.WithContext(vlog.WithTraceID(ctx.Request.Context(), id))
	return id
}
