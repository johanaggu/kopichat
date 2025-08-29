package handlers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/johanaggu/kopichat/internal/application/chatbot"
)

func New(chatbot chatbot.ChatBot) *Handler {
	return &Handler{
		chatbot: chatbot,
	}
}

type Handler struct {
	chatbot chatbot.ChatBot
}

func (h *Handler) Chat(ctx context.Context, event json.RawMessage) (Res, error) {
	var req Req
	if err := json.Unmarshal(event, &req); err != nil {
		log.Printf("failed to unmarshal event: %v", err)
		return Res{}, err
	}

	msg, err := h.chatbot.Discuss(ctx, req.ConversationID, req.Message)
	if err != nil {
		log.Printf("errocreating conversation: %v", err)
		return Res{}, err
	}

	var res Res
	res.ConversationID = msg[0].ChatID()
	for _, m := range msg {
		res.Message = append(res.Message, ConversationMessage{
			Role:    string(m.Role()),
			Message: m.Content(),
		})
	}

	return res, nil
}

type Req struct {
	ConversationID string `json:"conversation_id"`
	Message        string `json:"message"`
}

type Res struct {
	ConversationID string                `json:"conversation_id"`
	Message        []ConversationMessage `json:"message"`
}

type ConversationMessage struct {
	Role    string `json:"role"`
	Message string `json:"message"`
}
