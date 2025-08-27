package chat

import "context"

type Repository interface {
	SaveChat(ctx context.Context, c *Chat) error
	RetrieveChat(id string) (*Chat, error)
}