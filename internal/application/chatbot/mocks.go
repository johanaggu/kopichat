package chatbot

import (
	"context"

	"github.com/johanaggu/kopichat/internal/model/chat"
	"github.com/johanaggu/kopichat/internal/model/message"
)

// chatRepoMock
type chatRepoMock struct {
	SaveChatFunc    func(ctx context.Context, externalID string) (*chat.Chat, error)
	GetChatByIDFunc func(ctx context.Context, id string) (*chat.Chat, error)
}

func (m *chatRepoMock) SaveChat(ctx context.Context, externalID string) (*chat.Chat, error) {
	return m.SaveChatFunc(ctx, externalID)
}
func (m *chatRepoMock) GetChatByID(ctx context.Context, id string) (*chat.Chat, error) {
	return m.GetChatByIDFunc(ctx, id)
}

type msgRepoMock struct {
	SaveMessageFunc func(ctx context.Context, role message.Role, content, chatID string) (*message.Message, error)
}

func (m *msgRepoMock) SaveMessage(ctx context.Context, role message.Role, content, chatID string) (*message.Message, error) {
	return m.SaveMessageFunc(ctx, role, content, chatID)
}

type openAIMock struct {
	CreateConversationFunc func(ctx context.Context) (string, error)
	TalkFunc               func(ctx context.Context, conversationID, message string) (string, error)
	GetMessagesFunc        func(ctx context.Context, conversationID, chatID string) ([]*message.Message, error)
}

func (m *openAIMock) CreateConversation(ctx context.Context) (string, error) {
	return m.CreateConversationFunc(ctx)
}

func (m *openAIMock) Talk(ctx context.Context, conversationID, message string) (string, error) {
	return m.TalkFunc(ctx, conversationID, message)
}

func (m *openAIMock) GetMessages(ctx context.Context, conversationID, chatID string) ([]*message.Message, error) {
	return m.GetMessagesFunc(ctx, conversationID, chatID)
}
