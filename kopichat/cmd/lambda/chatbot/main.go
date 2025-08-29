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

func main() {
	ctx := context.Background()
	conn, err := db.NewConn(ctx, db.Conf{
		User: os.Getenv("MYSQL_USER"),
		Pass: os.Getenv("MYSQL_PASSWORD"),
		Host: os.Getenv("MYSQL_HOST"),
		Port: os.Getenv("MYSQL_PORT"),
		Name: os.Getenv("MYSQL_DATABASE"),
	})
	if err != nil {
		log.Println(err.Error(), os.Getenv("MYSQL_HOST"))
		os.Exit(1)

	}

	httpClient := &http.Client{}
	instructions := ""
	openaiCli := openai.NewClient(httpClient, os.Getenv("OENAI_API_URI"), os.Getenv("OENAI_API_KEY"), os.Getenv("OENAI_MODEL"), instructions)

	chatbotClient := chatbot.New(conn, conn, openaiCli)
	h := handlers.New(chatbotClient)

	lambda.Start(h.Chat)
}
