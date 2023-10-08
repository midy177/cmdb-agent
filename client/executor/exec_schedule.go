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

const execLogFilePrefix = "/var/log/cmdb-agent/exec_schedule/"
const cronLogFilePrefix = "/var/log/cmdb-agent/cron_schedule/"

var runningMap syncmap.Map

func init() {
	_ = os.MkdirAll(execLogFilePrefix, 0644)
	_ = os.MkdirAll(cronLogFilePrefix, 0644)
}

func (e *ExecReq) Run() (string, error) {
	var logFilePath string
	if e.IsCron {
		logFilePath = fmt.Sprintf("%s%d", cronLogFilePrefix, e.Uuid)
	} else {
		logFilePath = fmt.Sprintf("%s%d", execLogFilePrefix, e.Uuid)
	}
	// 保存脚本内容
	e.saveContent(logFilePath + "/content.sh")
	// 保存脚本运行日志
	writer, err := os.OpenFile(logFilePath+"/stat.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", err
	}

	go func(w *os.File, content string, keyName uint64) {
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
		runningMap.Store(keyName, exec)
		defer runningMap.Delete(keyName)
		_ = exec.Wait()
		_, _ = w.Write([]byte("\n-------The script finish running-------\n"))
	}(writer, e.Content, e.GroupId)

	return logFilePath, nil
}

func (e *ExecReq) saveContent(dstDir string) {
	writer, err := os.OpenFile(dstDir, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	defer writer.Close()
	if err != nil {
		logrus.Errorf("faild save content:err-> %s", err.Error())
		return
	}
	_, err = writer.WriteString(e.Content)
	if err != nil {
		logrus.Errorf("faild save content:err-> %s", err.Error())
	}
}

// GetAllOnRunningExec 获取在运行的脚本列表
func GetAllOnRunningExec() []uint64 {
	var list []uint64
	runningMap.Range(func(key, value any) bool {
		name, ok := key.(uint64)
		if ok {
			list = append(list, name)
		}
		return true
	})
	return list
}

// IsRunningExec 获取在运行的脚本列表
func IsRunningExec(keyName uint64) bool {
	_, ok := runningMap.Load(keyName)
	return ok
}

// StopRunningExec 关闭在运行的脚本，返回成功失败状态
func StopRunningExec(keyName uint64) (bool, error) {
	v, ok := runningMap.Load(keyName)
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
