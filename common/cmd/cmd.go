package cmd

import (
	"bytes"
	"github.com/fanghongbo/dlog"
	"os/exec"
	"syscall"
	"time"
)

func RunLocalCommand(name string, arg ...string) (string, error) {
	var (
		cmd *exec.Cmd
		out bytes.Buffer
		err error
	)

	cmd = exec.Command(name, arg...)
	cmd.Stdout = &out
	err = cmd.Run()
	return out.String(), err
}

func RunLocalCommandWithTimeout(cmd *exec.Cmd, timeout time.Duration) (error, bool) {
	var (
		err  error
		done chan error
	)

	done = make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(timeout):
		dlog.Infof("timeout, process:%s will be killed", cmd.Path)

		go func() {
			<-done // allow goroutine to exit
		}()

		// IMPORTANT: cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} is necessary before cmd.Start()
		err = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		if err != nil {
			dlog.Errorf("kill %s failed, error:", -cmd.Process.Pid, err)
		}

		return err, true
	case err = <-done:
		return err, false
	}
}
