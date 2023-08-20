package utils

import (
	"bytes"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"net"
	"os/exec"
	"strings"
)

func GetPublicAddr() (*PublicAddr, error) {
	client := resty.New()
	resp, err := client.R().SetHeader("User-Agent", "Mozilla").Get("https://api.ip.sb/geoip")
	if err != nil {
		return nil, err
	}
	info := new(PublicAddr)
	err = json.Unmarshal(resp.Body(), &info)
	if err != nil {
		return nil, err
	}
	return info, nil
}

type PublicAddr struct {
	Organization    string  `json:"organization"`
	Longitude       float64 `json:"longitude"`
	City            string  `json:"city"`
	Timezone        string  `json:"timezone"`
	Isp             string  `json:"isp"`
	Offset          int     `json:"offset"`
	Region          string  `json:"region"`
	Asn             int     `json:"asn"`
	AsnOrganization string  `json:"asn_organization"`
	Country         string  `json:"country"`
	Ip              string  `json:"ip"`
	Latitude        float64 `json:"latitude"`
	ContinentCode   string  `json:"continent_code"`
	CountryCode     string  `json:"country_code"`
	RegionCode      string  `json:"region_code"`
}

func GetPrivateAddr() []string {
	ifaces, err := net.Interfaces()
	if err != nil {
		logrus.Errorf("net.Interfaces err: %s", err.Error())
		return nil
	}
	var list []string
	for _, iFace := range ifaces {
		// 过滤掉容器网卡
		if strings.HasPrefix(iFace.Name, "docker") ||
			strings.HasPrefix(iFace.Name, "br-") ||
			strings.HasPrefix(iFace.Name, "cni") ||
			strings.HasPrefix(iFace.Name, "cilium") ||
			strings.HasPrefix(iFace.Name, "flannel") ||
			strings.HasPrefix(iFace.Name, "tun") ||
			strings.HasPrefix(iFace.Name, "weave") ||
			strings.HasPrefix(iFace.Name, "cali") ||
			(iFace.Flags&net.FlagUp == 0) ||
			(iFace.Flags&net.FlagLoopback != 0) ||
			(iFace.Flags&net.FlagPointToPoint != 0) {
			continue
		}
		var (
			addrArr []net.Addr
			ips     []string
		)
		addrArr, err = iFace.Addrs()
		if err != nil {
			logrus.Errorf("returns a list of unicast interface addresses for a specific interface. err: %s", err.Error())
			err = nil
			continue
		}
		for _, addr := range addrArr {
			ipAddr, ok := addr.(*net.IPNet)
			if ok && !ipAddr.IP.IsLoopback() && ipAddr.IP.To4() != nil {
				ips = append(ips, addr.String())
			}
			//ip, ipNet, err := net.ParseCIDR(addr.String())
			//if err != nil {
			//	fmt.Printf("Failed to parse address for interface %s: %v\n", iFace.Name, err)
			//	continue
			//}
			//fmt.Printf("Interface: %s, IP: %s, Subnet Mask: %s\n", iface.Name, ip, ipNet.Mask)
		}
		list = append(list, ips...)
	}
	return list
}

func GetNotFormatDisks() []string {
	cmd := exec.Command("bash", "-c", "lsblk -r --output NAME,MOUNTPOINT | awk -F \\/ '/sd/ { dsk=substr($1,1,3);dsks[dsk]+=1 } END { for ( i in dsks ) { if (dsks[i]==1) print i } }'\n")
	output, err := cmd.Output()
	if err != nil {
		logrus.Errorf("Failed to run lsblk command: %v\n", err)
		return nil
	}
	bytes.TrimLeft(output, "\n")
	bytes.TrimRight(output, "\n")
	bytes.TrimSpace(output)
	blks := bytes.Split(output, []byte("\n"))
	var list []string
	for _, v := range blks {
		if len(v) > 0 {
			list = append(list, "/dev/"+string(v))
		}
	}
	return list
}
