.DEFAULT_GOAL := help

.PHONY: help 
help:
	@echo "Hello code reviewer!"
	@echo "Here are the possible commands with the make tool"
	@echo ""
	@echo "Available targets:"
	@echo "  install  - Verify Docker CLI, Docker daemon, and Docker Compose v2 are available."
	@echo "  test     - Run unit tests with coverage, print coverage summary, then tear down test services."
	@echo "  run      - Build the chatbot image (buildx, linux/amd64) and start services with docker compose."
	@echo "             Requires env vars: OPENAI_API_KEY"
	@echo "  down     - Stop and remove services, networks, and volumes (removes orphans)."
	@echo "  clean    - Alias for 'down'."
	@echo ""
	@echo "Developer tools:"
	@echo "  tidy     - Run 'go mod tidy'."
	@echo "  fmt      - Run 'go fmt ./...'."
	@echo ""
	@echo "Usage: make <target>    e.g., 'make test'"

.PHONY:install
install:
	@echo "reviewing the necessary tools..."
	@type docker >/dev/null 2>&1 || (echo 'docker not found: install from https://docs.docker.com/get-docker/'; exit 1)
	@docker info >/dev/null 2>&1 || (echo 'docker daemon not running: start Docker Desktop or "sudo systemctl start docker"'; exit 1)
	@docker compose version >/dev/null 2>&1 || (echo 'docker compose v2 not found: install plugin or use Docker Desktop'; exit 1)


.PHONY: test
test:
	@echo "Running unit tests ..."
	@docker compose -f docker-compose.test.yml run --remove-orphans --rm unit_tests \
		test ./... -coverprofile=coverage.out 
	@docker compose -f docker-compose.test.yml run --remove-orphans --rm unit_tests \
		tool cover -func=coverage.out 
	@docker compose -f docker-compose.test.yml down -v --remove-orphans
		
.PHONY: run
run:
	@[ -n "$$OPENAI_API_KEY" ] || (echo 'OPENAI_API_KEY is required. Example: OPENAI_API_KEY=sk-xxx make run'; exit 1)
	@docker buildx build --platform linux/amd64 --provenance=false -t kopichatjag-bot:local -f cmd/lambda/chatbot/Dockerfile .
	@docker compose up 
	
.PHONY: down
down:
	@echo "teardown all running services..."
	@docker compose down -v --remove-orphans

.PHONY: clean
clean: down

.PHONY: tidy
tidy:
	@echo "Executing go mod tidy..."
	@go mod tidy

.PHONY: fmt
fmt:
	@echo "Executing go fmt..."
	@go fmt ./...
