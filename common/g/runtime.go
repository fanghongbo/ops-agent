package g

import (
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/ops-agent/utils"
	"runtime"
	"time"
)

func InitRuntime() {
	var osType string

	osType = runtime.GOOS
	switch osType {
	case "linux":
	default:
		dlog.Fatalf("the %s system platform is not supported\n", osType)
	}

	// 内存监控
	go MemMonitor()
}

func MemMonitor() {
	go func() {
		var (
			nowMemUsedMB uint64
			maxMemMB     uint64
			rate         uint64
		)

		for {
			time.Sleep(time.Second * 10)

			nowMemUsedMB = getMemUsedMB()
			maxMemMB = uint64(utils.CalculateMemLimit(config.MaxMemRate))
			rate = (nowMemUsedMB * 100) / maxMemMB

			if config.Debug {
				dlog.Infof("agent mem used: %dMB, percent: %d%%", nowMemUsedMB, rate)
			}

			// 若超50%限制，打印 warning
			// 超过100%，就退出了
			if rate > 50 {
				dlog.Warningf("agent heap memory used rate, current: %d%%", rate)
			}
			if rate > 100 {
				// 堆内存已超过限制，退出进程
				dlog.Fatalf("heap memory size over limit. quit process.[used:%dMB][limit:%dMB][rate:%d]", nowMemUsedMB, maxMemMB, rate)
			}
		}
	}()
}

func getMemUsedMB() uint64 {
	var (
		sts runtime.MemStats
		ret uint64
	)

	runtime.ReadMemStats(&sts)
	// 这里取了mem.Alloc
	ret = sts.HeapAlloc / 1024 / 1024
	return ret
}
