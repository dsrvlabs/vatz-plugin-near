package main

import (
	"context"
	pluginpb "github.com/xellos00/silver-bentonville/dist/proto/dsrv/api/node_manager/plugin"
	"log"
	"net"
	"os"
	"sync"
	rpc "vatz-plugin-near/protocol_status/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func initiateServer(ch <-chan os.Signal) {
	log.Println("Start vatz-near-plugin")

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, err := net.Listen("tcp", "0.0.0.0:9091")
	if err != nil {
		log.Println(err)
	}

	s := grpc.NewServer()
	serv := rpc.GrpcService{}
	pluginpb.RegisterManagerPluginServer(s, &serv)
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

func main() {
	ch := make(chan os.Signal, 1)
	initiateServer(ch)
}
