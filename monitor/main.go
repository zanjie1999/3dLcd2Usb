package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

func main() {
	for {
		// out := "U:" + strconv.FormatFloat(GetCpuPercent(), 'f', 2, 64) + "% "
		// out += "M:" + strconv.FormatFloat(GetVMemUsed(), 'f', 2, 64) + "/" + strconv.FormatFloat(GetSMemUsed(), 'f', 2, 64) + "G"
		// fmt.Println(out)
		out := fmt.Sprintf("%-10s", fmt.Sprintf("U:%.2f%%", GetCpuPercent()))
		out += fmt.Sprintf("%10s", fmt.Sprintf("M:%.1f/%.1f", GetVMemUsed(), GetSMemUsed()))
		fmt.Println(out)
	}
}

func GetCpuPercent() float64 {
	percent, _ := cpu.Percent(time.Second, false)
	return percent[0]
}

func GetVMemUsed() float64 {
	memInfo, _ := mem.VirtualMemory()
	return float64(memInfo.Used) / 1024 / 1024 / 1024
}

func GetSMemUsed() float64 {
	memInfo, _ := mem.SwapMemory()
	return float64(memInfo.Used) / 1024 / 1024 / 1024
}

func GetDiskPercent() float64 {
	parts, _ := disk.Partitions(true)
	diskInfo, _ := disk.Usage(parts[0].Mountpoint)
	return diskInfo.UsedPercent
}
