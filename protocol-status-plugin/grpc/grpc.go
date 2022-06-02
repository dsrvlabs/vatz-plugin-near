package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

const (
	grpcPort = 9092
)

type grpcService struct {
	pluginpb.UnimplementedPluginServer
}

func (s *grpcService) Verify(ctx context.Context, in *emptypb.Empty) (*pluginpb.VerifyInfo, error) {
	fmt.Println("Plugin Verify Has been Called: ALIVE 9091")
	resp := pluginpb.VerifyInfo{
		VerifyMsg: "UP",
	}
	return &resp, nil
}

func (s *grpcService) Execute(ctx context.Context, in *pluginpb.ExecuteRequest) (*pluginpb.ExecuteResponse, error) {
	// Method name 에 따라서 알아서 구현할 수 있도록 한다.
	fmt.Println("Plugin Execute Has been Called::: ALIVE")
	req := in.ExecuteInfo.AsMap()["execute_method"].(string)

	if req == "isUp" {
		return &pluginpb.ExecuteResponse{
			State:        pluginpb.STATE_FAILURE,
			Message:      "This is SUCCESS Message:9091 for Test : isUp ",
			AlertType:    []pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD},
			Severity:     pluginpb.SEVERITY_CRITICAL,
			ResourceType: "near-",
		}, nil
	} else if req == "getBlockHeight" {
		return &pluginpb.ExecuteResponse{
			State:        pluginpb.STATE_FAILURE,
			Message:      "This is FAIL Message for Test:9091 :getBlockHeight ",
			AlertType:    []pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD},
			Severity:     pluginpb.SEVERITY_CRITICAL,
			ResourceType: "near-",
		}, nil
	} else {
		return &pluginpb.ExecuteResponse{
			State:        pluginpb.STATE_FAILURE,
			Message:      "This is FAIL Message for Test:9091 : else",
			AlertType:    []pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD},
			Severity:     pluginpb.SEVERITY_CRITICAL,
			ResourceType: "near-",
		}, nil
	}
}

// StartServer try to start api manager.
func StartServer() error {
	s := grpc.NewServer()
	serv := grpcService{}

	pluginpb.RegisterPluginServer(s, &serv)
	reflection.Register(s)

	addr := fmt.Sprintf(":%d", grpcPort)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panic(err)
		return err
	}

	log.Println("listen ", addr)

	go func() {
		if err := s.Serve(l); err != nil {
			log.Panic(err)
		}
	}()

	log.Println("Pilot Plugin started")
	return nil
}
