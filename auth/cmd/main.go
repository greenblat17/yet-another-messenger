package main

import (
	"os"
	"sync"

	handler "github.com/greenblat17/yet-another-messenger/auth/internal/api/http"
	"github.com/greenblat17/yet-another-messenger/auth/internal/grpc"
)

func main() {
	probeHandler := handler.NewProbeHandler()

	server := grpc.NewGRPCServer(probeHandler)

	httpPort, ok := os.LookupEnv("AUTH_HTTP_PORT")
	if !ok {
		httpPort = "8080"
	}

	grpcPort, ok := os.LookupEnv("AUTH_GRPC_PORT")
	if !ok {
		grpcPort = "50050"
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
