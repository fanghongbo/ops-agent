package g

import (
	"github.com/fanghongbo/dlog"
	"time"
)

func InitLog() {
	var (
		backend *dlog.FileBackend
		err     error
	)

	if config.Log == nil {
		return
	}

	backend, err = dlog.NewFileBackend(config.Log.LogPath, config.Log.LogFileName)
	if err != nil {
		dlog.Fatalf("create log file backend err: %s", err)
	}

	dlog.SetLogging(config.Log.LogLevel, backend)

	// 日志rotate设置
	backend.SetKeepHours(uint(config.Log.LogKeepHours))
	backend.SetFlushDuration(1 * time.Second)
	backend.SetRotateByHour(true)

	if config.Debug {
		dlog.LogToStdout()
	}
}
