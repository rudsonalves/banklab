.PHONY: help \
	build test \
	api-build api-migrate-up api-migrate-down api-test api-stop \
	api-port-clean \
	api-run api-run-dev api-run-staging env-init bootstrap run-dev run-staging run-prod \
	dev staging prod \
	port-check port-clean \
	cloudflared-tunnel \
	mobile-test mobile-test-unit mobile-sync-ip \
	commit diff push pull gitlog

# =========================
# Variables
# =========================
API_ENV ?= dev
API_ENV_FILE=api/$(API_ENV).env
-include $(API_ENV_FILE)

DB_HOST ?= 127.0.0.1
DB_PORT ?= 5432
DB_NAME ?= bank
DB_USER ?= postgres
DB_PASSWORD ?= postgres
SERVER_PORT ?= 8080

DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
MIGRATIONS_PATH=api/migrations
BOOK_PT_DIR=api/docs/visao_geral
TEMPLATE=templates/eisvogel.latex
DOCKER_COMPOSE=docker compose --env-file $(API_ENV_FILE) -f docker-compose.yml -p banklab-$(API_ENV)

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
docker-up: env-init port-check ## Start the selected PostgreSQL container
	$(DOCKER_COMPOSE) up -d --build postgres

port-check: ## Fail early if DB_PORT is already in use by another process/container
	@project_pg_container="banklab-$(API_ENV)-postgres-1"; \
	if docker ps --format '{{.Names}}' | grep -qx "$$project_pg_container"; then \
		exit 0; \
	fi; \
	$(MAKE) --no-print-directory port-clean; \
	if lsof -nP -iTCP:$(DB_PORT) -sTCP:LISTEN > /dev/null 2>&1; then \
		echo "Port $(DB_PORT) is already allocated."; \
		echo "Listeners:"; \
		lsof -nP -iTCP:$(DB_PORT) -sTCP:LISTEN; \
		echo ""; \
		echo "Tip: stop the conflicting non-Docker process or change DB_PORT in api/$(API_ENV).env"; \
		exit 1; \
	fi

port-clean: ## Stop conflicting Docker containers bound to DB_PORT
	@conflict_containers=$$(docker ps --format '{{.Names}}\t{{.Ports}}' | awk -F'\t' -v p='$(DB_PORT)' '$$2 ~ ":"p"->" {print $$1}'); \
	if [ -n "$$conflict_containers" ]; then \
		echo "Port $(DB_PORT) is in use by Docker container(s). Running automatic cleanup:"; \
		echo "$$conflict_containers" | sed 's/^/ - /'; \
		echo "$$conflict_containers" | xargs docker stop > /dev/null; \
		echo "Cleanup done."; \
	fi

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
setup: env-init docker-check docker-up db-wait migrate-up api-run ## Full setup from scratch

run-dev: env-init mobile-sync-ip docker-check docker-up db-wait migrate-up api-run ## Start full development system

run-staging: ## Start the full staging environment
	$(MAKE) API_ENV=staging run-dev

run-prod: ## Start the full production environment
	$(MAKE) API_ENV=prod run-dev

reset: env-init docker-check docker-clean docker-up db-wait db-reset migrate-up api-run ## Hard reset environment

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
dev: run-dev ## Alias para desenvolvimento

staging: ## Start the full staging environment
	$(MAKE) run-staging

prod: ## Start the full production environment
	$(MAKE) run-prod

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

api-run: env-init api-port-clean ## Run the selected API on the host
	cd api && ENV_FILE=$(API_ENV).env go run ./cmd/api

api-port-clean: ## Stop stale API process if SERVER_PORT is already occupied
	@pid=$$(lsof -ti tcp:$(SERVER_PORT) -sTCP:LISTEN | head -n1); \
	if [ -n "$$pid" ]; then \
		comm=$$(ps -p $$pid -o comm= | tr -d '[:space:]'); \
		if [ "$$comm" = "api" ] || [ "$$comm" = "go" ] || echo "$$comm" | grep -Eq '(^|/)api$$|go-build'; then \
			echo "Port $(SERVER_PORT) is occupied by stale API process (pid=$$pid, comm=$$comm). Cleaning up..."; \
			kill $$pid; \
		else \
			echo "Port $(SERVER_PORT) is occupied by another process (pid=$$pid, comm=$$comm)."; \
			echo "Stop it manually or change SERVER_PORT in api/$(API_ENV).env"; \
			exit 1; \
		fi; \
	fi

api-run-dev: ## Run API with api/dev.env
	$(MAKE) mobile-sync-ip
	$(MAKE) API_ENV=dev api-run

api-run-staging: ## Run API with api/staging.env
	$(MAKE) API_ENV=staging api-run

api-stop: ## Stop the API listening on SERVER_PORT
	@if command -v fuser > /dev/null 2>&1; then \
		fuser -k $(SERVER_PORT)/tcp; \
	elif command -v lsof > /dev/null 2>&1; then \
		pid=$$(lsof -ti tcp:$(SERVER_PORT)); \
		[ -z "$$pid" ] || kill $$pid; \
	else \
		echo "Install psmisc (fuser) or lsof to use api-stop"; \
		exit 1; \
	fi

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
