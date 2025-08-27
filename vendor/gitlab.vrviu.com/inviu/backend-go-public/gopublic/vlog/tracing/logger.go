package tracing

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"

	"gitlab.vrviu.com/inviu/backend-go-public/go.uber.org/zap"
	"gitlab.vrviu.com/inviu/backend-go-public/go.uber.org/zap/zapcore"
	"go.opentelemetry.io/otel/trace"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var (
	KeyXRequestID = "X-Request-ID"
)

func NewTraceCtx(traceID string) context.Context {
	return WithTraceID(context.Background(), traceID)
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, KeyXRequestID, traceID)
}

func GetTraceID(ctx context.Context) string {
	// fasthttp 只支持 string 类型 key
	if val, ok := ctx.Value(KeyXRequestID).(string); ok {
		return val
	}
	return ""
}

type VLogger struct {
	log    *zap.Logger
	config Config
}

func newCore(level string, encoding string, w zapcore.WriteSyncer) (core zapcore.Core) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeDuration = zapcore.NanosDurationEncoder
	encoderConfig.TimeKey = "T"
	encoderConfig.LevelKey = "L"
	encoderConfig.MessageKey = "M"
	encoderConfig.CallerKey = "LFILE"

	var l zapcore.Level
	switch strings.ToLower(level) {
	case "debug":
		l = zap.DebugLevel
	case "info":
		l = zap.InfoLevel
	case "warn":
		l = zap.WarnLevel
	case "error":
		l = zap.ErrorLevel
	case "fatal":
		l = zap.FatalLevel
	case "panic":
		l = zap.PanicLevel
	default:
		l = zap.InfoLevel
	}

	var e zapcore.Encoder
	switch encoding {
	case "json":
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		e = zapcore.NewJSONEncoder(encoderConfig)
	case "console":
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		e = zapcore.NewConsoleEncoder(encoderConfig)
	default:
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		e = zapcore.NewConsoleEncoder(encoderConfig)
	}

	core = zapcore.NewCore(e, w, l)
	return
}

func NewLogger(config Config, opts ...Option) *VLogger {
	coreFlag := false
	var core zapcore.Core

	for _, opt := range opts {
		opt(&config)
	}

	if config.EnableFile {
		hook := lumberjack.Logger{
			Filename:   config.Filename,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			LocalTime:  config.LocalTime,
			Compress:   config.Compress,
		}
		w := zapcore.AddSync(&hook)
		filecore := newCore(config.FileLevel, config.FileEncodeing, w)

		if coreFlag {
			core = zapcore.NewTee(core, filecore)
		} else {
			core = filecore
			coreFlag = true
		}
	}

	if config.EnableConsole {
		w := zapcore.Lock(os.Stderr)
		consolecore := newCore(config.ConsoleLevel, config.ConsoleEncodeing, w)

		if coreFlag {
			core = zapcore.NewTee(core, consolecore)
		} else {
			core = consolecore
			coreFlag = true
		}
	}

	if !config.EnableFile && !config.EnableConsole {
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.AddSync(ioutil.Discard),
			zap.PanicLevel)
	}

	fields := []zap.Field{zap.String("LAPP", config.AppName)}
	if config.EnablePID {
		fields = append(fields, zap.Int("LPID", os.Getpid()))
	}

	if config.EnableSourceIP {
		fields = append(fields, zap.String("LIP", GetIP(config.SourceEth)))
	}

	core = core.With(fields)

	zapOption := []zap.Option{}
	if config.EnableCaller {
		zapOption = append(zapOption, zap.AddCaller(), zap.AddCallerSkip(config.GlobalCallerSkip+1))
	}

	l := zap.New(core, zapOption...)

	return &VLogger{l, config}
}

func GetIP(eth string) string {
	ifi, err := net.InterfaceByName(eth)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	addrs, err := ifi.Addrs()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return ""
}

