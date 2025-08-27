package vlog

import (
	"context"

	"gitlab.vrviu.com/inviu/backend-go-public/go.uber.org/zap"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic/vlog/tracing"
)

var (
	oldCtx = tracing.NewTraceCtx("no_trace")
)

func GetLogger() *tracing.VLogger {
	return tracing.GetLogger()
}

func AppName() string {
	return tracing.AppName()
}

func GetPrintLogger(err error) func(context.Context, ...interface{}) {
	return tracing.GetPrintLogger(err)
}

func GetPrintfLogger(err error) func(context.Context, string, ...interface{}) {
	return tracing.GetPrintfLogger(err)
}

func Init(filename string) {
	tracing.Init(filename)
}

func SetConfig(config Config) func() {
	return tracing.SetConfig(tracing.Config(config))
}

func ReplaceLogger(logger *VLogger) func() {
	return tracing.ReplaceLogger(logger.log)
}

// Print logs a message at level Debug on the VLogger.
func Print(args ...interface{}) {
	GetLogger().Print(oldCtx, args...)
}

func PrintContext(ctx context.Context, args ...interface{}) {
	GetLogger().Print(ctx, args...)
}

// Printf logs a message at level Debug on the VLogger.
func Printf(format string, args ...interface{}) {
	GetLogger().Printf(oldCtx, format, args...)
}

func PrintfContext(ctx context.Context, format string, args ...interface{}) {
	GetLogger().Printf(ctx, format, args...)
}

// Debug logs a message at level Debug on the VLogger.
func Debug(args ...interface{}) {
	GetLogger().Debug(oldCtx, args...)
}

func DebugContext(ctx context.Context, args ...interface{}) {
	GetLogger().Debug(ctx, args...)
}

// Debugf logs a message at level Debug on the VLogger.
func Debugf(format string, args ...interface{}) {
	GetLogger().Debugf(oldCtx, format, args...)
}

func DebugfContext(ctx context.Context, format string, args ...interface{}) {
	GetLogger().Debugf(ctx, format, args...)
}

// Info logs a message at level Info on the VLogger.
func Info(args ...interface{}) {
	GetLogger().Info(oldCtx, args...)
}

func InfoContext(ctx context.Context, args ...interface{}) {
	GetLogger().Info(ctx, args...)
}

// Infof logs a message at level Info on the VLogger.
func Infof(format string, args ...interface{}) {
	GetLogger().Infof(oldCtx, format, args...)
}

func InfofContext(ctx context.Context, format string, args ...interface{}) {
	GetLogger().Infof(ctx, format, args...)
}

// Warn logs a message at level Warn on the VLogger.
func Warn(args ...interface{}) {
	GetLogger().Warn(oldCtx, args...)
}

func WarnContext(ctx context.Context, args ...interface{}) {
	GetLogger().Warn(ctx, args...)
}

// Warnf logs a message at level Warn on the VLogger.
func Warnf(format string, args ...interface{}) {
	GetLogger().Warnf(oldCtx, format, args...)
}

func WarnfContext(ctx context.Context, format string, args ...interface{}) {
	GetLogger().Warnf(ctx, format, args...)
}

// Error logs a message at level Error on the VLogger.
func Error(args ...interface{}) {
	GetLogger().Error(oldCtx, args...)
}

func ErrorContext(ctx context.Context, args ...interface{}) {
	GetLogger().Error(ctx, args...)
}

// Errorf logs a message at level Error on the VLogger.
func Errorf(format string, args ...interface{}) {
	GetLogger().Errorf(oldCtx, format, args...)
}

func ErrorfContext(ctx context.Context, format string, args ...interface{}) {
	GetLogger().Errorf(ctx, format, args...)
}

// Fatal logs a message at level Fatal on the VLogger.
func Fatal(args ...interface{}) {
	GetLogger().Fatal(oldCtx, args...)
}

func FatalContext(ctx context.Context, args ...interface{}) {
	GetLogger().Fatal(ctx, args...)
}

// Fatalf logs a message at level Fatal on the VLogger.
func Fatalf(format string, args ...interface{}) {
	GetLogger().Fatalf(oldCtx, format, args...)
}

func FatalfContext(ctx context.Context, format string, args ...interface{}) {
	GetLogger().Fatalf(ctx, format, args...)
}

// Panic logs a message at level Panic on the VLogger.
func Panic(args ...interface{}) {
	GetLogger().Panic(oldCtx, args...)
}

func PanicContext(ctx context.Context, args ...interface{}) {
	GetLogger().Panic(ctx, args...)
}

// Panicf logs a message at level Panic on the VLogger.
func Panicf(format string, args ...interface{}) {
	GetLogger().Panicf(oldCtx, format, args...)
}

func PanicfContext(ctx context.Context, format string, args ...interface{}) {
	GetLogger().Panicf(ctx, format, args...)
}

// With return a logger with an extra field.
func With(key string, value interface{}) *VLogger {
	return &VLogger{log: GetLogger().With(key, value)}
}

// Withs return a logger with extra fields.
func Withs(fields map[string]interface{}) *VLogger {
	return &VLogger{log: GetLogger().Withs(fields)}
}

// WithField return a logger with extra zap fields.
func WithField(fields ...zap.Field) *VLogger {
	return &VLogger{log: GetLogger().WithField(fields...)}
}

// AddCallerSkip return a logger with new caller skip.
func AddCallerSkip(skip int) *VLogger {
	return &VLogger{log: GetLogger().AddCallerSkip(skip)}
}

// RedirectStdLog redirects output from the standard library's package-global
// logger to the supplied logger at InfoLevel. Since zap already handles caller
// annotations, timestamps, etc., it automatically disables the standard
// library's annotations and prefixing.
//
// It returns a function to restore the original prefix and flags and reset the
// standard library's output to os.Stderr.
func RedirectStdLog(l *VLogger) (func(), error) {
	return tracing.RedirectStdLog(l.log)
}

// RedirectStdLogAt redirects output from the standard library's package-global
// logger to the supplied logger at the specified level. Since zap already
// handles caller annotations, timestamps, etc., it automatically disables the
// standard library's annotations and prefixing.
//
// It returns a function to restore the original prefix and flags and reset the
// standard library's output to os.Stderr.
func RedirectStdLogAt(l *VLogger, level string) (func(), error) {
	return tracing.RedirectStdLogAt(l.log, level)
}

func ReplaceStdLog() (func(), error) {
	return tracing.ReplaceStdLog()
}
