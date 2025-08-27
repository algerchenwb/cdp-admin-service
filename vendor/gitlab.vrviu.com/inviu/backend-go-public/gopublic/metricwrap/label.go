package metricwrap

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	MetricLabelName     = "name"
	MetricLabelBizType  = "biz"
	MetricLabelAreaType = "area"
	MetricLabelIDC      = "idc"
	MetricLabelUGID     = "ugid"
	MetricLabelCode     = "code"
)

var (
	noopMetricLabelsBuilder = NewMetricLabelsBuilder()
)

func NewMetricLabelsBuilder(labels ...string) MetricLabelsBuilder {
	builder := &metricLabelsBuilder{
		labels: make(map[string]string),
	}
	for _, label := range labels {
		builder.labels[label] = ""
	}
	return builder
}

type metricLabelsBuilder struct {
	labels prometheus.Labels
}

func (b *metricLabelsBuilder) SetKV(key, val string) {
	if _, ok := b.labels[key]; !ok {
		return
	}
	b.labels[key] = val
}

func (b *metricLabelsBuilder) SetKVs(kvs map[string]string) {
	for k, v := range kvs {
		b.SetKV(k, v)
	}
}

func (b *metricLabelsBuilder) Build() prometheus.Labels {
	return b.labels
}
