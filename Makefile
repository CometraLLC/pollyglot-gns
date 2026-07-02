# Pollyglot — repo-wide developer commands
# Backend-specific targets live in backend/Makefile.

.PHONY: check check-backend check-frontend dev dev-down test help

## Run every gate CI runs: backend lint+tests, frontend lint+tests+build
check: check-backend check-frontend

check-backend:
	cd backend && golangci-lint run ./... && go test ./... -race -count=1

check-frontend:
	cd frontend && bun run lint && bun run test && bun run build

## Run all tests only (no lint/build)
test:
	cd backend && go test ./... -race -count=1
	cd frontend && bun run test

## Bring up the dev stack (Postgres, Redis, API on :8080) and Next on :3000
dev:
	cd backend && docker compose -f docker-compose.dev.yaml up --build -d
	cd frontend && bun run dev

dev-down:
	cd backend && docker compose -f docker-compose.dev.yaml down

help:
	@echo "make check          - run every gate CI runs (lint + tests + build)"
	@echo "make check-backend  - backend lint + tests"
	@echo "make check-frontend - frontend lint + tests + build"
	@echo "make test           - all tests, no lint/build"
	@echo "make dev            - docker stack + Next dev server"
	@echo "make dev-down       - stop the docker stack"
