.PHONY: help setup up down logs build build-dashboard test migrate dev-dashboard scale clean

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

setup: ## One-time setup (run before first `make up`)
	go mod tidy
	cd dashboard && npm install

up: ## Start all services (postgres + api + worker + migrate)
	docker compose up --build

down: ## Stop all services
	docker compose down

logs: ## Tail logs (optional: SERVICE=worker)
	docker compose logs -f $(SERVICE)

build-dashboard: ## Build the React dashboard into internal/web/static
	cd dashboard && npm install && npm run build
	rm -rf internal/web/static
	mkdir -p internal/web/static
	cp -r dashboard/dist/. internal/web/static/

build: build-dashboard ## Build the Go binary (with embedded dashboard)
	mkdir -p bin
	go build -o bin/jobqueue ./cmd/jobqueue

test: ## Run Go tests
	go test ./...

migrate: ## Apply migrations (against compose postgres)
	docker compose run --rm migrate migrate up

dev-dashboard: ## Run Vite dev server with HMR on :5173
	cd dashboard && npm run dev

scale: ## Scale workers (usage: make scale N=3)
	docker compose up --scale worker=$(N)

clean: ## Tear down + nuke postgres volume
	docker compose down -v
	rm -rf bin/ dashboard/dist/ internal/web/static/*
