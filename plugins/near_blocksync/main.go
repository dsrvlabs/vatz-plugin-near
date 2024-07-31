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
	// Default values
	defaultAddr      = "127.0.0.1"
	defaultPort      = 10002
	defaultTarget    = "localhost"
	defaultBlockDiff = 3
	pluginName       = "near_blocksync"
	methodName       = "NearGetBlockHeight"
)

var (
	addr           string
	target         string
	port           int
	blockDiff      int
	preBlockHeight = -1
)

func init() {
	flag.StringVar(&addr, "addr", defaultAddr, "IP Address(e.g. 0.0.0.0, 127.0.0.1)")
	flag.IntVar(&port, "port", defaultPort, "Port number, default: 10002")
	flag.IntVar(&blockDiff, "ticker", defaultBlockDiff, " for BlockHeight Difference value, default is 3 blocks")
	flag.StringVar(&target, "target", defaultTarget, "Target Node (e.g. 0.0.0.0, default localhost)")
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

	cmd := "curl -s " + target + ":3030/metrics | grep -e ^near_block_height_head"
	cmdOutput, err := runCommand(cmd)
	if err != nil {
		log.Error().Str(methodName, "Fail to get block height").Msg(pluginName)
		return createResponse(pluginpb.STATE_FAILURE, pluginpb.SEVERITY_ERROR, "Fail to get block height due to runCMD"), nil
	}

	bHeightCurrent := strings.Split(cmdOutput, " ")
	if len(bHeightCurrent) < 2 {
		log.Error().Str(methodName, "Unexpected format from block height command output").Msg(pluginName)
		return createResponse(pluginpb.STATE_FAILURE, pluginpb.SEVERITY_ERROR, "Unexpected format from block height command output"), nil
	}

	BHValInt, err := strconv.Atoi(bHeightCurrent[1])
	if err != nil {
		log.Error().Str(methodName, "Parsing Error from Current BlockHeight").Msg(pluginName)
		return createResponse(pluginpb.STATE_FAILURE, pluginpb.SEVERITY_ERROR, "Parsing Error from Current BlockHeight"), nil
	}

	var contentMSG string
	severity := pluginpb.SEVERITY_INFO

	if preBlockHeight == -1 {
		preBlockHeight = BHValInt
		contentMSG = "Setting checked first value of BlockHeight"
	} else {
		diff := BHValInt - preBlockHeight
		if diff < 1 {
			severity = pluginpb.SEVERITY_CRITICAL
			contentMSG = fmt.Sprintf("Block Height's increase has halted for the moment by (%d) > %d | %d", diff, preBlockHeight, BHValInt)
		} else if diff < blockDiff {
			severity = pluginpb.SEVERITY_ERROR
			contentMSG = fmt.Sprintf("Block Height is NOT increasing for the moment by (%d) > %d | %d", diff, preBlockHeight, BHValInt)
		} else {
			contentMSG = fmt.Sprintf("Block Height is increasing by (%d) from %d To %d", diff, preBlockHeight, BHValInt)
			log.Info().Str(methodName, contentMSG).Msg(pluginName)
		}
		preBlockHeight = BHValInt
	}

	return createResponse(pluginpb.STATE_SUCCESS, severity, contentMSG), nil
}

func createResponse(state pluginpb.STATE, severity pluginpb.SEVERITY, message string) sdk.CallResponse {
	return sdk.CallResponse{
		FuncName:   methodName,
		Message:    message,
		Severity:   severity,
		State:      state,
		AlertTypes: []pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD},
	}
}

func runCommand(cmd string) (string, error) {
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Error().Str(methodName, "Fail to get block height").Msg(pluginName)
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
