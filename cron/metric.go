package cron

import (
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/ops-agent/common/g"
	"github.com/fanghongbo/ops-agent/common/model"
	"strconv"
	"strings"
	"time"
)

func SyncBuiltinMetrics() {
	if g.Conf().Heartbeat == nil || !g.Conf().Heartbeat.Enabled {
		dlog.Warning("heartbeat is disable, builtin metric collector does not work")
		return
	}

	go syncBuiltinMetrics()
}

func syncBuiltinMetrics() {
	var (
		timestamp int64  = -1
		checksum  string = "nil"
	)

	for {
		var (
			interval time.Duration
			portData []int64
			pathData []string
			procData map[string]map[int]string
			urlData  map[string]string
			req      model.AgentHeartbeatRequest
			resp     model.BuiltinMetricResponse
		)

		interval = time.Duration(g.Conf().Heartbeat.Interval) * time.Second
		time.Sleep(interval)

		portData = []int64{}
		pathData = []string{}
		procData = make(map[string]map[int]string)
		urlData = make(map[string]string)

		hostname, err := g.Hostname()
		if err != nil {
			continue
		}

		req = model.AgentHeartbeatRequest{
			Hostname: hostname,
			Checksum: checksum,
		}

		err = g.HbsClient.Call("Agent.BuiltinMetrics", req, &resp)
		if err != nil {
			dlog.Errorf("call Agent.BuiltinMetrics err: %s", err)
			continue
		}

		if resp.Timestamp <= timestamp {
			continue
		}

		if resp.Checksum == checksum {
			continue
		}

		timestamp = resp.Timestamp
		checksum = resp.Checksum

		for _, metric := range resp.Metrics {
			if metric.Metric == g.UrlCheckHealth {
				var (
					urlOption     []string
					timeoutOption []string
					arr           []string
				)

				arr = strings.Split(metric.Tags, ",")
				for _, item := range arr {
					item = strings.TrimSpace(item)
					if strings.HasPrefix(item, "url=") {
						urlOption = strings.Split(item, "=")

					}
					if strings.HasPrefix(item, "timeout=") {
						timeoutOption = strings.Split(item, "=")
					}
				}

				if len(urlOption) != 2 {
					dlog.Errorf("%s url argument is missing", metric.String())
					continue
				}

				// init timeout default 30s
				if len(timeoutOption) != 2 {
					timeoutOption = []string{"timeout", "30"}
				}

				if n, err := strconv.ParseInt(timeoutOption[1], 10, 64); err == nil {
					// timeout gt 60s is not allow
					if n <= 0 || n > 60 {
						timeoutOption = []string{"timeout", "30"}
					}
					urlData[urlOption[1]] = timeoutOption[1]
				} else {
					dlog.Errorf("metric %s convert timeout string to int err: %s", g.UrlCheckHealth, err)
				}

			} else if metric.Metric == g.NetPortListen {
				// 端口监控
				var (
					portOption []string
					arr        []string
				)

				arr = strings.Split(metric.Tags, ",")

				for _, item := range arr {
					item = strings.TrimSpace(item)
					if strings.HasPrefix(item, "port=") {
						portOption = strings.Split(item, "=")
					}
				}

				if len(portOption) != 2 {
					dlog.Errorf("%s port argument is missing", metric.String())
					continue
				}

				if port, err := strconv.ParseInt(portOption[1], 10, 64); err == nil {
					portData = append(portData, port)
				} else {
					dlog.Errorf("metrics %s convert string to int err: %s", g.NetPortListen, err)
				}

			} else if metric.Metric == g.DuBs {
				// 目录监控
				var (
					pathOption []string
					arr        []string
				)

				arr = strings.Split(metric.Tags, ",")

				for _, item := range arr {
					item = strings.TrimSpace(item)
					if strings.HasPrefix(item, "path=") {
						pathOption = strings.Split(item, "=")
					}
				}

				if len(pathOption) != 2 {
					dlog.Errorf("%s path argument is missing", metric.String())
					continue
				}

				pathData = append(pathData, strings.TrimSpace(pathOption[1]))

			} else if metric.Metric == g.ProcNum {
				// 进程监控
				arr := strings.Split(metric.Tags, ",")
				tmpMap := make(map[int]string)

				for i := 0; i < len(arr); i++ {
					if strings.HasPrefix(arr[i], "name=") {
						tmpMap[1] = strings.TrimSpace(arr[i][5:])
					} else if strings.HasPrefix(arr[i], "cmdline=") {
						tmpMap[2] = strings.TrimSpace(arr[i][8:])
					}
				}

				procData[metric.Tags] = tmpMap

			} else {
				dlog.Errorf("invalid metric: %s tags: %s", metric.Metric, metric.Tags)
			}
		}

		g.SetReportUrlMeta(urlData)
		g.SetReportPortMeta(portData)
		g.SetReportProcMeta(procData)
		g.SetDuPathMeta(pathData)
	}
}
