package metrics

import (
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/nux"
	"github.com/fanghongbo/ops-agent/common/g"
	"github.com/fanghongbo/ops-agent/common/model"
	"strings"
)

func ProcMetrics() (L []*model.MetricValue) {
	var (
		reportProc map[string]map[int]string
		sz         int
		ps         []*nux.Proc
		psLen      int
		err        error
	)

	reportProc = g.ReportProcMeta()
	sz = len(reportProc)
	if sz == 0 {
		return
	}

	ps, err = nux.AllProcs()
	if err != nil {
		dlog.Errorf("get all proc err: %s", err)
		return
	}

	psLen = len(ps)

	for tags, m := range reportProc {
		cnt := 0
		for i := 0; i < psLen; i++ {
			if validator(ps[i], m) {
				cnt++
			}
		}

		L = append(L, GaugeValue(g.ProcNum, cnt, tags))
	}

	return
}

func validator(p *nux.Proc, m map[int]string) bool {
	// only one kv pair
	for key, val := range m {
		if key == 1 {
			// name
			if val != p.Name {
				return false
			}
		} else if key == 2 {
			// cmdline
			if !strings.Contains(p.Cmdline, val) {
				return false
			}
		}
	}
	return true
}
