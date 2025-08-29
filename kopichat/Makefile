.DEFAULT_GOAL := help

.PHONY: help 
help:
	@echo "Hello code reviewer!"
	@echo "Here are the possible commands with the make tool"

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
	@docker buildx build --platform linux/amd64 --provenance=false -t kopichatjag-bot:local -f cmd/lambda/chatbot/Dockerfile .
	@docker compose up -d
	
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
