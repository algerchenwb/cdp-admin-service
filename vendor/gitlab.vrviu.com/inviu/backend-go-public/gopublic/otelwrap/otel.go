package otelwrap

import (
	"context"
	"net/http"
	"strings"
	"sync"

	"go.opentelemetry.io/contrib/exporters/autoexport"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	_globalMu sync.RWMutex
	_config   = otelConfig{
		EnvGroup: make(map[string]string),
	}
)

func SetConfig(opts ...Option) {
	_globalMu.Lock()
	for _, opt := range opts {
		opt(&_config)
	}
	_config.InitEnv()
	_globalMu.Unlock()
}

func GetServiceName() string {
	return _config.ServiceName
}

func GetFilters() []Filter {
	return _config.Filters
}

func GetSpanNameFormatter() SpanNameFormatter {
	if _config.SpanNameFormatter != nil {
		return _config.SpanNameFormatter
	}
	return func(r *http.Request) string {
		return strings.Clone(r.URL.Path)
	}
}

// Init TracerProvider
//
//	https://opentelemetry.io/docs/concepts/sdk-configuration/otlp-exporter-configuration/
//
// # 1.1 disable exporter (default=none)
//   - OTEL_TRACES_EXPORTER = none
//
// # 1.2 enable http otel
//   - OTEL_TRACES_EXPORTER = otlp
//   - OTEL_EXPORTER_OTLP_PROTOCOL = http/protobuf
//   - OTEL_EXPORTER_OTLP_ENDPOINT = http://my-api-endpoint/
//
// # 1.3 enable grpc otel
//   - OTEL_TRACES_EXPORTER = otlp
//   - OTEL_EXPORTER_OTLP_PROTOCOL = grpc
//   - OTEL_EXPORTER_OTLP_ENDPOINT = https://my-api-endpoint:443
//
// # 2.1 sampler (default_sampler=parentbased_always_on)
//   - OTEL_TRACES_SAMPLER = [always_on|always_off|parentbased_always_on|parentbased_always_off]
//
// # 2.2 ratio sampler (default_sampler_arg=1.0)
//   - OTEL_TRACES_SAMPLER = [traceidratio|parentbased_traceidratio]
//   - OTEL_TRACES_SAMPLER_ARG = [0-1.0]
//
// # 3 resource service name
//   - OTEL_SERVICE_NAME = service-name
//   - OTEL_RESOURCE_ATTRIBUTES = key1=val1,key2=val2
func Init(ctx context.Context, name string, opts ...sdktrace.TracerProviderOption) (*sdktrace.TracerProvider, error) {
	exporter, err := autoexport.NewSpanExporter(ctx)
	if err != nil {
		return nil, err
	}

	res, err := resource.Merge(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceName(name)), resource.Environment())
	if err != nil {
		return nil, err
	}

	opts = append([]sdktrace.TracerProviderOption{
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	}, opts...)

	tp := sdktrace.NewTracerProvider(opts...)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}

const (
	Version    = "v1.8.15"
	tracerName = "gitlab.vrviu.com/inviu/backend-go-public/gopublic/otelwrap"
)

var (
	tracer     trace.Tracer
	tracerOnce sync.Once
)

func NewTracer(options ...trace.TracerOption) trace.Tracer {
	options = append(options, trace.WithInstrumentationVersion(Version))
	return otel.GetTracerProvider().Tracer(tracerName, options...)
}

func GetTracer() trace.Tracer {
	tracerOnce.Do(func() {
		tracer = NewTracer()
	})
	return tracer
}

func Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return GetTracer().Start(ctx, spanName, opts...)
}
