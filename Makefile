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
# Monorepo
# =========================
build: api-build ## Build backend binary

test: api-test mobile-test ## Run API and Mobile tests

# =========================
# API (Go)
# =========================
api-build: ## Build API binary into api/build/
	cd api && go build -o build/bank-api ./cmd/api

api-migrate-up: ## Run API database migrations
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

api-migrate-down: ## Rollback last API database migration
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down

api-test: ## Run API tests with coverage
	cd api && go test -cover ./...

# =========================
# Mobile (Flutter)
# =========================
mobile-test: ## Run all mobile tests
	cd mobile && flutter test

mobile-test-unit: ## Run mobile unit tests
	cd mobile && flutter test test/core

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
