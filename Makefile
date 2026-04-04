.PHONY: help migration commit diff push pull

# =========================
# Variables
# =========================
DB_URL=postgres://postgres:postgres@localhost:5432/bank?sslmode=disable
MIGRATIONS_PATH=migrations

# =========================
# Help
# =========================
help:
	@echo ""
	@echo "Available commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  %-12s %s\n", $$1, $$2}'
	@echo ""

# =========================
# Database
# =========================
migration: ## Run database migrations
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

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

tests: ## Run tests with coverage
	go test -cover ./...