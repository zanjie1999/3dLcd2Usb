package main

import (
	"fmt"

	"github.com/shirou/gopsutil/net"
)

func main() {
	netInfos, _ := net.IOCounters(true)
	fmt.Println(netInfos)
}
