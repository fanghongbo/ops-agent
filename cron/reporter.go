package cron

import (
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/ops-agent/common/g"
	"github.com/fanghongbo/ops-agent/common/model"
	"github.com/fanghongbo/ops-agent/common/rpc"
	"github.com/fanghongbo/ops-agent/http"
	"time"
)

func ReportAgentStatus() {
	if g.Conf().Heartbeat.Enabled {
		go reportAgentStatus()
	} else {
		dlog.Warning("heartbeat is disable, agent status does not sent")
	}
}

func reportAgentStatus() {
	for {
		var (
			hostname string
			req      model.AgentReportRequest
			resp     rpc.SimpleRpcResponse
			err      error
			interval time.Duration
			hash     string
		)

		interval = time.Duration(g.Conf().Heartbeat.Interval) * time.Second
		time.Sleep(interval)

		hostname, err = g.Hostname()
		if err != nil {
			dlog.Errorf("get hostname err: %s", err)
			continue
		}

		if !g.Conf().Plugin.Enabled {
			hash, err = http.GetPluginVersion()
			if err != nil {
				dlog.Errorf("get plugin version err: %s", err)
			}
		}

		req = model.AgentReportRequest{
			Hostname:      hostname,
			IP:            g.IP(),
			AgentVersion:  g.Version,
			PluginVersion: hash,
		}

		err = g.HbsClient.Call("Agent.ReportStatus", req, &resp)
		if err != nil || resp.Code != 0 {
			dlog.Errorf("call Agent.ReportStatus fail: %s Request: %s Response: %s", err, req, resp)
			continue
		}
	}
}
