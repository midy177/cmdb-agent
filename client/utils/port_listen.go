package utils

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
)

var (
	listenIsRun   bool
	listenCancel  context.CancelFunc
	listenMux     sync.Mutex
	tcpListenList sync.Map
	udpListenList sync.Map
)

func ListenPort(b []byte) error {
	if listenIsRun {
		return fmt.Errorf("已经运行，请先关闭")
	}
	listenMux.Lock()
	listenIsRun = true
	listenMux.Unlock()
	ctx, cancel := context.WithCancel(context.Background())
	scan := new(ScanInstance)
	ports, err := scan.getAllPort(string(b))
	if err != nil {
		log.Fatalln(err)
	}
	g := sync.WaitGroup{}
	for _, v := range ports {
		g.Add(2)
		go func(p int, c context.Context) {
			runTcp(p, c)
			g.Done()
		}(v, ctx)
		go func(p int, c context.Context) {
			runUdp(p, c)
			g.Done()
		}(v, ctx)
	}
	g.Wait()
	cancel()
	return nil
}

func StopListen() error {
	listenMux.Lock()
	defer listenMux.Unlock()
	if listenIsRun {
		listenCancel()
	} else {
		return fmt.Errorf("没有在监听端口")
	}
	return nil
}

func GetListenList() *OnListenList {
	listenMux.Lock()
	defer listenMux.Unlock()
	if !listenIsRun {
		return nil
	}
	list := new(OnListenList)
	tcpListenList.Range(func(key, value any) bool {
		list.TCP = append(list.TCP, key.(int))
		return true
	})
	udpListenList.Range(func(key, value any) bool {
		list.UDP = append(list.UDP, key.(int))
		return true
	})
	return list
}

func runTcp(p int, ctx context.Context) {
	tcpListenList.Store(p, struct{}{})
	defer tcpListenList.Delete(p)
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:"+strconv.Itoa(p))
	if err != nil {
		fmt.Println("Error resolving TCP address:", err)
		return
	}
	conn, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Println("Error listening on TCP address:", err)
		return
	}
	defer conn.Close()
	log.Println("listening on TCP address:" + addr.String())
	for {
		select {
		case <-ctx.Done():
			return
		default:
			cli, cerr := conn.Accept()
			if cerr != nil {
				fmt.Println("Error reading TCP packet: ", cerr)
				continue
			}
			_, _ = cli.Write([]byte("Hello: " + cli.RemoteAddr().String() + ",Accessible tcp port: " + strconv.Itoa(p) + "\n"))
			_ = cli.Close()
			fmt.Printf("TCP访问可达-> 远端地址：%s 本地端口: %d\n", cli.RemoteAddr().String(), p)
		}
	}
}

func runUdp(p int, ctx context.Context) {
	udpListenList.Store(p, struct{}{})
	defer udpListenList.Delete(p)
	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:"+strconv.Itoa(p))
	if err != nil {
		fmt.Printf("Error resolving UDP address: %s", err)
		return
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Printf("Error listening on UDP address: %s", err)
		return
	}
	defer conn.Close()
	log.Println("listening on UDP address:" + addr.String())
	buffer := make([]byte, 1024)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, cli, cerr := conn.ReadFromUDP(buffer)
			if cerr != nil {
				fmt.Printf("Error reading UDP packet: %s", cerr)
				continue
			}
			_, _ = conn.WriteToUDP([]byte("Hello: "+cli.String()+",Accessible udp port: "+strconv.Itoa(p)+"\n"), cli)
			fmt.Printf("UDP访问可达-> 远端地址：%s 本地端口: %d\n", cli.String(), p)
		}
	}
}

type OnListenList struct {
	TCP []int `json:"tcp"`
	UDP []int `json:"udp"`
}
