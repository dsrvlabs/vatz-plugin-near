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
	defaultAddr = "127.0.0.1"
	defaultPort = 9091
	pluginName  = "near-metric-blockheight"
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
	diffValue := 3
	timeTicker := 5

	cmd := "curl -s localhost:3030/metrics | grep -e ^near_block_height_head"
	cmdOutputFirst := runCommand(cmd)
	contentMSG := ""

	if cmdOutputFirst == "" {
		state = pluginpb.STATE_FAILURE
		severity = pluginpb.SEVERITY_ERROR
		contentMSG = "Fail to get Block Height"

	} else {

		time.Sleep(time.Duration(timeTicker) * time.Second)
		cmdOutputSecond := runCommand(cmd)
		if cmdOutputSecond == "" {
			state = pluginpb.STATE_FAILURE
			severity = pluginpb.SEVERITY_ERROR
			contentMSG = "Fail to get Block Height"

		} else {
			bHeightVal2 := strings.Split(cmdOutputSecond, " ")
			BHValInt2, err2 := strconv.Atoi(bHeightVal2[1])
			if err2 != nil {
				log.Error().
					Str("NearGetBlockheight", "Parsing Error on Second block height").
					Msg("near-metric-blockheight")
			}

			bHeightVal1 := strings.Split(cmdOutputFirst, " ")
			BHValInt1, err1 := strconv.Atoi(bHeightVal1[1])
			if err1 != nil {
				log.Error().
					Str("NearGetBlockheight", "Parsing Error on First block height").
					Msg("near-metric-blockheight")
			}
			diff := BHValInt2 - BHValInt1
			if diff < diffValue {
				state = pluginpb.STATE_FAILURE
				severity = pluginpb.SEVERITY_ERROR
				contentMSG = "Block Height is NOT increasing for the Moment by (" + fmt.Sprintf("%d", diff) + ") in " + fmt.Sprintf("%d", timeTicker) + " Seconds "
			} else {
				contentMSG = "Block Height is increasing by (" + fmt.Sprintf("%d", diff) + ") in " + fmt.Sprintf("%d", timeTicker) + " Seconds\n" + fmt.Sprintf("%d", BHValInt2) + " | " + fmt.Sprintf("%d", BHValInt1)
				log.Info().
					Str("SuccessToGetBlockHeight", contentMSG).
					Msg("near-metric-blockheight")
			}
		}
	}

	ret := sdk.CallResponse{
		FuncName:   "NearBlockHeight",
		Message:    contentMSG,
		Severity:   severity,
		State:      state,
		AlertTypes: []pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD},
	}
	return ret, nil
}

func runCommand(cmd string) (output string) {
	stdOutput := ""
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Error().
			Str("NearBlockHeight", "Fail to get block height").
			Msg("near-metric-blockheight")
	}
	outputFinal := strings.TrimSpace(string(out))
	stdOutput = outputFinal
	return stdOutput

}
