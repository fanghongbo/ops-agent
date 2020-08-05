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

	backend, err = dlog.NewFileBackend(Conf().Log.LogPath, Conf().Log.LogFileName)
	if err != nil {
		dlog.Fatalf("create log file backend err: %s", err)
	}

	dlog.SetLogging(Conf().Log.LogLevel, backend)

	// 日志rotate设置
	backend.SetKeepHours(uint(Conf().Log.LogKeepHours))
	backend.SetFlushDuration(1 * time.Second)
	backend.SetRotateByHour(true)

	if Conf().Debug {
		dlog.LogToStdout()
	}
}
