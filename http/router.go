package http

import (
	"encoding/json"
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/ops-agent/common/g"
	"github.com/fanghongbo/ops-agent/common/model"
	"github.com/fanghongbo/ops-agent/metrics"
	"net/http"
	"os"
	"time"
)

func init() {
	http.HandleFunc("/v1/push", func(w http.ResponseWriter, req *http.Request) {
		var (
			metric  []*model.MetricValue
			err     error
			decoder *json.Decoder
		)

		if req.ContentLength == 0 {
			w.WriteHeader(http.StatusBadRequest)
			RenderJson(w, map[string]interface{}{
				"success": false,
				"msg":     "body is blank",
				"data":    nil,
			})
			return
		}

		decoder = json.NewDecoder(req.Body)
		err = decoder.Decode(&metric)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			RenderJson(w, map[string]interface{}{
				"success": false,
				"msg":     err.Error(),
				"data":    nil,
			})
			return
		}

		g.SendToTransfer(metric)

		RenderJson(w, map[string]interface{}{
			"success": true,
			"msg":     "push success",
			"data":    nil,
		})
	})

	http.HandleFunc("/config/reload", func(w http.ResponseWriter, r *http.Request) {
		if IsLocalRequest(r) {
			if err := g.ReloadConfig(); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				RenderJson(w, map[string]interface{}{
					"success": false,
					"msg":     err.Error(),
					"data":    nil,
				})
			} else {
				RenderJson(w, map[string]interface{}{
					"success": true,
					"msg":     "reload success",
					"data":    nil,
				})
			}
		} else {
			w.WriteHeader(http.StatusForbidden)
			RenderJson(w, map[string]interface{}{
				"success": false,
				"msg":     "no privilege",
				"data":    nil,
			})
		}
	})

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		if IsLocalRequest(r) {
			data := GetMetrics()
			RenderJson(w, map[string]interface{}{
				"success": true,
				"msg":     "query success",
				"data":    data,
			})
		} else {
			w.WriteHeader(http.StatusForbidden)
			RenderJson(w, map[string]interface{}{
				"success": false,
				"msg":     "no privilege",
				"data":    nil,
			})
		}
	})

	http.HandleFunc("/metric/check", func(w http.ResponseWriter, r *http.Request) {
		if IsLocalRequest(r) {
			data := metrics.CheckCollector()
			RenderJson(w, map[string]interface{}{
				"success": true,
				"msg":     "query success",
				"data":    data,
			})
		} else {
			w.WriteHeader(http.StatusForbidden)
			RenderJson(w, map[string]interface{}{
				"success": false,
				"msg":     "no privilege",
				"data":    nil,
			})
		}
	})

	http.HandleFunc("/plugins", func(w http.ResponseWriter, r *http.Request) {
		if IsLocalRequest(r) {
			RenderJson(w, map[string]interface{}{
				"success": true,
				"msg":     "query success",
				"data":    g.Plugins,
			})
		} else {
			w.WriteHeader(http.StatusForbidden)
			RenderJson(w, map[string]interface{}{
				"success": false,
				"msg":     "no privilege",
				"data":    nil,
			})
		}
	})

	http.HandleFunc("/plugin/update", func(w http.ResponseWriter, r *http.Request) {
		if IsLocalRequest(r) {
			var (
				err error
			)

			if err = UpdatePlugin(); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				RenderJson(w, map[string]interface{}{
					"success": false,
					"msg":     err.Error(),
					"data":    nil,
				})
			} else {
				RenderJson(w, map[string]interface{}{
					"success": true,
					"msg":     "update success",
					"data":    nil,
				})
			}

		} else {
			w.WriteHeader(http.StatusForbidden)
			RenderJson(w, map[string]interface{}{
				"success": false,
				"msg":     "no privilege",
				"data":    nil,
			})
		}
	})

	http.HandleFunc("/plugin/version", func(w http.ResponseWriter, r *http.Request) {
		if IsLocalRequest(r) {
			var (
				hash string
				err  error
			)

			if hash, err = GetPluginVersion(); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				RenderJson(w, map[string]interface{}{
					"success": false,
					"msg":     err.Error(),
					"data":    nil,
				})
			} else {
				RenderJson(w, map[string]interface{}{
					"success": true,
					"msg":     "query success",
					"data":    hash,
				})
			}

		} else {
			w.WriteHeader(http.StatusForbidden)
			RenderJson(w, map[string]interface{}{
				"success": false,
				"msg":     "no privilege",
				"data":    nil,
			})
		}
	})

	http.HandleFunc("/exit", func(w http.ResponseWriter, r *http.Request) {
		if IsLocalRequest(r) {
			RenderJson(w, map[string]interface{}{
				"success": true,
				"msg":     "exited success",
				"data":    nil,
			})
			go func() {
				time.Sleep(time.Second)
				dlog.Warning("exited..")
				os.Exit(0)
			}()
		} else {
			w.WriteHeader(http.StatusForbidden)
			RenderJson(w, map[string]interface{}{
				"success": false,
				"msg":     "no privilege",
				"data":    nil,
			})
		}
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		RenderJson(w, map[string]interface{}{
			"success": true,
			"msg":     "query success",
			"data":    "ok",
		})
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		RenderJson(w, map[string]interface{}{
			"success": true,
			"msg":     "query success",
			"data":    g.VersionInfo(),
		})
	})
}
