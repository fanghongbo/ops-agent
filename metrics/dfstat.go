package metrics

import (
	"fmt"
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/nux"
	"github.com/fanghongbo/ops-agent/common/g"
	"github.com/fanghongbo/ops-agent/common/model"
)

func DeviceMetrics() (L []*model.MetricValue) {
	var (
		mountPoints   [][3]string
		err           error
		myMountPoints map[string]bool = make(map[string]bool)
		diskTotal     uint64          = 0
		diskUsed      uint64          = 0
	)

	mountPoints, err = nux.ListMountPoint()
	if err != nil {
		dlog.Error("collect device metrics fail: %s", err)
		return
	}

	if len(g.Conf().Collector.MountPoint) > 0 {
		for _, mp := range g.Conf().Collector.MountPoint {
			myMountPoints[mp] = true
		}
	}

	for idx := range mountPoints {
		fsSpec, fsFile, fsVfsType := mountPoints[idx][0], mountPoints[idx][1], mountPoints[idx][2]
		if len(myMountPoints) > 0 {
			if _, ok := myMountPoints[fsFile]; !ok {
				dlog.Debug("mount point not matched with config", fsFile, "ignored.")
				continue
			}
		}

		var du *nux.DeviceUsage
		du, err = nux.BuildDeviceUsage(fsSpec, fsFile, fsVfsType)
		if err != nil {
			dlog.Errorf("Generate disk usage err: %s", err)
			continue
		}

		if du.BlocksAll == 0 {
			continue
		}

		diskTotal += du.BlocksAll
		diskUsed += du.BlocksUsed

		tags := fmt.Sprintf("mount=%s,fstype=%s", du.FsFile, du.FsVfstype)
		L = append(L, GaugeValue("df.bytes.total", du.BlocksAll, tags))
		L = append(L, GaugeValue("df.bytes.used", du.BlocksUsed, tags))
		L = append(L, GaugeValue("df.bytes.free", du.BlocksFree, tags))
		L = append(L, GaugeValue("df.bytes.used.percent", du.BlocksUsedPercent, tags))
		L = append(L, GaugeValue("df.bytes.free.percent", du.BlocksFreePercent, tags))

		if du.InodesAll == 0 {
			continue
		}

		L = append(L, GaugeValue("df.inodes.total", du.InodesAll, tags))
		L = append(L, GaugeValue("df.inodes.used", du.InodesUsed, tags))
		L = append(L, GaugeValue("df.inodes.free", du.InodesFree, tags))
		L = append(L, GaugeValue("df.inodes.used.percent", du.InodesUsedPercent, tags))
		L = append(L, GaugeValue("df.inodes.free.percent", du.InodesFreePercent, tags))

	}

	if len(L) > 0 && diskTotal > 0 {
		L = append(L, GaugeValue("df.statistics.total", float64(diskTotal)))
		L = append(L, GaugeValue("df.statistics.used", float64(diskUsed)))
		L = append(L, GaugeValue("df.statistics.used.percent", float64(diskUsed)*100.0/float64(diskTotal)))
	}

	return
}

func DeviceMetricsCheck() bool {
	var (
		mountPoints [][3]string
		err         error
	)

	mountPoints, err = nux.ListMountPoint()

	if err != nil {
		dlog.Error("collect device metrics fail: %s", err)
		return false
	}

	if len(mountPoints) <= 0 {
		return false
	}

	return true
}
