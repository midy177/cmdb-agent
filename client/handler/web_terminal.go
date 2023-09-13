package handler

import (
	"encoding/json"
	"github.com/creack/pty"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"time"
)

type windowSize struct {
	Resize bool   `json:"resize"`
	Rows   uint16 `json:"rows"`
	Cols   uint16 `json:"cols"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源的请求
	},
}

func (cc *CustomContext) WebsocketTerminal(c echo.Context) error {
	l := logrus.WithField("remoteAddr", c.Request().RemoteAddr)
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		l.WithError(err).Error("Unable to upgrade connection")
		return err
	}
	cmd := exec.Command("/bin/bash", "-l")
	usr, err := user.Current()
	if err != nil {
		l.Error("Failed to get current user:", err)
		cmd.Dir = "/root"
	} else {
		cmd.Dir = usr.HomeDir
	}
	cmd.Env = append(os.Environ(), "TERM="+getTerm(), "HOME="+cmd.Dir)
	tty, err := pty.Start(cmd)
	if err != nil {
		l.WithError(err).Error("Unable to start pty/cmd")
		_ = conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		return err
	}
	defer func() {
		tty.Write([]byte("exit\n"))
		time.Sleep(time.Second)
		cmd.Process.Kill()
		cmd.Process.Wait()
		tty.Close()
		conn.Close()
	}()
	pty.Setsize(tty, &pty.Winsize{
		Rows: 30,
		Cols: 80,
	})
	go func() {
		for {
			buf := make([]byte, 1024)
			read, err := tty.Read(buf)
			if err != nil {
				_ = conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
				// 处理退出
				_ = conn.Close()
				return
			}
			_ = conn.WriteMessage(websocket.BinaryMessage, buf[:read])
		}
	}()

	for {
		messageType, reader, err := conn.NextReader()
		if err != nil {
			//l.WithError(err).Error("Unable to grab next reader")
			return nil
		}

		if messageType == websocket.TextMessage {
			decoder := json.NewDecoder(reader)
			resizeMessage := windowSize{}
			err = decoder.Decode(&resizeMessage)
			if err == nil && resizeMessage.Resize {
				err = pty.Setsize(tty, &pty.Winsize{
					Rows: resizeMessage.Rows,
					Cols: resizeMessage.Cols,
				})
				if err != nil {
					l.WithError(err).Error("Unable to resize terminal")
				}
			} else {
				_ = conn.WriteMessage(websocket.PongMessage, []byte("00"))
			}
		} else {
			copied, err := io.Copy(tty, reader)
			if err != nil {
				l.WithError(err).Errorf("Error after copying %d bytes", copied)
			}
		}
	}
}

func getTerm() (term string) {
	if term = os.Getenv("xterm"); term == "" {
		term = "xterm-256color"
	}
	return
}
