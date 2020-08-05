package metrics

import (
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/nux"
	"github.com/fanghongbo/ops-agent/common/model"
)

func LoadAvgMetrics() []*model.MetricValue {
	var (
		load *nux.Loadavg
		err  error
	)

	load, err = nux.LoadAvg()
	if err != nil {
		dlog.Errorf("get load avg err: %s", err)
		return nil
	}

	return []*model.MetricValue{
		GaugeValue("load.1min", load.Avg1min),
		GaugeValue("load.5min", load.Avg5min),
		GaugeValue("load.15min", load.Avg15min),
	}

}
