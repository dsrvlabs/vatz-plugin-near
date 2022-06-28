package main

import (
	"flag"
	"fmt"
	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/sdk"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/structpb"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	// Default values.
	defaultAddr = "127.0.0.1"
	defaultPort = 9091
	pluginName  = "near-metric-up"
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	cmd := "curl -s localhost:3030/metrics"

	c, b := exec.CommandContext(ctx, "bash", "-c", cmd), new(strings.Builder)
	c.Stdout = b
	c.Run()

	cancel()
	contentMSG := ""
	if len(b.String()) > 0 {
		contentMSG = "NEAR Process is UP"
		log.Info().
			Str("NearUP", contentMSG).
			Msg("near-metric-up")
	} else {
		contentMSG = "NEAR Process is DOWN"
		severity = pluginpb.SEVERITY_CRITICAL
		log.Error().
			Str("NearUP", contentMSG).
			Msg("near-metric-up")
	}

	ret := sdk.CallResponse{
		FuncName:   "NearUP",
		Message:    contentMSG,
		Severity:   severity,
		State:      state,
		AlertTypes: []pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD},
	}

	return ret, nil
}
