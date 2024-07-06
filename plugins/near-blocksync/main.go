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
	defaultAddr     = "127.0.0.1"
	defaultPort     = 10002
	defaultTarget   = "localhost"
	defaultTimeTick = 5
	pluginName      = "near-blocksync"
	methodName      = "NearGetBlockHeight"
)

var (
	addr   string
	target string
	port   int
	ticker int
)

func init() {
	flag.StringVar(&addr, "addr", defaultAddr, "IP Address(e.g. 0.0.0.0, 127.0.0.1)")
	flag.IntVar(&port, "port", defaultPort, "Port number, default: 10002")
	flag.IntVar(&ticker, "ticker", defaultTimeTick, "Time ticker for block chain Diff, default is 5 (seconds)")
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
	diffValue := 3

	cmd := "curl -s " + target + ":3030/metrics | grep -e ^near_block_height_head"
	contentMSG := ""
	cmdOutput1st, err1st := runCommand(cmd)

	if err1st != nil {
		state = pluginpb.STATE_FAILURE
		severity = pluginpb.SEVERITY_ERROR
		contentMSG = "Fail to get Block Height 1st in Diff"

	} else {
		time.Sleep(time.Duration(ticker) * time.Second)
		cmdOutput2nd, err2nd := runCommand(cmd)
		if err2nd != nil {
			state = pluginpb.STATE_FAILURE
			severity = pluginpb.SEVERITY_ERROR
			contentMSG = "Fail to get Block Height 2nd in Diff"
		}

		if state == pluginpb.STATE_SUCCESS {

			bHeightVal2 := strings.Split(cmdOutput2nd, " ")
			BHValInt2, err2 := strconv.Atoi(bHeightVal2[1])

			if err2 != nil {
				state = pluginpb.STATE_FAILURE
				severity = pluginpb.SEVERITY_ERROR
				log.Error().
					Str(methodName, "Parsing Error on 2nd block height for diff").
					Msg(pluginName)
			}

			bHeightVal1 := strings.Split(cmdOutput1st, " ")
			BHValInt1, err1 := strconv.Atoi(bHeightVal1[1])

			if err1 != nil {
				state = pluginpb.STATE_FAILURE
				severity = pluginpb.SEVERITY_ERROR
				log.Error().
					Str(methodName, "Parsing Error on 1st block height for diff").
					Msg(pluginName)
			}

			if state == pluginpb.STATE_SUCCESS {
				diff := BHValInt2 - BHValInt1
				if diff < 1 {
					severity = pluginpb.SEVERITY_CRITICAL
					contentMSG = "Block Height's increase has halted for the moment by (" + fmt.Sprintf("%d", diff) + ") in " + fmt.Sprintf("%d", ticker) + " Seconds "
				} else if diff < diffValue {
					severity = pluginpb.SEVERITY_ERROR
					contentMSG = "Block Height is NOT increasing for the moment by (" + fmt.Sprintf("%d", diff) + ") in " + fmt.Sprintf("%d", ticker) + " Seconds "
				} else {
					contentMSG = "Block Height is increasing by (" + fmt.Sprintf("%d", diff) + ") in " + fmt.Sprintf("%d", ticker) + " Seconds\n" + fmt.Sprintf("%d", BHValInt2) + " | " + fmt.Sprintf("%d", BHValInt1)
					log.Info().
						Str(methodName, contentMSG).
						Msg(pluginName)
				}
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
			Str(methodName, "Fail to get block height").
			Msg(pluginName)
		return stdOutput, err
	}
	outputFinal := strings.TrimSpace(string(out))
	stdOutput = outputFinal
	return stdOutput, nil
}
