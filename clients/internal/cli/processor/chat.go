package processor

import (
	"context"

	"github.com/greenblat17/yet-another-messenger/clients/pkg/clients/api/proto/chat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (c *CommandProcessor) GetChatHistory(ctx context.Context, args map[string]string) (any, error) {
	req := &chat.GetChatHistoryRequest{
		ConversationId: args["conversation_id"],
	}

	resp, err := c.chatClient.GetChatHistory(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *CommandProcessor) SendMessage(ctx context.Context, args map[string]string) (any, error) {
	stream, err := c.chatClient.SendMessage(ctx)
	if err != nil {
		return nil, err
	}

	req := &chat.ChatMessage{
		ConversationId: args["conversation_id"],
		SenderId:       args["sender_id"],
		Text:           args["text"],
		Timestamp:      timestamppb.Now(),
	}
	if err := stream.Send(req); err != nil {
		return nil, err
	}

	resp, err := stream.Recv()
	if err != nil {
		return nil, err
	}

	return chat.ChatMessage{
		ConversationId: resp.ConversationId,
		SenderId:       resp.SenderId,
		Text:           resp.Text,
		Timestamp:      resp.Timestamp,
	}, nil
}
