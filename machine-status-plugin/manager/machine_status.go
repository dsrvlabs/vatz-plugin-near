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
	percent, _ := cpu.Percent(time.Second, true)
	totalUsed := 0.0
	totalUsed = percent[cpu.CPUser] + percent[cpu.CPNice] + percent[cpu.CPSys] + percent[cpu.CPIntr] + percent[cpu.CPIdle] + percent[cpu.CPUStates]
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
