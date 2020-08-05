# ops-agent

基于open-falcon二次开发的linux监控客户端

## 特性

- 代码重构、移除无用的接口和配置信息;
- 重写 url metric 的 probeUrl方法, 使用 go http 库替代使用系统 curl 命令来探测 url;
- 重写 /plugin/update 接口, 使用 go-git 库替代使用系统的git命令来更新插件仓库;
- 新增从私有git仓库同步插件功能;
- 新增运行日志配置, 支持日志滚动;
- 新增 /metrics 接口, 支持查看当前监控的所有 metric;
- 新增 /metric/check 接口, 支持查看当前系统 metric 依赖环境;
- 新增 cpu核心绑定、内存阈值配置; 当 agent 内存达到阈值的50%时, 打印告警信息；当内存达到阈值的100%, 程序直接退出;


## 编译

it is a golang classic project

``` shell
cd $GOPATH/src/github.com/fanghongbo/ops-agent/
./control build
./control start
```

## 配置
Refer to `cfg.example.json`, modify the file name to `cfg.json` :

```config
{
  "debug": false,
  "hostname": "",
  "ip": "",
  "log": {
    "log_level": "INFO",
    "log_path": "./logs",
    "log_file_name": "run.log",
    "log_keep_hours": 3
  },
  "plugin": {
    "enabled": true,
    "dir": "./plugin",
    "git": "https://github.com/open-falcon/plugin.git",
    "username": "",
    "password": "",
    "logs": "./logs"
  },
  "heartbeat": {
    "enabled": true,
    "addr": "127.0.0.1:6030",
    "interval": 60,
    "timeout": 1000
  },
  "transfer": {
    "enabled": true,
    "addrs": [
      "127.0.0.1:8433"
    ],
    "interval": 60,
    "timeout": 1000
  },
  "http": {
    "enabled": true,
    "listen": ":1988"
  },
  "collector": {
    "ifacePrefix": [
      "eth",
      "em",
      "ens"
    ],
    "mountPoint": []
  },
  "default_tags": {
  },
  "ignore": {
    "cpu.busy": true,
    "df.bytes.free": true,
    "df.bytes.total": true,
    "df.bytes.used": true,
    "df.bytes.used.percent": true,
    "df.inodes.total": true,
    "df.inodes.free": true,
    "df.inodes.used": true,
    "df.inodes.used.percent": true,
    "mem.memtotal": true,
    "mem.memused": true,
    "mem.memused.percent": true,
    "mem.memfree": true,
    "mem.swaptotal": true,
    "mem.swapused": true,
    "mem.swapfree": true
  },
  "max_cpu_rate": 0.2,
  "max_mem_rate": 0.3
}
```

## License

This software is licensed under the Apache License. See the LICENSE file in the top distribution directory for the full license text.
