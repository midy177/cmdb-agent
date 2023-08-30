package utils

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/net"
	"log"
	"testing"
)

func TestName(t *testing.T) {
	connections, err := net.Connections("udp")
	if err != nil {
		log.Fatal(err)
	}

	count := 0
	for _, conn := range connections {
		if conn.Status == "ESTABLISHED" {
			count++
		}
	}
	fmt.Println(len(connections), count)
}
