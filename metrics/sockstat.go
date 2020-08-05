package metrics

import (
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/nux"
	"github.com/fanghongbo/ops-agent/common/model"
)

func SocketStatSummaryMetrics() (L []*model.MetricValue) {
	var (
		ssMap map[string]uint64
		err   error
	)

	ssMap, err = nux.SocketStatSummary()
	if err != nil {
		dlog.Errorf("get socket status summary err: %s", err)
		return
	}

	for k, v := range ssMap {
		L = append(L, GaugeValue("ss."+k, v))
	}

	return
}
