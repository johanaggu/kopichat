package openai

import (
	"context"
	"net/http"

	"github.com/johanaggu/kopichat/internal/model/message"
)

// OpenAIClient defines actions for interacting with the OpenAI API.
type OpenAIClient interface {
	// CreateConversation initiates a new conversation and returns its ID,
	// conversation it is messages registry in a OpenAI API.
	CreateConversation(ctx context.Context) (string, error)
	// Talk sends a message to the OpenAI API and returns the message response.
	Talk(ctx context.Context, conversationID, message string) (string, error)
	// GetMessages return
	GetMessages(ctx context.Context, conversationID, chatID string) ([]*message.Message, error)
}

func NewClient(httpClient *http.Client, uriBase, apiKEY, model, instructions string) OpenAIClient {
	return &client{
		httpClient:   httpClient,
		uriBase:      uriBase,
		apiKEY:       apiKEY,
		model:        model,
		instructions: instructions,
	}
}
