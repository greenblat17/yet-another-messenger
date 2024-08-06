package processor

import (
	"context"

	"github.com/greenblat17/yet-another-messenger/clients/pkg/clients/api/proto/auth"
)

func (p *CommandProcessor) RegisterUser(ctx context.Context, args map[string]string) (any, error) {
	req := &auth.RegisterUserRequest{
		Email:    args["email"],
		Password: args["password"],
	}

	resp, err := p.authClient.RegisterUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (p *CommandProcessor) LoginUser(ctx context.Context, args map[string]string) (any, error) {
	req := &auth.LoginUserRequest{
		Email:    args["email"],
		Password: args["password"],
	}

	resp, err := p.authClient.LoginUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (p *CommandProcessor) Logout(ctx context.Context, args map[string]string) (any, error) {
	req := &auth.LogoutRequest{}

	resp, err := p.authClient.Logout(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
