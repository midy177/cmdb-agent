package utils

import (
	"context"
	"github.com/creativeprojects/go-selfupdate"
	"os"
)

func UpgradeMyself(upgradeUrl string) error {
	assetName := "cmdb-agent"
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	return selfupdate.UpdateTo(context.Background(), upgradeUrl, assetName, exe)
}
