.PHONY: build
build:
	docker buildx build --platform linux/amd64 --provenance=false -t kopichatjag-bot:local -f cmd/lambda/chatbot/Dockerfile .
	docker compose up
