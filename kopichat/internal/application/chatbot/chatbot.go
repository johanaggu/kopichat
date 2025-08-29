package chatbot

import (
	"context"
	"fmt"

	"github.com/johanaggu/kopichat/internal/infrastructure/openai"
	"github.com/johanaggu/kopichat/internal/model/chat"
	"github.com/johanaggu/kopichat/internal/model/message"
)

// ChatBot defines the chatbot service interface, represent a use pricpiple usecase.
type ChatBot interface {
	// Discuss handles the chat discussion flow.
	// Returns the updated the last five messages in the conversation.
	Discuss(ctx context.Context, chatID, message string) ([]*message.Message, error)
}

// chatBot implements the ChatBot interface.
type chatBot struct {
	chatRepo  chat.Repository
	msgRepo   message.Repository
	openAICli openai.OpenAIClient
}

// New creates a new instance of object that implements the ChatBot interface.
func New(chatRepo chat.Repository, msgRepo message.Repository, openAICli openai.OpenAIClient) ChatBot {
	return &chatBot{
		chatRepo:  chatRepo,
		msgRepo:   msgRepo,
		openAICli: openAICli,
	}
}

// Discuss handles the chat discussion flow, receives a chat ID and a message.
// If the chat ID is empty, it creates a new conversation and saves it to the repository.
// It then sends the message to the OpenAI client and returns the last five messages in the conversation.
func (c *chatBot) Discuss(ctx context.Context, chatID, msg string) ([]*message.Message, error) {
	var ch *chat.Chat
	var err error

	if chatID != "" {
		ch, err = c.chatRepo.GetChatByID(ctx, chatID)
		if err != nil {
			return nil, fmt.Errorf("failed to get chat: %w", err)
		}
	}

	if chatID == "" {
		externalID, err := c.openAICli.CreateConversation(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create conversation: %w", err)
		}

		ch, err = c.chatRepo.SaveChat(ctx, externalID)
		if err != nil {
			return nil, fmt.Errorf("failed to save chat: %w", err)
		}
	}

	_, err = c.msgRepo.SaveMessage(ctx, message.UserRol, msg, ch.ID())
	if err != nil {
		return nil, fmt.Errorf("error creating user message")
	}

	botMsg, err := c.openAICli.Talk(ctx, ch.ExternalID(), msg)
	if err != nil {
		return nil, fmt.Errorf("error talking with bot")
	}

	_, err = c.msgRepo.SaveMessage(ctx, message.BotRole, botMsg, ch.ID())
	if err != nil {
		return nil, fmt.Errorf("error creating bot message")
	}

	messages, err := c.openAICli.GetMessages(ctx, ch.ExternalID(), ch.ID())
	if err != nil {
		return nil, fmt.Errorf("error retrieve messages from openai api")
	}

	if len(messages) > 5 {
		return messages[:5], nil
	}

	return messages, nil
}
