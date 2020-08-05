package metrics

import (
	"github.com/fanghongbo/nux"
	"github.com/fanghongbo/ops-agent/common/cmd"
)

func CheckCollector() map[string]bool {
	var (
		output map[string]bool
	)

	output = make(map[string]bool)
	_, procStatErr := nux.CurrentProcStat()
	_, listDiskErr := nux.ListDiskStats()
	ports, listeningPortsErr := nux.ListeningPorts()
	proc, psErr := nux.AllProcs()

	_, duErr := cmd.RunLocalCommand("du", "--help")

	output["kernel"] = len(KernelMetrics()) > 0
	output["df.bytes"] = DeviceMetricsCheck()
	output["net.if"] = len(CoreNetMetrics([]string{})) > 0
	output["loadavg"] = len(LoadAvgMetrics()) > 0
	output["cpustat"] = procStatErr == nil
	output["disk.io"] = listDiskErr == nil
	output["memory"] = len(MemMetrics()) > 0
	output["netstat"] = len(NetStatMetrics()) > 0
	output["ss -s"] = len(SocketStatSummaryMetrics()) > 0
	output["ss -tln"] = listeningPortsErr == nil && len(ports) > 0
	output["ps aux"] = psErr == nil && len(proc) > 0
	output["du -bs"] = duErr == nil

	return output
}
