package grpc

import (
	"context"

	"github.com/greenblat17/yet-another-messenger/pkg/api/proto/chat/v1/chat/v1"
	"github.com/greenblat17/yet-another-messenger/pkg/api/proto/friendship/v1/friendship/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FriendshipService struct {
	friendship.UnimplementedFriendshipServiceServer
}

func NewFriendshipService() *FriendshipService {
	return &FriendshipService{}
}

func (s *FriendshipService) SendMessage(req chat.ChatService_SendMessageServer) error {
	return status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}

func (s *FriendshipService) SendFriendRequest(ctx context.Context, req *friendship.FriendRequest) (*friendship.FriendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendFriendRequest not implemented")
}

func (s *FriendshipService) AcceptFriendRequest(ctx context.Context, req *friendship.FriendRequest) (*friendship.FriendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AcceptFriendRequest not implemented")
}

func (s *FriendshipService) RejectFriendRequest(ctx context.Context, req *friendship.FriendRequest) (*friendship.FriendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RejectFriendRequest not implemented")
}

func (s *FriendshipService) RemoveFriend(ctx context.Context, req *friendship.RemoveFriendRequest) (*friendship.RemoveFriendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveFriend not implemented")
}

func (s *FriendshipService) GetFriends(ctx context.Context, req *friendship.GetFriendsRequest) (*friendship.GetFriendsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFriends not implemented")
}
