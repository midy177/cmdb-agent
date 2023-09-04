package remotedialerx

import (
	"context"
	"errors"
	"github.com/rancher/remotedialer"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

type RemoteDialerX struct {
	Ctx     context.Context
	Session *remotedialer.Session
	mux     sync.RWMutex
}

func NewRemoteDialerStruct() *RemoteDialerX {
	return new(RemoteDialerX)
}

func (r *RemoteDialerX) NewRemoteDialerX(ctx context.Context, serverUrl string, headers http.Header) error {
	return remotedialer.ClientConnect(ctx, serverUrl+"/cmdb-api/cmdb_dialer/connect", headers, nil, func(proto, address string) bool {
		logrus.Infof("remotedialer: %s %s", proto, address)
		return true
	}, r.onConnect)
}

func (r *RemoteDialerX) onConnect(ctx context.Context, session *remotedialer.Session) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.Ctx = ctx
	r.Session = session
	return nil
}

func (r *RemoteDialerX) Close() {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.Ctx = nil
	r.Session = nil
}

func (r *RemoteDialerX) GetRemoteDialer() (remotedialer.Dialer, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()
	if r.Session != nil {
		return r.Session.Dial, nil
	}
	return nil, errors.New("remote dialer is close")
}

func (r *RemoteDialerX) HttpClient() (*http.Client, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()
	if r.Session != nil {
		return &http.Client{
			Transport: &http.Transport{
				DialContext: r.Session.Dial,
				//DialContext: func(c context.Context, network, addr string) (net.Conn, error) {
				//	return r.Session.Dial(ctx, "unix", "/var/run/cmdb-server.sock")
				//},
			},
		}, nil
	}
	return nil, errors.New("remote dialer is close")

}
