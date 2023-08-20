package executor

import (
	"errors"
	"fmt"
	execstreamer "github.com/midy177/exec-streamer"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/syncmap"
	"os"
	"os/exec"
	"os/user"
)

type ExecReq struct {
	// ID
	Id uint64 `json:"id"`
	// Name
	Name string `json:"name" validate:"required"`
	// Content
	Content string `json:"content" validate:"required"`
	WithEnd bool   `json:"withEnd"`
}

const logFilePrefix = "/var/log/exec_schedule/"

var runningMap syncmap.Map

func init() {
	_ = os.MkdirAll(logFilePrefix, 0644)
}

func (e *ExecReq) Run() (string, error) {
	logFilePath := fmt.Sprintf("%s%s", logFilePrefix, e.Name)
	writer, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", err
	}
	go func(w *os.File, content string, path string) {
		defer w.Close()
		home := "/root"
		usr, err := user.Current()
		if err == nil {
			home = usr.HomeDir
		}
		streamer, err := execstreamer.NewExecStreamerBuilder().
			ExecutorName("bash").
			Exe(content).
			Env(append(os.Environ(), "TERM=xterm-256color", "HOME="+home)...).
			Dir(home).
			Writers(writer).
			AutoFlush().
			Build()
		if err != nil {
			logrus.Error(err)
			return
		}
		_, _ = w.Write([]byte("-------The script start running-------\n"))
		exec, err := streamer.StartExec()
		if err != nil {
			logrus.Error(err)
			return
		}
		runningMap.Store(path, exec)
		defer runningMap.Delete(path)
		_ = exec.Wait()
		_, _ = w.Write([]byte("\n-------The script finish running-------\n"))
	}(writer, e.Content, logFilePath)

	return logFilePath, nil
}

// GetOnRunningExec 获取在运行的脚本列表
func GetOnRunningExec() []string {
	var list []string
	runningMap.Range(func(key, value any) bool {
		name, ok := key.(string)
		if ok {
			list = append(list, name)
		}
		return true
	})
	return list
}

// IsRunningExec 获取在运行的脚本列表
func IsRunningExec(name string) bool {
	_, ok := runningMap.Load(name)
	return ok
}

// StopRunningExec 关闭在运行的脚本，返回成功失败状态
func StopRunningExec(name string) (bool, error) {
	v, ok := runningMap.Load(name)
	if !ok {
		return false, errors.New("脚本没有在运行")
	}
	val, vOk := v.(*exec.Cmd)
	if !vOk {
		return false, errors.New("类型不是*exec.Cmd")
	}
	err := val.Process.Kill()
	if err != nil {
		return false, err
	}
	return true, nil
}
