# ==============================================================================
# SkillNODE Makefile - Developer Convenience Commands
# ==============================================================================

.PHONY: help install setup dev dev-client dev-server build build-client build-server \
        db-up db-down db-logs db-generate db-push db-studio \
        lint test audit ci clean clean-all rebase push

# Default target when running `make`
.DEFAULT_GOAL := help

# ------------------------------------------------------------------------------
# HELP / DOCUMENTATION
# ------------------------------------------------------------------------------
help:
	@echo "================================================================================"
	@echo " SkillNODE Development & Automation Commands"
	@echo "================================================================================"
	@echo "Setup & Dependencies:"
	@echo "  make install        - Install root and workspace npm dependencies"
	@echo "  make setup          - Full project setup (install + prisma client generation)"
	@echo ""
	@echo "Development:"
	@echo "  make dev-client     - Launch frontend Go dev server (Go SSG + WebAssembly)"
	@echo "  make dev-server     - Launch backend Node dev server (Nodemon + Express)"
	@echo ""
	@echo "Build:"
	@echo "  make build          - Build client and verify server syntax"
	@echo "  make build-client   - Build client production bundle (Go SSG + WebAssembly)"
	@echo "  make build-server   - Validate server syntax and modules"
	@echo ""
	@echo "Database & Docker:"
	@echo "  make db-up          - Start PostgreSQL container in background (docker-compose)"
	@echo "  make db-down        - Stop PostgreSQL container"
	@echo "  make db-logs        - Stream logs from PostgreSQL container"
	@echo "  make db-generate    - Generate Prisma Client for server"
	@echo "  make db-push        - Push Prisma schema migrations to database"
	@echo "  make db-studio      - Open Prisma Studio GUI"
	@echo ""
	@echo "Quality & Security:"
	@echo "  make test           - Run project test suites"
	@echo "  make lint           - Run ESLint checks"
	@echo "  make audit          - Run npm security vulnerability audit"
	@echo "  make ci             - Run complete local CI check (build + audit)"
	@echo ""
	@echo "Git & Automation:"
	@echo "  make rebase         - Safe git pull --rebase on origin/main"
	@echo "  make push msg=\"...\" - Auto rebase, stage all changes, commit, and push"
	@echo ""
	@echo "Cleanup:"
	@echo "  make clean          - Remove build artifacts and caches"
	@echo "  make clean-all      - Remove build artifacts and node_modules"
	@echo "================================================================================"

# ------------------------------------------------------------------------------
# SETUP & DEPENDENCIES
# ------------------------------------------------------------------------------
install:
	npm ci

setup: install db-generate

# ------------------------------------------------------------------------------
# DEVELOPMENT
# ------------------------------------------------------------------------------
dev-client:
	cd client && go run build.go -serve

dev-server:
	npm run dev:server

# ------------------------------------------------------------------------------
# BUILD
# ------------------------------------------------------------------------------
build: build-client

build-client:
	cd client && go run build.go

build-server:
	npm run build -w server

# ------------------------------------------------------------------------------
# DATABASE & DOCKER
# ------------------------------------------------------------------------------
db-up:
	docker compose up -d postgres

db-down:
	docker compose down

db-logs:
	docker compose logs -f postgres

db-generate:
	npm run prisma:generate -w server

db-push:
	npm run db:push

db-studio:
	npx prisma studio --schema=server/prisma/schema.prisma

# ------------------------------------------------------------------------------
# TESTING & SECURITY
# ------------------------------------------------------------------------------
test:
	npm test -w client

lint:
	npx eslint .

audit:
	npm audit --audit-level=high

ci: build audit

# ------------------------------------------------------------------------------
# CLEANUP
# ------------------------------------------------------------------------------
clean:
	@echo "Cleaning build artifacts..."
	@if exist client\dist rmdir /s /q client\dist
	@if exist client\.vite rmdir /s /q client\.vite
	@echo "Done."

clean-all: clean
	@echo "Removing node_modules..."
	@if exist node_modules rmdir /s /q node_modules
	@if exist client\node_modules rmdir /s /q client\node_modules
	@if exist server\node_modules rmdir /s /q server\node_modules
	@echo "Done."

# ------------------------------------------------------------------------------
# GIT & AUTOMATION WORKFLOWS
# ------------------------------------------------------------------------------
MSG ?= $(msg)
ifeq ($(MSG),)
  MSG := $(m)
endif

rebase:
	@echo "==> Pulling and rebasing on origin/main (with autostash)..."
	git pull --rebase --autostash origin main

push:
ifeq ($(strip $(MSG)),)
	@echo "Error: Commit message required. Usage: make push msg=\"your message\" or make push m=\"your message\""
	@exit 1
else
	@echo "==> Executing automated safe rebase, stage, commit and push..."
	git pull --rebase --autostash origin main
	git add -A
	git commit -m "$(MSG)"
	git push origin HEAD
endif

