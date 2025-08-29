package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/johanaggu/kopichat/internal/model/message"
)

var _ OpenAIClient = (*client)(nil)

type client struct {
	httpClient   *http.Client
	uriBase      string
	apiKEY       string
	model        string
	instructions string
}

type conversationRes struct {
	ID string `json:"id"`
}

func (c *client) CreateConversation(ctx context.Context) (string, error) {
	conversationURI := fmt.Sprintf("%s/conversations", c.uriBase)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, conversationURI, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.apiKEY))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var conversation conversationRes
	err = json.Unmarshal(body, &conversation)
	if err != nil {
		return "", err
	}

	return conversation.ID, nil
}

type talkReq struct {
	Model          string `json:"model"`
	ConversationID string `json:"conversation"`
	Instructions   string `json:"instructions"`
	Store          bool   `json:"store"`
	Message        string `json:"input"`
}

type talkRes struct {
	Output []struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	} `json:"output"`
}

func (c *client) Talk(ctx context.Context, conversationID, message string) (string, error) {
	responsesURI := fmt.Sprintf("%s/responses", c.uriBase)

	talkReq := talkReq{
		Model:          c.model,
		ConversationID: conversationID,
		Instructions:   c.instructions,
		Store:          true,
		Message:        message,
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(talkReq)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, responsesURI, &buf)
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.apiKEY))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var res talkRes
	err = json.Unmarshal(body, &res)
	if err != nil {
		return "", err
	}

	return res.Output[0].Content[0].Text, nil
}

type msgRes struct {
	Data []struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
		Role string `json:"role"`
	} `json:"data"`
}

func (c *client) GetMessages(ctx context.Context, conversationID, chatID string) ([]*message.Message, error) {
	responsesURI := fmt.Sprintf("%s/conversations/%s/items", c.uriBase, conversationID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, responsesURI, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.apiKEY))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var items msgRes
	err = json.Unmarshal(body, &items)
	if err != nil {
		return nil, err
	}

	var messages []*message.Message
	for _, item := range items.Data {
		var role message.Role
		if item.Role == "assistant" {
			role = message.BotRole
		} else {
			role = message.UserRol
		}

		m := message.NewMessage(0, role, item.Content[0].Text, chatID)
		messages = append(messages, m)
	}

	return messages, nil
}
