.DEFAULT_GOAL := help

.PHONY: help 
description:
	@echo "Hello code reviewer!"
	@echo "Here are the possible commands with the make tool"


.PHONY: test
test:
	@echo "Running unit tests ..."
	@go test ./... -shuffle=on -short -v
	ocker compose -f docker-compose.test.yml run --remove-orphans test test ./... 

.PHONY: build # Running aplication to local 
build:
	docker buildx build --platform linux/amd64 --provenance=false -t kopichatjag-bot:local -f cmd/lambda/chatbot/Dockerfile .
	docker compose up

	
.PHONY: tidy
tidy:
	@echo "Executing go mod tidy..."
	@go mod tidy

.PHONY: fmt
fmt:
	@echo "Executing go fmt..."
	@go fmt ./...
