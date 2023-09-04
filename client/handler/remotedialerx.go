package handler

import (
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func (cc *CustomContext) CircleRemoteDialerX(ctx context.Context, url string, headers http.Header) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Error("Recovered CircleRemoteDialerX:", r)
		}
	}()
	err := cc.RemoteDialerX.NewRemoteDialerX(ctx, url, headers)
	if err != nil {
		logrus.Error(err)
	}
	cc.RemoteDialerX.Close()
	time.Sleep(10 * time.Second)
}
