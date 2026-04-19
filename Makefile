.PHONY: help \
	build test \
	api-build api-migrate-up api-migrate-down api-test \
	mobile-test mobile-test-unit \
	commit diff push pull gitlog

# =========================
# Variables
# =========================
DB_URL=postgres://postgres:postgres@localhost:5432/bank?sslmode=disable
MIGRATIONS_PATH=api/migrations

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
	docker compose up -d

docker-down: ## Stop and remove Docker containers
	docker compose down

docker-logs: ## Follow Docker container logs
	docker compose logs -f

docker-clean: ## Remove containers and volumes
	docker compose down -v

docker-check: ## Check if Docker is running
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

bootstrap: setup ## Alias semântico
dev: run ## Alias para desenvolvimento

# =========================
# Monorepo
# =========================
build: api-build ## Build backend binary

test: api-test mobile-test ## Run API and Mobile tests

# =========================
# API (Go)
# =========================
api-build: ## Build API binary into api/build/
	cd api && go build -o build/bank-api ./cmd/api

api-test: ## Run API tests with coverage
	cd api && go test -cover ./...

api-run: ## Run API server
	cd api && go run ./cmd/api

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
mobile-test: ## Run all mobile tests
	cd mobile && flutter test

mobile-test-unit: ## Run mobile unit tests
	cd mobile && flutter test test/core

tests: ## Run all tests
	make api-test mobile-test

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

# =========================
# Flutter specific
# =========================
fclean: ## Clean Flutter build and get dependencies
	cd mobile && flutter clean && flutter pub get

fbuild: ## Build Flutter app for release
	cd mobile && flutter build apk --release

fadd: ## Add a Flutter package (make fadd pkg=package_name)
	cd mobile && flutter pub add $(pkg)