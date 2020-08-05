package metrics

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/ops-agent/common/cmd"
	"github.com/fanghongbo/ops-agent/common/g"
	"github.com/fanghongbo/ops-agent/common/model"
)

var timeout = 30

func DuMetrics() (L []*model.MetricValue) {
	var (
		paths     []string
		result    chan *model.MetricValue
		resultLen int
		wg        sync.WaitGroup
	)

	paths = g.DuPathMeta()
	result = make(chan *model.MetricValue, len(paths))

	for _, path := range paths {
		wg.Add(1)
		go func(path string) {
			var err error
			defer func() {
				if err != nil {
					dlog.Errorf("get disk status err: %s", err)
					result <- GaugeValue(g.DuBs, -1, "path="+path)
				}
				wg.Done()
			}()
			//tips:osx  does not support -b.
			command := exec.Command("du", "-bs", path)
			var stdout bytes.Buffer
			command.Stdout = &stdout
			var stderr bytes.Buffer
			command.Stderr = &stderr
			err = command.Start()
			if err != nil {
				return
			}
			err, isTimeout := cmd.RunLocalCommandWithTimeout(command, time.Duration(timeout)*time.Second)
			if isTimeout {
				err = errors.New(fmt.Sprintf("exec cmd : du -bs %s timeout", path))
				return
			}

			errStr := stderr.String()
			if errStr != "" {
				err = errors.New(errStr)
				return
			}

			if err != nil {
				err = errors.New(fmt.Sprintf("du -bs %s failed: %s", path, err.Error()))
				return
			}

			arr := strings.Fields(stdout.String())
			if len(arr) < 2 {
				err = errors.New(fmt.Sprintf("du -bs %s failed: %s", path, "return fields < 2"))
				return
			}

			size, err := strconv.ParseUint(arr[0], 10, 64)
			if err != nil {
				err = errors.New(fmt.Sprintf("cannot parse du -bs %s output", path))
				return
			}
			result <- GaugeValue(g.DuBs, size, "path="+path)
		}(path)
	}
	wg.Wait()

	resultLen = len(result)
	for i := 0; i < resultLen; i++ {
		L = append(L, <-result)
	}
	return
}
