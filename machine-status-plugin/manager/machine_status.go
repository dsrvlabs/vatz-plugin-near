package manager

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"syscall"
	"time"
)

type DiskStatus struct {
	All  uint64 `json:"All"`
	Used uint64 `json:"Used"`
	Free uint64 `json:"Free"`
}

type machineStatus struct {
}

func DiskInfo(path string) (disk DiskStatus) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	return
}

func (ms machineStatus) GetDiskInfo(path string) (DiskStatus, error) {
	diskInfo := DiskInfo(path)
	return diskInfo, nil
}

func (ms machineStatus) GetMemoryUsage() (float64, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println(err)
		//os.Exit(-1)
	}
	return vmStat.UsedPercent, nil
}

func (ms machineStatus) GetCPUUsage() (float64, error) {
	_, err := cpu.Info()
	if err != nil {
		fmt.Printf("get cpu info failed, err:%v", err)
	}

	totalUsed := 0.0
	percent, _ := cpu.Percent(3*time.Second, false)
	//fmt.Printf("cpu percent:%v\n", percent)
	for _, numb := range percent {
		totalUsed += numb
	}

	return totalUsed, nil
}

type MachineStatus interface {
	GetDiskInfo(path string) (DiskStatus, error)
	GetMemoryUsage() (float64, error)
	GetCPUUsage() (float64, error)
}

func NewMachineStatus() MachineStatus {
	return &machineStatus{}
}
