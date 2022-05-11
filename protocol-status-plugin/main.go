package main

import (
	"log"
	"time"

	//"time"

	plugin_grpc "protocol-status-plugin/grpc"
)

const (
	servName = "Protocol Status Plugin"
)

func main() {
	log.Println("Start Server: ", servName)

	plugin_grpc.StartServer()

	for {
		// TODO: Handle SIGTERM, Shutdown gracefully.
		time.Sleep(time.Second * 10)
	}
}
