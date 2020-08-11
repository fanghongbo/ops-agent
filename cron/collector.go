package cron

import (
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/ops-agent/common/g"
	"github.com/fanghongbo/ops-agent/common/model"
	"github.com/fanghongbo/ops-agent/metrics"
	"time"
)

func Collect() {
	if g.Conf().Transfer == nil || !g.Conf().Transfer.Enabled {
		dlog.Warning("transfer is disable, metric collector does not work")
		return
	}

	mappers := metrics.InitMetricFuncMappers()
	for _, v := range mappers {
		go collect(int64(v.Interval), v.Fs)
	}
}

func collect(sec int64, fns []func() []*model.MetricValue) {
	t := time.NewTicker(time.Second * time.Duration(sec))
	defer t.Stop()

	for {
		var (
			hostname      string
			mvs           []*model.MetricValue
			ignoreMetrics map[string]bool
			err           error
		)

		<-t.C

		hostname, err = g.Hostname()
		if err != nil {
			dlog.Errorf("get hostname err: %s", err)
			continue
		}

		ignoreMetrics = g.Conf().IgnoreMetrics

		for _, fn := range fns {
			items := fn()
			if items == nil {
				continue
			}

			if len(items) == 0 {
				continue
			}

			for _, mv := range items {
				if b, ok := ignoreMetrics[mv.Metric]; ok && b {
					continue
				} else {
					mvs = append(mvs, mv)
				}
			}
		}

		now := time.Now().Unix()
		for j := 0; j < len(mvs); j++ {
			mvs[j].Step = sec
			mvs[j].Endpoint = hostname
			mvs[j].Timestamp = now
		}

		g.SendToTransfer(mvs)
	}
}
