// 3D LCD to USB 上位机
// 2021/01/19
// v2.0

package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/albenik/go-serial/v2"
	"github.com/albenik/go-serial/v2/enumerator"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	//"golang.org/x/exp/shiny/driver/internal/win32"
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
			port, err := serial.Open(portName, serial.WithBaudrate(9600))
			if err != nil {
				log.Print(err)
			} else {
				log.Print(portName)
				for {
					out := Screen1()
					// fmt.Println(out)
					n, err := port.Write([]byte(out))
					if err != nil {
						log.Print(err)
						break
					}
					buf := make([]byte, 128)
					n, err = port.Read(buf)
					if err != nil {
						log.Print(err)
						break
					}
					if n != 0 {
						log.Printf("%q", buf[0])
						Button(string(buf[0]))
					}
				}
			}
		}
		time.Sleep(time.Second * 10)
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
		// fmt.Println(memInfo.Used, lastVRAMUse, memInfo.Used-lastVRAMUse)
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
			if runtime.GOOS == "windows" {
				if strings.HasPrefix(i.Name, "本地连接") {
					sent = i.BytesSent
					recv = i.BytesRecv
					break
				}
			} else if strings.HasPrefix(i.Name, "en") {
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

// h:mm:ss
func TimeHMS() string {
	return time.Now().Format("3:04:05")
}

// 寻找串口并初始化
func findSerialPort() string {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		log.Print(err)
	}
	if len(ports) == 0 {
		log.Print("No serial ports found!")
	}
	for _, port := range ports {
		if port.IsUSB && port.VID == serialVid && port.PID == serialPid {
			return port.Name
		}
	}
	return ""
}

// 第1屏
func Screen1() string {
	cpuPre := CPUPercent()
	out := "["
	out += fmt.Sprintf("%-18s", cpuBar[:int(cpuPre*0.19)])
	out += "]\n"
	if cpuPre >= 100 {
		out += fmt.Sprintf("%-12s", fmt.Sprintf("Cpu:%.2f%%", cpuPre))
	} else if cpuPre >= 10 {
		out += fmt.Sprintf("%-12s", fmt.Sprintf("Cpu:%.3f%%", cpuPre))
	} else {
		out += fmt.Sprintf("%-12s", fmt.Sprintf("Cpu: %.3f%%", cpuPre))
	}
	out += fmt.Sprintf("%8s", TimeHMS())
	out += "\n"
	send, recv := NetworkSpeed(0.5)
	out += fmt.Sprintf("%-11s", fmt.Sprintf("Men:%.3fG", VMemUsed()))
	if send >= 10 {
		out += fmt.Sprintf("%9s", fmt.Sprintf("^:%.3f", send))
	} else {
		out += fmt.Sprintf("%9s", fmt.Sprintf("^: %.3f", send))
	}
	out += "\n"
	out += fmt.Sprintf("%-11s", fmt.Sprintf("Swp:%.3fG", SMemUsed()))
	if recv >= 10 {
		out += fmt.Sprintf("%9s", fmt.Sprintf("v:%.3f", recv))
	} else {
		out += fmt.Sprintf("%9s", fmt.Sprintf("v: %.3f", recv))
	}
	out += "\r"
	return out
}

func Button(in string) {
	switch in {
	case "1":
		if runtime.GOOS == "windows" {
			// 关闭屏幕
			//win32.SendMessage(-1, 0x0112, 0xF170, 2)
		}
	case "0":

	case "+":

	case "-":

	default:
	}
}
