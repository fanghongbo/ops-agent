package model

import (
	"fmt"
)

type AgentReportRequest struct {
	Hostname      string
	IP            string
	AgentVersion  string
	PluginVersion string
}

func (u *AgentReportRequest) String() string {
	return fmt.Sprintf(
		"<Hostname:%s, IP:%s, AgentVersion:%s, PluginVersion:%s>",
		u.Hostname,
		u.IP,
		u.AgentVersion,
		u.PluginVersion,
	)
}

type AgentUpdateInfo struct {
	LastUpdate    int64
	ReportRequest *AgentReportRequest
}

type AgentHeartbeatRequest struct {
	Hostname string
	Checksum string
}

func (u *AgentHeartbeatRequest) String() string {
	return fmt.Sprintf(
		"<Hostname: %s, Checksum: %s>",
		u.Hostname,
		u.Checksum,
	)
}

type AgentPluginsResponse struct {
	Plugins   []string
	Timestamp int64
}

func (u *AgentPluginsResponse) String() string {
	return fmt.Sprintf(
		"<Plugins:%v, Timestamp:%v>",
		u.Plugins,
		u.Timestamp,
	)
}

// e.g. net.port.listen or proc.num
type BuiltinMetric struct {
	Metric string
	Tags   string
}

func (u *BuiltinMetric) String() string {
	return fmt.Sprintf(
		"%s/%s",
		u.Metric,
		u.Tags,
	)
}

type BuiltinMetricResponse struct {
	Metrics   []*BuiltinMetric
	Checksum  string
	Timestamp int64
}

func (u *BuiltinMetricResponse) String() string {
	return fmt.Sprintf(
		"<Metrics:%v, Checksum:%s, Timestamp:%v>",
		u.Metrics,
		u.Checksum,
		u.Timestamp,
	)
}

type BuiltinMetricSlice []*BuiltinMetric

func (u BuiltinMetricSlice) Len() int {
	return len(u)
}

func (u BuiltinMetricSlice) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

func (u BuiltinMetricSlice) Less(i, j int) bool {
	return u[i].String() < u[j].String()
}
