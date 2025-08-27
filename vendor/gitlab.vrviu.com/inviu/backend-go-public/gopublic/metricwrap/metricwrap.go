package metricwrap

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
)

type MetricLabelsBuilder interface {
	SetKV(key, val string)
	SetKVs(kvs map[string]string)
	Build() prometheus.Labels
}

type MetricReport interface {
	SetMetricLabels(b MetricLabelsBuilder)
}

type metricContextKeyType int

const (
	metricLabelsBuilderKey metricContextKeyType = 0
)

// ContextWithMetricLabelsBuilder returns a copy of parent with span set as the metric labels.
func ContextWithMetricLabelsBuilder(parent context.Context, builder MetricLabelsBuilder) context.Context {
	return context.WithValue(parent, metricLabelsBuilderKey, builder)
}

// MetricLabelsBuilderFromContext returns a copy of parent with span set as the metric labels.
func MetricLabelsBuilderFromContext(ctx context.Context) MetricLabelsBuilder {
	if ctx == nil {
		return noopMetricLabelsBuilder
	}
	if metricLabelsBuilder, ok := ctx.Value(metricLabelsBuilderKey).(MetricLabelsBuilder); ok {
		return metricLabelsBuilder
	}
	return noopMetricLabelsBuilder
}
