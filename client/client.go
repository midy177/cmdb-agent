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

//func Run(ctx context.Context) error {
//	customContext := new(handler.CustomContext)
//	customContext.RemoteDialerX = new(remotedialerx.RemoteDialerX)
//	cfg, err := pkg.GetRemoteDialerConnectConfig()
//	if err != nil {
//		logrus.Fatalln(err)
//	}
//	// 创建 Unix socket 监听器
//	socket := "/var/run/cmdb-agent.sock"
//	_ = os.Remove(socket) // 检查并删除现有的 socket 文件
//	listener, err := net.Listen("unix", socket)
//	if err != nil {
//		logrus.Fatalln(err)
//	}
//	defer listener.Close()
//	e := echo.New()
//	e.Use(middleware.Recover())
//	e.Validator = midd.NewValidate()
//	e.GET("/web_terminal", customContext.WebsocketTerminal)
//	e.POST("/upgrading", customContext.Upgrading)
//	e.POST("/exec_backend/run", customContext.ExecOnBackend)
//	e.GET("/exec_backend/list", customContext.GetExecOnBackendList)
//	e.POST("/exec_backend/status", customContext.ExecOnBackendIsRunning)
//	e.POST("/exec_backend/stop", customContext.StopExecOnBackendList)
//	e.GET("/exec_terminal/:id", customContext.WebsocketTerminalOnScript)
//	e.GET("/remove_log/:id", customContext.RemoveExecLog)
//	e.POST("/save_config", customContext.SaveConfig)
//	e.POST("/exec_online", customContext.ExecOnline)
//	e.POST("/tail_file", customContext.TailLogStream)
//	e.POST("/upload_file", customContext.UploadFile)
//	e.POST("/download_file", customContext.DownloadFile)
//	e.GET("/get_stat", customContext.GetStat)
//	// 启动 HTTP 服务器
//	e.Listener = listener // 将 Echo 的 Listener 设置为 Unix socket 监听器
//	go func() {
//		logrus.Fatal(e.Start(""))
//	}()
//	utils.SetHostId(cfg.Id)
//	// 启动上报服务
//	go customContext.ReportCircle(ctx)
//	headers := http.Header{
//		"tunnel-token": []string{cfg.Auth},
//		"client-key":   []string{strconv.FormatUint(cfg.Id, 10)},
//	}
//	for {
//		CircleRemoteDialerX(customContext, ctx, cfg.Server, headers)
//	}
//}
//
//func CircleRemoteDialerX(customContext *handler.CustomContext, ctx context.Context, url string, headers http.Header) {
//	defer func() {
//		if r := recover(); r != nil {
//			logrus.Error("Recovered CircleRemoteDialerX:", r)
//		}
//	}()
//	err := customContext.RemoteDialerX.NewRemoteDialerX(ctx, url, headers)
//	if err != nil {
//		logrus.Error(err)
//	}
//	customContext.RemoteDialerX.Close()
//	time.Sleep(10 * time.Second)
//}

func NewClient(ctx context.Context) error {
	cusCtx := new(handler.CustomContext)
	cusCtx.RemoteDialerX = new(remotedialerx.RemoteDialerX)
	cfg, err := pkg.GetRemoteDialerConnectConfig()
	if err != nil {
		logrus.Fatalln(err)
	}
	utils.SetHostId(cfg.Id)
	remotedialer.ClientRouter.Use(middleware.Recover())
	remotedialer.ClientRouter.Validator = midd.NewValidate()
	remotedialer.ClientRouter.GET("/web_terminal", cusCtx.WebsocketTerminal)
	remotedialer.ClientRouter.POST("/upgrading", cusCtx.Upgrading)
	remotedialer.ClientRouter.POST("/exec_backend/run", cusCtx.ExecOnBackend)
	remotedialer.ClientRouter.GET("/exec_backend/list", cusCtx.GetExecOnBackendList)
	remotedialer.ClientRouter.POST("/exec_backend/status", cusCtx.ExecOnBackendIsRunning)
	remotedialer.ClientRouter.POST("/exec_backend/stop", cusCtx.StopExecOnBackendList)
	remotedialer.ClientRouter.GET("/exec_terminal/:id", cusCtx.WebsocketTerminalOnScript)
	remotedialer.ClientRouter.GET("/remove_log/:id", cusCtx.RemoveExecLog)
	remotedialer.ClientRouter.POST("/save_config", cusCtx.SaveConfig)
	remotedialer.ClientRouter.POST("/exec_online", cusCtx.ExecOnline)
	remotedialer.ClientRouter.POST("/tail_file", cusCtx.TailLogStream)
	remotedialer.ClientRouter.POST("/upload_file", cusCtx.UploadFile)
	remotedialer.ClientRouter.POST("/download_file", cusCtx.DownloadFile)
	remotedialer.ClientRouter.GET("/get_stat", cusCtx.GetStat)
	go cusCtx.ReportCircle(ctx)
	headers := http.Header{
		"tunnel-token": []string{cfg.Auth},
		"client-key":   []string{strconv.FormatUint(cfg.Id, 10)},
	}
	for {
		cusCtx.CircleRemoteDialerX(ctx, cfg.Server, headers)
	}
}
