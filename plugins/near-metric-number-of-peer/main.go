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
	"strconv"
	"strings"
	"time"
)

const (
	// Default values.
	defaultAddr   = "127.0.0.1"
	defaultPort   = 10004
	defaultTarget = "localhost"
	pluginName    = "near-metric-number-of-peer"
	methodName    = "NearGetNumberOfPeer"
)

var (
	addr   string
	target string
	port   int
)

func init() {
	flag.StringVar(&addr, "addr", defaultAddr, "IP Address(e.g. 0.0.0.0, 127.0.0.1)")
	flag.StringVar(&target, "target", defaultTarget, "Target Node (e.g. 0.0.0.0, default localhost)")
	flag.IntVar(&port, "port", defaultPort, "Port number, default 10004")
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

	cmd := "curl -s " + target + ":3030/metrics | grep -e ^near_peer_connections_total"
	contentMSG := ""
	cmdOutput, cmdError := runCommand(cmd)
	if cmdError != nil {
		state = pluginpb.STATE_FAILURE
		severity = pluginpb.SEVERITY_ERROR
		log.Error().
			Str(methodName, "Error to get Connected number of peers").
			Msg(pluginName)
	}

	f := strings.Split(cmdOutput, " ")
	if len(f) > 1 {
		numOfPeers, numErr := strconv.Atoi(f[1])
		if numErr != nil {
			state = pluginpb.STATE_FAILURE
			severity = pluginpb.SEVERITY_ERROR
			numOfPeers = 0
			log.Error().
				Str(methodName, "Parsing Connected number of peers").
				Msg(pluginName)
		}
		if state == pluginpb.STATE_SUCCESS {
			if numOfPeers == 0 {
				severity = pluginpb.SEVERITY_CRITICAL
				contentMSG = "Number of Peer is 0."
				log.Warn().
					Str(methodName, "CRITICAL: "+contentMSG).
					Msg(pluginName)
			} else {
				contentMSG = "Number of Peer is " + fmt.Sprintf("%d", numOfPeers) + "."
				log.Info().
					Str(methodName, "INFO: "+contentMSG).
					Msg(pluginName)
			}
		}
	}

	ret := sdk.CallResponse{
		FuncName:   methodName,
		Message:    contentMSG,
		Severity:   severity,
		State:      state,
		AlertTypes: []pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD},
	}
	return ret, nil
}

func runCommand(cmd string) (string, error) {
	stdOutput := ""
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Error().
			Str(methodName, "Fail to get connected number of peers").
			Msg(pluginName)
		return stdOutput, err
	}
	outputFinal := strings.TrimSpace(string(out))
	stdOutput = outputFinal
	return stdOutput, nil

}
