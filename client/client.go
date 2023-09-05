package client

import (
	"cmdb-agent/client/handler"
	"cmdb-agent/client/midd"
	"cmdb-agent/client/pkg"
	"cmdb-agent/client/remotedialerx"
	"cmdb-agent/client/utils"
	"context"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rancher/remotedialer"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func NewClient(ctx context.Context) error {
	cusCtx := new(handler.CustomContext)
	cusCtx.RemoteDialerX = new(remotedialerx.RemoteDialerX)
	cfg, err := pkg.GetRemoteDialerConnectConfig()
	if err != nil {
		logrus.Fatalln(err)
	}
	utils.SetHostId(cfg.Id)
	remotedialer.ClientHandler.Use(middleware.Recover())
	remotedialer.ClientHandler.Validator = midd.NewValidate()
	remotedialer.ClientHandler.GET("/web_terminal", cusCtx.WebsocketTerminal)
	remotedialer.ClientHandler.POST("/upgrading", cusCtx.Upgrading)
	remotedialer.ClientHandler.POST("/exec_backend/run", cusCtx.ExecOnBackend)
	remotedialer.ClientHandler.GET("/exec_backend/list", cusCtx.GetExecOnBackendList)
	remotedialer.ClientHandler.POST("/exec_backend/status", cusCtx.ExecOnBackendIsRunning)
	remotedialer.ClientHandler.POST("/exec_backend/stop", cusCtx.StopExecOnBackendList)
	remotedialer.ClientHandler.GET("/exec_terminal/:id", cusCtx.WebsocketTerminalOnScript)
	remotedialer.ClientHandler.GET("/remove_log/:id", cusCtx.RemoveExecLog)
	remotedialer.ClientHandler.POST("/save_config", cusCtx.SaveConfig)
	remotedialer.ClientHandler.POST("/exec_online", cusCtx.ExecOnline)
	remotedialer.ClientHandler.POST("/tail_file", cusCtx.TailLogStream)
	remotedialer.ClientHandler.POST("/upload_file", cusCtx.UploadFile)
	remotedialer.ClientHandler.POST("/download_file", cusCtx.DownloadFile)
	remotedialer.ClientHandler.GET("/get_stat", cusCtx.GetStat)
	go cusCtx.ReportCircle(ctx)
	headers := http.Header{
		"tunnel-token": []string{cfg.Auth},
		"client-key":   []string{strconv.FormatUint(cfg.Id, 10)},
	}
	for {
		cusCtx.CircleRemoteDialerX(ctx, cfg.Server, headers)
	}
}
