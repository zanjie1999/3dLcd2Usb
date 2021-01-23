// 3D LCD to USB 上位机
// 2021/01/19
// v1

package main

import (
	"runtime"
	"strings"
	"time"
	"log"
	"fmt"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/albenik/go-serial/v2"
	"github.com/albenik/go-serial/v2/enumerator"
)

var lastSent uint64 = 0
var lastRecv uint64 = 0
var lastVRAMUse uint64 = 0
var cpuBar = "=================="
var useSerial = true
var serialVid = "04D9"
var serialPid = "B534"

func main() {
	for {
		portName := findSerialPort()
		if portName != "" {
			port, err := serial.Open(portName, serial.WithBaudrate(115200))
			if err != nil {
				log.Fatal(err)
			} else {
				log.Println(portName)
				for {
					out := screen1()
					// fmt.Println(out)
					n, err := port.Write([]byte(out))
					if err != nil {
						log.Fatal(err)
					}
					buf := make([]byte, 128)
					n, err = port.Read(buf)
					if err != nil {
						log.Fatal(err)
					}
					if n!=0 {
						log.Printf("%q", buf[0])
					}
				}
			}
		}
	}
}

// cpu使用率
func CPUPercent() float64 {
	percent, _ := cpu.Percent(time.Second/2, false)
	return percent[0]
}

// 物理内存使用
func VMemUsed() float64 {
	memInfo, _ := mem.VirtualMemory()
	lastVRAMUse = memInfo.Used
	return float64(memInfo.Used) / 1024 / 1024 / 1024
}

// 虚拟内存使用
func SMemUsed() float64 {
	memInfo, _ := mem.SwapMemory()
	// win拿到的是已提交 需要减掉物理内存使用
	if memInfo.Used != 0 && runtime.GOOS == "windows" {
		if lastVRAMUse == 0 {
			VMemUsed()
		}
		return float64(memInfo.Used-lastVRAMUse) / 1024 / 1024 / 1024
	}
	return float64(memInfo.Used) / 1024 / 1024 / 1024
}

// 磁盘占用率
func DiskPercent() float64 {
	parts, _ := disk.Partitions(true)
	diskInfo, _ := disk.Usage(parts[0].Mountpoint)
	return diskInfo.UsedPercent
}

// 网速				倍数		上传	下载
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

// hh:mm:ss
func TimeHMS() string {
	return time.Now().Format("03:04:05")
}

// 寻找串口并初始化
func findSerialPort() string {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}
	for _, port := range ports {
		if port.IsUSB && port.VID == serialVid && port.PID == serialPid {
			return port.Name
		}
	} 
	return ""
}

// 第1屏
func screen1() string {
	cpuPre := CPUPercent()
	out := "["
	out += fmt.Sprintf("%-18s", cpuBar[:int(cpuPre*0.18)])
	out += "]\n"
	if cpuPre >= 100 {
		out += fmt.Sprintf("%-12s", fmt.Sprintf("Cpu:%.2f%%", cpuPre))
	} else {
		out += fmt.Sprintf("%-12s", fmt.Sprintf("Cpu:%.3f%%", cpuPre))
	}
	out += TimeHMS()
	out += "\n"
	send, recv := NetworkSpeed(0.5)
	out += fmt.Sprintf("%-11s", fmt.Sprintf("Men:%.3fG", VMemUsed()))
	out += fmt.Sprintf("%9s", fmt.Sprintf("^:%.3f", send))
	out += "\n"
	out += fmt.Sprintf("%-11s", fmt.Sprintf("Swp:%.3fG", SMemUsed()))
	out += fmt.Sprintf("%9s", fmt.Sprintf("v:%.3f", recv))
	out += "\r"
	return out
}
