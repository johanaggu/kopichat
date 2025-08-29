package chat

import "context"

type Repository interface {
	SaveChat(ctx context.Context, externalID string) (*Chat, error)
	GetChatByID(ctx context.Context, id string) (*Chat, error)
}
