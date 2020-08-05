package g

import "sync"

const ProcNum = "proc.num"

var (
	reportProc     map[string]map[int]string
	reportProcLock = new(sync.RWMutex)
)

func ReportProcMeta() map[string]map[int]string {
	reportProcLock.RLock()
	defer reportProcLock.RUnlock()
	return reportProc
}

func SetReportProcMeta(proc map[string]map[int]string) {
	reportProcLock.Lock()
	defer reportProcLock.Unlock()
	reportProc = proc
}
