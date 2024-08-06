package processor

import (
	"context"

	"github.com/greenblat17/yet-another-messenger/clients/pkg/clients/api/proto/user"
)

func (p *CommandProcessor) GetUser(ctx context.Context, args map[string]string) (any, error) {
	req := &user.GetUserRequest{
		Username: args["username"],
	}

	resp, err := p.userClient.GetUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (p *CommandProcessor) UpdateUser(ctx context.Context, args map[string]string) (any, error) {
	req := &user.UpdateUserRequest{
		Username: args["username"],
		Bio:      args["bio"],
	}

	resp, err := p.userClient.UpdateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
