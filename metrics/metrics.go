package metrics

import (
	"github.com/fanghongbo/ops-agent/common/g"
	"github.com/fanghongbo/ops-agent/common/model"
)

type MetricFunc struct {
	Fs       []func() []*model.MetricValue
	Interval int
}

func InitMetricFuncMappers() []MetricFunc {
	var (
		interval int
		Mappers  []MetricFunc
	)

	interval = g.Conf().Transfer.Interval
	Mappers = []MetricFunc{
		{
			Fs: []func() []*model.MetricValue{
				AgentMetrics,
				CpuMetrics,
				NetMetrics,
				KernelMetrics,
				LoadAvgMetrics,
				MemMetrics,
				DiskIOMetrics,
				IOStatsMetrics,
				NetStatMetrics,
				ProcMetrics,
				UdpMetrics,
			},
			Interval: interval,
		},
		{
			Fs: []func() []*model.MetricValue{
				DeviceMetrics,
			},
			Interval: interval,
		},
		{
			Fs: []func() []*model.MetricValue{
				PortMetrics,
				SocketStatSummaryMetrics,
			},
			Interval: interval,
		},
		{
			Fs: []func() []*model.MetricValue{
				DuMetrics,
			},
			Interval: interval,
		},
		{
			Fs: []func() []*model.MetricValue{
				UrlMetrics,
			},
			Interval: interval,
		},
		{
			Fs: []func() []*model.MetricValue{
				GpuMetrics,
			},
			Interval: interval,
		},
	}

	return Mappers
}
