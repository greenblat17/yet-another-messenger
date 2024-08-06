package processor

import (
	"context"

	"github.com/greenblat17/yet-another-messenger/clients/pkg/clients/api/proto/friendship"
)

func (p *CommandProcessor) SendFriendRequest(ctx context.Context, args map[string]string) (any, error) {
	req := &friendship.FriendRequest{
		UserId:   args["user_id"],
		FriendId: args["friend_id"],
	}

	resp, err := p.friendshipClient.SendFriendRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (p *CommandProcessor) AcceptFriendRequest(ctx context.Context, args map[string]string) (any, error) {
	req := &friendship.FriendRequest{
		UserId:   args["user_id"],
		FriendId: args["friend_id"],
	}

	resp, err := p.friendshipClient.AcceptFriendRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (p *CommandProcessor) RejectFriendRequest(ctx context.Context, args map[string]string) (any, error) {
	req := &friendship.FriendRequest{
		UserId:   args["user_id"],
		FriendId: args["friend_id"],
	}

	resp, err := p.friendshipClient.RejectFriendRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (p *CommandProcessor) RemoveFriend(ctx context.Context, args map[string]string) (any, error) {
	req := &friendship.RemoveFriendRequest{
		UserId:   args["user_id"],
		FriendId: args["friend_id"],
	}

	resp, err := p.friendshipClient.RemoveFriend(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (p *CommandProcessor) GetFriends(ctx context.Context, args map[string]string) (any, error) {
	req := &friendship.GetFriendsRequest{
		UserId: args["user_id"],
	}

	resp, err := p.friendshipClient.GetFriends(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
