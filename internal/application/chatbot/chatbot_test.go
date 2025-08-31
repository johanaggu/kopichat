package chatbot

import (
	"context"
	"errors"
	"testing"

	"github.com/johanaggu/kopichat/internal/model/chat"
	"github.com/johanaggu/kopichat/internal/model/message"
	"github.com/stretchr/testify/assert"
)

func TestDiscuss(t *testing.T) {
	chatRepo := &chatRepoMock{}
	msgRepo := &msgRepoMock{}
	openAICli := &openAIMock{}

	openaiRes := "world"

	chatID := "4f28c18b-8494-119h-bb77-fa023a94cfd2"
	externalChatID := "external_id"
	chatRepo.GetChatByIDFunc = func(ctx context.Context, id string) (*chat.Chat, error) {
		return chat.New(chatID, externalChatID), nil
	}

	msgRepo.SaveMessageFunc = func(ctx context.Context, role message.Role, content, chatID string) (*message.Message, error) {
		return message.NewMessage(1, message.BotRole, "hello", chatID), nil
	}

	openAICli.TalkFunc = func(ctx context.Context, conversationID, message string) (string, error) {
		return openaiRes, nil
	}

	openAICli.GetMessagesFunc = func(ctx context.Context, conversationID, chatID string) ([]*message.Message, error) {
		return []*message.Message{
			message.NewMessage(1, message.UserRol, "hello", chatID),
			message.NewMessage(2, message.UserRol, "world", chatID),
		}, nil
	}

	service := New(chatRepo, msgRepo, openAICli)

	msgs, err := service.Discuss(context.Background(), chatID, "hello")
	assert.Len(t, msgs, 2)
	assert.NoError(t, err)
}

func TestDiscuss_ErrInvalidChatID(t *testing.T) {
	chatRepo := &chatRepoMock{}
	msgRepo := &msgRepoMock{}
	openAICli := &openAIMock{}

	chatID := "4f28c18b-8494-119h-bb77-fa023a94cfd2"

	chatRepo.GetChatByIDFunc = func(ctx context.Context, id string) (*chat.Chat, error) {
		if id != chatID {
			return nil, errors.New("invalid chat id")
		}
		t.Fatal("same chat id")
		return nil, nil
	}

	service := New(chatRepo, msgRepo, openAICli)

	ms, err := service.Discuss(context.Background(), "1", "hello")
	assert.Error(t, err)
	t.Log(ms)
}
