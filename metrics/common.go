package metrics

import (
	"github.com/fanghongbo/ops-agent/common/model"
	"strings"
)

func NewMetricValue(metric string, val interface{}, dataType string, tags ...string) *model.MetricValue {
	var (
		mv   model.MetricValue
		size int
	)

	mv = model.MetricValue{
		Metric:      metric,
		Value:       val,
		CounterType: dataType,
	}

	size = len(tags)
	if size > 0 {
		mv.Tags = strings.Join(tags, ",")
	}

	return &mv
}

func GaugeValue(metric string, val interface{}, tags ...string) *model.MetricValue {
	return NewMetricValue(metric, val, "GAUGE", tags...)
}

func CounterValue(metric string, val interface{}, tags ...string) *model.MetricValue {
	return NewMetricValue(metric, val, "COUNTER", tags...)
}
