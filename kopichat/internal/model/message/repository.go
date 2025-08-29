package message

import "context"

type Repository interface {
	SaveMessage(ctx context.Context, role Role, content, chatID string) (*Message, error)
}