func (l *VLogger) getFields(ctx context.Context) (fields []zap.Field) {
	if !l.config.DisableTraceID {
		span := trace.SpanContextFromContext(ctx)
		if span.HasTraceID() {
			fields = append(fields, zap.String("TRACE_ID", span.TraceID().String()))
			fields = append(fields, zap.String("SPAN_ID", span.SpanID().String()))
		} else {
			fields = append(fields, zap.String("TRACE_ID", GetTraceID(ctx)))
		}
	}
	return
}

func (l VLogger) AppName() string {
	return l.config.AppName
}

// Print logs a message at level Debug on the VLogger.
func (l VLogger) Print(ctx context.Context, args ...interface{}) {
	l.log.Debug(fmt.Sprint(args...), l.getFields(ctx)...)
}

// Printf logs a message at level Debug on the VLogger.
func (l VLogger) Printf(ctx context.Context, format string, args ...interface{}) {
	l.log.Debug(fmt.Sprintf(format, args...), l.getFields(ctx)...)
}

// Debug logs a message at level Debug on the VLogger.
func (l VLogger) Debug(ctx context.Context, args ...interface{}) {
	l.log.Debug(fmt.Sprint(args...), l.getFields(ctx)...)
}

// Debugf logs a message at level Debug on the VLogger.
func (l VLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.log.Debug(fmt.Sprintf(format, args...), l.getFields(ctx)...)
}

// Info logs a message at level Info on the VLogger.
func (l VLogger) Info(ctx context.Context, args ...interface{}) {
	l.log.Info(fmt.Sprint(args...), l.getFields(ctx)...)
}

// Infof logs a message at level Info on the VLogger.
func (l VLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	l.log.Info(fmt.Sprintf(format, args...), l.getFields(ctx)...)
}

// Warn logs a message at level Warn on the VLogger.
func (l VLogger) Warn(ctx context.Context, args ...interface{}) {
	l.log.Warn(fmt.Sprint(args...), l.getFields(ctx)...)
}

// Warnf logs a message at level Warn on the VLogger.
func (l VLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	l.log.Warn(fmt.Sprintf(format, args...), l.getFields(ctx)...)
}

// Error logs a message at level Error on the VLogger.
func (l VLogger) Error(ctx context.Context, args ...interface{}) {
	l.log.Error(fmt.Sprint(args...), l.getFields(ctx)...)
}

// Errorf logs a message at level Error on the VLogger.
func (l VLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.log.Error(fmt.Sprintf(format, args...), l.getFields(ctx)...)
}

// Fatal logs a message at level Fatal on the VLogger.
func (l VLogger) Fatal(ctx context.Context, args ...interface{}) {
	l.log.Fatal(fmt.Sprint(args...), l.getFields(ctx)...)
}

// Fatalf logs a message at level Fatal on the VLogger.
func (l VLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	l.log.Fatal(fmt.Sprintf(format, args...), l.getFields(ctx)...)
}

// Panic logs a message at level Painc on the VLogger.
func (l VLogger) Panic(ctx context.Context, args ...interface{}) {
	l.log.Panic(fmt.Sprint(args...), l.getFields(ctx)...)
}

// Panicf logs a message at level Painc on the VLogger.
func (l VLogger) Panicf(ctx context.Context, format string, args ...interface{}) {
	l.log.Panic(fmt.Sprintf(format, args...), l.getFields(ctx)...)
}

func (l *VLogger) zapFields(fields map[string]interface{}) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return zapFields
}

// With return a logger with an extra field.
func (l *VLogger) With(key string, value interface{}) *VLogger {
	return l.WithField(zap.Any(key, value))
}

// Withs return a logger with extra fields.
func (l *VLogger) Withs(fields map[string]interface{}) *VLogger {
	return l.WithField(l.zapFields(fields)...)
}

// WithField return a logger with extra zap fields.
func (l *VLogger) WithField(fields ...zap.Field) *VLogger {
	return &VLogger{l.log.With(fields...), l.config}
}

// AddCallerSkip return a logger with new caller skip.
func (l *VLogger) AddCallerSkip(skip int) *VLogger {
	return &VLogger{l.log.WithOptions(zap.AddCallerSkip(skip)), l.config}
}
