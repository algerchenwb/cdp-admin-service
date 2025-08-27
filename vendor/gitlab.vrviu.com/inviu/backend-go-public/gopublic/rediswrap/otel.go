package rediswrap

import (
	"context"

	"github.com/go-redis/redis/extra/redisotel/v8"
	"github.com/go-redis/redis/v8"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic/otelwrap"
	"go.opentelemetry.io/otel/trace"
)

type TracingHook struct {
	*redisotel.TracingHook
}

func NewTracingHook(opts ...redisotel.Option) *TracingHook {
	return &TracingHook{TracingHook: redisotel.NewTracingHook(opts...)}
}

func (th *TracingHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	if otelwrap.IsSkip(ctx) {
		return ctx, nil
	}

	return th.TracingHook.BeforeProcess(ctx, cmd)
}

func (th *TracingHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if otelwrap.IsSkip(ctx) {
		return nil
	}

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(ResultKey.String(cmd.String()))
	return th.TracingHook.AfterProcess(ctx, cmd)
}

func (th *TracingHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	if otelwrap.IsSkip(ctx) {
		return ctx, nil
	}

	return th.TracingHook.BeforeProcessPipeline(ctx, cmds)
}

func (th *TracingHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	if otelwrap.IsSkip(ctx) {
		return nil
	}

	span := trace.SpanFromContext(ctx)
	show := []string{}
	for _, cmd := range cmds {
		show = append(show, cmd.String())
	}
	span.SetAttributes(ResultKey.StringSlice(show))
	return th.TracingHook.AfterProcessPipeline(ctx, cmds)
}
