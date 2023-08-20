package handler

import (
	"cmdb-agent/client/echox"
	"cmdb-agent/client/executor"
	"encoding/base64"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"strings"
)

func (cc *CustomContext) ExecOnBackend(c echo.Context) error {
	req := new(executor.ExecReq)
	err := c.Bind(&req)
	if err != nil {
		return echox.Response{Code: http.StatusOK, Message: err.Error()}.JSON(c)
	}
	err = c.Validate(req)
	if err != nil {
		return echox.Response{Code: http.StatusOK, Message: err.Error()}.JSON(c)
	}
	runLogPath, err := req.Run()
	if err != nil {
		return echox.Response{Code: http.StatusOK, Message: err.Error()}.JSON(c)
	}
	return echox.Response{Code: http.StatusOK, Data: runLogPath}.JSON(c)
}

func (cc *CustomContext) ExecOnBackendIsRunning(c echo.Context) error {
	req := new(StopOnRunningExec)
	err := c.Bind(&req)
	if err != nil {
		return echox.Response{Code: http.StatusOK, Message: err.Error()}.JSON(c)
	}
	err = c.Validate(req)
	if err != nil {
		return echox.Response{Code: http.StatusOK, Message: err.Error()}.JSON(c)
	}
	ok := executor.IsRunningExec(req.Name)
	return echox.Response{Code: http.StatusOK, Data: ok}.JSON(c)
}

func (cc *CustomContext) GetExecOnBackendList(c echo.Context) error {
	list := executor.GetOnRunningExec()
	return echox.Response{Code: http.StatusOK, Data: list}.JSON(c)
}

func (cc *CustomContext) StopExecOnBackendList(c echo.Context) error {
	req := new(StopOnRunningExec)
	err := c.Bind(&req)
	if err != nil {
		return echox.Response{Code: http.StatusOK, Message: err.Error()}.JSON(c)
	}
	err = c.Validate(req)
	if err != nil {
		return echox.Response{Code: http.StatusOK, Message: err.Error()}.JSON(c)
	}
	ok, err := executor.StopRunningExec(req.Name)
	if err != nil {
		return echox.Response{Code: http.StatusOK, Message: err.Error()}.JSON(c)
	}
	return echox.Response{Code: http.StatusOK, Data: ok}.JSON(c)
}

func (cc *CustomContext) RemoveExecLog(c echo.Context) error {
	enScript := c.Param("id")
	decodeScript, err := base64.StdEncoding.DecodeString(enScript)
	if err != nil {
		return echox.Response{Code: http.StatusOK, Message: err.Error()}.JSON(c)
	}
	if strings.HasPrefix(string(decodeScript), "/var/log/exec_schedule/") {
		err = os.RemoveAll(string(decodeScript))
		if err != nil {
			return echox.Response{Code: http.StatusOK, Message: err.Error()}.JSON(c)
		}
	}
	return echox.Response{Code: http.StatusOK, Data: true}.JSON(c)
}
