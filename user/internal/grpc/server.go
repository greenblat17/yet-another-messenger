package grpc

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	grpcservice "github.com/greenblat17/yet-another-messenger/user/internal/api/grpc"
	user "github.com/greenblat17/yet-another-messenger/user/pkg/user/api/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type ProbeHandler interface {
	StartUp(w http.ResponseWriter, r *http.Request, pathParams map[string]string)
	Live(w http.ResponseWriter, r *http.Request, pathParams map[string]string)
	Ready(w http.ResponseWriter, r *http.Request, pathParams map[string]string)
}

type UserServer struct {
	Server       *grpc.Server
	probeHandler ProbeHandler
}

func NewGRPCServer(probeHandler ProbeHandler) *UserServer {
	kasp := keepalive.ServerParameters{
		MaxConnectionIdle:     30 * time.Minute, // максимальное время бездействия соединения
		MaxConnectionAge:      30 * time.Minute, // максимальное время соединения
		MaxConnectionAgeGrace: 10 * time.Minute, // время на завершение активных соединений после достижения MaxConnectionAge
		Time:                  10 * time.Minute, // время ожидания перед первым PING
		Timeout:               5 * time.Minute,  // время ожидания PING от клиента
	}

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(kasp),
	)

	user.RegisterUserServiceServer(grpcServer, grpcservice.NewUserService())

	return &UserServer{
		Server:       grpcServer,
		probeHandler: probeHandler,
	}
}

func (s *UserServer) RunGRPCServer(port string) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		log.Printf("Starting gRPC grpc on port %s...", port)
		if err := s.Server.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC grpc: %v", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down gRPC grpc...")
	s.Server.GracefulStop()
}

func (s *UserServer) RunProxyServer(port string) {
	grpcServerEndpoint := flag.String("grpc-endpoint", fmt.Sprintf("localhost:%s", port), "gRPC endpoint")

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Minute,
			Timeout:             5 * time.Minute,
			PermitWithoutStream: true,
		}),
	}

	// http
	err := mux.HandlePath(http.MethodGet, "/start-up", s.probeHandler.StartUp)
	if err != nil {
		log.Fatalf("failed to RegisterUserServiceHandlerFromEndpoint: %v", err)
	}

	err = mux.HandlePath(http.MethodGet, "/live", s.probeHandler.Live)
	if err != nil {
		log.Fatalf("failed to RegisterUserServiceHandlerFromEndpoint: %v", err)
	}

	err = mux.HandlePath(http.MethodGet, "/ready", s.probeHandler.Ready)
	if err != nil {
		log.Fatalf("failed to RegisterUserServiceHandlerFromEndpoint: %v", err)
	}

	// grpc
	err = user.RegisterUserServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	if err != nil {
		log.Fatalf("failed to RegisterUserServiceHandlerFromEndpoint: %v", err)
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	go func() {
		log.Printf("Starting proxy grpc on port %s...", port)
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Error starting proxy grpc: %v", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down proxy grpc...")
}
