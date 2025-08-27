package metricwrap

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
)

type RequestMetricFunc func()

var (
	NoopRequestMetricFunc = func() {}
)

var (
	RequestDefaultNamespace       = "unknown"
	RequestDefaultSubsystem       = "unknown"
	RequestDefaultLabelNames      = []string{MetricLabelName}
	RequestDefaultDurationBuckets = []float64{10, 20, 50, 75, 100, 150, 200, 250, 300, 400, 500, 750, 1000, 1250, 1500, 1750, 2000, 2500, 3000}
	RequestDefaultObjectives      = map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.95: 0.005, 0.99: 0.001}
)

func NewMetricConfig() *MetricConfig {
	return &MetricConfig{
		Disable:         false,
		Namespace:       RequestDefaultNamespace,
		Subsystem:       RequestDefaultSubsystem,
		LabelNames:      RequestDefaultLabelNames,
		EnableHistogram: false,
		DurationBuckets: RequestDefaultDurationBuckets,
		EnableSummary:   true,
		Objectives:      RequestDefaultObjectives,
	}
}

type MetricConfig struct {
	Disable         bool
	Namespace       string
	Subsystem       string
	LabelNames      []string
	EnableHistogram bool
	DurationBuckets []float64
	EnableSummary   bool
	Objectives      map[float64]float64
}

func NewRequestMetricHandler(cfg *MetricConfig) *RequestMetricHandler {
	if cfg == nil {
		cfg = NewMetricConfig()
	}
	if !gopublic.StringInArray(MetricLabelName, cfg.LabelNames) {
		cfg.LabelNames = append(cfg.LabelNames, MetricLabelName)
	}
	return &RequestMetricHandler{
		disable:         cfg.Disable,
		enableHistogram: cfg.EnableHistogram,
		enableSummary:   cfg.EnableSummary,
		labelNames:      cfg.LabelNames,
		requestCounterVec: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: cfg.Namespace,
			Subsystem: cfg.Subsystem,
			Name:      "request_total",
			Help:      "Counter for request.",
		}, cfg.LabelNames),
		requestDurationVec: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: cfg.Namespace,
			Subsystem: cfg.Subsystem,
			Name:      "request_duration_ms",
			Help:      "Histogram of latencies for request duration.",
			Buckets:   cfg.DurationBuckets,
		}, cfg.LabelNames),
		requestDurationSummaryVec: promauto.NewSummaryVec(prometheus.SummaryOpts{
			Namespace:  cfg.Namespace,
			Subsystem:  cfg.Subsystem,
			Name:       "request_duration_ms_summary",
			Help:       "Summary of latencies for request duration.",
			Objectives: cfg.Objectives,
		}, cfg.LabelNames),
	}
}

type RequestMetricHandler struct {
	disable                   bool
	enableHistogram           bool
	enableSummary             bool
	labelNames                []string
	requestCounterVec         *prometheus.CounterVec
	requestDurationVec        *prometheus.HistogramVec
	requestDurationSummaryVec *prometheus.SummaryVec
}

func (h RequestMetricHandler) Metric(ctx context.Context, name string) (context.Context, RequestMetricFunc) {
	if h.disable {
		return ctx, NoopRequestMetricFunc
	}

	startAt := time.Now()

	builder := h.builder()
	builder.SetKV(MetricLabelName, name)
	ctx = ContextWithMetricLabelsBuilder(ctx, builder)

	return ctx, func() {
		h.requestCounterVec.With(builder.Build()).Inc()

		cost := float64(time.Since(startAt).Milliseconds())
		if h.enableHistogram {
			h.requestDurationVec.With(builder.Build()).Observe(cost)
		}
		if h.enableSummary {
			h.requestDurationSummaryVec.With(builder.Build()).Observe(cost)
		}
	}
}

func (h RequestMetricHandler) builder() MetricLabelsBuilder {
	return NewMetricLabelsBuilder(h.labelNames...)
}
