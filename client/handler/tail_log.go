package handler

import (
	"cmdb-agent/client/echox"
	"github.com/hpcloud/tail"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

func (cc *CustomContext) TailLogStream(c echo.Context) error {
	req := new(TailLogReq)
	err := c.Bind(&req)
	if err != nil {
		return echox.Response{Code: http.StatusOK, Message: err.Error()}.JSON(c)
	}
	err = c.Validate(req)
	if err != nil {
		return echox.Response{Code: http.StatusOK, Message: err.Error()}.JSON(c)
	}
	var tailConfig = tail.Config{
		// 文件夹被移除或被或被打包，需要重新打开
		ReOpen: false,
		//实时跟踪
		Follow: false,
		//支持文件不存在
		MustExist: false,
		Poll:      true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 1},
	}
	if req.Follow {
		tailConfig.ReOpen = true
		tailConfig.Follow = true
	}
	if req.SeekEnd {
		// 如果出现移除，保存上次读取位置，避免重新读取
		tailConfig.Location = &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd}
	}
	tailFile, err := tail.TailFile(req.LogPath, tailConfig)
	if err != nil {
		return echox.Response{Code: http.StatusOK, Message: err.Error()}.JSON(c)
	}
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
	ctx := c.Request().Context()
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			msg, ok := <-tailFile.Lines
			if !ok {
				logrus.Infof("tail file close open, filename: %s", tailFile.Filename)
				time.Sleep(100 * time.Millisecond)
				if req.Follow {
					continue
				}
				return nil
			}
			// 传输日志
			_, err2 := c.Response().Write([]byte(msg.Text))
			if err2 != nil {
				return err2
			}
			c.Response().Flush()
		}
	}
}
