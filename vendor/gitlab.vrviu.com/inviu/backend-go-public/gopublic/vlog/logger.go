package vlog

import (
	"context"
	"fmt"

	"gitlab.vrviu.com/inviu/backend-go-public/go.uber.org/zap"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic/vlog/tracing"
)

type VLogger struct {
	log *tracing.VLogger
}

func NewLogger(config Config, opts ...Option) *VLogger {
	nopts := []tracing.Option{}
	for _, opt := range opts {
		nopts = append(nopts, tracing.Option(opt))
	}
	logger := tracing.NewLogger(tracing.Config(config), nopts...)
	return &VLogger{log: logger}
}

func GetIP(eth string) string {
	return tracing.GetIP(eth)
}

func NewTraceCtx(traceID string) context.Context {
	return tracing.NewTraceCtx(traceID)
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return tracing.WithTraceID(ctx, traceID)
}

func GetTraceID(ctx context.Context) string {
	return tracing.GetTraceID(ctx)
}

// Print logs a message at level Debug on the VLogger.
func (l VLogger) Print(args ...interface{}) {
	l.log.Print(oldCtx, fmt.Sprint(args...))
}

func (l VLogger) PrintContext(ctx context.Context, args ...interface{}) {
	l.log.Print(ctx, args...)
}

// Printf logs a message at level Debug on the VLogger.
func (l VLogger) Printf(format string, args ...interface{}) {
	l.log.Printf(oldCtx, format, args...)
}

func (l VLogger) PrintfContext(ctx context.Context, format string, args ...interface{}) {
	l.log.Printf(ctx, format, args...)
}

// Debug logs a message at level Debug on the VLogger.
func (l VLogger) Debug(args ...interface{}) {
	l.log.Debug(oldCtx, args...)
}

func (l VLogger) DebugContext(ctx context.Context, args ...interface{}) {
	l.log.Debug(ctx, args...)
}

// Debugf logs a message at level Debug on the VLogger.
func (l VLogger) Debugf(format string, args ...interface{}) {
	l.log.Debugf(oldCtx, format, args...)
}

func (l VLogger) DebugfContext(ctx context.Context, format string, args ...interface{}) {
	l.log.Debugf(ctx, format, args...)
}

// Info logs a message at level Info on the VLogger.
func (l VLogger) Info(args ...interface{}) {
	l.log.Info(oldCtx, args...)
}

func (l VLogger) InfoContext(ctx context.Context, args ...interface{}) {
	l.log.Info(ctx, args...)
}

// Infof logs a message at level Info on the VLogger.
func (l VLogger) Infof(format string, args ...interface{}) {
	l.log.Infof(oldCtx, format, args...)
}

func (l VLogger) InfofContext(ctx context.Context, format string, args ...interface{}) {
	l.log.Infof(ctx, format, args...)
}

// Warn logs a message at level Warn on the VLogger.
func (l VLogger) Warn(args ...interface{}) {
	l.log.Warn(oldCtx, args...)
}

func (l VLogger) WarnContext(ctx context.Context, args ...interface{}) {
	l.log.Warn(ctx, args...)
}

// Warnf logs a message at level Warn on the VLogger.
func (l VLogger) Warnf(format string, args ...interface{}) {
	l.log.Warnf(oldCtx, format, args...)
}

func (l VLogger) WarnfContext(ctx context.Context, format string, args ...interface{}) {
	l.log.Warnf(ctx, format, args...)
}

// Error logs a message at level Error on the VLogger.
func (l VLogger) Error(args ...interface{}) {
	l.log.Error(oldCtx, args...)
}

func (l VLogger) ErrorContext(ctx context.Context, args ...interface{}) {
	l.log.Error(ctx, args...)
}

// Errorf logs a message at level Error on the VLogger.
func (l VLogger) Errorf(format string, args ...interface{}) {
	l.log.Errorf(oldCtx, format, args...)
}

func (l VLogger) ErrorfContext(ctx context.Context, format string, args ...interface{}) {
	l.log.Errorf(ctx, format, args...)
}

// Fatal logs a message at level Fatal on the VLogger.
func (l VLogger) Fatal(args ...interface{}) {
	l.log.Fatal(oldCtx, args...)
}

func (l VLogger) FatalContext(ctx context.Context, args ...interface{}) {
	l.log.Fatal(ctx, args...)
}

// Fatalf logs a message at level Fatal on the VLogger.
func (l VLogger) Fatalf(format string, args ...interface{}) {
	l.log.Fatalf(oldCtx, format, args...)
}

func (l VLogger) FatalfContext(ctx context.Context, format string, args ...interface{}) {
	l.log.Fatalf(ctx, format, args...)
}

// Panic logs a message at level Painc on the VLogger.
func (l VLogger) Panic(args ...interface{}) {
	l.log.Panic(oldCtx, args...)
}

func (l VLogger) PanicContext(ctx context.Context, args ...interface{}) {
	l.log.Panic(ctx, args...)
}

// Panicf logs a message at level Painc on the VLogger.
func (l VLogger) Panicf(format string, args ...interface{}) {
	l.log.Panicf(oldCtx, format, args...)
}

func (l VLogger) PanicfContext(ctx context.Context, format string, args ...interface{}) {
	l.log.Panicf(ctx, format, args...)
}

// With return a logger with an extra field.
func (l *VLogger) With(key string, value interface{}) *VLogger {
	return &VLogger{log: l.log.With(key, value)}
}

// Withs return a logger with extra fields.
func (l *VLogger) Withs(fields map[string]interface{}) *VLogger {
	return &VLogger{log: l.log.Withs(fields)}
}

// WithField return a logger with extra zap fields.
func (l *VLogger) WithField(fields ...zap.Field) *VLogger {
	return &VLogger{log: l.log.WithField(fields...)}
}

// AddCallerSkip return a logger with new caller skip.
func (l *VLogger) AddCallerSkip(skip int) *VLogger {
	return &VLogger{log: l.log.AddCallerSkip(skip)}
}
