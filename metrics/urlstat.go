package metrics

import (
	"context"
	"fmt"
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/ops-agent/common/g"
	"github.com/fanghongbo/ops-agent/common/model"
	"net/http"
	"strconv"
	"time"
)

func UrlMetrics() (L []*model.MetricValue) {
	var (
		reportUrls map[string]string
		hostname   string
		err        error
	)

	reportUrls = g.ReportUrlMeta()
	sz := len(reportUrls)
	if sz == 0 {
		return
	}

	hostname, err = g.Hostname()
	if err != nil {
		hostname = "None"
		return
	}

	for furl, timeout := range reportUrls {
		tags := fmt.Sprintf("url=%v,timeout=%v,src=%v", furl, timeout, hostname)
		if ok := probeUrl(furl, timeout); !ok {
			L = append(L, GaugeValue(g.UrlCheckHealth, 0, tags))
			continue
		}
		L = append(L, GaugeValue(g.UrlCheckHealth, 1, tags))
	}
	return
}

func probeUrl(url string, t string) bool {
	var (
		err    error
		ctx    context.Context
		cancel context.CancelFunc
		req    *http.Request
		resp   *http.Response
	)

	timeout, err = strconv.Atoi(t)
	if err != nil {
		dlog.Errorf("convert %s timeout %s string to int err: %s", url, t, err)
		return false
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		dlog.Errorf("create new http request err: %s", err)
		return false
	}

	req = req.WithContext(ctx)

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		dlog.Errorf("request err: %s", err)
		return false
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		dlog.Errorf("get %s status code: %d", url, resp.StatusCode)
		return false
	}

	return true
}
