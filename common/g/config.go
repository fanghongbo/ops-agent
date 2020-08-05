package g

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/ops-agent/utils"
	"io/ioutil"
	"os"
	"runtime"
	"sync"
)

var (
	cfg            = flag.String("c", "./config/cfg.json", "specify config file")
	v              = flag.Bool("v", false, "show version")
	vv             = flag.Bool("vv", false, "show version detail")
	ConfigFile     string
	configFileLock = new(sync.RWMutex)
	config         *GlobalConfig
)

type LogConfig struct {
	LogPath      string `json:"log_path"`
	LogLevel     string `json:"log_level"`
	LogFileName  string `json:"log_file_name"`
	LogKeepHours int    `json:"log_keep_hours"`
}

type PluginConfig struct {
	Enabled  bool   `json:"enabled"`
	Dir      string `json:"dir"`
	Git      string `json:"git"`
	Username string `json:"username"`
	Password string `json:"password"`
	LogDir   string `json:"logs"`
}

type HeartbeatConfig struct {
	Enabled  bool   `json:"enabled"`
	Addr     string `json:"addr"`
	Interval int    `json:"interval"`
	Timeout  int    `json:"timeout"`
}

type TransferConfig struct {
	Enabled  bool     `json:"enabled"`
	Addrs    []string `json:"addrs"`
	Interval int      `json:"interval"`
	Timeout  int      `json:"timeout"`
}

type HttpConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type CollectorConfig struct {
	IfacePrefix []string `json:"ifacePrefix"`
	MountPoint  []string `json:"mountPoint"`
}

type GlobalConfig struct {
	Debug         bool              `json:"debug"`
	Hostname      string            `json:"hostname"`
	IP            string            `json:"ip"`
	Log           *LogConfig        `json:"log"`
	Plugin        *PluginConfig     `json:"plugin"`
	Heartbeat     *HeartbeatConfig  `json:"heartbeat"`
	Transfer      *TransferConfig   `json:"transfer"`
	Http          *HttpConfig       `json:"http"`
	Collector     *CollectorConfig  `json:"collector"`
	DefaultTags   map[string]string `json:"default_tags"`
	IgnoreMetrics map[string]bool   `json:"ignore"`
	MaxCPURate    float64           `json:"max_cpu_rate"`
	MaxMemRate    float64           `json:"max_mem_rate"`
}

func Conf() *GlobalConfig {
	configFileLock.RLock()
	defer configFileLock.RUnlock()

	return config
}

func InitConfig() {
	flag.Parse()

	if *v {
		fmt.Println(VersionInfo())
		os.Exit(0)
	}

	if *vv {
		fmt.Println(AgentInfo())
		os.Exit(0)
	}

	cfgFile := *cfg
	ConfigFile = cfgFile

	if cfgFile == "" {
		dlog.Fatal("config file not specified: use -c $filename")
	}

	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		dlog.Fatalf("config file specified not found: %s", cfgFile)
	} else {
		dlog.Infof("use config file: %s", ConfigFile)
	}

	if bs, err := ioutil.ReadFile(cfgFile); err != nil {
		dlog.Fatalf("read config file failed: %s", err.Error())
	} else {
		if err := json.Unmarshal(bs, &config); err != nil {
			dlog.Fatalf("decode config file failed: %s", err.Error())
		} else {
			dlog.Infof("load config success from %s", cfgFile)
		}
	}

	if err := Validator(); err != nil {
		dlog.Errorf("validator config file fail: %s", err)
		os.Exit(127)
	}

	// 最大使用内存
	maxMemMB := utils.CalculateMemLimit(config.MaxMemRate)

	// 最大cpu核数
	maxCPUNum := utils.GetCPULimitNum(config.MaxCPURate)

	dlog.Infof("bind [%d] cpu core", maxCPUNum)
	runtime.GOMAXPROCS(maxCPUNum)

	dlog.Infof("memory limit: %d MB", maxMemMB)
}

func ReloadConfig() error {
	if _, err := os.Stat(ConfigFile); os.IsNotExist(err) {
		dlog.Fatalf("config file specified not found: %s", ConfigFile)
		return err
	} else {
		dlog.Infof("reload config file: %s", ConfigFile)
	}

	if bs, err := ioutil.ReadFile(ConfigFile); err != nil {
		dlog.Fatalf("reload config file failed: %s", err)
		return err
	} else {
		configFileLock.RLock()
		defer configFileLock.RUnlock()

		if err := json.Unmarshal(bs, &config); err != nil {
			dlog.Fatalf("decode config file failed: %s", err)
			return err
		} else {
			dlog.Infof("reload config success from %s", ConfigFile)
		}
	}

	if err := Validator(); err != nil {
		dlog.Errorf("validator config file fail: %s", err)
		return err
	}

	return nil
}

func Validator() error {
	// 设置默认日志路径为 ./logs
	if config.Log.LogPath == "" {
		config.Log.LogPath = "./logs"
	}

	// 设置默认日志文件名称为 run.log
	if config.Log.LogFileName == "" {
		config.Log.LogFileName = "run.log"
	}

	// 设置默认日志级别为 LogLevel
	if config.Log.LogLevel == "" {
		config.Log.LogLevel = "INFO"
	}

	// 设置默认保留24小时的日志
	if config.Log.LogKeepHours == 0 {
		config.Log.LogKeepHours = 24
	}

	// 插件设置
	if config.Plugin.Enabled {
		if config.Plugin.Dir == "" {
			return errors.New("plugin dir is empty")
		}

		if config.Plugin.LogDir == "" {
			return errors.New("plugin log dir is empty")
		}

		if !utils.ValidGitUrl(config.Plugin.Git) {
			return errors.New("plugin git repo must be used the web url")
		}
	}

	// 心跳设置
	if config.Heartbeat.Enabled {
		if config.Heartbeat.Addr == "" {
			return errors.New("heartbeat addr is empty")
		}

		if config.Heartbeat.Interval == 0 {
			config.Heartbeat.Interval = 60
		}

		if config.Heartbeat.Timeout == 0 {
			config.Heartbeat.Timeout = 1000
		}
	}

	// transfer 设置
	if config.Transfer.Enabled {
		if len(config.Transfer.Addrs) == 0 {
			return errors.New("transfer addrs is empty")
		}

		if config.Transfer.Interval == 0 {
			config.Transfer.Interval = 60
		}

		if config.Transfer.Timeout == 0 {
			config.Transfer.Timeout = 1000
		}
	}

	// http 设置
	if config.Http.Enabled {
		if config.Http.Listen == "" {
			return errors.New("local listen addr is empty")
		}
	}

	// MaxCPURate
	if config.MaxCPURate < 0 || config.MaxCPURate > 1 {
		return errors.New("max_cpu_rate is range 0 to 1")
	}

	// MaxMemRate
	if config.MaxMemRate < 0 || config.MaxMemRate > 1 {
		return errors.New("max_mem_rate is range 0 to 1")
	}

	return nil
}
