package cron

import (
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/ops-agent/common/g"
	"github.com/fanghongbo/ops-agent/common/model"
	"github.com/fanghongbo/ops-agent/utils"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func SyncMinePlugins() {
	if g.Conf().Plugin == nil || !g.Conf().Plugin.Enabled {
		dlog.Warning("plugin is disable, plugin does not work")
		return
	}

	if g.Conf().Heartbeat == nil || !g.Conf().Heartbeat.Enabled {
		dlog.Warning("heartbeat is disable, plugin does not work")
		return
	}

	go syncMinePlugins()
}

func syncMinePlugins() {
	for {
		var (
			interval    time.Duration
			timestamp   int64 = -1
			pluginDirs  []string
			desiredAll  map[string]*g.Plugin
			re          *regexp.Regexp
			scriptFiles [][]string
			scriptDirs  []string
			resp        model.AgentPluginsResponse
		)

		interval = time.Duration(g.Conf().Heartbeat.Interval) * time.Second
		time.Sleep(interval)

		hostname, err := g.Hostname()
		if err != nil {
			dlog.Errorf("get hostname err: %s", err)
			continue
		}

		req := model.AgentHeartbeatRequest{
			Hostname: hostname,
		}

		err = g.HbsClient.Call("Agent.MinePlugins", req, &resp)
		if err != nil {
			dlog.Errorf("call Agent.MinePlugin fail: %s", err)
			continue
		}

		if resp.Timestamp <= timestamp {
			continue
		}

		pluginDirs = resp.Plugins
		timestamp = resp.Timestamp

		if g.Conf().Debug {
			dlog.Infof("call Agent.MinePlugin: %v", resp)
		}

		if len(pluginDirs) == 0 {
			g.ClearAllPlugins()
			continue
		}

		desiredAll = make(map[string]*g.Plugin)
		re = regexp.MustCompile(`(.*)\((.*)\)`)

		for _, scriptItem := range pluginDirs {
			// 插件配置项, 约定是目录文件、脚本路径或者是带参数的脚本路径
			// 比如： sys/ntp/60_ntp.py(arg1,arg2) 或者 sys/ntp/60_ntp.py 或者 sys/ntp
			// 1. 参数只对单个脚本文件生效，目录不支持参数
			// 2. 如果某个目录下的某个脚本被单独绑定到某个机器，那么再次绑定该目录时，该文件会不会再次执行
			var args string = ""

			scriptItem = strings.TrimSpace(scriptItem)
			if scriptItem == "" {
				continue
			}

			matchArgs := re.FindAllStringSubmatch(scriptItem, -1)
			if matchArgs != nil {
				scriptItem = matchArgs[0][1]
				args = matchArgs[0][2]
			}

			absPath := filepath.Join(g.Conf().Plugin.Dir, scriptItem)
			if !utils.FileIsExist(absPath) {
				dlog.Errorf("%s is not exist", absPath)
				continue
			}

			// 对脚本文件和目录进行归类
			if utils.IsFile(absPath) {
				scriptFiles = append(scriptFiles, []string{scriptItem, args})
			} else {
				scriptDirs = append(scriptDirs, scriptItem)
			}
		}

		taken := make(map[string]struct{})
		for _, scriptFile := range scriptFiles {
			var (
				cycle    int
				err      error
				absPath  string
				fileName string
				arr      []string
				fi       os.FileInfo
			)

			absPath = filepath.Join(g.Conf().Plugin.Dir, scriptFile[0])
			_, fileName = filepath.Split(absPath)
			arr = strings.Split(fileName, "_")

			cycle, err = strconv.Atoi(arr[0])
			if err == nil {
				fi, _ = os.Stat(absPath)
				plugin := &g.Plugin{FilePath: scriptFile[0], MTime: fi.ModTime().Unix(), Cycle: cycle, Args: scriptFile[1]}
				desiredAll[scriptFile[0]+"("+scriptFile[1]+")"] = plugin
			}

			//针对某个 host group 绑定了单个脚本后，再绑定该脚本的目录时，会忽略目录中的该文件
			taken[scriptFile[0]] = struct{}{}
		}

		for _, scriptDir := range scriptDirs {
			ps := g.ListPlugins(strings.Trim(scriptDir, "/"))
			for k, p := range ps {
				if _, ok := taken[k]; ok {
					continue
				}
				desiredAll[k] = p
			}
		}

		g.DelNoUsePlugins(desiredAll)
		g.AddNewPlugins(desiredAll)

		if g.Conf().Debug {
			dlog.Infof("current plugins: %v", g.Plugins)
		}
	}
}
