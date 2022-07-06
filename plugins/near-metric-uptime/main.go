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
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	// Default values.
	defaultAddr    = "127.0.0.1"
	defaultPort    = 10005
	defaultTarget  = "localhost"
	defaultNetwork = "mainnet"
	pluginName     = "near-metric-uptime"
	methodName     = "NearGetUptime"
)

var (
	addr    string
	target  string
	network string
	port    int
)

func init() {
	flag.StringVar(&addr, "addr", defaultAddr, "IP Address(e.g. 0.0.0.0, 127.0.0.1)")
	flag.StringVar(&network, "network", defaultNetwork, "network, default mainnet")
	flag.StringVar(&target, "target", defaultTarget, "Target Node (e.g. 0.0.0.0, default localhost)")
	flag.IntVar(&port, "port", defaultPort, "Port number, default 10005")
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
	contentMSG := ""
	account := "dsrvlabs.poolv1.near"
	if network != "mainnet" {
		account = "dsrvlabs.pool.f863973.m0"
	}
	fmt.Println("account:", account)
	cmdBlockProduced := "curl -s " + target + ":3030/metrics | grep -e ^near_validators_blocks_produced{account_id='\"" + account + "\"'}"

	cmdOutput1, err1 := runCommand(cmdBlockProduced)
	if err1 != nil {
		state = pluginpb.STATE_FAILURE
		severity = pluginpb.SEVERITY_ERROR
		contentMSG = contentMSG + "Error to get validators chunks produced\n"
	}

	cmdBlockExpected := "curl -s " + target + ":3030/metrics | grep -e ^near_validators_blocks_expected{account_id='\"" + account + "\"'}"

	cmdOutput2, err2 := runCommand(cmdBlockExpected)
	if err2 != nil {
		state = pluginpb.STATE_FAILURE
		severity = pluginpb.SEVERITY_ERROR
		contentMSG = contentMSG + "Error to get validators expected Block\n"
	}

	producedVal := strings.Split(cmdOutput1, " ")
	producedRate, errPR := strconv.Atoi(producedVal[1])
	if errPR != nil {
		state = pluginpb.STATE_FAILURE
		severity = pluginpb.SEVERITY_ERROR
		log.Error().
			Str(methodName, "Parsing Error on block produced Block").
			Msg(pluginName)
	}
	expectedVal := strings.Split(cmdOutput2, " ")
	expectedRate, errER := strconv.Atoi(expectedVal[1])
	if errER != nil {
		state = pluginpb.STATE_FAILURE
		severity = pluginpb.SEVERITY_ERROR
		log.Error().
			Str(methodName, "Parsing Error on block expected Rate").
			Msg(pluginName)
	}

	f := producedRate / expectedRate
	chunkProducedRate := math.Round(float64(f * 100))
	if state == pluginpb.STATE_SUCCESS {
		if chunkProducedRate < 50 {
			severity = pluginpb.SEVERITY_CRITICAL
			contentMSG = "Node's uptime Rate is (" + fmt.Sprintf("%.2f", chunkProducedRate) + "), which is way too lower than normal rate(95%)."
			log.Warn().
				Str(methodName, "CRITICAL: "+contentMSG).
				Msg(pluginName)
		} else if chunkProducedRate < 94 {
			severity = pluginpb.SEVERITY_WARNING
			contentMSG = "Node's uptime Rate is (" + fmt.Sprintf("%.2f", chunkProducedRate) + "), which is lower than normal rate(95%)."
			log.Warn().
				Str(methodName, "WARNING: "+contentMSG).
				Msg(pluginName)
		} else {
			contentMSG = "Node's uptime Rate is Normal as (" + fmt.Sprintf("%.2f", chunkProducedRate) + " %)"
			log.Info().
				Str(methodName, "INFO: "+contentMSG).
				Msg(pluginName)
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
			Str(methodName, "Fail to get").
			Msg(pluginName)
		return stdOutput, err
	}
	outputFinal := strings.TrimSpace(string(out))
	stdOutput = outputFinal
	return stdOutput, nil

}
