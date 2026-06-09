.PHONY: help \
	build test \
	api-build api-migrate-up api-migrate-down api-test api-stop api-kill-local \
	api-run api-run-dev api-run-staging api-run-prod env-init bootstrap \
	dev staging prod \
	cloudflared-tunnel \
	docker-db-up \
	mobile-test mobile-test-unit mobile-sync-ip \
	commit diff push pull gitlog

# =========================
# Variables
# =========================
API_ENV ?= dev
API_ENV_FILE=api/$(API_ENV).env
-include $(API_ENV_FILE)

DB_HOST ?= postgres
DB_PORT ?= 5432
DB_PUBLISHED_HOST ?= 127.0.0.1
DB_PUBLISHED_PORT ?= 5432
DB_NAME ?= bank
DB_USER ?= postgres
DB_PASSWORD ?= postgres

DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_PUBLISHED_HOST):$(DB_PUBLISHED_PORT)/$(DB_NAME)?sslmode=disable
MIGRATIONS_PATH=api/migrations
BOOK_PT_DIR=api/docs/visao_geral
TEMPLATE=templates/eisvogel.latex
DOCKER_COMPOSE=API_ENV_FILE=$(API_ENV_FILE) docker compose --env-file $(API_ENV_FILE) -f docker-compose.yml -p banklab-$(API_ENV)

# =========================
# Help
# =========================
help: ## List available commands
	@echo ""
	@echo "Available commands:"
	@echo ""
	@grep -hE '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'
	@echo ""

# =========================
# Docker
# =========================
docker-up: env-init ## Start Docker containers in detached mode
	$(DOCKER_COMPOSE) up -d --build

docker-db-up: env-init ## Start only the selected PostgreSQL container
	$(DOCKER_COMPOSE) up -d postgres

docker-down: env-init ## Stop and remove Docker containers
	$(DOCKER_COMPOSE) down

docker-logs: env-init ## Follow Docker container logs
	$(DOCKER_COMPOSE) logs -f

docker-clean: env-init ## Remove containers and volumes
	$(DOCKER_COMPOSE) down -v

docker-check: ## Check if Docker is running
	@if command -v colima > /dev/null 2>&1; then \
		echo "Colima detected. Ensuring daemon is running..."; \
		colima start; \
	fi
	@docker info > /dev/null 2>&1 || (echo "Docker is not running" && exit 1)

cloudflared-tunnel: ## Run the banklab Cloudflare tunnel
	cloudflared tunnel run banklab

# =========================
# Bootstrap and Reset
# =========================
setup: env-init docker-check docker-db-up db-wait migrate-up api-run ## Full setup from scratch

run: env-init mobile-sync-ip docker-check docker-db-up db-wait migrate-up api-run ## Start full development system

reset: env-init docker-check docker-clean docker-db-up db-wait db-reset migrate-up api-run ## Hard reset environment

db-reset: env-init ## Reset only the database
	$(DOCKER_COMPOSE) exec -T postgres psql -v ON_ERROR_STOP=1 -U $(DB_USER) -d postgres -c "DROP DATABASE IF EXISTS $(DB_NAME) WITH (FORCE);"
	$(DOCKER_COMPOSE) exec -T postgres psql -v ON_ERROR_STOP=1 -U $(DB_USER) -d postgres -c "CREATE DATABASE $(DB_NAME);"

db-wait: env-init ## Wait for the database to be ready
	@echo "Waiting for database..."
	@for i in $$(seq 1 30); do \
		if $(DOCKER_COMPOSE) exec -T postgres pg_isready -U $(DB_USER) > /dev/null 2>&1; then \
			echo "Database is ready"; \
			exit 0; \
		fi; \
		sleep 1; \
	done; \
	echo "Database not ready after timeout"; \
	exit 1

bootstrap: setup ## Full bootstrap from scratch
dev: run ## Alias para desenvolvimento

staging: ## Start the full staging environment
	$(MAKE) API_ENV=staging docker-check
	$(MAKE) API_ENV=staging docker-db-up
	$(MAKE) API_ENV=staging db-wait
	$(MAKE) API_ENV=staging migrate-up
	$(MAKE) API_ENV=staging api-run

