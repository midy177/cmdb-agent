package handler

import (
	"cmdb-agent/client/echox"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func (cc *CustomContext) UploadFile(c echo.Context) error {
	path := c.FormValue("path")
	filename := c.FormValue("filename")
	if !filepath.IsAbs(path) {
		return echox.Response{Code: http.StatusOK, Message: "目标文件夹必须是绝对路径"}.JSON(c)
	}
	_, err := os.Stat(path)
	if err != nil {
		return echox.Response{Code: http.StatusOK, Message: err.Error()}.JSON(c)
	}
	file, err := c.FormFile("file")
	if err != nil {
		return echox.Response{Code: http.StatusOK, Message: err.Error()}.JSON(c)
	}
	src, err := file.Open()
	if err != nil {
		return echox.Response{Code: http.StatusOK, Message: err.Error()}.JSON(c)
	}
	defer src.Close()

	fullFilePath := filepath.Join(path, filename)
	dst, err := os.Create(fullFilePath)
	if err != nil {
		return echox.Response{Code: http.StatusOK, Message: err.Error()}.JSON(c)
	}
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		return echox.Response{Code: http.StatusOK, Message: err.Error()}.JSON(c)
	}
	return echox.Response{Code: http.StatusOK, Data: "上传成功"}.JSON(c)
}
