.PHONY: setup dev api web build deploy

HOST ?= stan@34.55.165.112
SITE ?= https://stanmao.com

setup:
	cd api && go mod tidy
	cd web && pnpm install

api:
	cd api && go run ./cmd/server

web:
	cd web && pnpm dev

dev:
	$(MAKE) -j2 api web

build:
	cd api && go build -o bin/server ./cmd/server
	cd web && pnpm build

deploy:
	cd web && pnpm build
	cd api && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -o bin/linux-amd64/server ./cmd/server
	ssh $(HOST) 'mkdir -p apps/dashboard/bin apps/nginx/html/dashboard'
	rsync -az --delete-after --exclude .DS_Store web/dist/ $(HOST):apps/nginx/html/dashboard/
	rsync -az deploy/compose.yaml $(HOST):apps/dashboard/compose.yaml
	rsync -az api/bin/linux-amd64/server $(HOST):apps/dashboard/bin/server
	ssh $(HOST) 'cd apps/dashboard && docker compose up -d --force-recreate'
	curl -fsS --retry 5 --retry-all-errors --retry-delay 2 $(SITE)/api/health
