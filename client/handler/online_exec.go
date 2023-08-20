package handler

import (
	"cmdb-agent/client/echox"
	"cmdb-agent/client/executor"
	"context"
	"encoding/base64"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	execstreamer "github.com/midy177/exec-streamer"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/user"
)

func (cc *CustomContext) ExecOnline(c echo.Context) error {
	req := new(executor.ExecReq)
	err := c.Bind(&req)
	if err != nil {
		return echox.Response{Code: http.StatusOK, Message: err.Error()}.JSON(c)
	}
	err = c.Validate(req)
	if err != nil {
		return echox.Response{Code: http.StatusOK, Message: err.Error()}.JSON(c)
	}
	streamWriter := &StreamWriter{W: c.Response()}
	home := "/root"
	usr, err := user.Current()
	if err == nil {
		home = usr.HomeDir
	}
	streamer, err := execstreamer.NewExecStreamerBuilder().
		ExecutorName("bash").
		Exe(req.Content).
		Dir(home).
		Env(append(os.Environ(), "TERM=xterm-256color", "HOME="+home)...).
		Writers(streamWriter).
		AutoFlush().
		Build()
	if err != nil {
		return echox.Response{Code: http.StatusOK, Message: err.Error()}.JSON(c)
	}
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextPlain)
	c.Response().WriteHeader(http.StatusOK)
	exec, err := streamer.StartExec()
	if err != nil {
		_, _ = c.Response().Write([]byte(err.Error()))
		return nil
	}
	ctx, cancel := context.WithCancel(c.Request().Context())
	go func() {
		_ = exec.Wait()
		select {
		case <-ctx.Done():
			return
		default:
			if req.WithEnd {
				_, _ = streamWriter.Write([]byte("-------The script finish running-------\n"))
			}
			cancel()
		}
	}()
	<-ctx.Done()
	_ = exec.Process.Kill()
	return nil
}

func (cc *CustomContext) WebsocketTerminalOnScript(c echo.Context) error {
	enScript := c.Param("id")
	decodeScript, err := base64.StdEncoding.DecodeString(enScript)
	if err != nil {
		return echox.Response{Code: http.StatusOK, Message: err.Error()}.JSON(c)
	}
	l := logrus.WithField("remoteAddr", c.Request().RemoteAddr)
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		l.WithError(err).Error("Unable to upgrade connection")
		return err
	}
	home := "/root"
	usr, err := user.Current()
	if err == nil {
		home = usr.HomeDir
	}
	streamer, err := execstreamer.NewExecStreamerBuilder().
		ExecutorName("bash").
		Exe(string(decodeScript)).
		Env(append(os.Environ(), "TERM=xterm-256color", "HOME="+home)...).
		Dir(home).
		Writers(&execWriter{W: conn}).
		AutoFlush().
		Build()
	exec, err := streamer.StartExec()
	if err != nil {
		_, _ = c.Response().Write([]byte(err.Error()))
		return nil
	}
	ctx, cancel := context.WithCancel(c.Request().Context())
	go func() {
		_ = exec.Wait()
		select {
		case <-ctx.Done():
			return
		default:
			cancel()
		}
	}()
	<-ctx.Done()
	_ = exec.Process.Kill()
	return nil
}

type execWriter struct {
	W *websocket.Conn
}

func (s *execWriter) Write(p []byte) (n int, err error) {
	err = s.W.WriteMessage(websocket.BinaryMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), err
}
