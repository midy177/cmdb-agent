package handler

import (
	"cmdb-agent/client/echox"
	"cmdb-agent/client/executor"
	"github.com/creack/pty"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os/exec"
	"os/user"
)

func (cc *CustomContext) OnlineExecPty(c echo.Context) error {
	req := new(executor.ExecReq)
	err := c.Bind(&req)
	if err != nil {
		return echox.Response{Code: http.StatusOK, Message: err.Error()}.JSON(c)
	}
	err = c.Validate(req)
	if err != nil {
		return echox.Response{Code: http.StatusOK, Message: err.Error()}.JSON(c)
	}
	cmd := exec.Command("bash", "-c", req.Content)
	usr, err := user.Current()
	if err != nil {
		cmd.Dir = "/root"
	} else {
		cmd.Dir = usr.HomeDir
	}
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return echox.Response{Code: http.StatusOK, Message: err.Error()}.JSON(c)
	}
	defer func() {
		cmd.Process.Kill()
		cmd.Process.Wait()
		_ = ptmx.Close()
	}()
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
	streamWriter := &StreamWriter{W: c.Response()}
	_, err = io.Copy(streamWriter, ptmx)
	return err
}
