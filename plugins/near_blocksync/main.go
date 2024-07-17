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

	cmd := "curl -s " + target + ":3030/metrics | grep -e ^near_block_height_head"
	contentMSG := ""
	cmdOutput, err := runCommand(cmd)
	if err != nil {
		state = pluginpb.STATE_FAILURE
		severity = pluginpb.SEVERITY_ERROR
		contentMSG = "Fail to get Block Height 1st in Diff"

	}
	bHeightCurrent := strings.Split(cmdOutput, " ")
	BHValInt, errParse := strconv.Atoi(bHeightCurrent[1])
	if errParse != nil {
		state = pluginpb.STATE_FAILURE
		severity = pluginpb.SEVERITY_ERROR
		log.Error().
			Str(methodName, "Parsing Error from Current BlockHeight").
			Msg(pluginName)
	}

	if state == pluginpb.STATE_SUCCESS {
		if preBlockHeight == -1 {
			preBlockHeight = BHValInt
			contentMSG = "Setting checked first value of BlockHeight"
		} else {
			diff := BHValInt - preBlockHeight
			if diff < 1 {
				severity = pluginpb.SEVERITY_CRITICAL
				contentMSG = "Block Height's increase has halted for the moment by (" + fmt.Sprintf("%d", diff) + ") > " + fmt.Sprintf("%d", preBlockHeight) + " | " + fmt.Sprintf("%d", BHValInt)
			} else if diff < blockDiff {
				severity = pluginpb.SEVERITY_ERROR
				contentMSG = "Block Height is NOT increasing for the moment by (" + fmt.Sprintf("%d", diff) + ") > " + fmt.Sprintf("%d", preBlockHeight) + " | " + fmt.Sprintf("%d", BHValInt)
			} else {
				contentMSG = "Block Height is increasing by (" + fmt.Sprintf("%d", diff) + ") from " + fmt.Sprintf("%d", preBlockHeight) + " To " + fmt.Sprintf("%d", BHValInt)
				log.Info().
					Str(methodName, contentMSG).
					Msg(pluginName)
			}
			preBlockHeight = BHValInt
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
			Str(methodName, "Fail to get block height").
			Msg(pluginName)
		return stdOutput, err
	}
	outputFinal := strings.TrimSpace(string(out))
	stdOutput = outputFinal
	return stdOutput, nil
}
