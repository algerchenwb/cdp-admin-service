package fasthttpmetric

import (
	"sync/atomic"

	"github.com/valyala/fasthttp"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic/metricwrap"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic/otelwrap/fasthttpotel"
)

var (
	requestMetricWrapHandler atomic.Value
)

func init() {
	UpdateHandler(nil)
}

func GetHandler() *RequestMetricWrapHandler {
	return requestMetricWrapHandler.Load().(*RequestMetricWrapHandler)
}

func UpdateHandler(cfg *metricwrap.MetricConfig) {
	requestMetricWrapHandler.Store(NewRequestMetricWrapHandler(cfg))
}

func RequestMetricWrap(fn fasthttp.RequestHandler) fasthttp.RequestHandler {
	return GetHandler().RequestMetricWrap(fn)
}

func NewRequestMetricWrapHandler(cfg *metricwrap.MetricConfig) *RequestMetricWrapHandler {
	if cfg == nil {
		cfg = metricwrap.NewMetricConfig()
	}
	if !gopublic.StringInArray(metricwrap.MetricLabelCode, cfg.LabelNames) {
		cfg.LabelNames = append(cfg.LabelNames, metricwrap.MetricLabelCode)
	}
	return &RequestMetricWrapHandler{
		requestMetric: metricwrap.NewRequestMetricHandler(cfg),
	}
}

type RequestMetricWrapHandler struct {
	requestMetric *metricwrap.RequestMetricHandler
}

func (h RequestMetricWrapHandler) RequestMetricWrap(fn fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		iCtx := fasthttpotel.GetTraceCtx(ctx)
		iCtx, handle := h.requestMetric.Metric(iCtx, string(ctx.Path()))
		defer handle()

		fasthttpotel.SetTraceCtx(ctx, iCtx)

		fn(ctx)
	}
}
