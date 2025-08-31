package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
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

func (h *Handler) Chat(ctx context.Context, req Req) (Res, error) {
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

func (h *Handler) ChatDispacher(ctx context.Context, e events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var raw []byte
	if e.IsBase64Encoded {
		b, err := base64.StdEncoding.DecodeString(e.Body)
		if err != nil {
			log.Printf("error decoding base64 request: %v", err)
			return events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusBadRequest,
				Headers:    genericHeaders(),
			}, nil
		}
		raw = b
	} else {
		raw = []byte(e.Body)
	}

	var req Req
	if err := json.Unmarshal(raw, &req); err != nil || req.Message == "" {
		log.Printf("error unmarshalling request: %v", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    genericHeaders(),
		}, nil
	}

	res, err := h.Chat(ctx, req)
	if err != nil {
		log.Printf("error handling chat: %v", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    genericHeaders(),
		}, nil
	}

	b, err := json.Marshal(res)
	if err != nil {
		log.Printf("error marshalling response: %v", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    genericHeaders(),
		}, nil
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Headers:    genericHeaders(),
		Body:       string(b),
	}, nil

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

func genericHeaders() map[string]string {
	json.Marshal(`{"message":"hello"}`)
	return map[string]string{
		"Content-Type":                 "application/json",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Content-Type,Authorization",
	}
}

type ErrRes struct {
	Message string `json:"message"`
}
