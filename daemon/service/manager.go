// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package service

import (
	"cmdb-agent/daemon"
	"context"
	"github.com/kardianos/service"
)

var noContext = context.Background()

// a manager manages the service lifecycle.
type manager struct {
	cancel context.CancelFunc
}

// Start starts the service in a separate go routine.
func (m *manager) Start(service.Service) error {
	ctx, cancel := context.WithCancel(noContext)
	m.cancel = cancel
	go func() {
		_ = daemon.Run(ctx)
	}()
	return nil
}

// Stop stops the service.
func (m *manager) Stop(service.Service) error {
	m.cancel()
	return nil
}
