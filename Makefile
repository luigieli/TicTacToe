# Makefile for Ultimate Tic Tac Toe Infrastructure

# Variables
DC = docker compose

# Default target
.PHONY: help
help:
	@echo "Usage:"
	@echo "  make build         Build Docker images"
	@echo "  make up            Start services in detached mode"
	@echo "  make down          Stop and remove containers"
	@echo "  make test-infra    Verify connectivity to containers"

.PHONY: build
build:
	$(DC) build

.PHONY: up
up:
	$(DC) up -d

.PHONY: down
down:
	$(DC) down

.PHONY: test-infra
test-infra:
	@echo "Testing Frontend (8082)..."
	@curl -s -o /dev/null -w "%{http_code}" http://localhost:8082 | grep 200 > /dev/null && echo "Frontend is UP" || (echo "Frontend is DOWN" && exit 1)
	@echo "Testing Backend (8083)..."
	@curl -s http://localhost:8083/games -X POST | grep "game_id" > /dev/null && echo "Backend is UP" || (echo "Backend is DOWN" && exit 1)
