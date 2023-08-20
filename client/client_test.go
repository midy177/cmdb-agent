package client

import (
	"cmdb-agent/client/handler"
	"cmdb-agent/client/midd"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"os"
	"os/user"
	"testing"
	//"cmdb-agent/client/svc"
)

func TestGetEnv(t *testing.T) {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Failed to get current user:", err)
		return
	}
	fmt.Println("Home directory: ", usr.HomeDir, "<->", os.Getenv("HOME"))

	envVars := os.Environ()

	for _, envVar := range envVars {
		fmt.Println(envVar)
	}
}

func TestName(t *testing.T) {
	customContext := new(handler.CustomContext)
	e := echo.New()
	e.Use(middleware.Recover())
	e.Validator = midd.NewValidate()
	e.GET("/web_terminal", customContext.WebsocketTerminal)
	e.POST("/exec_backend/run", customContext.ExecOnBackend)
	e.GET("/exec_backend/list", customContext.GetExecOnBackendList)
	e.POST("/exec_backend/status", customContext.ExecOnBackendIsRunning)
	e.POST("/exec_backend/stop", customContext.StopExecOnBackendList)
	e.POST("/exec_online", customContext.ExecOnline)
	e.POST("/exec_online_pty", customContext.OnlineExecPty)
	e.GET("/upgrading", customContext.Upgrading)
	//e.POST("/tail_file", handler.TailLogStream)
	e.POST("/upload_file", customContext.UploadFile)
	e.POST("/download_file", customContext.DownloadFile)
	// 启动 HTTP 服务器
	//e.Listener = listener // 将 Echo 的 Listener 设置为 Unix socket 监听器
	logrus.Fatal(e.Start("0.0.0.0:8080"))
}
