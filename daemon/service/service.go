// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package service

import (
	"os"
	"runtime"

	"github.com/kardianos/service"
)

// Config configures the service.
type Config struct {
	Name     string // service name
	Desc     string // service description
	Username string // service username (windows only)
	Password string // service password (windows only)
	Config   string
}

// New creates and configures a new service.
func New(conf Config) (service.Service, error) {
	config := &service.Config{
		Name:        conf.Name,
		DisplayName: conf.Name,
		Description: conf.Desc,
		Arguments:   []string{"service", "run"},
	}

	switch runtime.GOOS {
	case "darwin":
		config.Option = service.KeyValue{
			"KeepAlive":   true,
			"RunAtLoad":   true,
			"UserService": os.Getuid() != 0,
		}
	case "windows":
		if conf.Username != "" {
			config.UserName = conf.Username
			config.Option = service.KeyValue{
				"Password": conf.Password,
			}
		}
	}

	m := new(manager)
	return service.New(m, config)
}
