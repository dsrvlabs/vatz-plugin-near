package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"
	"time"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/sdk"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	// Default values.
	defaultAddr = "127.0.0.1"
	defaultPort = 9091
	pluginName  = "machine-status-disk"
)

var (
	addr string
	port int
)

func init() {
	flag.StringVar(&addr, "addr", defaultAddr, "IP Address(e.g. 0.0.0.0, 127.0.0.1)")
	flag.IntVar(&port, "port", defaultPort, "Port number, defulat 9091")
	flag.Parse()
}

type DiskStatus struct {
	All  uint64 `json:"All"`
	Used uint64 `json:"Used"`
	Free uint64 `json:"Free"`
}

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

func getLargeInteger(arr []int) int {
	for j := 1; j < len(arr); j++ {
		if arr[0] < arr[j] {
			arr[0] = arr[j]
		}
	}
	return arr[0]
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

func main() {
	p := sdk.NewPlugin(pluginName)
	p.Register(pluginFeature)

	ctx := context.Background()
	if err := p.Start(ctx, addr, port); err != nil {
		fmt.Println("exit")
	}
}

func pluginFeature(info, option map[string]*structpb.Value) (sdk.CallResponse, error) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	state := pluginpb.STATE_SUCCESS
	severity := pluginpb.SEVERITY_INFO

	var flags []int
	diskPaths := [2]string{"/", "/mnt/near"}

	message := ""

	for idx, diskPath := range diskPaths {
		DiskInfo := DiskInfo(diskPath)
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

	log.Info().
		Str("GetDiskUsage", message).
		Msg("machine-status-disk")

	if message == "" {
		state = pluginpb.STATE_FAILURE
	}

	if state == pluginpb.STATE_SUCCESS {
		if severityScale > 3 {
			severity = pluginpb.SEVERITY_CRITICAL
		} else if severityScale > 2 {
			severity = pluginpb.SEVERITY_WARNING
		}
	}

	// TODO: Fill here.
	ret := sdk.CallResponse{
		FuncName:   "GetDiskUsage",
		Message:    message,
		Severity:   severity,
		State:      state,
		AlertTypes: []pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD},
	}

	return ret, nil
}
