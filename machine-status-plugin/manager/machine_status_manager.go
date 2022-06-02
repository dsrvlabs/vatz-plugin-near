package manager

import (
	"fmt"
	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
)

var (
	MSInstance MachineStatus
	MSManager  machineStatusManager
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

func init() {
	MSInstance = NewMachineStatus()
}

type machineStatusManager struct {
}

func getLargeInteger(arr []int) int {
	for j := 1; j < len(arr); j++ {
		if arr[0] < arr[j] {
			arr[0] = arr[j]
		}
	}
	return arr[0]
}

func (ms *machineStatusManager) GetDiskUsage() (*pluginpb.ExecuteResponse, error) {
	var flags []int
	//diskPaths := [2]string{"/", "/mnt/near"}
	diskPaths := [1]string{"/"}
	message := ""

	for idx, diskPath := range diskPaths {
		DiskInfo, _ := MSInstance.GetDiskInfo(diskPath)
		all := float64(DiskInfo.All) / float64(GB)
		used := float64(DiskInfo.Used) / float64(GB)
		defaultUsage := used / all * 100

		currentUsage := fmt.Sprintf("%.2f", defaultUsage)

		flagStatus := 0
		if defaultUsage < 60 {
			flagStatus = 1
		} else if defaultUsage < 85 {
			flagStatus = 2
		} else if defaultUsage < 95 {
			flagStatus = 3
		} else if defaultUsage >= 95 {
			flagStatus = 4
		}

		flags = append(flags, flagStatus)
		if idx > 0 && message != "" {
			message += "| "
		}
		usedSTR := fmt.Sprintf("%f", used)
		allSTR := fmt.Sprintf("%f", all)
		message += "Mounted on " + diskPath + " with Usage: " + currentUsage + "%" + " (" + usedSTR + "/" + allSTR + ")"
	}

	severityScale := getLargeInteger(flags)
	state := pluginpb.STATE_FAILURE
	severity := pluginpb.SEVERITY_INFO

	fmt.Println("message: ", message)

	if message != "" {
		state = pluginpb.STATE_SUCCESS
	}

	if severityScale > 3 {
		severity = pluginpb.SEVERITY_CRITICAL
	} else if severityScale > 2 {
		severity = pluginpb.SEVERITY_WARNING
	}

	return &pluginpb.ExecuteResponse{
		State:        state,
		Message:      message,
		AlertType:    []pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD},
		Severity:     severity,
		ResourceType: "near-machine-status-plugin",
	}, nil
}

func (ms *machineStatusManager) GetCPUUsage() (*pluginpb.ExecuteResponse, error) {
	totalUsage, _ := MSInstance.GetCPUUsage()
	cpuScale := 0
	if totalUsage < 50.0 {
		cpuScale = 1
	} else if totalUsage < 65 {
		cpuScale = 2
	} else if totalUsage < 90 {
		cpuScale = 3
	} else if totalUsage >= 90 {
		cpuScale = 4
	}

	state := pluginpb.STATE_FAILURE
	severity := pluginpb.SEVERITY_INFO

	if totalUsage > 0.0 {
		state = pluginpb.STATE_SUCCESS
	}

	if cpuScale > 3 {
		severity = pluginpb.SEVERITY_CRITICAL
	} else if cpuScale > 2 {
		severity = pluginpb.SEVERITY_WARNING
	}

	fmt.Println("Total CPU Usage: ", totalUsage, "%")
	s := "Total CPU Usage: " + fmt.Sprintf("%.2f", totalUsage) + "%"

	return &pluginpb.ExecuteResponse{
		State:        state,
		Message:      s,
		AlertType:    []pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD},
		Severity:     severity,
		ResourceType: "near-machine-status-plugin",
	}, nil
}

func (ms *machineStatusManager) GetMemoryUsage() (*pluginpb.ExecuteResponse, error) {
	totalUsage, _ := MSInstance.GetMemoryUsage()
	memoryScale := 0
	if totalUsage < 60 {
		memoryScale = 1
	} else if totalUsage < 80 {
		memoryScale = 2
	} else if totalUsage < 90 {
		memoryScale = 3
	} else if totalUsage >= 90 {
		memoryScale = 4
	}

	state := pluginpb.STATE_FAILURE
	severity := pluginpb.SEVERITY_INFO

	if totalUsage > 0.0 {
		state = pluginpb.STATE_SUCCESS
	}

	if memoryScale > 3 {
		severity = pluginpb.SEVERITY_CRITICAL
	} else if memoryScale > 2 {
		severity = pluginpb.SEVERITY_WARNING
	}

	fmt.Println("Total Memory Usage: ", totalUsage, "%")
	s := "Total Memory Usage: " + fmt.Sprintf("%.2f", totalUsage) + "%"

	return &pluginpb.ExecuteResponse{
		State:        state,
		Message:      s,
		AlertType:    []pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD},
		Severity:     severity,
		ResourceType: "near-machine-status-plugin",
	}, nil
}
