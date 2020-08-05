package metrics

import (
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/nux"
	"github.com/fanghongbo/ops-agent/common/model"
)

var USES = map[string]struct{}{
	"PruneCalled":        {},
	"LockDroppedIcmps":   {},
	"ArpFilter":          {},
	"TW":                 {},
	"DelayedACKLocked":   {},
	"ListenOverflows":    {},
	"ListenDrops":        {},
	"TCPPrequeueDropped": {},
	"TCPTSReorder":       {},
	"TCPDSACKUndo":       {},
	"TCPLoss":            {},
	"TCPLostRetransmit":  {},
	"TCPLossFailures":    {},
	"TCPFastRetrans":     {},
	"TCPTimeouts":        {},
	"TCPSchedulerFailed": {},
	"TCPAbortOnMemory":   {},
	"TCPAbortOnTimeout":  {},
	"TCPAbortFailed":     {},
	"TCPMemoryPressures": {},
	"TCPSpuriousRTOs":    {},
	"TCPBacklogDrop":     {},
	"TCPMinTTLDrop":      {},
}

func NetStatMetrics() (L []*model.MetricValue) {
	var (
		tcpExtList map[string]uint64
		cnt        int
		err        error
	)

	tcpExtList, err = nux.Netstat("TcpExt")

	if err != nil {
		dlog.Errorf("get net stats err: %s", err)
		return
	}

	cnt = len(tcpExtList)
	if cnt == 0 {
		return
	}

	for key, val := range tcpExtList {
		if _, ok := USES[key]; !ok {
			continue
		}
		L = append(L, CounterValue("TcpExt."+key, val))
	}

	return
}
