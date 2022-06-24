package main

import (
	"fmt"
	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/sdk"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/mem"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/structpb"
	"os"
	"time"
)

const (
	addr       = "0.0.0.0"
	port       = 9093
	pluginName = "machine-status-memory"
)

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

	vmStat, err := mem.VirtualMemory()

	if err != nil {
		state = pluginpb.STATE_FAILURE
		severity = pluginpb.SEVERITY_ERROR
	}

	totalUsage := vmStat.UsedPercent

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

	if state == pluginpb.STATE_SUCCESS {
		if memoryScale > 3 {
			severity = pluginpb.SEVERITY_CRITICAL
		} else if memoryScale > 2 {
			severity = pluginpb.SEVERITY_WARNING
		}
	}

	contentMSG := "Total Memory Usage: " + fmt.Sprintf("%.2f", totalUsage) + "%"

	log.Info().
		Str("GetMemoryUsage", contentMSG).
		Msg("machine-status-memory")

	ret := sdk.CallResponse{
		FuncName:   "GetMemoryUsage",
		Message:    contentMSG,
		Severity:   severity,
		State:      state,
		AlertTypes: []pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD},
	}

	return ret, nil
}
