package processor

import (
	"github.com/greenblat17/yet-another-messenger/clients/pkg/clients/api/proto/auth"
	"github.com/greenblat17/yet-another-messenger/clients/pkg/clients/api/proto/chat"
	"github.com/greenblat17/yet-another-messenger/clients/pkg/clients/api/proto/friendship"
	"github.com/greenblat17/yet-another-messenger/clients/pkg/clients/api/proto/user"
)

type CommandProcessor struct {
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
) *CommandProcessor {
	return &CommandProcessor{
		authClient:       authClient,
		userClient:       userClient,
		friendshipClient: friendshipClient,
		chatClient:       chatClient,
	}
}
