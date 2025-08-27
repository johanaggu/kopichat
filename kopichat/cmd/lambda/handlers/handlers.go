package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type RickAPIRes struct {
	Characters string `json:"characters"`
	Locations  string `json:"locations"`
	Episodes   string `json:"episodes"`
}

type Message struct {
	ConversationID string `json:"conversation_id"`
	Message        string `json:"message"`
}

func Chat(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp, err := http.Get("https://rickandmortyapi.com/api")
	if err != nil {
		log.Printf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed read body: %v", err)
	}

	var rickAPIRes RickAPIRes
	err = json.Unmarshal(body, &rickAPIRes)
	if err != nil {
		log.Printf("Failed unmarshall body: %v", err)
	}

	responseBody , err := json.Marshal(rickAPIRes)
	if err != nil {
		log.Printf("Failed marshall response body: %v", err)
	}
	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
	}
	return response, nil
}


func Migrate() {
	
}