package grpc

import (
	"context"

	"github.com/greenblat17/yet-another-messenger/pkg/api/proto/chat/v1/chat/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ChatService struct {
	chat.UnimplementedChatServiceServer
}

func NewChatService() *ChatService {
	return &ChatService{}
}

func (s *ChatService) SendMessage(req chat.ChatService_SendMessageServer) error {
	return status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}

func (s *ChatService) GetChatHistory(ctx context.Context, req *chat.GetChatHistoryRequest) (*chat.GetChatHistoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChatHistory not implemented")
}
