package main

import (
	"os"
	"sync"

	handler "github.com/greenblat17/yet-another-messenger/auth/api/http"
	"github.com/greenblat17/yet-another-messenger/auth/grpc"
)

func main() {
	probeHandler := handler.NewProbeHandler()

	server := grpc.NewGRPCServer(probeHandler)

	httpPort, ok := os.LookupEnv("CHAT_HTTP_PORT")
	if !ok {
		httpPort = "8081"
	}

	grpcPort, ok := os.LookupEnv("CHAT_GRPC_PORT")
	if !ok {
		grpcPort = "50051"
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		server.RunGRPCServer(grpcPort)
	}()

	go func() {
		defer wg.Done()
		server.RunProxyServer(httpPort)
	}()

	wg.Wait()
}
