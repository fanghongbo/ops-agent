package http

import (
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/ops-agent/common/g"
	"net/http"
)

func Start() {
	var (
		addr string
		s    *http.Server
		err  error
	)

	if g.Conf().Http == nil || !g.Conf().Http.Enabled {
		dlog.Warning("http is disable")
		return
	}

	addr = g.Conf().Http.Listen
	if addr == "" {
		return
	}

	s = &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 30,
	}

	dlog.Infof("listening %s", addr)

	if err = s.ListenAndServe(); err != nil {
		dlog.Fatalf("start http server err: %s", err)
	}
}
