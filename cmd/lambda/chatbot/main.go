package main

import (
	"context"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/johanaggu/kopichat/cmd/lambda/chatbot/handlers"
	"github.com/johanaggu/kopichat/internal/application/chatbot"
	db "github.com/johanaggu/kopichat/internal/infrastructure/mysql"
	"github.com/johanaggu/kopichat/internal/infrastructure/openai"
)

type environment string

const (
	Local      environment = "local"
	Production environment = "production"
)

func main() {
	env := environment(os.Getenv("ENVIRONMENT"))

	userDB := os.Getenv("MYSQL_USER")
	passDB := os.Getenv("MYSQL_PASSWORD")
	hostDB := os.Getenv("MYSQL_HOST")
	portDB := os.Getenv("MYSQL_PORT")
	nameDB := os.Getenv("MYSQL_DATABASE")

	openaiAPIURI := os.Getenv("OPENAI_API_URI")
	openaiAPIKey := os.Getenv("OPENAI_API_KEY")
	openaiModel := os.Getenv("OPENAI_MODEL")
	ctx := context.Background()
	conn, err := db.NewConn(ctx, db.Conf{
		User: userDB,
		Pass: passDB,
		Host: hostDB,
		Port: portDB,
		Name: nameDB,
	})
	if err != nil {
		log.Println("error connecting with mysql")
		os.Exit(1)

	}

	httpClient := &http.Client{}
	instructions := `
	`
	openaiCli := openai.NewClient(httpClient, openaiAPIURI, openaiAPIKey, openaiModel, instructions)

	chatbotClient := chatbot.New(conn, conn, openaiCli)
	h := handlers.New(chatbotClient)

	if env == Local {
		lambda.Start(h.Chat)
	}

	lambda.Start(h.ChatDispacher)
}
