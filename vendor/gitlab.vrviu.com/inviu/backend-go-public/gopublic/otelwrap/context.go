package otelwrap

import (
	"context"

	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
	vlog "gitlab.vrviu.com/inviu/backend-go-public/gopublic/vlog/tracing"
)

type ContextKeyType int

const (
	CtxKeyOtelFlag ContextKeyType = iota
)

type OtelFlag int

const (
	OtelFlagSkip OtelFlag = 1 << iota
)

var (
	ctxOtelBackend = context.Background()
	ctxOtelTODO    = context.TODO()
	ctxOtelSkip    = NewOtelCtx(OtelFlagSkip)
)

func Background() context.Context {
	return ctxOtelBackend
}

func TODO() context.Context {
	return ctxOtelTODO
}

func Skip() context.Context {
	return ctxOtelSkip
}

func NewOtelCtx(flag OtelFlag) context.Context {
	return WithOtelFlag(Background(), flag)
}

func GetOtelFlag(ctx context.Context) OtelFlag {
	flag, ok := ctx.Value(CtxKeyOtelFlag).(OtelFlag)
	if !ok {
		return 0
	}
	return flag
}

func WithOtelFlag(ctx context.Context, flag OtelFlag) context.Context {
	return context.WithValue(ctx, CtxKeyOtelFlag, GetOtelFlag(ctx)|flag)
}

func WithoutOtelFlag(ctx context.Context, flag OtelFlag) context.Context {
	return context.WithValue(ctx, CtxKeyOtelFlag, GetOtelFlag(ctx)&(^flag))
}

func CheckOtelFlag(ctx context.Context, flag OtelFlag) bool {
	return GetOtelFlag(ctx)&flag == flag
}

func IsSkip(ctx context.Context) bool {
	return CheckOtelFlag(ctx, OtelFlagSkip)
}

func NewTraceCtx(name string) context.Context {
	return vlog.WithTraceID(Background(), gopublic.RandonStringWithPrefix(name, 10))
}

func NewSkipTraceCtx(name string) context.Context {
	return vlog.WithTraceID(Skip(), gopublic.RandonStringWithPrefix(name, 10))
}
