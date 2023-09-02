package handler

import (
	"bytes"
	"cmdb-agent/client/echox"
	"cmdb-agent/client/utils"
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func (cc *CustomContext) GetStat(c echo.Context) error {
	hostStat := utils.GetHostStat()
	return echox.Response{Code: http.StatusOK, Data: hostStat}.JSON(c)
}

func (cc *CustomContext) ReportCircle(ctx context.Context) {
	for {
		cc.reportFun(ctx)
	}
}

func (cc *CustomContext) reportFun(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Error("Recovered reportFun:", r)
		}
	}()
	time.Sleep(time.Minute * 4)
	client, err := cc.RemoteDialerX.HttpClient(ctx)
	if err != nil {
		logrus.Errorf("can't get http.client-> %s", err.Error())
		return
	}
	hostStat := utils.GetHostStat()
	marshalData, err := jsoniter.Marshal(hostStat)
	if err != nil {
		logrus.Errorf("can't Marshal hostStat data -> %s", err.Error())

	}
	resp, err := client.Post("http://unix/report_host_stat", "application/json; charset=utf-8", bytes.NewBuffer(marshalData))
	if err != nil {
		logrus.Errorf("can't Post hostStat data to server -> %s", err.Error())
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		result, err := echox.ParseResponse(resp)
		if err != nil {
			logrus.Errorf("can't ParseResponse -> %s", err.Error())
			return
		}
		logrus.Errorf("Post err info -> %s", result.Message)
	}
}
