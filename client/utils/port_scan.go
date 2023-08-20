package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ScanInfo struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

func ScanPort(b []byte) {
	var scanInfo ScanInfo
	err := json.Unmarshal(b, &scanInfo)
	if err != nil {
		return
	}
	scan := NewScanInstance(300, 200, true)
	scan.ScanIpOpenPort("tcp", scanInfo.IP, scanInfo.Port)
	scan.ScanIpOpenPort("udp", scanInfo.IP, scanInfo.Port)
}

// ScanInstance ip 扫描
type ScanInstance struct {
	debug   bool
	timeout int
	process int
	mux     *sync.RWMutex
}

func NewScanInstance(timeout int, process int, debug bool) *ScanInstance {
	return &ScanInstance{
		debug:   debug,
		timeout: timeout,
		process: process,
		mux:     new(sync.RWMutex),
	}
}

// ScanIpOpenPort 获取开放端口号
func (s *ScanInstance) ScanIpOpenPort(network, ip string, port string) {
	var (
		total     int
		pageCount int
		num       int
		//openPorts []int
		//mutex     sync.Mutex
	)
	ports, _ := s.getAllPort(port)
	total = len(ports)
	if total < s.process {
		pageCount = total
	} else {
		pageCount = s.process
	}
	InitScanLogFile(network)
	num = int(math.Ceil(float64(total) / float64(pageCount)))
	startLog := fmt.Sprintf("【%v】需要扫描%s端口总数:%v 个，总协程:%v 个，并发:%v 个，超时:%d 毫秒\n", ip, network, total, pageCount, num, s.timeout)
	s.debugLog(startLog)
	s.saveScanLog([]byte(startLog), network)
	start := time.Now()
	all := map[int][]int{}
	for i := 1; i <= pageCount; i++ {
		for j := 0; j < num; j++ {
			tmp := (i-1)*num + j
			if tmp < total {
				all[i] = append(all[i], ports[tmp])
			}
		}
	}
	wg := sync.WaitGroup{}
	for k, v := range all {
		wg.Add(1)
		go func(value []int, key int) {
			defer wg.Done()
			//var tmpPorts []int
			for i := 0; i < len(value); i++ {
				opened := s.isOpen(network, ip, value[i])
				if opened {
					//tmpPorts = append(tmpPorts, value[i])
					s.debugLog(fmt.Sprintf("【%v】%s端口:%v ...... 开放 ......", ip, network, value[i]))
				} else {
					s.debugLog(fmt.Sprintf("【%v】%s端口:%v ...... 未开放 ......", ip, network, value[i]))
					s.saveScanLog([]byte(fmt.Sprintf("【%v】%s端口:%v ...... 未开放 ......\n", ip, network, value[i])), network)
				}
			}
			//mutex.Lock()
			//openPorts = append(openPorts, tmpPorts...)
			//mutex.Unlock()
		}(v, k)
	}
	wg.Wait()
	s.saveScanLog([]byte(fmt.Sprintf("【%v】^_^扫描结束，执行时长%.3fs\n", ip, time.Since(start).Seconds())), network)
}

// GetAllIp 获取所有ip
func (s *ScanInstance) GetAllIp(ip string) ([]string, error) {
	var (
		ips []string
	)
	ipTmp := strings.Split(ip, "-")
	firstIp, err := net.ResolveIPAddr("ip", ipTmp[0])
	if err != nil {
		return ips, errors.New(ipTmp[0] + "域名解析失败" + err.Error())
	}
	if net.ParseIP(firstIp.String()) == nil {
		return ips, errors.New(ipTmp[0] + " ip地址有误~")
	}
	//域名转化成ip再塞回去
	ipTmp[0] = firstIp.String()
	ips = append(ips, ipTmp[0]) //最少有一个ip地址

	if len(ipTmp) == 2 {
		//以切割第一段ip取到最后一位
		ipTmp2 := strings.Split(ipTmp[0], ".")
		startIp, _ := strconv.Atoi(ipTmp2[3])
		endIp, err := strconv.Atoi(ipTmp[1])
		if err != nil || endIp < startIp {
			endIp = startIp
		}
		if endIp > 255 {
			endIp = 255
		}
		totalIp := endIp - startIp + 1
		for i := 1; i < totalIp; i++ {
			ips = append(ips, fmt.Sprintf("%s.%s.%s.%d", ipTmp2[0], ipTmp2[1], ipTmp2[2], startIp+i))
		}
	}
	return ips, nil
}

// 记录日志
func (s *ScanInstance) debugLog(str string) {
	if s.debug == true {
		fmt.Println(str)
	}
}

// 记录日志
func (s *ScanInstance) saveScanLog(content []byte, network string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	filename := "scan-" + network + ".log"
	// 检查文件是否存在，不存在则创建
	var _, err = os.Stat(filename)
	if os.IsNotExist(err) {
		var file, err = os.Create(filename)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
	}
	// 追加内容到文件末尾
	err = os.WriteFile(filename, content, os.ModeAppend)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func InitScanLogFile(network string) {
	filename := "scan-" + network + ".log"
	_ = os.WriteFile(filename, nil, 0666)
}

// 获取所有端口
func (s *ScanInstance) getAllPort(port string) ([]int, error) {
	var ports []int
	//处理 ","号 如 80,81,88 或 80,88-100
	portArr := strings.Split(strings.Trim(port, ","), ",")
	for _, v := range portArr {
		portArr2 := strings.Split(strings.Trim(v, "-"), "-")
		startPort, err := s.filterPort(portArr2[0])
		if err != nil {
			continue
		}
		//第一个端口先添加
		ports = append(ports, startPort)
		if len(portArr2) > 1 {
			//添加第一个后面的所有端口
			endPort, _ := s.filterPort(portArr2[1])
			if endPort > startPort {
				for i := 1; i <= endPort-startPort; i++ {
					ports = append(ports, startPort+i)
				}
			}
		}
	}
	//去重复
	ports = s.arrayUnique(ports)

	return ports, nil
}

// 端口合法性过滤
func (s *ScanInstance) filterPort(str string) (int, error) {
	port, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	if port < 1 || port > 65535 {
		return 0, errors.New("端口号范围超出")
	}
	return port, nil
}

// 查看端口号是否打开
func (s *ScanInstance) isOpen(network, ip string, port int) bool {
	conn, err := net.DialTimeout(network, fmt.Sprintf("%s:%d", ip, port), time.Millisecond*time.Duration(s.timeout))
	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			fmt.Println("连接数超出系统限制！" + err.Error())
			os.Exit(1)
		}
		return false
	}
	_ = conn.Close()
	return true
}

// 数组去重
func (s *ScanInstance) arrayUnique(arr []int) []int {
	var newArr []int
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return newArr
}
