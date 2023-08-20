package pkg

import (
	jsoniter "github.com/json-iterator/go"
	"os"
)

var hostInfo HostInfo

func GetRemoteDialerConnectConfig() (*HostInfo, error) {
	onfigBuffer, err := os.ReadFile("/etc/cmdb-agent/config.json")
	if err != nil {
		return &hostInfo, err
	}
	err = jsoniter.Unmarshal(onfigBuffer, &hostInfo)
	if err != nil {
		return &hostInfo, err
	}
	return &hostInfo, err
}
