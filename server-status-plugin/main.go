package main

import (
	"log"
	"time"

	//"time"

	plugin_grpc "server-status-plugin/grpc"
)

const (
	servName = "Server Status Plugin"
)

func main() {
	log.Println("Start Server: ", servName)

	plugin_grpc.StartServer()

	for {
		// TODO: Handle SIGTERM, Shutdown gracefully.
		time.Sleep(time.Second * 10)
	}
}
