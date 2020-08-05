package metrics

import (
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/nux"
	"github.com/fanghongbo/ops-agent/common/model"
)

func KernelMetrics() (L []*model.MetricValue) {
	var (
		maxFiles      uint64
		maxProc       uint64
		allocateFiles uint64
		err           error
	)

	maxFiles, err = nux.KernelMaxFiles()
	if err != nil {
		dlog.Errorf("get kernel.maxfiles err: %s", err)
		return
	}

	L = append(L, GaugeValue("kernel.maxfiles", maxFiles))

	maxProc, err = nux.KernelMaxProc()
	if err != nil {
		dlog.Errorf("get kernel.maxproc err: %s", err)
		return
	}

	L = append(L, GaugeValue("kernel.maxproc", maxProc))

	allocateFiles, err = nux.KernelAllocateFiles()
	if err != nil {
		dlog.Errorf("get kernel.files.allocated err: %s", err)
		return
	}

	L = append(L, GaugeValue("kernel.files.allocated", allocateFiles))
	L = append(L, GaugeValue("kernel.files.left", maxFiles-allocateFiles))
	return
}
