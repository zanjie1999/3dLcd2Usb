package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

var lastSent uint64 = 0
var lastRecv uint64 = 0
var cpuBar = "=================="

func main() {
	for {
		cpuPre := CpuPercent()
		out := "["
		out += fmt.Sprintf("%-18s", cpuBar[:int(cpuPre*0.18)])
		out += "]\n"
		out += fmt.Sprintf("%-12s", fmt.Sprintf("U:%.3f%%", CpuPercent()))
		out += TimeHMS()
		out += "\n"
		send, recv := NetworkSpeed(0.5)
		out += fmt.Sprintf("%-10s", fmt.Sprintf("Men:%.2fG", VMemUsed()))
		out += fmt.Sprintf("%10s", fmt.Sprintf("^:%.3f", send))
		out += "\n"
		out += fmt.Sprintf("%-10s", fmt.Sprintf("Swap:%.1fG", SMemUsed()))
		out += fmt.Sprintf("%10s", fmt.Sprintf("v:%.3f", recv))
		out += "\n"
		fmt.Println(out)
		// time.Sleep(time.Second)
	}
}

func CpuPercent() float64 {
	percent, _ := cpu.Percent(time.Second/2, false)
	return percent[0]
}

func VMemUsed() float64 {
	memInfo, _ := mem.VirtualMemory()
	return float64(memInfo.Used) / 1024 / 1024 / 1024
}

func SMemUsed() float64 {
	memInfo, _ := mem.SwapMemory()
	return float64(memInfo.Used) / 1024 / 1024 / 1024
}

func DiskPercent() float64 {
	parts, _ := disk.Partitions(true)
	diskInfo, _ := disk.Usage(parts[0].Mountpoint)
	return diskInfo.UsedPercent
}

func NetworkSpeed(s float64) (float64, float64) {
	netInfos, _ := net.IOCounters(true)
	var sent uint64 = 0
	var recv uint64 = 0
	if len(netInfos) > 1 {
		for _, i := range netInfos {
			if strings.HasPrefix(i.Name, "en") {
				sent = i.BytesSent
				recv = i.BytesRecv
				break
			}
		}
	}
	if sent == 0 && recv == 0 {
		sent = netInfos[0].BytesSent
		recv = netInfos[0].BytesRecv

	}
	sendF64 := float64(sent-lastSent) / 1024 / 1024 / s
	recvF64 := float64(recv-lastRecv) / 1024 / 1024 / s
	lastSent = sent
	lastRecv = recv
	return sendF64, recvF64
}

func TimeHMS() string {
	return time.Now().Format("03:04:05")
}
