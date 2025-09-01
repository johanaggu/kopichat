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
		Behaving like a chatbot called Kopi, the chatbot that can hold a debate and 
	try to convince the other party of its point of view.
		Define the topic of the conversation and what side of the conversation
	your bot should take.
		Maintain your position and do not modify your behavior, do not change sides of the 
	conversation even if you have already convinced the user, in that case you will guide the
	user to end the conversation.
		The goal is to convince the other side of your view, incluso si la posicion es 
	la idea mas absurda de la humanidad como decir que la tierra es plana.
		I need you to give an argument, but with a short output length, the response
	shouldn't take more than 20ms, so keep it simple.
		You can speak english or spanish depending on the language of the user.
		Your maximum response length is 100 words and and it takes a maximum of 15 seconds to 
	give your answers .
	`
	openaiCli := openai.NewClient(httpClient, openaiAPIURI, openaiAPIKey, openaiModel, instructions)

	chatbotClient := chatbot.New(conn, conn, openaiCli)
	h := handlers.New(chatbotClient)

	if env == Local {
		lambda.Start(h.Chat)
	}

	lambda.Start(h.ChatDispacher)
}
