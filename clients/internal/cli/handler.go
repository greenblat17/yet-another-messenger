package cli

import (
	"github.com/greenblat17/yet-another-messenger/pkg/api/proto/auth/v1/auth/v1"
	"github.com/greenblat17/yet-another-messenger/pkg/api/proto/chat/v1/chat/v1"
	"github.com/greenblat17/yet-another-messenger/pkg/api/proto/friendship/v1/friendship/v1"
	"github.com/greenblat17/yet-another-messenger/pkg/api/proto/user/v1/user/v1"
)

type CommandClient struct {
	authClient       auth.AuthServiceClient
	userClient       user.UserServiceClient
	friendshipClient friendship.FriendshipServiceClient
	chatClient       chat.ChatServiceClient
}

func NewCommandClient(
	authClient auth.AuthServiceClient,
	userClient user.UserServiceClient,
	friendshipClient friendship.FriendshipServiceClient,
	chatClient chat.ChatServiceClient,
) *CommandClient {
	return &CommandClient{
		authClient:       authClient,
		userClient:       userClient,
		friendshipClient: friendshipClient,
		chatClient:       chatClient,
	}
}
