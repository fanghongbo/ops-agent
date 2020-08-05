package g

import "sync"

const NetPortListen  = "net.port.listen"

var (
	reportPorts     []int64
	reportPortsLock = new(sync.RWMutex)
)

func ReportPortMeta() []int64 {
	reportPortsLock.RLock()
	defer reportPortsLock.RUnlock()
	return reportPorts
}

func SetReportPortMeta(ports []int64) {
	reportPortsLock.Lock()
	defer reportPortsLock.Unlock()
	reportPorts = ports
}
