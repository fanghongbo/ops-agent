package metrics

import (
	"fmt"
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/nux"
	"github.com/fanghongbo/ops-agent/common/g"
	"github.com/fanghongbo/ops-agent/common/model"
	"github.com/fanghongbo/ops-agent/utils"
)

func PortMetrics() (L []*model.MetricValue) {
	var (
		reportPorts []int64
		allTcpPorts []int64
		allUdpPorts []int64
		sz          int
		err         error
	)

	reportPorts = g.ReportPortMeta()
	sz = len(reportPorts)
	if sz == 0 {
		return
	}

	allTcpPorts, err = nux.TcpPorts()
	if err != nil {
		dlog.Errorf("get tcp port err: %s", err)
		return
	}

	allUdpPorts, err = nux.UdpPorts()
	if err != nil {
		dlog.Errorf("get udp port err: %s", err)
		return
	}

	for i := 0; i < sz; i++ {
		tags := fmt.Sprintf("port=%d", reportPorts[i])
		if utils.ContainsInt64(allTcpPorts, reportPorts[i]) || utils.ContainsInt64(allUdpPorts, reportPorts[i]) {
			L = append(L, GaugeValue(g.NetPortListen, 1, tags))
		} else {
			L = append(L, GaugeValue(g.NetPortListen, 0, tags))
		}
	}

	return
}
