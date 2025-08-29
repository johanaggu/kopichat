package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/johanaggu/kopichat/cmd/lambda/migrations/handlers"
)

func main() {
	lambda.Start(handlers.Migrate)
}