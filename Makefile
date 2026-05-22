.PHONY: help \
	build test \
	api-build api-migrate-up api-migrate-down api-test api-stop env-init bootstrap \
	mobile-test mobile-test-unit mobile-sync-ip \
	commit diff push pull gitlog

# =========================
# Variables
# =========================
DB_URL=postgres://postgres:postgres@localhost:5432/bank?sslmode=disable
MIGRATIONS_PATH=api/migrations
BOOK_PT_DIR=api/docs/visao_geral
TEMPLATE=templates/eisvogel.latex

# =========================
# Help
# =========================
help: ## List available commands
	@echo ""
	@echo "Available commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'
	@echo ""

# =========================
# Docker
# =========================
docker-up: ## Start Docker containers in detached mode
	docker compose up -d --no-recreate

docker-down: ## Stop and remove Docker containers
	docker compose down

docker-logs: ## Follow Docker container logs
	docker compose logs -f

docker-clean: ## Remove containers and volumes
	docker compose down -v

docker-check: ## Check if Docker is running
	@if command -v colima > /dev/null 2>&1; then \
		echo "Colima detected. Ensuring daemon is running..."; \
		colima start; \
	fi
	@docker info > /dev/null 2>&1 || (echo "Docker is not running" && exit 1)

# =========================
# Bootstrap and Reset
# =========================
setup: docker-check docker-up db-wait migrate-up ## Full setup from scratch

run: docker-check docker-up db-wait migrate-up api-run ## Start full system

reset: docker-check docker-clean docker-up db-wait db-reset migrate-up ## Hard reset environment

db-reset: ## Reset only the database
	docker exec -i bank-postgres psql -U postgres -c "DROP DATABASE IF EXISTS bank;"
	docker exec -i bank-postgres psql -U postgres -c "CREATE DATABASE bank;"

db-wait: ## Wait for the database to be ready
	@echo "Waiting for database..."
	@for i in $$(seq 1 30); do \
		if docker exec bank-postgres pg_isready -U postgres > /dev/null 2>&1; then \
			echo "Database is ready"; \
			exit 0; \
		fi; \
		sleep 1; \
	done; \
	echo "Database not ready after timeout"; \
	exit 1

bootstrap: env-init docker-up db-wait migrate-up api-run ## Full bootstrap from scratch
dev: run ## Alias para desenvolvimento

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

api-run: mobile-sync-ip ## Run API server
	cd api && go run ./cmd/api

api-stop: ## Stop API server running on port 8080
	@pid=$$(lsof -ti tcp:8080); \
	if [ -z "$$pid" ]; then \
		echo "API is not running on port 8080"; \
	else \
		echo "Stopping API process $$pid"; \
		kill $$pid; \
	fi

# =========================
# Database Migrations
# =========================
migrate-up: ## Run API database migrations
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

migrate-down: ## Rollback last API database migration
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down

dbschema: ## Export database schema to schema.sql
	docker exec -t bank-postgres pg_dump -U postgres -d bank > schema.sql

# =========================
# Mobile (Flutter)
# =========================
mobile-tests: ## Run all mobile tests
	cd mobile && flutter test

mobile-test-unit: ## Run mobile unit tests
	cd mobile && flutter test test/core

mobile-sync-ip: ## Update mobile .env BASE_URL with current host LAN IP
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