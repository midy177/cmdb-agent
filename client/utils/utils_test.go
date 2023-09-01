package utils

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/disk"
	"testing"
)

func TestName(t *testing.T) {
	// Disk
	partitions, err := disk.Partitions(false)
	if err == nil {
		for _, partition := range partitions {
			usageStat, err := disk.Usage(partition.Mountpoint)
			if err == nil && isPhysicalDisk(partition.Device) {
				fmt.Printf("%+v %+v\n", partition, usageStat)
			}
		}
	}
}
