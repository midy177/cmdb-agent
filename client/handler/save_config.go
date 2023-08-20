package handler

import (
	"cmdb-agent/client/echox"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"path/filepath"
)

type ConfigFile struct {
	FilePath string `json:"filePath" validate:"required"`
	Filename string `json:"filename" validate:"required"`
	Content  string `json:"content" validate:"required"`
}

func (cc *CustomContext) SaveConfig(c echo.Context) error {
	req := new(ConfigFile)
	err := c.Bind(&req)
	if err != nil {
		return echox.Response{Code: http.StatusOK, Message: err.Error()}.JSON(c)
	}
	err = c.Validate(req)
	if err != nil {
		return echox.Response{Code: http.StatusOK, Message: err.Error()}.JSON(c)
	}
	err = os.MkdirAll(req.FilePath, 0666)
	if err != nil {
		return echox.Response{Code: http.StatusOK, Message: err.Error()}.JSON(c)
	}
	fullpath := filepath.Join(req.FilePath, req.Filename)
	err = os.WriteFile(fullpath, []byte(req.Content), 0666)
	if err != nil {
		return echox.Response{Code: http.StatusOK, Message: err.Error()}.JSON(c)
	}
	return echox.Response{Code: http.StatusOK, Data: "上传成功"}.JSON(c)
}
