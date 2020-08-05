package utils

import (
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/nux"
	"math"
	"runtime"
)

func GetCPULimitNum(maxCPURate float64) int {
	var cpuLimit int

	cpuLimit = int(math.Ceil(float64(runtime.NumCPU()) * maxCPURate))
	if cpuLimit < 1 {
		cpuLimit = 1
	}
	return cpuLimit
}

func CalculateMemLimit(maxMemRate float64) int {
	var (
		memTotal, memLimit int
		m                  *nux.Mem
		err                error
	)

	m, err = nux.MemInfo()
	if err != nil {
		dlog.Error("failed to get mem.total:", err)
		memLimit = 512
	} else {
		memTotal = int(m.MemTotal / (1024 * 1024))
		memLimit = int(float64(memTotal) * maxMemRate)
	}

	if memLimit < 512 {
		memLimit = 512
	}

	return memLimit
}
