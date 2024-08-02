package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/greenblat17/yet-another-messenger/chat/internal/api/http"
	"github.com/greenblat17/yet-another-messenger/chat/internal/grpc"
)

func main() {
	httpPort, ok := os.LookupEnv("CHAT_HTTP_PORT")
	if !ok {
		httpPort = "8083"
	}

	wsPort, ok := os.LookupEnv("CHAT_GRPC_PORT")
	if !ok {
		wsPort = "8090"
	}

	grpcPort, ok := os.LookupEnv("CHAT_WS_PORT")
	if !ok {
		grpcPort = "50053"
	}

	probeHandler := http.NewProbeHandler()
	portConfig := grpc.NewPortConfig(grpcPort, httpPort, wsPort)
	grpcServerEndpoint := flag.String(
		"grpc-endpoint",
		fmt.Sprintf("localhost:%s", portConfig.GRPCPort()),
		"gRPC endpoint",
	)

	server := grpc.NewGRPCServer(probeHandler, portConfig, grpcServerEndpoint)

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		server.RunGRPCServer()
	}()

	go func() {
		defer wg.Done()
		server.RunHTTPProxyServer()
	}()

	go func() {
		defer wg.Done()
		server.RunWebSocketProxyServer()
	}()

	wg.Wait()
}
