package metrics

import (
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/nux"
	"github.com/fanghongbo/ops-agent/common/model"
)

func MemMetrics() []*model.MetricValue {
	var (
		m   *nux.Mem
		err error
	)

	m, err = nux.MemInfo()
	if err != nil {
		dlog.Errorf("get mem info err: %s", err)
		return nil
	}

	memFree := m.MemFree + m.Buffers + m.Cached
	if m.MemAvailable > 0 {
		memFree = m.MemAvailable
	}
	memUsed := m.MemTotal - memFree

	freeMemRate := 0.0
	usedMemRate := 0.0
	if m.MemTotal != 0 {
		freeMemRate = float64(memFree) * 100.0 / float64(m.MemTotal)
		usedMemRate = float64(memUsed) * 100.0 / float64(m.MemTotal)
	}

	freeSwapRate := 0.0
	usedSwapRate := 0.0
	if m.SwapTotal != 0 {
		freeSwapRate = float64(m.SwapFree) * 100.0 / float64(m.SwapTotal)
		usedSwapRate = float64(m.SwapUsed) * 100.0 / float64(m.SwapTotal)
	}

	return []*model.MetricValue{
		GaugeValue("mem.memtotal", m.MemTotal),
		GaugeValue("mem.memused", memUsed),
		GaugeValue("mem.memfree", memFree),
		GaugeValue("mem.swaptotal", m.SwapTotal),
		GaugeValue("mem.swapused", m.SwapUsed),
		GaugeValue("mem.swapfree", m.SwapFree),
		GaugeValue("mem.memfree.percent", freeMemRate),
		GaugeValue("mem.memused.percent", usedMemRate),
		GaugeValue("mem.swapfree.percent", freeSwapRate),
		GaugeValue("mem.swapused.percent", usedSwapRate),
	}

}
