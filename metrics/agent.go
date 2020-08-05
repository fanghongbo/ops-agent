package metrics

import (
	"github.com/fanghongbo/ops-agent/common/model"
)

func AgentMetrics() []*model.MetricValue {
	return []*model.MetricValue{GaugeValue("agent.alive", 1)}
}
