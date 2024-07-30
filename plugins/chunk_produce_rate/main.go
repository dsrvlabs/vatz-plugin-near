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
	defaultPort    = 10003
	defaultNetwork = "mainnet"
	defaultTarget  = "localhost"
	defaultNode    = "dsrvlabs.poolv1.near"
	pluginName     = "chunk_produce_rate"
	methodName     = "NearGetChunkProduceRate"
)

var (
	addr    string
	target  string
	node    string
	network string
	port    int
)

func init() {
	flag.StringVar(&addr, "addr", defaultAddr, "IP Address(e.g. 0.0.0.0, 127.0.0.1)")
	flag.StringVar(&network, "network", defaultNetwork, "network, default: mainnet")
	flag.StringVar(&node, "node", defaultNode, "network, default: dsrvlabs.poolv1.near")
	flag.StringVar(&target, "target", defaultTarget, "Target Node (e.g. 0.0.0.0, default localhost)")
	flag.IntVar(&port, "port", defaultPort, "Port number, default 10003")
	flag.Parse()
}

func main() {
	p := sdk.NewPlugin(pluginName)
	if err := p.Register(pluginFeature); err != nil {
		log.Fatal().Err(err).Msg("Failed to register plugin feature")
	}
	ctx := context.Background()
	if err := p.Start(ctx, addr, port); err != nil {
		fmt.Println("exit")
	}
}

func pluginFeature(info, option map[string]*structpb.Value) (sdk.CallResponse, error) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	state := pluginpb.STATE_SUCCESS
	severity := pluginpb.SEVERITY_INFO
	expRate := math.NaN()
	prdRate := math.NaN()
	contentMSG := ""

	networkAddr := "dsrvlabs.poolv1.near"
	if network != "mainnet" {
		networkAddr = "dsrvlabs.pool.f863973.m0"
	}
	fmt.Println("networkAddr:", networkAddr)
	cmdChunkProduced := "curl -s " + target + ":3030/metrics | grep -e ^near_validators_chunks_produced{account_id='\"" + networkAddr + "\"'}"

	cmdOutput1, err1 := runCommand(cmdChunkProduced)
	if err1 != nil {
		state = pluginpb.STATE_FAILURE
		severity = pluginpb.SEVERITY_ERROR
		contentMSG = contentMSG + "Error to get validators chunks produced\n"
	}

	cmdChunkExpected := "curl -s " + target + ":3030/metrics | grep -e ^near_validators_chunks_expected{account_id='\"" + networkAddr + "\"'}"

	cmdOutput2, err2 := runCommand(cmdChunkExpected)
	if err2 != nil {
		state = pluginpb.STATE_FAILURE
		severity = pluginpb.SEVERITY_ERROR
		contentMSG = contentMSG + "Error to get validators chunks expected\n"
	}

	producedVal := strings.Split(cmdOutput1, " ")
	expectedVal := strings.Split(cmdOutput2, " ")

	if len(producedVal) > 1 && len(expectedVal) > 1 {
		producedRate, errPR := strconv.Atoi(producedVal[1])
		if errPR != nil {
			state = pluginpb.STATE_FAILURE
			severity = pluginpb.SEVERITY_ERROR
			log.Error().
				Str(methodName, "Parsing Error on chunk produced Rate").
				Msg(pluginName)
		}

		expectedRate, errER := strconv.Atoi(expectedVal[1])
		if errER != nil {
			state = pluginpb.STATE_FAILURE
			severity = pluginpb.SEVERITY_ERROR
			log.Error().
				Str(methodName, "Parsing Error on chunk expected Rate").
				Msg(pluginName)
		}

		if state == pluginpb.STATE_SUCCESS {
			prdRate = float64(producedRate)
			expRate = float64(expectedRate)
		}
	}

	if !math.IsNaN(prdRate) && !math.IsNaN(expRate) {
		chunkProducedRate := math.Round(prdRate / expRate * 100)
		log.Info().
			Str("Calculated Chunk Produce Rate: ", contentMSG).
			Msg(pluginName)
		if state == pluginpb.STATE_SUCCESS {
			if chunkProducedRate < 51 {
				severity = pluginpb.SEVERITY_CRITICAL
				contentMSG = "Chunk Produced Rate is (" + fmt.Sprintf("%.2f", chunkProducedRate) + "), which is under normal rate(50%)."
				log.Warn().
					Str(methodName, "CRITICAL: "+contentMSG).
					Msg(pluginName)
			} else if chunkProducedRate < 95 {
				severity = pluginpb.SEVERITY_WARNING
				contentMSG = "Chunk Produced Rate is (" + fmt.Sprintf("%.2f", chunkProducedRate) + "), which is under normal rate(95%)."
				log.Warn().
					Str(methodName, "WARNING: "+contentMSG).
					Msg(pluginName)
			} else {
				contentMSG = "Chunk Produced Rate is (" + fmt.Sprintf("%.2f", chunkProducedRate) + " %)"
				log.Info().
					Str(methodName, "SUCCESS: "+contentMSG).
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
			Str(methodName, "Fail to get Chunk produce Rate").
			Msg(pluginName)
		return stdOutput, err
	}
	outputFinal := strings.TrimSpace(string(out))
	stdOutput = outputFinal
	return stdOutput, nil

}
