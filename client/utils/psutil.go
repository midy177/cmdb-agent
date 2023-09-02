package utils

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"strings"
	"time"
)

type HostStat struct {
	CpuStat struct {
		Core        int     `json:"core"`
		UsedPercent float64 `json:"usedPercent"` // 每30 s
	} `json:"cpuStat"`
	MemoryStat struct {
		Total       uint64  `json:"total"` // MB
		UsedPercent float64 `json:"usedPercent"`
	} `json:"memoryStat"`
	DiskStat []PartitionStat `json:"diskStat"`
	Uptime   uint64          `json:"uptime"`   // 正常运行时间
	BootTime uint64          `json:"bootTime"` // 开机时间
	HostId   uint64          `json:"hostId"`
}

type PartitionStat struct {
	Device      string  `json:"device"`
	MountPoint  string  `json:"mountPoint"`
	FsType      string  `json:"fsType"`
	Total       uint64  `json:"total"` // GB
	UsedPercent float64 `json:"usedPercent"`
}

var hostStat HostStat

func SetHostId(hostId uint64) {
	hostStat.HostId = hostId
}

func GetHostStat() *HostStat {
	// Cpu
	cInfo, err := cpu.Info()
	if err == nil {
		hostStat.CpuStat.Core = len(cInfo)
	}
	c, err := cpu.Percent(time.Minute, false)

	if err == nil && len(c) > 0 {
		hostStat.CpuStat.UsedPercent = c[0]
	}

	// Memory
	v, err := mem.VirtualMemory()
	if err == nil {
		hostStat.MemoryStat.Total = v.Total / 1024 / 1024
		hostStat.MemoryStat.UsedPercent = v.UsedPercent
	}

	// Disk
	partitions, err := disk.Partitions(false)
	if err == nil {
		var diskStat []PartitionStat
		for _, partition := range partitions {
			if isPhysicalDisk(partition.Device) {
				usageStat, err := disk.Usage(partition.Mountpoint)
				if err == nil {
					diskStat = append(diskStat, PartitionStat{
						Device:      partition.Device,
						MountPoint:  partition.Mountpoint,
						FsType:      partition.Fstype,
						Total:       usageStat.Total / 1024 / 1024 / 1024, // GB
						UsedPercent: usageStat.UsedPercent,
					})
				}
			}

		}
		hostStat.DiskStat = diskStat
	}
	// Uptime
	hInfo, err := host.Info()
	if err == nil {
		hostStat.Uptime = hInfo.Uptime
		hostStat.BootTime = hInfo.BootTime
	}

	return &hostStat
}

func isPhysicalDisk(device string) bool {
	if strings.HasPrefix(device, "/dev/") && !strings.Contains(device, "loop") {
		return true
	}
	return false
}
