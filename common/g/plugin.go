package g

import (
	"bytes"
	"encoding/json"
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/ops-agent/common/cmd"
	"github.com/fanghongbo/ops-agent/common/model"
	"github.com/fanghongbo/ops-agent/utils"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	Plugins              = make(map[string]*Plugin)
	PluginsWithScheduler = make(map[string]*PluginScheduler)
)

type Plugin struct {
	FilePath string `json:"filepath"`
	MTime    int64  `json:"mtime"`
	Cycle    int    `json:"cycle"`
	Args     string `json:"args"`
}

func DelNoUsePlugins(newPlugins map[string]*Plugin) {
	for currKey, currPlugin := range Plugins {
		newPlugin, ok := newPlugins[currKey]
		if !ok || currPlugin.MTime != newPlugin.MTime {
			deletePlugin(currKey)
		}
	}
}

func AddNewPlugins(newPlugins map[string]*Plugin) {
	for filePath, newPlugin := range newPlugins {
		if _, ok := Plugins[filePath]; ok && newPlugin.MTime == Plugins[filePath].MTime {
			continue
		}

		Plugins[filePath] = newPlugin
		newScheduler := NewPluginScheduler(newPlugin)
		PluginsWithScheduler[filePath] = newScheduler
		newScheduler.Schedule()
	}
}

func ClearAllPlugins() {
	for k := range Plugins {
		deletePlugin(k)
	}
}

func deletePlugin(key string) {
	v, ok := PluginsWithScheduler[key]
	if ok {
		v.Stop()
		delete(PluginsWithScheduler, key)
	}
	delete(Plugins, key)
}

type PluginScheduler struct {
	Ticker *time.Ticker
	Plugin *Plugin
	Quit   chan struct{}
}

func NewPluginScheduler(p *Plugin) *PluginScheduler {
	var scheduler PluginScheduler

	scheduler = PluginScheduler{Plugin: p}
	scheduler.Ticker = time.NewTicker(time.Duration(p.Cycle) * time.Second)
	scheduler.Quit = make(chan struct{})
	return &scheduler
}

func (u *PluginScheduler) Schedule() {
	go func() {
		for {
			select {
			case <-u.Ticker.C:
				PluginRun(u.Plugin)
			case <-u.Quit:
				u.Ticker.Stop()
				return
			}
		}
	}()
}

func (u *PluginScheduler) Stop() {
	close(u.Quit)
}

func PluginArgsParse(rawArgs string) []string {
	var (
		out  [][]string
		ss   []string
		ret  []string
		tail string
	)

	ss = strings.Split(rawArgs, "\\,")
	out = [][]string{}

	for _, s := range ss {
		var (
			cleanArgs []string
		)

		cleanArgs = []string{}
		for _, arg := range strings.Split(s, ",") {
			arg = strings.Trim(arg, " ")
			arg = strings.Trim(arg, "\"")
			arg = strings.Trim(arg, "'")
			cleanArgs = append(cleanArgs, arg)
		}
		out = append(out, cleanArgs)
	}

	ret = []string{}
	tail = ""

	for _, x := range out {
		for j, y := range x {
			if j == 0 {
				if tail != "" {
					ret = append(ret, tail+","+y)
					tail = ""
				} else {
					ret = append(ret, y)
				}
			} else if j == len(x)-1 {
				tail = y
			} else {
				ret = append(ret, y)
			}
		}
	}

	if tail != "" {
		ret = append(ret, tail)
	}

	return ret
}

func PluginRun(plugin *Plugin) {
	var (
		timeout   int
		filePath  string
		args      string
		command   *exec.Cmd
		stdout    bytes.Buffer
		stderr    bytes.Buffer
		err       error
		isTimeout bool
		errStr    string
		metrics   []*model.MetricValue
	)

	timeout = plugin.Cycle*1000 - 500
	filePath = filepath.Join(config.Plugin.Dir, plugin.FilePath)
	args = plugin.Args

	if !utils.FileIsExist(filePath) {
		dlog.Infof("no such plugin: %s args: (%s)", filePath, args)
		return
	}

	debug := config.Debug
	if debug {
		dlog.Infof("plugin: %s args: (%s) running...", filePath, args)
	}

	if args == "" {
		command = exec.Command(filePath)
	} else {
		argList := PluginArgsParse(args)
		command = exec.Command(filePath, argList...)
	}

	command.Stdout = &stdout
	command.Stderr = &stderr

	// 将进程GID设置成与PID相同的值
	command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	err = command.Start()
	if err != nil {
		dlog.Errorf("plugin: %s args: (%s) start fail, error: %s", filePath, args, err)
		return
	}

	if debug {
		dlog.Infof("plugin: %s args: (%s) started", filePath, args)
	}

	err, isTimeout = cmd.RunLocalCommandWithTimeout(command, time.Duration(timeout)*time.Millisecond)

	errStr = stderr.String()
	if errStr != "" {
		logFile := filepath.Join(config.Plugin.LogDir, plugin.FilePath+"("+plugin.Args+")"+".stderr.log")
		if _, err = utils.WriteString(logFile, errStr); err != nil {
			dlog.Printf("write log to %s fail, error: %s", logFile, err)
		}
	}

	if isTimeout {
		// has be killed
		if err == nil && debug {
			dlog.Errorf("timeout and kill plugin: %s args: (%s) successfully", filePath, args)
		}

		if err != nil {
			dlog.Errorf("kill plugin: %s args: (%s) occur error: %s", filePath, args, err)
		}

		return
	}

	if err != nil {
		dlog.Errorf("exec plugin: %s args: (%s) fail, error: %s", filePath, args, err)
		return
	}

	// exec successfully
	data := stdout.Bytes()
	if len(data) == 0 {
		if debug {
			dlog.Infof("stdout of plugin: %s args: (%s) is blank", filePath, args)
		}
		return
	}

	err = json.Unmarshal(data, &metrics)
	if err != nil {
		dlog.Errorf("json decode stdout of plugin: %s args: (%s) fail. error:%s stdout: %s", filePath, args, err, stdout.String())
		return
	}

	SendToTransfer(metrics)
}

func ListPlugins(scriptPath string) map[string]*Plugin {
	var (
		ret     map[string]*Plugin
		absPath string
		err     error
		fs      []os.FileInfo
	)

	ret = make(map[string]*Plugin)

	if scriptPath == "" {
		return ret
	}

	absPath = filepath.Join(config.Plugin.Dir, scriptPath)
	fs, err = ioutil.ReadDir(absPath)
	if err != nil {
		dlog.Errorf("can not list files under %s", absPath)
		return ret
	}

	for _, f := range fs {
		var (
			filename string
			arr      []string
			cycle    int
		)

		if f.IsDir() {
			continue
		}

		filename = f.Name()
		arr = strings.Split(filename, "_")
		if len(arr) < 2 {
			continue
		}

		// filename should be: $cycle_$xx
		cycle, err = strconv.Atoi(arr[0])
		if err != nil {
			continue
		}

		filePath := filepath.Join(scriptPath, filename)
		plugin := &Plugin{FilePath: filePath, MTime: f.ModTime().Unix(), Cycle: cycle, Args: ""}
		ret[filePath] = plugin
	}

	return ret
}
