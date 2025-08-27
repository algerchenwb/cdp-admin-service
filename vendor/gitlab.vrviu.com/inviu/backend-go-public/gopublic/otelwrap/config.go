package otelwrap

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	vlog "gitlab.vrviu.com/inviu/backend-go-public/gopublic/vlog/tracing"
)

type otelConfig struct {
	ServiceName       string
	Filters           []Filter
	SpanNameFormatter SpanNameFormatter
	EnvGroup          map[string]string
}

func (c *otelConfig) InitEnv() {
	for k, v := range c.EnvGroup {
		err := os.Setenv(k, v)
		if err != nil {
			panic(fmt.Errorf("os.Setenv[%s=%s] err[%s]", k, v, err))
		}
		vlog.Infof(vlog.NewTraceCtx("init_otel"), "os.Setenv[%s=%s]", k, v)
	}
}

// Filter is a predicate used to determine whether a given http.request should
// be traced. A Filter must return true if the request should be traced.
type Filter func(*http.Request) bool

// SpanNameFormatter is used to set span name by http.request.
type SpanNameFormatter func(r *http.Request) string

type Option func(*otelConfig)

// WithFilter adds a filter to the list of filters used by the handler.
// If any filter indicates to exclude a request then the request will not be
// traced. All filters must allow a request to be traced for a Span to be created.
// If no filters are provided then all requests are traced.
// Filters will be invoked for each processed request, it is advised to make them
// simple and fast.
func WithFilter(f ...Filter) Option {
	return func(c *otelConfig) {
		c.Filters = append(c.Filters, f...)
	}
}

// WithSpanNameFormatter takes a function that will be called on every
// request and the returned string will become the Span Name.
func WithSpanNameFormatter(f func(r *http.Request) string) Option {
	return func(c *otelConfig) {
		c.SpanNameFormatter = f
	}
}

func WithSpanNameRegexFormatter(regexRules ...string) Option {
	return func(c *otelConfig) {
		res := make([]*regexp.Regexp, len(regexRules))
		for i, rule := range regexRules {
			res[i] = regexp.MustCompile(rule)
		}

		c.SpanNameFormatter = func(r *http.Request) string {
			for _, re := range res {
				if re.Match([]byte(r.URL.Path)) {
					return re.String()
				}
			}
			return strings.Clone(r.URL.Path)
		}
	}
}

const (
	TracesExporterKey             = "OTEL_TRACES_EXPORTER"
	TracesExporterOTLPProtocolKey = "OTEL_EXPORTER_OTLP_PROTOCOL"
	TracesExporterOTLPEndpointKey = "OTEL_EXPORTER_OTLP_ENDPOINT"

	TracesExporterNone             = "none"
	TracesExporterOTLP             = "otlp"
	TracesExporterOTLPProtocolHTTP = "http/protobuf"
	TracesExporterOTLPProtocolGRPC = "grpc"
)

func WithTracesExporter(exporter, protocol, endpoint string) Option {
	return func(c *otelConfig) {
		switch exporter {
		case TracesExporterNone:
			c.EnvGroup[TracesExporterKey] = exporter
		case TracesExporterOTLP:
			c.EnvGroup[TracesExporterKey] = exporter
			c.EnvGroup[TracesExporterOTLPProtocolKey] = protocol
			c.EnvGroup[TracesExporterOTLPEndpointKey] = endpoint
		default:
			c.EnvGroup[TracesExporterKey] = TracesExporterNone
			vlog.Warnf(vlog.NewTraceCtx("init_otel"), "unknown exporter type[%s], set exporter[none]", exporter)
		}
	}
}

const (
	TracesSamplerKey    = "OTEL_TRACES_SAMPLER"
	TracesSamplerArgKey = "OTEL_TRACES_SAMPLER_ARG"

	TracesSamplerAlwaysOn                = "always_on"
	TracesSamplerAlwaysOff               = "always_off"
	TracesSamplerTraceIDRatio            = "traceidratio"
	TracesSamplerParentBasedAlwaysOn     = "parentbased_always_on"
	TracesSamplerParsedBasedAlwaysOff    = "parentbased_always_off"
	TracesSamplerParentBasedTraceIDRatio = "parentbased_traceidratio"
)

func WithTracesSampler(sampler, arg string) Option {
	return func(c *otelConfig) {
		c.EnvGroup[TracesSamplerKey] = sampler
		c.EnvGroup[TracesSamplerArgKey] = arg
	}
}

const (
	ServiceNameKey        = "OTEL_SERVICE_NAME"
	ResourceAttributesKey = "OTEL_RESOURCE_ATTRIBUTES"
)

func WithServiceName(serviceName string) Option {
	return func(c *otelConfig) {
		if len(serviceName) == 0 {
			vlog.Warnf(vlog.NewTraceCtx("init_otel"), "service name is empty")
			return
		}
		c.ServiceName = serviceName
		c.EnvGroup[ServiceNameKey] = serviceName
	}
}

func WithResourceAttributes(attributes string) Option {
	return func(c *otelConfig) {
		if len(attributes) == 0 {
			vlog.Warnf(vlog.NewTraceCtx("init_otel"), "resource attributes is empty")
			return
		}
		c.EnvGroup[ResourceAttributesKey] = attributes
	}
}

func WithEnv(key, val string) Option {
	return func(c *otelConfig) {
		c.EnvGroup[key] = val
	}
}
