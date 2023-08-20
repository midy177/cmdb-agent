package handler

import (
	"cmdb-agent/client/echox"
	"github.com/labstack/echo/v4"
	"net/http"
	"path/filepath"
)

func (cc *CustomContext) DownloadFile(c echo.Context) error {
	req := new(DownloadFileReq)
	err := c.Bind(&req)
	if err != nil {
		return echox.Response{Code: http.StatusOK, Message: err.Error()}.JSON(c)
	}
	err = c.Validate(req)
	if err != nil {
		return echox.Response{Code: http.StatusOK, Message: err.Error()}.JSON(c)
	}
	c.Response().Header().Set("Content-Type", "application/octet-stream")
	return c.Attachment(req.Filepath, filepath.Base(req.Filepath))
}