prod: ## Start the full production environment
	$(MAKE) API_ENV=prod docker-check
	$(MAKE) API_ENV=prod docker-db-up
	$(MAKE) API_ENV=prod db-wait
	$(MAKE) API_ENV=prod migrate-up
	$(MAKE) API_ENV=prod api-run

env-init: ## Create API and Mobile .env files if they do not exist
	bash infra/scripts/ensure-env-files.sh

# =========================
# Monorepo
# =========================
build: api-build ## Build backend binary

# =========================
# API (Go)
# =========================
api-build: ## Build API binary into api/build/
	cd api && go build -o build/bank-api ./cmd/api

api-tests: ## Run API tests with coverage
	cd api && go test -cover ./...

api-kill-local: ## Stop any running local bank-api process
	@pids=$$(pgrep -f '(^|/)bank-api( |$$)|api/build/bank-api' || true); \
	if [ -n "$$pids" ]; then \
		echo "Stopping local bank-api process(es): $$pids"; \
		kill $$pids; \
	fi

api-run: env-init api-kill-local ## Build and run the selected API in Docker
	$(DOCKER_COMPOSE) up -d --build api

api-run-dev: ## Run API with api/dev.env
	$(MAKE) mobile-sync-ip
	$(MAKE) API_ENV=dev api-run

api-run-staging: ## Run API with api/staging.env
	$(MAKE) API_ENV=staging api-run

api-run-prod: ## Run API with api/prod.env
	$(MAKE) API_ENV=prod api-run

api-stop: env-init ## Stop the selected API container
	$(DOCKER_COMPOSE) stop api

# =========================
# Database Migrations
# =========================
migrate-up: env-init ## Run API database migrations
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

migrate-down: env-init ## Rollback last API database migration
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down

dbschema: env-init ## Export database schema to schema.sql
	$(DOCKER_COMPOSE) exec -T postgres pg_dump -U $(DB_USER) -d $(DB_NAME) > schema.sql

# =========================
# Mobile (Flutter)
# =========================
mobile-tests: ## Run all mobile tests
	cd mobile && flutter test

mobile-test-unit: ## Run mobile unit tests
	cd mobile && flutter test test/core

mobile-sync-ip: ## Update only mobile/dev.env BASE_URL with current host LAN IP
	bash infra/scripts/update-mobile-env-ip.sh

tests: ## Run all tests
	make api-tests mobile-tests

# =========================
# Git
# =========================
commit: ## Commit using predefined message file
	git add .
	git commit -F ~/commit.md

diff: ## Show staged diff and line count
	git add .
	git diff --cached > ~/diff
	wc -l ~/diff

push: ## Push current branch or specified branch (make push branch=xxx)
	@branch=$${branch:-$$(git branch --show-current)}; \
	if [ -z "$$branch" ]; then \
		echo "Branch not found"; \
		exit 1; \
	fi; \
	git push origin $$branch

pull: ## Pull current branch or specified branch (make pull branch=xxx)
	@branch=$${branch:-$$(git branch --show-current)}; \
	if [ -z "$$branch" ]; then \
		echo "Branch not found"; \
		exit 1; \
	fi; \
	git pull origin $$branch

gitlog: ## Show git log in one line format
	git log --oneline

pullmain: ## Pull latest changes from main branch
	git checkout main && git pull origin main

# =========================
# Flutter specific
# =========================
fclean: ## Clean Flutter build and get dependencies
	cd mobile && flutter clean && flutter pub get

fbuild: ## Build Flutter app for release
	cd mobile && flutter build apk --release

fadd: ## Add a Flutter package (make fadd pkg=package_name)
	cd mobile && flutter pub add $(pkg)

# =========================
# Documentation
# =========================
book-pt: ## Generate the Portuguese book PDF from markdown visao_geral
	pandoc \
		$(BOOK_PT_DIR)/chapters/*.md \
		--filter pandoc-mermaid \
		--metadata-file=metadata.yaml \
		--template=$(TEMPLATE) \
		--pdf-engine=xelatex \
		--toc \
		-V toc-depth=2 \
		-V book=true \
		-V titlepage=true \
		-V top-level-division=chapter \
		--resource-path=.:$(BOOK_PT_DIR) \
		-o banklab_BR.pdf
