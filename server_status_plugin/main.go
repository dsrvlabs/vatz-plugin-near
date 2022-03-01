package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"vatz-plugin-near/server_status/plugin"
	"vatz-plugin-near/server_status/policy"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	executor policy.Executor
)

func init() {
	executor = policy.NewExecutor()
}

type pluginServer struct {
	plugin.UnimplementedManagerPluginServer
}

func (s *pluginServer) Init(context.Context, *emptypb.Empty) (*plugin.PluginInfo, error) {
	// TODO: TBD
	return nil, nil
}

func (s *pluginServer) Verify(context.Context, *emptypb.Empty) (*plugin.VerifyInfo, error) {
	// TODO: TBD
	return nil, nil
}

func (s *pluginServer) Execute(ctx context.Context, req *plugin.ExecuteRequest) (*plugin.ExecuteResponse, error) {
	log.Println("pluginServer.Execute")

	resp := &plugin.ExecuteResponse{
		State:   plugin.ExecuteResponse_SUCCESS,
		Message: "OK",
	}

	fmt.Printf("ExecuteInfo %+v\n", req.ExecuteInfo)
	fmt.Printf("Fields %+v\n", req.ExecuteInfo.Fields)

	val, ok := req.ExecuteInfo.Fields["function"]
	if !ok {
		resp.State = plugin.ExecuteResponse_FAILURE
		resp.Message = "no valid function"
		return resp, nil
	}

	funcName := val.GetStringValue()

	fmt.Println("Function is ", funcName)
	if funcName == "" {
		resp.State = plugin.ExecuteResponse_FAILURE
		resp.Message = "no valid function"
		return resp, nil
	}

	switch funcName {
	case "IsNearUp":
		isUp, err := executor.IsNearUp()
		if err != nil {
			return nil, err
		}

		if !isUp {
			resp.Message = "dead"
		}

	default:
		log.Println("No selection")
		resp.Message = "No function"
	}

	return resp, nil
}

func main() {
	ch := make(chan os.Signal, 1)
	startServer(ch)
}

func startServer(ch <-chan os.Signal) {
	log.Println("Start vatz-near-plugin")

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, err := net.Listen("tcp", "0.0.0.0:9091")
	if err != nil {
		log.Println(err)
	}

	s := grpc.NewServer()

	serv := pluginServer{}
	plugin.RegisterManagerPluginServer(s, &serv)

	reflection.Register(s)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		_ = <-ch
		cancel()
		s.GracefulStop()
		wg.Done()
	}()

	if err := s.Serve(c); err != nil {
		log.Panic(err)
	}

	wg.Wait()
}
