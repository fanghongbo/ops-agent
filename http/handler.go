package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fanghongbo/ops-agent/common/g"
	"github.com/fanghongbo/ops-agent/common/git"
	"github.com/fanghongbo/ops-agent/metrics"
	"net/http"
	"strings"
)

func IsLocalRequest(r *http.Request) bool {
	addr := strings.Split(r.RemoteAddr, ":")
	if addr[0] == "127.0.0.1" {
		return true
	}
	return false
}

func RenderJson(w http.ResponseWriter, v interface{}) {
	bs, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, _ = w.Write(bs)
}

func GetMetrics() []map[string]interface{} {
	var data []map[string]interface{}

	mappers := metrics.InitMetricFuncMappers()
	for _, item := range mappers {
		for _, fs := range item.Fs {
			res := fs()
			for _, metric := range res {
				if b, ok := g.Conf().IgnoreMetrics[metric.Metric]; ok && b {
					continue
				}

				if metric.Tags == "" {
					data = append(data, map[string]interface{}{
						"metric": metric.Metric,
						"value":  metric.Value,
					})
				} else {
					data = append(data, map[string]interface{}{
						"metric": fmt.Sprintf("%s/%s", metric.Metric, metric.Tags),
						"value":  metric.Value,
					})
				}
			}
		}
	}

	return data
}

func UpdatePlugin() error {
	var (
		repo git.NewGitClient
		err  error
	)

	if !g.Conf().Plugin.Enabled {
		return errors.New("plugin is disable")
	}

	if g.Conf().Plugin.Username == "" || g.Conf().Plugin.Password == "" {
		repo = git.NewGitClient{
			Url:      g.Conf().Plugin.Git,
			Path:     g.Conf().Plugin.Dir,
			RepoType: 0,
		}
	} else {
		repo = git.NewGitClient{
			Url:      g.Conf().Plugin.Git,
			Path:     g.Conf().Plugin.Dir,
			Username: g.Conf().Plugin.Username,
			Password: g.Conf().Plugin.Password,
			RepoType: 1,
		}
	}

	if err = repo.Update(); err != nil {
		return err
	}

	return nil
}

func GetPluginVersion() (string, error) {
	var (
		repo git.NewGitClient
		hash string
		err  error
	)

	if !g.Conf().Plugin.Enabled {
		return "", errors.New("plugin is disable")
	}

	if g.Conf().Plugin.Username == "" || g.Conf().Plugin.Password == "" {
		repo = git.NewGitClient{
			Url:      g.Conf().Plugin.Git,
			Path:     g.Conf().Plugin.Dir,
			RepoType: 0,
		}
	} else {
		repo = git.NewGitClient{
			Url:      g.Conf().Plugin.Git,
			Path:     g.Conf().Plugin.Dir,
			Username: g.Conf().Plugin.Username,
			Password: g.Conf().Plugin.Password,
			RepoType: 1,
		}
	}

	if hash, err = repo.Head(); err != nil {
		return "", err
	} else {
		return hash, nil
	}
}
