package handler

import (
	"cmdb-agent/client/echox"
	"context"
	selfupdate "github.com/creativeprojects/go-selfupdate"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"sync"
	"time"
)

var (
	upgraded = false
	mux      sync.Mutex
)

func (cc *CustomContext) Upgrading(c echo.Context) error {
	if upgraded {
		return echox.Response{Code: http.StatusOK, Data: "Have upgraded,wait to restart"}.JSON(c)
	} else {
		mux.Lock()
		defer mux.Unlock()
		upgraded = true
		// 获取请求
		req := new(UpgradeReq)
		err := c.Bind(&req)
		if err != nil {
			return echox.Response{Code: http.StatusBadRequest, Message: err.Error()}.JSON(c)
		}
		err = c.Validate(req)
		if err != nil {
			return echox.Response{Code: http.StatusBadRequest, Message: err.Error()}.JSON(c)
		}
		streamWriter := &StreamWriter{W: c.Response()}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextPlain)
		c.Response().WriteHeader(http.StatusOK)
		_, err = streamWriter.Write([]byte("Start updating...\n"))
		if err != nil {
			return err
		}
		assetURL := req.UpgradeUrl
		assetName := "cmdb-agent"
		exe, err := os.Executable()
		if err != nil {
			_, err = streamWriter.Write([]byte("err: could not locate executable path.\n"))
			return err
		}
		if err := selfupdate.UpdateTo(context.Background(), assetURL, assetName, exe); err != nil {
			_, err = streamWriter.Write([]byte("error occurred while updating binary: " + err.Error() + "\n"))
			return err
		}
		_, err = streamWriter.Write([]byte("Successfully updated,wait restart program.\n"))
		if err != nil {
			return err
		}
		_, err = streamWriter.Write([]byte("Start to restart program....\n"))
		if err != nil {
			return err
		}
		time.Sleep(time.Second * 2)
		// 在新进程中重启动自己
		cmd := exec.Command("/usr/local/bin/cmdb-agent", "service", "restart")
		usr, err := user.Current()
		if err != nil {
			logrus.Error("Failed to get current user:", err)
			cmd.Dir = "/root"
		} else {
			cmd.Dir = usr.HomeDir
		}
		cmd.Env = append(os.Environ(), "TERM="+getTerm(), "HOME="+cmd.Dir)
		if err := cmd.Run(); err != nil {
			logrus.Errorf("Error restarting: %s", err)
		}
		return err
	}
}
