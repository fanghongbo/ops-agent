package metrics

import (
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/nux"
	"github.com/fanghongbo/ops-agent/common/g"
	"github.com/fanghongbo/ops-agent/common/model"
)

func NetMetrics() []*model.MetricValue {
	return CoreNetMetrics(g.Conf().Collector.IfacePrefix)
}

func CoreNetMetrics(interfacePrefix []string) []*model.MetricValue {
	var (
		netIfs []*nux.NetIf
		ret    []*model.MetricValue
		cnt    int
		err    error
	)

	netIfs, err = nux.NetIfs(interfacePrefix)
	if err != nil {
		dlog.Errorf("get net interface err: %s", err)
		return []*model.MetricValue{}
	}

	cnt = len(netIfs)
	ret = make([]*model.MetricValue, cnt*26)

	for idx, netIf := range netIfs {
		interfaceStr := "iface=" + netIf.Iface
		ret[idx*26+0] = CounterValue("net.if.in.bytes", netIf.InBytes, interfaceStr)
		ret[idx*26+1] = CounterValue("net.if.in.packets", netIf.InPackages, interfaceStr)
		ret[idx*26+2] = CounterValue("net.if.in.errors", netIf.InErrors, interfaceStr)
		ret[idx*26+3] = CounterValue("net.if.in.dropped", netIf.InDropped, interfaceStr)
		ret[idx*26+4] = CounterValue("net.if.in.fifo.errs", netIf.InFifoErrs, interfaceStr)
		ret[idx*26+5] = CounterValue("net.if.in.frame.errs", netIf.InFrameErrs, interfaceStr)
		ret[idx*26+6] = CounterValue("net.if.in.compressed", netIf.InCompressed, interfaceStr)
		ret[idx*26+7] = CounterValue("net.if.in.multicast", netIf.InMulticast, interfaceStr)
		ret[idx*26+8] = CounterValue("net.if.out.bytes", netIf.OutBytes, interfaceStr)
		ret[idx*26+9] = CounterValue("net.if.out.packets", netIf.OutPackages, interfaceStr)
		ret[idx*26+10] = CounterValue("net.if.out.errors", netIf.OutErrors, interfaceStr)
		ret[idx*26+11] = CounterValue("net.if.out.dropped", netIf.OutDropped, interfaceStr)
		ret[idx*26+12] = CounterValue("net.if.out.fifo.errs", netIf.OutFifoErrs, interfaceStr)
		ret[idx*26+13] = CounterValue("net.if.out.collisions", netIf.OutCollisions, interfaceStr)
		ret[idx*26+14] = CounterValue("net.if.out.carrier.errs", netIf.OutCarrierErrs, interfaceStr)
		ret[idx*26+15] = CounterValue("net.if.out.compressed", netIf.OutCompressed, interfaceStr)
		ret[idx*26+16] = CounterValue("net.if.total.bytes", netIf.TotalBytes, interfaceStr)
		ret[idx*26+17] = CounterValue("net.if.total.packets", netIf.TotalPackages, interfaceStr)
		ret[idx*26+18] = CounterValue("net.if.total.errors", netIf.TotalErrors, interfaceStr)
		ret[idx*26+19] = CounterValue("net.if.total.dropped", netIf.TotalDropped, interfaceStr)
		ret[idx*26+20] = GaugeValue("net.if.speed.bits", netIf.SpeedBits, interfaceStr)
		ret[idx*26+21] = CounterValue("net.if.in.percent", netIf.InPercent, interfaceStr)
		ret[idx*26+22] = CounterValue("net.if.out.percent", netIf.OutPercent, interfaceStr)
		ret[idx*26+23] = CounterValue("net.if.in.bits", netIf.InBytes*8, interfaceStr)
		ret[idx*26+24] = CounterValue("net.if.out.bits", netIf.OutBytes*8, interfaceStr)
		ret[idx*26+25] = CounterValue("net.if.total.bits", netIf.TotalBytes*8, interfaceStr)
	}
	return ret
}
