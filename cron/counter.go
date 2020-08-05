package cron

import (
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/ops-agent/metrics"
	"time"
)

const CollectInterval = time.Second

func InitCounterData() {
	for {
		initCpuStatCounter()
		initDiskStatCounter()
		time.Sleep(CollectInterval)
	}
}

func initCpuStatCounter() {
	go func() {
		if err := metrics.UpdateCpuStat(); err != nil {
			dlog.Errorf("update cpu stats err: %s", err)
		}
	}()
}

func initDiskStatCounter() {
	go func() {
		if err := metrics.UpdateDiskStats(); err != nil {
			dlog.Errorf("update disk stats err: %s", err)
		}
	}()
}
