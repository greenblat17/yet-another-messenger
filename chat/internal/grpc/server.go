package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	grpcservice "github.com/greenblat17/yet-another-messenger/chat/internal/api/grpc"
	chat "github.com/greenblat17/yet-another-messenger/chat/pkg/chat/api/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type ProbeHandler interface {
	StartUp(w http.ResponseWriter, r *http.Request, pathParams map[string]string)
	Live(w http.ResponseWriter, r *http.Request, pathParams map[string]string)
	Ready(w http.ResponseWriter, r *http.Request, pathParams map[string]string)
}

type ChatServer struct {
	Server             *grpc.Server
	probeHandler       ProbeHandler
	portConfig         *PortConfig
	grpcServerEndpoint *string
}

func NewGRPCServer(probeHandler ProbeHandler, portConfig *PortConfig, grpcServerEndpoint *string) *ChatServer {
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

	chat.RegisterChatServiceServer(grpcServer, grpcservice.NewChatService())

	return &ChatServer{
		Server:             grpcServer,
		probeHandler:       probeHandler,
		portConfig:         portConfig,
		grpcServerEndpoint: grpcServerEndpoint,
	}
}

func (s *ChatServer) RunGRPCServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.portConfig.grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		log.Printf("Starting gRPC grpc on port %s...", s.portConfig.grpcPort)
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

func (s *ChatServer) RunHTTPProxyServer() {

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
		log.Fatalf("failed to RegisterOrderHandlerFromEndpoint: %v", err)
	}

	err = mux.HandlePath(http.MethodGet, "/live", s.probeHandler.Live)
	if err != nil {
		log.Fatalf("failed to RegisterOrderHandlerFromEndpoint: %v", err)
	}

	err = mux.HandlePath(http.MethodGet, "/ready", s.probeHandler.Ready)
	if err != nil {
		log.Fatalf("failed to RegisterOrderHandlerFromEndpoint: %v", err)
	}

	// grpc
	err = chat.RegisterChatServiceHandlerFromEndpoint(ctx, mux, *s.grpcServerEndpoint, opts)
	if err != nil {
		log.Fatalf("failed to RegisterOrderHandlerFromEndpoint: %v", err)
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", s.portConfig.httpProxyPort),
		Handler: mux,
	}

	go func() {
		log.Printf("Starting proxy grpc on port %s...", s.portConfig.httpProxyPort)
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Error starting proxy grpc: %v", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down proxy grpc...")
}

func (s *ChatServer) RunWebSocketProxyServer() {
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

	// grpc
	err := chat.RegisterChatServiceHandlerFromEndpoint(ctx, mux, *s.grpcServerEndpoint, opts)
	if err != nil {
		log.Fatalf("failed to RegisterOrderHandlerFromEndpoint: %v", err)
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", s.portConfig.wsProxyPort),
		Handler: wsproxy.WebsocketProxy(mux),
	}

	// TODO: не работает
	go func() {
		log.Printf("Starting ws proxy grpc on port %s...", s.portConfig.wsProxyPort)
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Error starting proxy grpc: %v", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down ws proxy grpc...")
}
