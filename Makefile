.PHONY: help build run-api run-notifier dev lint format migrate-up migrate-down docker-up docker-down generate-docs

help:
	@echo "Available commands:"
	@echo "  build - Build the application"
	@echo "  run-api - Run the API"
	@echo "  run-notifier - Run the notifier"
	@echo "  dev - Run the application in development mode"
	@echo "  lint - Lint the application"
	@echo "  format - Format the application"
	@echo "  migrate-up - Run database migrations"
	@echo "  migrate-down - Rollback database migrations"
	@echo "  docker-up - Start the application in development mode"
	@echo "  docker-down - Stop the application"
	@echo "  generate-docs - Generate Swagger documentation"

build:
	go run ./scripts/build.go

run-api: 
	go run ./cmd/api

run-notifier: 
	go run ./cmd/notifier

dev:
	go run ./cmd/api

lint: format
	golangci-lint run ./...

format:
	gofmt -s -w .

migrate-up:
	migrate -path db/migrations -database "postgresql://postgres:password@localhost:5432/ecommerce_shop?sslmode=disable" up

migrate-down:
	migrate -path db/migrations -database "postgresql://postgres:password@localhost:5432/ecommerce_shop?sslmode=disable" down

docker-up:
	docker-compose -f docker/docker-compose.yml up -d

docker-down:
	docker-compose -f docker/docker-compose.yml down


# Generate Swagger documentation
# This command requires the swag CLI tool to be installed
# Run: go install github.com/swaggo/swag/cmd/swag@latest
# Then run: make generate-docs
generate-docs:
# 	mkdir -p docs
	swag init -g cmd/api/main.go --o docs --parseDependency --parseInternal --exclude .git,docs,docker,db