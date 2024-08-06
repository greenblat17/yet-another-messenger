package main

import (
	"os"
	"sync"

	"github.com/greenblat17/yet-another-messenger/chat/internal/api/http"
	"github.com/greenblat17/yet-another-messenger/chat/internal/grpc"
)

func main() {
	probeHandler := http.NewProbeHandler()

	server := grpc.NewGRPCServer(probeHandler)

	httpPort, ok := os.LookupEnv("CHAT_HTTP_PORT")
	if !ok {
		httpPort = "8083"
	}

	grpcPort, ok := os.LookupEnv("CHAT_GRPC_PORT")
	if !ok {
		grpcPort = "50053"
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
