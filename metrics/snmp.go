package metrics

import (
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/nux"
	"github.com/fanghongbo/ops-agent/common/model"
)

func UdpMetrics() []*model.MetricValue {
	var (
		udp   map[string]int64
		count int
		ret   []*model.MetricValue
		err   error
	)

	udp, err = nux.Snmp("Udp")
	if err != nil {
		dlog.Errorf("read snmp err: %s", err)
		return []*model.MetricValue{}
	}

	count = len(udp)
	ret = make([]*model.MetricValue, count)
	i := 0
	for key, val := range udp {
		ret[i] = CounterValue("snmp.Udp."+key, val)
		i++
	}

	return ret
}
