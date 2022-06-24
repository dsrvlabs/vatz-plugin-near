package main

import (
	"fmt"
	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/sdk"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/cpu"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/structpb"
	"os"
	"time"
)

const (
	addr       = "0.0.0.0"
	port       = 9091
	pluginName = "machine-status-cpu"
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

	_, err := cpu.Info()

	if err != nil {
		fmt.Printf("get cpu info failed, err:%v", err)
		state = pluginpb.STATE_FAILURE
		severity = pluginpb.SEVERITY_ERROR
	}

	totalUsed := 0.0
	percent, _ := cpu.Percent(3*time.Second, false)

	for _, numb := range percent {
		totalUsed += numb
	}

	cpuScale := 0
	if totalUsed < 50.0 {
		cpuScale = 1
	} else if totalUsed < 65 {
		cpuScale = 2
	} else if totalUsed < 90 {
		cpuScale = 3
	} else if totalUsed >= 90 {
		cpuScale = 4
	}

	if state == pluginpb.STATE_SUCCESS {
		if cpuScale > 3 {
			severity = pluginpb.SEVERITY_CRITICAL
		} else if cpuScale > 2 {
			severity = pluginpb.SEVERITY_WARNING
		}
	}
	contentMSG := "Total CPU Usage: " + fmt.Sprintf("%.2f", totalUsed) + "%"

	log.Info().
		Str("GetCPUUsage", contentMSG).
		Msg("machine-status-cpu")

	ret := sdk.CallResponse{
		FuncName:   "GetCPUUsage",
		Message:    contentMSG,
		Severity:   severity,
		State:      state,
		AlertTypes: []pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD},
	}

	return ret, nil
}
